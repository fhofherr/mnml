// Code generated by "stringer -type TokenType -trimprefix TokenType -output tokentype.string.go"; DO NOT EDIT.

package agmi

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[tokenTypeUnknown-0]
	_ = x[TokenTypeModeline-1]
	_ = x[TokenTypeLineBreak-2]
	_ = x[TokenTypeParSep-3]
	_ = x[TokenTypeQuoteMod-4]
	_ = x[TokenTypePreFmtMod-5]
	_ = x[TokenTypeLinkMod-6]
	_ = x[TokenTypeLinkURI-7]
	_ = x[TokenTypeBulletPoint-8]
	_ = x[TokenTypeIndent-9]
	_ = x[TokenTypeText-10]
}

const _TokenType_name = "tokenTypeUnknownModelineLineBreakParSepQuoteModPreFmtModLinkModLinkURIBulletPointIndentText"

var _TokenType_index = [...]uint8{0, 16, 24, 33, 39, 47, 56, 63, 70, 81, 87, 91}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
