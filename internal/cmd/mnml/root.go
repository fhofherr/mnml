package mnml

import "github.com/spf13/cobra"

// New creates a new instance of mnml ready for use by the command line.
func New() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "mnml",
		Short: "A minimalistic Gemini and Gopher site generator.",
	}
	rootCmd.AddCommand(newAGMI2GMICmd())
	rootCmd.AddCommand(newVersionCmd())

	return rootCmd
}
