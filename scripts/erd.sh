#!/usr/bin/env bash

echo "Installing Elrond CLI (erd)"
curl -LOs http://tools.elrond.com.s3.amazonaws.com/release/linux-x86_64/erd && chmod u+x erd
mkdir -p configs && curl -o configs/economics.toml https://raw.githubusercontent.com/ElrondNetwork/elrond-config/master/economics.toml
echo "Elrond CLI is now ready to use!"
echo "Invoke it using ./erd - see ./erd --help for all available commands and options"
