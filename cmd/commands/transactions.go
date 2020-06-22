package commands

import (
	"fmt"
	"time"

	cmdConfig "github.com/SebastianJ/elrond-cli/config/cmd"
	sdkAPI "github.com/SebastianJ/elrond-sdk/api"
	sdkTransactions "github.com/SebastianJ/elrond-sdk/transactions"
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
			return sendTransactions(cmd)
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
	cmdTx.Flags().StringVar(&cmdConfig.Tx.Proxy, "proxy", "", "Proxy address to use for interacting with the REST API")
	cmdTx.Flags().Int64Var(&cmdConfig.Tx.GasLimit, "gas-limit", -1, "Gas limit")
	cmdTx.Flags().Int64Var(&cmdConfig.Tx.GasPrice, "gas-price", -1, "Gas price")
	cmdTx.Flags().Int64Var(&cmdConfig.Tx.Count, "count", -1, "How many transactions to send")
	cmdTx.Flags().Int64Var(&cmdConfig.Tx.Sleep, "sleep", -1, "How long the CLI should sleep after sending a transaction")
	cmdTx.Flags().StringVar(&cmdConfig.Tx.ConfigPath, "config", "./configs/economics.toml", "The economics configuration file to load")
	cmdTx.Flags().BoolVar(&cmdConfig.Tx.ForceAPINonceLookups, "force-api-nonce-lookups", false, "Force the usage of https://wallet-api.elrond.com for checking nonces when using local node endpoints")

	RootCmd.AddCommand(cmdTx)
}

func sendTransactions(cmd *cobra.Command) error {
	if cmdConfig.Tx.To == "" {
		return errors.New("please provide a valid receiver address using --to ADDRESS")
	}

	client := sdkAPI.Client{
		Host:                 cmdConfig.Persistent.Endpoint,
		ForceAPINonceLookups: cmdConfig.Tx.ForceAPINonceLookups,
	}

	if cmdConfig.Tx.Proxy != "" {
		client.Proxy = cmdConfig.Tx.Proxy
	}
	client.Initialize()

	defaultGasParams, err := sdkTransactions.ParseGasSettings(cmdConfig.Tx.ConfigPath)
	if err != nil {
		return err
	}

	// Make a copy of the default gas params that can be modified when processing the tx
	gasParams := defaultGasParams
	if cmdConfig.Tx.GasLimit != -1 {
		gasParams.GasLimit = uint64(cmdConfig.Tx.GasLimit)
	}
	if cmdConfig.Tx.GasPrice != -1 {
		gasParams.GasPrice = uint64(cmdConfig.Tx.GasPrice)
	}

	wallet, err := sdkWallet.Decrypt(cmdConfig.Tx.WalletPath)
	if err != nil {
		return err
	}

	if cmdConfig.Tx.Count <= 0 {
		sendSingleTransaction(wallet, gasParams, client)
	} else {
		sendMultipleTransactions(wallet, gasParams, client)
	}

	return nil
}

func sendSingleTransaction(wallet sdkWallet.Wallet, gasParams sdkTransactions.GasParams, client sdkAPI.Client) error {
	_, txHash, err := sdkTransactions.SendTransaction(
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

func sendMultipleTransactions(wallet sdkWallet.Wallet, gasParams sdkTransactions.GasParams, client sdkAPI.Client) error {
	txs := []*sdkAPI.TransactionData{}

	account, err := client.GetAccount(wallet.Address)
	if err != nil {
		return err
	}

	nonce := int64(account.Nonce)

	for i := int64(0); i < cmdConfig.Tx.Count; i++ {
		fmt.Printf("Nonce is now: %d\n", nonce)

		tx, err := sdkTransactions.GenerateAndSignTransaction(
			wallet,
			cmdConfig.Tx.To,
			cmdConfig.Tx.Amount,
			cmdConfig.Tx.MaximumAmount,
			nonce,
			cmdConfig.Tx.Data,
			gasParams,
			client,
		)
		if err != nil {
			return err
		}

		fmt.Printf("Will attempt to send tx. Receiver: %s, amount: %f, nonce: %d, tx hash %s\n", cmdConfig.Tx.To, cmdConfig.Tx.Amount, nonce, tx.TxHash)
		txs = append(txs, tx.APIData)
		nonce++
	}

	response, err := client.SendMultipleTransactions(txs)
	if err != nil {
		return err
	}

	fmt.Printf("Sent a total of %d transactions!\n", response.TxsSent)
	/*for _, txHash := range response.TxsHashes {
		fmt.Printf("Sent tx: %s\n", txHash)
	}*/

	return nil
}
