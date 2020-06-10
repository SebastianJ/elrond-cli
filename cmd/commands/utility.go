package commands

import (
	"fmt"

	"github.com/SebastianJ/elrond-cli/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var keys []string

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

	cmdConvert := &cobra.Command{
		Use:   "convert",
		Short: "Conversion functions",
		Long:  "Conversion functions",
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

	cmdConvert.AddCommand(cmdConvertToBech32)
	cmdConvert.AddCommand(cmdConvertFromBech32)
	cmdUtility.AddCommand(cmdConvert)
	RootCmd.AddCommand(cmdUtility)
}

func convertKeysToBech32(cmd *cobra.Command) error {
	if len(keys) == 0 {
		return errors.New("please provide keys to convert using --keys")
	}

	for _, key := range keys {
		bech32, err := utils.PublicKeyToBech32(key)
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
		key, err := utils.Bech32ToPublicKey(bech32)
		if err != nil {
			return err
		}
		fmt.Printf("Bech32: %s - key: %s\n", bech32, key)
	}

	return nil
}
