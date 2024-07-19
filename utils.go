package main

import (
	"encoding/json"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	RED    = "#FF0000"
	YELLOW = "#FFFF00"
	GREEN  = "#00FF00"
)

func MustCreateRandomAccount() Account {
	// Create a new account
	acc := Account{}
	var err error
	if acc.PrivKey, err = crypto.GenerateKey(); err != nil {
		log.Fatal(err)
	}
	acc.Address = crypto.PubkeyToAddress(acc.PrivKey.PublicKey)
	return acc
}

// MustBeautifyBlock formats and prints an Ethereum block in a readable JSON format
func MustBeautifyBlock(block *RpcBlock) string {
	blockJSON, err := json.MarshalIndent(block, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal block: %v", err)
	}
	return string(blockJSON)
}

// MustBeautifyReceipt formats and prints an Ethereum receipt in a readable JSON format
func MustBeautifyReceipt(receipt *types.Receipt) string {
	receiptJSON, err := json.MarshalIndent(receipt, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal receipt: %v", err)
	}
	return string(receiptJSON)
}

// MustBeautifyReceipts formats and prints a list of Ethereum receipts in a readable JSON format
func MustBeautifyReceipts(receipts types.Receipts) string {
	receiptsJSON, err := json.MarshalIndent(receipts, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal receipts: %v", err)
	}
	return string(receiptsJSON)
}

// MustBeautifyTransaction formats and prints an Ethereum transaction in a readable JSON format
func MustBeautifyTransaction(tx *types.Transaction) string {
	txJSON, err := tx.MarshalJSON()
	if err != nil {
		log.Fatalf("Failed to marshal transaction: %v", err)
	}
	return string(txJSON)
}
