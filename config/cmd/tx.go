package cmd

// Tx is a collection of tx related options
var Tx TxFlags

// TxFlags represents the tx flags
type TxFlags struct {
	WalletPath           string
	Password             string
	To                   string
	Amount               float64
	MaximumAmount        bool
	Nonce                int64
	Data                 string
	Sleep                int64
	ConfigPath           string
	Count                int64
	GasPrice             int64
	GasLimit             int64
	Proxy                string
	ForceAPINonceLookups bool
}
