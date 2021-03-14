package mnml

import (
	"fmt"

	"github.com/fhofherr/mnml/internal/version"
	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(*cobra.Command, []string) {
			fmt.Println(version.String())
		},
	}
}
