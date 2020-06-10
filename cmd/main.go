package main

import (
	"fmt"
	"os"
	"path"
	"runtime"

	cmd "github.com/SebastianJ/elrond-cli/cmd/commands"
	"github.com/spf13/cobra"
)

func main() {
	// Force usage of Go's own DNS implementation
	os.Setenv("GODEBUG", "netdns=go")

	cmd.VersionWrap = fmt.Sprintf("%s/%s-%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	cmd.RootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Show version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(os.Stderr, "SebastianJ (C) 2020. %v, version %s/%s-%s\n", path.Base(os.Args[0]), runtime.Version(), runtime.GOOS, runtime.GOARCH)
			os.Exit(0)
			return nil
		},
	})
	cmd.Execute()
}
