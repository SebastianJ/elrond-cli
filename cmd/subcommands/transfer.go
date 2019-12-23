package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/SebastianJ/elrond-cli/transactions"
	"github.com/urfave/cli"
)

// TransferCommand sets up the CLI functionality for performing transfers
func TransferCommand() cli.Command {
	return cli.Command{
		Name:    "transfer",
		Aliases: []string{"tx"},
		Usage:   "transfer tokens from one address to another",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "wallet",
				Usage: "Which wallet to use for sending transactions",
				Value: "./initialBalancesSk.pem",
			},
			cli.StringFlag{
				Name:  "to",
				Usage: "Which address to send tokens to",
				Value: "",
			},
			cli.Float64Flag{
				Name:  "amount",
				Usage: "How many tokens to send per transaction",
				Value: 1.0,
			},
			cli.BoolFlag{
				Name:  "maximum",
				Usage: "Send the maximum available amount of tokens",
			},
			cli.Int64Flag{
				Name:  "nonce",
				Usage: "What nonce to use for sending the transaction",
				Value: -1,
			},
			cli.StringFlag{
				Name:  "data",
				Usage: "Transaction data to use for sending the transaction",
				Value: "",
			},
			cli.Int64Flag{
				Name:  "sleep",
				Usage: "How long the CLI should sleep after sending a transaction",
				Value: -1,
			},
			cli.StringFlag{
				Name:  "config",
				Usage: "The economics configuration file to load",
				Value: "./config/economics.toml",
			},
		},
		Action: func(ctx *cli.Context) error {
			return sendTransactionCommand(ctx)
		},
	}
}

func sendTransactionCommand(ctx *cli.Context) error {
	certKeyPath := ctx.String("wallet")

	encodedKey, err := core.LoadSkFromPemFile(certKeyPath, 0)

	if err != nil {
		return errors.New("you need to provide a valid path to a wallet/initialBalancesSk.pem file using --wallet PATH_TO_PEM_FILE")
	}

	receiver := ctx.String("to")

	if receiver == "" || len(receiver) < 64 {
		return errors.New("please provide a valid receiver address using --to ADDRESS")
	}

	maximum := ctx.Bool("maximum")
	amount := ctx.Float64("amount")
	txData := ctx.String("data")
	apiHost := ctx.GlobalString("api-endpoint")
	nonce := ctx.Int64("nonce")
	sleep := ctx.Int64("sleep")

	configPath := ctx.String("config")
	gasPrice, gasLimit, err := parseGasSettings(configPath)

	if err != nil {
		return err
	}

	txHexHash, err := transactions.SendTransaction(encodedKey, receiver, amount, maximum, nonce, txData, gasPrice, gasLimit, apiHost)

	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Success! Your pending transaction hash is: %s", txHexHash))

	if sleep > 0 {
		time.Sleep(time.Duration(sleep) * time.Second)
	}

	return nil
}

func parseGasSettings(configPath string) (uint64, uint64, error) {
	var gasPrice uint64
	var gasLimit uint64

	economicsConfig, err := loadEconomicsConfig(configPath)
	if err == nil {
		gasPrice, err := strconv.ParseInt(economicsConfig.FeeSettings.MinGasPrice, 10, 64)

		if err != nil {
			return 0, 0, err
		}

		gasLimit, err := strconv.ParseInt(economicsConfig.FeeSettings.MinGasLimit, 10, 64)

		if err != nil {
			return 0, 0, err
		}

		return uint64(gasPrice), uint64(gasLimit), nil
	}

	gasPrice = 100000000000000
	gasLimit = 100000

	return gasPrice, gasLimit, nil
}

func loadEconomicsConfig(filepath string) (*config.ConfigEconomics, error) {
	cfg := &config.ConfigEconomics{}
	err := core.LoadTomlFile(cfg, filepath)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
