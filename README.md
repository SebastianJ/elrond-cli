# Elrond CLI

Elrond CLI is a CLI tool that aims users to interact with Elrond's blockchain using a set of CLI commands.

## Installation

Install the CLI using:

`bash <(curl -s -S -L https://raw.githubusercontent.com/SebastianJ/elrond-cli/master/scripts/erd.sh)`

The CLI tool should now be installed as `erd`

You can invoke it using `./erd COMMAND ARGUMENTS` - see `./erd --help` for available commands and arguments.

## Usage

### Sending transactions

If you already have a `initialBalancesSk.pem` file in the same directory as the CLI:

`./erd transfer --to RECEIVER_ADDRESS --amount AMOUNT`

If you want to specify a custom wallet/initialBalancesSk.pem file:

`./erd transfer --wallet PATH_TO_CUSTOM_INITIALBALANCESSK_FILE.pem --to RECEIVER_ADDRESS --amount AMOUNT`

If --amount AMOUNT isn't specified 1 ERD will be sent by default.
