package agmi

import (
	"fmt"
	"io"
)

// ConverterState is a function that processes the current token of the input.
//
// Any ConverterState may use the next to make decisions on how to process the
// current Token or to update the Converter state. However, it is absolutely
// necessary that the current token is processed one way or the other. After
// cur is passed to a ConverterState function it will not be passed again
// during the remainder of the formatting run.
type ConverterState func(c *Converter, cur, next Token)

// Converter is a type that helps implementing the conversion from
// Almost Gemtext to some other format.
//
// The zero value of Converter is not usable. Use the NewConverter
// function to obtain a working instance.
type Converter struct {
	State ConverterState // Current state of the Converter. Update when transitioning.
	Err   error          // Set this field if an error occurs while processing a token.

	scanner *Scanner
	out     io.Writer
}

// NewConverter creates a new Converter that writes its input in
// converted form to out.
//
// The start ConverterState is the state the converter helper starts with when
// processing begins.
func NewConverter(in io.Reader, out io.Writer, start ConverterState) Converter {
	return Converter{
		State:   start,
		scanner: NewScanner(in),
		out:     out,
	}
}

// Convert processes the input and writes it in a converted form to the
// output.
func (c *Converter) Convert() error {
	const op = "agmi/Converter.Format"

	if c.State == nil {
		return fmt.Errorf("%s: initial state not set", op)
	}

	// Initialize current token
	if !c.scanner.Scan() {
		if c.scanner.Err() != nil {
			return fmt.Errorf("%s: %v", op, c.scanner.Err())
		}
		return fmt.Errorf("%s: no input", op)
	}
	cur := c.scanner.Token()

	for c.scanner.Scan() {
		next := c.scanner.Token()

		c.State(c, cur, next)
		if c.Err != nil {
			break
		}
		cur = next
	}
	if c.scanner.Err() != nil {
		return fmt.Errorf("%s: %v", op, c.scanner.Err())
	}

	// Finish the last token
	c.State(c, cur, Token{})
	if c.Err != nil {
		return fmt.Errorf("%s: %v", op, c.Err)
	}
	return nil
}

// Write writes the string s to the output.
func (c *Converter) Write(s string) {
	const op = "agmi/Converter.Write"

	if _, err := c.out.Write([]byte(s)); err != nil {
		c.Err = fmt.Errorf("%s: %v", op, err)
	}
}
