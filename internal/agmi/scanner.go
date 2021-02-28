package agmi

import (
	"bufio"
	"bytes"
	"io"
)

// TokenType defines the type of an Almost Gemtext Token.
type TokenType int

const (
	tokenTypeUnknown TokenType = iota

	// TokenTypeModeline marks the Token as a modeline.
	TokenTypeModeline

	// TokenTypeLineBreak identifies the token as a simple line break.
	TokenTypeLineBreak

	// TokenTypeParSep identifies the token as a separator of paragraphs.
	TokenTypeParSep

	// TokenTypeQuoteMod identifies the next line of text as being a quote or
	// part of a quote.
	//
	// Consecutive occurrences of TokenTypeQuoteMod are to be treated as
	// interspersed with empty lines of text.
	TokenTypeQuoteMod

	// TokenTypePreFmtMod signals either the beginning or the end of
	// pre-formatted text.
	TokenTypePreFmtMod

	// TokenTypeLinkMod signals the beginning of a link.
	TokenTypeLinkMod

	// TokenTypeLinkURI is the URI part of a link.
	TokenTypeLinkURI

	// TokenTypeBulletPoint signals that the token is the bullet point
	// belonging to a list item. The immediately following Tokens of
	// TokenTypeLine and TokenTypeIndent constitute the list item for this
	// bullet point.
	TokenTypeBulletPoint

	// TokenTypeIndent signals that the Token's text was used to indent the
	// text of the following token.
	TokenTypeIndent

	// TokenTypeText marks the token as plain text. Tokens that occurred
	// earlier may have an influence on how the line is treated.
	TokenTypeText
)

//go:generate stringer -type TokenType -trimprefix TokenType -output tokentype.string.go

// Token represents a token of an Almost Gemtext document.
type Token struct {
	Type TokenType // Type of the token
	Text string    // Token text as read from the input.
}

// Scanner scans a single Almost Gemtext document and creates a stream of
// tokens for further processing.
//
// The zero value of Scanner is not usable. Use NewScanner to obtain a working
// instance.
type Scanner struct {
	scanner *bufio.Scanner
	state   bufio.SplitFunc
	token   Token
}

// NewScanner creates a new scanner for an Almost Gemtext document.
//
// The document is read from r.
func NewScanner(r io.Reader) *Scanner {
	var sc Scanner

	sc.scanner = bufio.NewScanner(r)
	sc.scanner.Split(sc.splitFunc)
	sc.state = sc.scanLine

	return &sc
}

// Scan advances the input until it was able to read a Token or an error occurs.
//
// Scan returns false if the scan stops. This may be the case if an error
// occurred or if the end of the input was reached. In the case of an error the
// Err method contains the error that occurred.
func (sc *Scanner) Scan() bool {
	sc.token = Token{} // reset token detected by last scan
	return sc.scanner.Scan()
}

// Token returns the Token read during the previous call to Scan.
func (sc *Scanner) Token() Token {
	sc.token.Text = sc.scanner.Text()
	return sc.token
}

// Err returns any error that occurred during the last call to scan.
//
// If the input was completely consumed by Scan Err returns nil instead of
// io.EOF.
func (sc *Scanner) Err() error {
	return sc.scanner.Err()
}

func (sc *Scanner) splitFunc(data []byte, atEOF bool) (int, []byte, error) {
	if len(data) == 0 {
		return 0, nil, nil
	}
	return sc.state(data, atEOF)
}

func (sc *Scanner) scanLine(data []byte, atEOF bool) (int, []byte, error) {
	switch data[0] {
	case '<':
		return sc.goToState(sc.scanModeLine, data, atEOF)
	case '>':
		return sc.goToState(sc.scanQuote, data, atEOF)
	case '\n':
		return sc.goToState(sc.scanParSep, data, atEOF)
	case '`':
		return sc.goToState(sc.scanFmtMod, data, atEOF)
	case ' ', '\t':
		return sc.goToState(sc.scanIndent, data, atEOF)
	case '*':
		return sc.goToState(sc.scanBulletPoint, data, atEOF)
	case '=':
		return sc.goToState(sc.scanLinkMod, data, atEOF)
	default:
		// Cannot decide on the type of line. Treat it as text.
		return sc.goToState(sc.scanText, data, atEOF)
	}
}

func (sc *Scanner) scanModeLine(data []byte, atEOF bool) (int, []byte, error) {
	if len(data) < 4 {
		if atEOF {
			// The current line is not a modeline, since it does not contain
			// enough data for an opening modeline identifier and there is no more
			// input left. Treat it as arbitrary text.
			return sc.goToState(sc.scanText, data, atEOF)
		}
		return 0, nil, nil
	}
	if string(data[0:4]) != "<!--" {
		// The line did not start with a modeline identifier. It must be an
		// arbitrary text
		return sc.goToState(sc.scanText, data, atEOF)
	}
	sc.tokenFound(TokenTypeModeline, sc.scanLine)
	return sc.scanFunc(data, atEOF, isVSpace)
}

