package config

import (
    "log"
    "github.com/hashgraph/hedera-sdk-go/v2"
)

type Config struct {
    AccountID  hedera.AccountID
    PrivateKey hedera.PrivateKey
    Network    string
}

func LoadConfig() Config {
    accountID, err := hedera.AccountIDFromString("0.0.4668422")
    if err != nil {
        log.Fatalf("Failed to parse Account ID: %v", err)
    }

    privateKey, err := hedera.PrivateKeyFromString("7674d8b8e535ed9182c2ac6d2d1dde0de18e7c1906cb20fc2663b2028fbf7fc7")
    if err != nil {
        log.Fatalf("Failed to parse Private Key: %v", err)
    }

    return Config{
        AccountID:  accountID,
        PrivateKey: privateKey,
        Network:    "testnet", 
    }
}
