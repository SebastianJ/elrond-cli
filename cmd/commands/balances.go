package commands

import (
	"fmt"

	cmdConfig "github.com/SebastianJ/elrond-cli/config/cmd"
	"github.com/SebastianJ/elrond-sdk/api"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	includeAddress bool
	addresses      []string
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
	cmdBalance.Flags().BoolVar(&includeAddress, "include-address", false, "Include the address in the balance output")

	cmdBalances := &cobra.Command{
		Use:   "balances",
		Short: "Check addresses balances",
		Long:  "Check addresses balances",
		RunE: func(cmd *cobra.Command, args []string) error {
			return checkBalances(cmd, args)
		},
	}
	cmdBalances.Flags().StringSliceVar(&addresses, "addresses", []string{}, "Addresses to check balances for")
	cmdBalances.Flags().BoolVar(&includeAddress, "include-address", false, "Include the address in the balance output")

	RootCmd.AddCommand(cmdBalance)
	RootCmd.AddCommand(cmdBalances)
}

func checkBalance(cmd *cobra.Command, args []string) error {
	address := args[0]

	if address == "" {
		return errors.New("please provide a valid address")
	}

	client := api.Client{Host: cmdConfig.Persistent.Endpoint}
	client.Initialize()

	account, err := client.GetBalance(address)
	if err != nil {
		return errors.New("failed to retrieve balance")
	}

	if includeAddress {
		fmt.Println(fmt.Sprintf("Address: %s, balance: %f", account.Address, account.Balance))
	} else {
		fmt.Println(fmt.Sprintf("%f", account.Balance))
	}

	return nil
}

func checkBalances(cmd *cobra.Command, args []string) error {
	client := api.Client{Host: cmdConfig.Persistent.Endpoint}
	client.Initialize()

	accounts, err := getBalances(client)
	if err != nil {
		return err
	}

	for _, account := range accounts {
		if includeAddress {
			fmt.Println(fmt.Sprintf("Address: %s, balance: %f", account.Address, account.Balance))
		} else {
			fmt.Println(fmt.Sprintf("%f", account.Balance))
		}
	}

	return nil
}

func getBalances(client api.Client) ([]api.Account, error) {
	accounts := []api.Account{}

	for _, address := range addresses {
		account, err := client.GetBalance(address)
		if err != nil {
			return accounts, errors.Errorf("failed to retrieve balance for address %s", address)
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}