func (sc *Scanner) scanParSep(data []byte, atEOF bool) (int, []byte, error) {
	// Assume we are dealing with a paragraph separator and find the index
	// of the first rune which is not a line break
	sc.tokenFound(TokenTypeParSep, sc.scanLine)
	i, tok, err := sc.scanFunc(data, atEOF, notIsRune('\n'))
	if err != nil {
		// Not wrapping the error is ok. We are only interested in the error
		// returned by the scanner.
		return i, tok, err
	}
	if i == 0 {
		// Read more data
		return 0, nil, nil
	}
	if i == 1 {
		// Our assumption was wrong. The input contained only a normal line
		// break.
		sc.tokenFound(TokenTypeLineBreak, sc.scanLine)
	}
	return i, tok, nil
}

func (sc *Scanner) scanText(data []byte, atEOF bool) (int, []byte, error) {
	sc.tokenFound(TokenTypeText, sc.scanLine)
	return sc.scanFunc(data, atEOF, isVSpace)
}

func (sc *Scanner) scanQuote(data []byte, atEOF bool) (int, []byte, error) {
	sc.tokenFound(TokenTypeQuoteMod, sc.scanLine)
	i, tok, err := sc.scanFunc(data, atEOF, notIsHSpace)
	return i, tok, err
}
func (sc *Scanner) scanFmtMod(data []byte, atEOF bool) (int, []byte, error) {
	if len(data) < 3 {
		if atEOF {
			// This can't be a format modifier. Read it as normal line.
			return sc.goToState(sc.scanText, data, atEOF)
		}
		// Read more data
		return 0, nil, nil
	}
	if string(data[0:3]) != "```" {
		// Not a format modifier. Read it as normal line.
		return sc.goToState(sc.scanText, data, atEOF)
	}
	sc.tokenFound(TokenTypePreFmtMod, sc.scanLine)
	return 3, data[0:3], nil
}

func (sc *Scanner) scanIndent(data []byte, atEOF bool) (int, []byte, error) {
	sc.tokenFound(TokenTypeIndent, sc.scanLine)
	return sc.scanFunc(data, atEOF, notIsHSpace)
}

func (sc *Scanner) scanBulletPoint(data []byte, atEOF bool) (int, []byte, error) {
	sc.tokenFound(TokenTypeBulletPoint, sc.scanLine)
	i, tok, err := sc.scanFunc(data, atEOF, notIsHSpace)
	return i, tok, err
}

func (sc *Scanner) scanLinkMod(data []byte, atEOF bool) (int, []byte, error) {
	if len(data) < 2 {
		if atEOF {
			// Not enough input left. This can't be be a link modifier.
			return sc.goToState(sc.scanText, data, atEOF)
		}
		return 0, nil, nil // read more data
	}
	if data[1] != '>' {
		// Not a link modifier
		return sc.goToState(sc.scanText, data, atEOF)
	}
	if atEOF {
		return 2, data[0:2], nil
	}
	i, _, err := sc.scanFunc(data[2:], atEOF, notIsHSpace)
	if err != nil {
		return 0, nil, err
	}
	if i == 0 {
		return 0, nil, nil
	}
	sc.tokenFound(TokenTypeLinkMod, sc.scanLinkURI)
	return i + 2, data[0 : i+2], nil // We offset our data by 2
}

func (sc *Scanner) scanLinkURI(data []byte, atEOF bool) (int, []byte, error) {
	i, tok, err := sc.scanFunc(data, atEOF, isHSpace)
	if err != nil {
		return 0, nil, err
	}
	if i == 0 {
		return 0, nil, nil // Read more data
	}
	sc.tokenFound(TokenTypeLinkURI, sc.scanLinkText)
	return i, tok, nil
}

func (sc *Scanner) scanLinkText(data []byte, atEOF bool) (int, []byte, error) {
	i, tok, err := sc.scanFunc(data, atEOF, isVSpace)
	if err != nil {
		return 0, nil, err
	}
	if i == 0 {
		return 0, nil, nil // Read more data
	}
	sc.tokenFound(TokenTypeText, sc.scanLine)
	return i, tok, nil
}

func (sc *Scanner) scanFunc(data []byte, atEOF bool, f func(r rune) bool) (int, []byte, error) {
	if len(data) == 1 && atEOF {
		// Return what we have got.
		return 1, data[0:1], nil
	}
	// Find the index of the first rune that matches f, excluding data[0].
	// Add 1 since we offset our data by 1.
	i := bytes.IndexFunc(data[1:], f) + 1
	if i == 0 {
		if !atEOF {
			return 0, nil, nil // Read more data
		}
		i = len(data)
	}
	return i, data[0:i], nil
}

func (sc *Scanner) goToState(state bufio.SplitFunc, data []byte, atEOF bool) (int, []byte, error) {
	sc.state = state
	return sc.state(data, atEOF)
}

// tokenFound sets the type of token that was found during the current call
// to Scan. nextState is the state the next call to Scan will start with.
func (sc *Scanner) tokenFound(typ TokenType, nextState bufio.SplitFunc) {
	sc.token.Type = typ
	sc.state = nextState
}
