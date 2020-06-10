package commands

import (
	"fmt"

	cmdConfig "github.com/SebastianJ/elrond-cli/config/cmd"
	"github.com/SebastianJ/elrond-sdk/api"
	"github.com/SebastianJ/elrond-sdk/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	cmdBalance := &cobra.Command{
		Use:   "balance",
		Short: "Check address balance",
		Long:  "Check address balance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return checkBalance(cmd, args)
		},
	}

	RootCmd.AddCommand(cmdBalance)
}

func checkBalance(cmd *cobra.Command, args []string) error {
	address := args[0]

	if address == "" {
		return errors.New("please provide a valid address")
	}

	client := api.Client{Host: cmdConfig.Persistent.Endpoint}
	client.Initialize()

	accountData, err := client.GetBalance(address)

	if err != nil {
		return errors.New("failed to retrieve balance")
	}

	balance := accountData.Balance
	converted, _ := utils.ConvertNumeralStringToBigFloat(balance)

	fmt.Println(fmt.Sprintf("%f", converted))

	return nil
}
