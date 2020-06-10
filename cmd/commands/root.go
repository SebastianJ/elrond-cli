package commands

import (
	"fmt"
	"os"

	cmd "github.com/SebastianJ/elrond-cli/config/cmd"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// VersionWrap - version displayed in case of errors
	VersionWrap = ""

	// RootCmd - main entry point for Cobra commands
	RootCmd = &cobra.Command{
		Use:          "erd",
		Short:        "Elrond CLI",
		Long:         "Elrond CLI",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Help()
			return nil
		},
	}
)

func init() {
	cmd.Persistent = cmd.PersistentFlags{}

	RootCmd.PersistentFlags().StringVar(&cmd.Persistent.Endpoint, "endpoint", "https://wallet-api.elrond.com", "Which API endpoint to use for API commands")
}

// Execute starts the actual app
func Execute() {
	RootCmd.SilenceErrors = true
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(errors.Wrapf(err, "commit: %s, error", VersionWrap).Error())
		os.Exit(1)
	}
}
