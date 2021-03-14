package mnml

import (
	"fmt"
	"os"

	"github.com/fhofherr/mnml/gemtext"
	"github.com/spf13/cobra"
)

func newAGMI2GMICmd() *cobra.Command {
	var outFile string

	agmi2gmi := &cobra.Command{
		Use:   "agmi2gmi",
		Short: "Transform Almost Gemtext to Gemtext",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inFile := args[0] // The ExactArgs ensures this is always there.

			in, err := os.Open(inFile)
			if err != nil {
				return fmt.Errorf("open input: %v", err)
			}
			defer in.Close()

			out := os.Stdout
			if outFile != "" {
				var err error

				out, err = os.Create(outFile)
				if err != nil {
					return fmt.Errorf("open output: %v", err)
				}
				defer out.Close()
			}

			if err := gemtext.FromAlmostGemtext(in, out); err != nil {
				return fmt.Errorf("convert %s to Gemtext: %v", inFile, err)
			}
			return nil
		},
	}
	agmi2gmi.Flags().StringVarP(
		&outFile, "output", "o", "", "Write the converted text to this file. Defaults to stdout if missing.")

	return agmi2gmi
}
