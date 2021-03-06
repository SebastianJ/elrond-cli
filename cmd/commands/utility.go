package commands

import (
	"encoding/hex"
	"fmt"
	"strings"

	sdkTransactions "github.com/SebastianJ/elrond-sdk/transactions"
	sdkUtils "github.com/SebastianJ/elrond-sdk/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	keys           []string
	numberOfShards uint32 = 2
)

func init() {
	cmdUtility := &cobra.Command{
		Use:   "utility",
		Short: "Utility functions",
		Long:  "Utility functions",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Help()
			return nil
		},
	}

	cmdConvertToBech32 := &cobra.Command{
		Use:   "to-bech32",
		Short: "Convert from public keys to Bech32",
		Long:  "Convert from public keys to Bech32",
		RunE: func(cmd *cobra.Command, args []string) error {
			return convertKeysToBech32(cmd)
		},
	}
	cmdConvertToBech32.Flags().StringSliceVar(&keys, "keys", []string{}, "Public keys to convert to Bech32, separated by a comma")

	cmdConvertFromBech32 := &cobra.Command{
		Use:   "from-bech32",
		Short: "Convert from Bech32 to public keys",
		Long:  "Convert from Bech32 to public keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			return convertKeysFromBech32(cmd)
		},
	}
	cmdConvertFromBech32.Flags().StringSliceVar(&keys, "keys", []string{}, "Bech32 keys to convert to public keys, separated by a comma")

	cmdShardForAddress := &cobra.Command{
		Use:   "shard-for-address",
		Short: "Detect shard for a given address",
		Long:  "Detect shard for a given address",
		RunE: func(cmd *cobra.Command, args []string) error {
			return detectShardForAddress(cmd)
		},
	}
	cmdShardForAddress.Flags().StringSliceVar(&keys, "keys", []string{}, "Public keys to check shard for")

	cmdUtility.AddCommand(cmdConvertToBech32)
	cmdUtility.AddCommand(cmdConvertFromBech32)
	cmdUtility.AddCommand(cmdShardForAddress)
	RootCmd.AddCommand(cmdUtility)
}

func convertKeysToBech32(cmd *cobra.Command) error {
	if len(keys) == 0 {
		return errors.New("please provide keys to convert using --keys")
	}

	for _, key := range keys {
		bech32, err := sdkUtils.PublicKeyToBech32(key)
		if err != nil {
			return err
		}
		fmt.Printf("Key: %s - bech32: %s\n", key, bech32)
	}

	return nil
}

func convertKeysFromBech32(cmd *cobra.Command) error {
	if len(keys) == 0 {
		return errors.New("please provide keys to convert using --keys")
	}

	for _, bech32 := range keys {
		key, err := sdkUtils.Bech32ToPublicKey(bech32)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", key)
	}

	return nil
}

func detectShardForAddress(cmd *cobra.Command) (err error) {
	if len(keys) == 0 {
		return errors.New("please provide keys to convert using --keys")
	}

	for _, key := range keys {
		var addressBytes []byte

		if strings.HasPrefix(key, "erd") {
			addressBytes, err = sdkUtils.Bech32ToPublicKeyBytes(key)
			if err != nil {
				return err
			}
		} else {
			addressBytes, err = hex.DecodeString(key)
			if err != nil {
				return err
			}
		}

		shardID := sdkTransactions.CalculateShardForAddress(addressBytes, numberOfShards)
		if err != nil {
			return err
		}
		fmt.Printf("Address: %s - shard: %d\n", key, shardID)
	}

	return nil
}
