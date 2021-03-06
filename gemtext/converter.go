package gemtext

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/fhofherr/mnml/internal/agmi"
)

// FromAlmostGemtext creates a Gemtext document of the Almost Gemtext
// document read from in and writes it to out.
func FromAlmostGemtext(in io.Reader, out io.Writer) error {
	const op = "gemtext/FromAlmostGemtext"

	c := agmi.NewConverter(in, out, fmtAGMIToken)
	if err := c.Convert(); err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}
	return nil
}

func fmtAGMIToken(c *agmi.Converter, cur, next agmi.Token) {
	switch cur.Type {
	case agmi.TokenTypeModeline:
		c.State = skipEmptyLines
	case agmi.TokenTypeQuoteMod:
		c.Write("> ")
		c.State = fmtQuoteLines
	case agmi.TokenTypePreFmtMod:
		c.Write("```")
		c.State = fmtPreFmt
	case agmi.TokenTypeIndent:
		if isSpaceIndent(cur.Text) && len(cur.Text) != 4 || (isTabIndent(cur.Text) && len(cur.Text) != 1) {
			// Cannot be pre-formatted text since it would need to be indented
			// by a single tab or exactly four spaces.
			c.Write(cur.Text)
			return
		}
		c.Write("```\n")
		c.State = fmtPreFmtByIndent
	case agmi.TokenTypeBulletPoint:
		c.Write("* ")
		c.State = fmtListItem
	case agmi.TokenTypeLineBreak:
		joinLines(c, next)
	default:
		c.Write(cur.Text)
	}
}

func fmtListItem(c *agmi.Converter, cur, next agmi.Token) {
	switch cur.Type {
	case agmi.TokenTypeIndent:
		if len(cur.Text) > 2 {
			c.Write(cur.Text[2:])
		}
	case agmi.TokenTypeParSep:
		// End of list
		c.State = fmtAGMIToken
	case agmi.TokenTypeLineBreak:
		if next.Type == agmi.TokenTypeBulletPoint {
			// Another list item is directly following the current one. Treat
			// it as end of list. The fmtAGMIToken function will know how to handle
			// it and delegate back to here.
			c.State = fmtAGMIToken
			c.Write("\n")
			return
		}
		joinLines(c, next)
	default:
		c.Write(cur.Text)
	}
}

func fmtPreFmtByIndent(c *agmi.Converter, cur, next agmi.Token) {
	if cur.Type == agmi.TokenTypeIndent {
		// Skip leading indent if it is only four spaces or a single tab.
		// Otherwise reduce it by four spaces or a single tab and write it
		// to the output.
		if isSpaceIndent(cur.Text) && utf8.RuneCountInString(cur.Text) > 4 {
			c.Write(strings.TrimPrefix(cur.Text, "    "))
		}
		if isTabIndent(cur.Text) && utf8.RuneCountInString(cur.Text) > 1 {
			c.Write(strings.TrimPrefix(cur.Text, "\t"))
		}
		return
	}
	if next.IsZero() || (cur.Type == agmi.TokenTypeParSep && next.Type != agmi.TokenTypeIndent) {
		// We reached the end of the pre-formatted block.
		c.State = fmtAGMIToken
		c.Write("\n```")
	}
	c.Write(cur.Text)
}

func fmtPreFmt(c *agmi.Converter, cur, _ agmi.Token) {
	c.Write(cur.Text)
	if cur.Type == agmi.TokenTypePreFmtMod {
		c.State = fmtAGMIToken
	}
}

func fmtQuoteLines(c *agmi.Converter, cur, next agmi.Token) {
	switch cur.Type {
	case agmi.TokenTypeQuoteMod:
		// The initial > was already written by fmtAGMIToken. Skip any
		// further > in the input.
		return
	case agmi.TokenTypeLineBreak:
		joinLines(c, next)
	case agmi.TokenTypeParSep:
		// We reached the end ouf our multi line quote.
		c.Write(cur.Text)
		c.State = fmtAGMIToken
	default:
		// Anything else is part of the quote's text.
		c.Write(cur.Text)
	}
}

func skipEmptyLines(c *agmi.Converter, cur, next agmi.Token) {
	if next.Type != agmi.TokenTypeLineBreak && next.Type != agmi.TokenTypeParSep {
		c.State = fmtAGMIToken // Set next state.
	}
}

func joinLines(c *agmi.Converter, next agmi.Token) {
	if next.IsZero() {
		c.Write("\n")
		return
	}
	c.Write(" ")
}

func isTabIndent(s string) bool {
	return strings.HasPrefix(s, "\t")
}

func isSpaceIndent(s string) bool {
	return strings.HasPrefix(s, " ")
}
