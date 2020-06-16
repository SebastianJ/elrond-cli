package commands

import (
	"fmt"
	"time"

	cmdConfig "github.com/SebastianJ/elrond-cli/config/cmd"
	"github.com/SebastianJ/elrond-sdk/api"
	"github.com/SebastianJ/elrond-sdk/transactions"
	sdkWallet "github.com/SebastianJ/elrond-sdk/wallet"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	cmdTx := &cobra.Command{
		Use:   "transfer",
		Short: "Send transaction",
		Long:  "Send transaction",
		RunE: func(cmd *cobra.Command, args []string) error {
			return sendTransaction(cmd)
		},
	}

	cmdConfig.Tx = cmdConfig.TxFlags{}
	cmdTx.Flags().StringVar(&cmdConfig.Tx.WalletPath, "wallet", "./keys/initialBalancesSk.pem", "Wallet PEM file to use for sending transactions")
	cmdTx.Flags().StringVar(&cmdConfig.Tx.Password, "password", "", "Wallet password")
	cmdTx.Flags().StringVar(&cmdConfig.Tx.To, "to", "", "Which address to send tokens to")
	cmdTx.Flags().Float64Var(&cmdConfig.Tx.Amount, "amount", 1.0, "How many tokens to send")
	cmdTx.Flags().BoolVar(&cmdConfig.Tx.MaximumAmount, "maximum-amount", false, "Send the maximum available amount of tokens")
	cmdTx.Flags().Int64Var(&cmdConfig.Tx.Nonce, "nonce", -1, "What nonce to use for sending the transaction")
	cmdTx.Flags().StringVar(&cmdConfig.Tx.Data, "data", "", "Transaction data to use for sending the transaction")
	cmdTx.Flags().Int64Var(&cmdConfig.Tx.Sleep, "sleep", -1, "How long the CLI should sleep after sending a transaction")
	cmdTx.Flags().StringVar(&cmdConfig.Tx.ConfigPath, "config", "./configs/economics.toml", "The economics configuration file to load")
	cmdTx.Flags().BoolVar(&cmdConfig.Tx.ForceAPINonceLookups, "force-api-nonce-lookups", false, "Force the usage of https://wallet-api.elrond.com for checking nonces when using local node endpoints")

	RootCmd.AddCommand(cmdTx)
}

func sendTransaction(cmd *cobra.Command) error {
	if cmdConfig.Tx.To == "" {
		return errors.New("please provide a valid receiver address using --to ADDRESS")
	}

	client := api.Client{
		Host:                 cmdConfig.Persistent.Endpoint,
		ForceAPINonceLookups: cmdConfig.Tx.ForceAPINonceLookups,
	}
	client.Initialize()

	defaultGasParams, err := transactions.ParseGasSettings(cmdConfig.Tx.ConfigPath)
	if err != nil {
		return err
	}

	// Make a copy of the default gas params that can be modified when processing the tx
	gasParams := defaultGasParams

	wallet, err := sdkWallet.Decrypt(cmdConfig.Tx.WalletPath)
	if err != nil {
		return err
	}

	_, txHash, err := transactions.SendTransaction(
		wallet,
		cmdConfig.Tx.To,
		cmdConfig.Tx.Amount,
		cmdConfig.Tx.MaximumAmount,
		cmdConfig.Tx.Nonce,
		cmdConfig.Tx.Data,
		gasParams,
		client,
	)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Success! Tx hash: %s", txHash))

	if cmdConfig.Tx.Sleep > 0 {
		time.Sleep(time.Duration(cmdConfig.Tx.Sleep) * time.Second)
	}

	return nil
}
