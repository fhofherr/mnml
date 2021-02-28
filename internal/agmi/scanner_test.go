package agmi_test

import (
	"bytes"
	"testing"

	"github.com/fhofherr/mnml/internal/agmi"
	"github.com/stretchr/testify/assert"
)

func TestScanner_Scan(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []agmi.Token
	}{
		{
			name:  "Scan modeline with close comment symbol",
			input: "<!-- vim: set tw=72 ft=markdown: -->\n",
			expected: []agmi.Token{
				{
					Type: agmi.TokenTypeModeline,
					Text: "<!-- vim: set tw=72 ft=markdown: -->",
				},
				{
					Type: agmi.TokenTypeLineBreak,
					Text: "\n",
				},
			},
		},
		{
			name:  "Scan modeline without close comment symbol",
			input: "<!-- vim: set tw=72 ft=markdown:\n",
			expected: []agmi.Token{
				{
					Type: agmi.TokenTypeModeline,
					Text: "<!-- vim: set tw=72 ft=markdown:",
				},
				{
					Type: agmi.TokenTypeLineBreak,
					Text: "\n",
				},
			},
		},
		// TODO Headings: currently require no special treatment. We can delay that.
		{
			name:  "Scan two paragraphs",
			input: "The first paragraph.\n\nThe second paragraph.",
			expected: []agmi.Token{
				{
					Type: agmi.TokenTypeText,
					Text: "The first paragraph.",
				},
				{
					Type: agmi.TokenTypeParSep,
					Text: "\n\n",
				},
				{
					Type: agmi.TokenTypeText,
					Text: "The second paragraph.",
				},
			},
		},
		{
			name:  "Scan single line of arbitrary text with trailing newline",
			input: "This is a line of text.\n",
			expected: []agmi.Token{
				{
					Type: agmi.TokenTypeText,
					Text: "This is a line of text.",
				},
				{
					Type: agmi.TokenTypeLineBreak,
					Text: "\n",
				},
			},
		},
		{
			name:  "Scan single line of arbitrary text without trailing newline",
			input: "This is a line of text.",
			expected: []agmi.Token{
				{
					Type: agmi.TokenTypeText,
					Text: "This is a line of text.",
				},
			},
		},
		{
			name:  "Scan two consecutive lines of arbitrary text",
			input: "This is a line of text.\nThis is another one.\n",
			expected: []agmi.Token{
				{
					Type: agmi.TokenTypeText,
					Text: "This is a line of text.",
				},
				{
					Type: agmi.TokenTypeLineBreak,
					Text: "\n",
				},
				{
					Type: agmi.TokenTypeText,
					Text: "This is another one.",
				},
				{
					Type: agmi.TokenTypeLineBreak,
					Text: "\n",
				},
			},
		},
		{
			name:  "Single quote line",
			input: "> This is a quote",
			expected: []agmi.Token{
				{
					Type: agmi.TokenTypeQuoteMod,
					Text: "> ",
				},
				{
					Type: agmi.TokenTypeText,
					Text: "This is a quote",
				},
			},
		},
		{
			name:  "Multiple quote lines",
			input: "> First line of quote.\n> Second line of quote.",
			expected: []agmi.Token{
				{
					Type: agmi.TokenTypeQuoteMod,
					Text: "> ",
				},
				{
					Type: agmi.TokenTypeText,
					Text: "First line of quote.",
				},
				{
					Type: agmi.TokenTypeLineBreak,
					Text: "\n",
				},
				{
					Type: agmi.TokenTypeQuoteMod,
					Text: "> ",
				},
				{
					Type: agmi.TokenTypeText,
					Text: "Second line of quote.",
				},
			},
		},
		{
			name:  "Multi paragraph quote",
			input: "> First line of quote.\n>\n> Second line of quote.",
			expected: []agmi.Token{
				{
					Type: agmi.TokenTypeQuoteMod,
					Text: "> ",
				},
				{
					Type: agmi.TokenTypeText,
					Text: "First line of quote.",
				},
				{
					Type: agmi.TokenTypeLineBreak,
					Text: "\n",
				},
				{
					Type: agmi.TokenTypeQuoteMod,
					Text: ">", // Note the missing space!
				},
				{
					Type: agmi.TokenTypeLineBreak,
					Text: "\n",
				},
				{
					Type: agmi.TokenTypeQuoteMod,
					Text: "> ",
				},
				{
					Type: agmi.TokenTypeText,
					Text: "Second line of quote.",
				},
			},
		},
		{
			name:  "Pre-formatted text with ```",
			input: "```\nThis is pre-formatted.\nThis as well.\n```",
			expected: []agmi.Token{
				{
					Type: agmi.TokenTypePreFmtMod,
					Text: "```",
				},
				{
					Type: agmi.TokenTypeLineBreak,
					Text: "\n",
				},
				{
					Type: agmi.TokenTypeText,
					Text: "This is pre-formatted.",
				},
				{
					Type: agmi.TokenTypeLineBreak,
					Text: "\n",
				},
				{
					Type: agmi.TokenTypeText,
					Text: "This as well.",
				},
				{
					Type: agmi.TokenTypeLineBreak,
					Text: "\n",
				},
				{
					Type: agmi.TokenTypePreFmtMod,
					Text: "```",
				},
			},
		},
		{
			name:  "Pre-formatted with space indent",
			input: "    This is pre-formatted.\n    This as well.\n",
			expected: []agmi.Token{
				{
					Type: agmi.TokenTypeIndent,
					Text: "    ",
				},
				{
					Type: agmi.TokenTypeText,
					Text: "This is pre-formatted.",
				},
				{
					Type: agmi.TokenTypeLineBreak,
					Text: "\n",
				},
				{
					Type: agmi.TokenTypeIndent,
					Text: "    ",
				},
				{
					Type: agmi.TokenTypeText,
					Text: "This as well.",
				},
				{
					Type: agmi.TokenTypeLineBreak,
					Text: "\n",
				},
			},
		},
		{
			name:  "Pre-formatted with tab indent",
			input: "\tThis is pre-formatted.\n\tThis as well.\n",
			expected: []agmi.Token{
				{
					Type: agmi.TokenTypeIndent,
					Text: "\t",
				},
				{
					Type: agmi.TokenTypeText,
					Text: "This is pre-formatted.",
				},
				{
					Type: agmi.TokenTypeLineBreak,
					Text: "\n",
				},
				{
					Type: agmi.TokenTypeIndent,
					Text: "\t",
				},
				{
					Type: agmi.TokenTypeText,
					Text: "This as well.",
				},
				{
					Type: agmi.TokenTypeLineBreak,
					Text: "\n",
				},
			},
		},
		{
			name:  "List items",
			input: "* First list item\n*Second list item",
			expected: []agmi.Token{
				{
					Type: agmi.TokenTypeBulletPoint,
					Text: "* ",
				},
				{
					Type: agmi.TokenTypeText,
					Text: "First list item",
				},
				{
					Type: agmi.TokenTypeLineBreak,
					Text: "\n",
				},
				{
					Type: agmi.TokenTypeBulletPoint,
					Text: "*", // Note the missing space.
				},
				{
					Type: agmi.TokenTypeText,
					Text: "Second list item",
				},
			},
		},
		{
			name:  "Multi-line list items",
			input: "* First list item\n  With a second line\n*Second list item\n  With another line",
			expected: []agmi.Token{
				{
					Type: agmi.TokenTypeBulletPoint,
					Text: "* ",
				},
				{
					Type: agmi.TokenTypeText,
					Text: "First list item",
				},
				{
					Type: agmi.TokenTypeLineBreak,
					Text: "\n",
				},
				{
					Type: agmi.TokenTypeIndent,
					Text: "  ",
				},
				{
					Type: agmi.TokenTypeText,
					Text: "With a second line",
				},
				{
					Type: agmi.TokenTypeLineBreak,
					Text: "\n",
				},
				{
					Type: agmi.TokenTypeBulletPoint,
					Text: "*", // Note the missing space.
				},
				{
					Type: agmi.TokenTypeText,
					Text: "Second list item",
				},
				{
					Type: agmi.TokenTypeLineBreak,
					Text: "\n",
				},
				{
					Type: agmi.TokenTypeIndent,
					Text: "  ",
				},
				{
					Type: agmi.TokenTypeText,
					Text: "With another line",
				},
			},
		},
		{
			name:  "link to HTTP",
			input: "=> gemini://example.com This is a link\n",
			expected: []agmi.Token{
				{
					Type: agmi.TokenTypeLinkMod,
					Text: "=> ",
				},
				{
					Type: agmi.TokenTypeLinkURI,
					Text: "gemini://example.com",
				},
				{
					Type: agmi.TokenTypeText,
					Text: " This is a link",
				},
				{
					Type: agmi.TokenTypeLineBreak,
					Text: "\n",
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			sc := agmi.NewScanner(bytes.NewBufferString(tt.input))

			nTokens := len(tt.expected)
			for i, etok := range tt.expected {
				assert.Truef(t, sc.Scan(), "Scan returned false for token %d/%d", i, nTokens)
				if !assert.NoErrorf(t, sc.Err(), "Scan failed for token %d/%d", i, nTokens) {
					return
				}
				tok := sc.Token()
				assert.Equal(t, etok, tok)
			}
			assert.False(t, sc.Scan(), "Not all tokens consumed")
		})
	}
}
