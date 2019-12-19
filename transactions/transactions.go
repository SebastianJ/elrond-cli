package transactions

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"
	"strings"
	"time"

	"github.com/ElrondNetwork/elrond-go/crypto"
	"github.com/ElrondNetwork/elrond-go/crypto/signing"
	"github.com/ElrondNetwork/elrond-go/crypto/signing/kyber"
	"github.com/ElrondNetwork/elrond-go/crypto/signing/kyber/singlesig"
	"github.com/ElrondNetwork/elrond-go/data/transaction"
	"github.com/SebastianJ/elrond-cli/api"
)

// SendTransaction - broadcast a transaction to the blockchain
func SendTransaction(encodedKey []byte, receiver string, amount float64, txData string, gasPrice uint64, gasLimit uint64, apiHost string) (string, error) {
	signer, privKey, pubKey, err := generateCryptoSuite(encodedKey)

	if err != nil {
		return "", err
	}

	pubKeyBytes, err := pubKey.ToByteArray()

	if err != nil {
		return "", err
	}

	sender := hex.EncodeToString(pubKeyBytes)
	hexSender, err := hex.DecodeString(sender)

	if err != nil {
		return "", err
	}

	hexReceiver, err := hex.DecodeString(receiver)

	if err != nil {
		return "", err
	}

	accountData, err := api.GetAccount(sender, apiHost)

	if err != nil {
		return "", errors.New("failed to retrieve account data")
	}

	nonce := accountData.Nonce

	realAmount := convertAmountToBigInt(amount)
	gasLimit = gasLimit + uint64(len(txData))

	tx := transaction.Transaction{
		Nonce:    nonce,
		SndAddr:  hexSender,
		RcvAddr:  hexReceiver,
		Value:    realAmount,
		Data:     txData,
		GasPrice: gasPrice,
		GasLimit: gasLimit,
	}

	txBuff, _ := json.Marshal(&tx)
	signature, _ := signer.Sign(privKey, txBuff)

	txHexHash, txError := api.SendTransaction(nonce, sender, receiver, realAmount.String(), gasPrice, gasLimit, txData, signature, apiHost)

	if txError != nil {
		// If we've sent an invalid nonce - sleep 3 seconds and then retry again using a fresh nonce
		if strings.Contains(txError.Error(), "transaction generation failed: invalid nonce") {
			time.Sleep(3 * time.Second)
			return SendTransaction(encodedKey, receiver, amount, txData, gasPrice, gasLimit, apiHost)
		} else {
			return "", txError
		}
	}

	return txHexHash, nil
}

func convertAmountToBigInt(amount float64) *big.Int {
	bigAmount := new(big.Float).SetFloat64(amount)
	base := new(big.Float).SetInt(big.NewInt(1000000000000000000))
	bigAmount.Mul(bigAmount, base)
	realAmount := new(big.Int)
	bigAmount.Int(realAmount)

	return realAmount
}

func generateCryptoSuite(encodedKey []byte) (signer *singlesig.SchnorrSigner, privKey crypto.PrivateKey, pubKey crypto.PublicKey, err error) {
	signer = &singlesig.SchnorrSigner{}
	suite := kyber.NewBlakeSHA256Ed25519()
	decodedKey, err := hex.DecodeString(string(encodedKey))

	keyGen := signing.NewKeyGenerator(suite)

	privKey, err = keyGen.PrivateKeyFromByteArray(decodedKey)
	if err != nil {
		return nil, nil, nil, err
	}

	pubKey = privKey.GeneratePublic()

	return signer, privKey, pubKey, err
}
