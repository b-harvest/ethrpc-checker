package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
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
	// First, use the default MarshalJSON to get the serialized data
	txJSON, err := tx.MarshalJSON()
	if err != nil {
		log.Fatalf("Failed to marshal transaction: %v", err)
	}

	// Then, unmarshal the serialized data into a map
	var txMap map[string]interface{}
	if err := json.Unmarshal(txJSON, &txMap); err != nil {
		log.Fatalf("Failed to unmarshal transaction: %v", err)
	}

	// Finally, marshal the map with indentation
	indentedTxJSON, err := json.MarshalIndent(txMap, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal indented transaction: %v", err)
	}

	return string(indentedTxJSON)
}

func MustCalculateSlotKey(rCtx *RpcContext, slotIndex uint64) common.Hash {
	addressTy, err := abi.NewType("address", "", nil)
	if err != nil {
		log.Fatalf("Failed to create address type: %v", err)
	}
	uint256Ty, err := abi.NewType("uint256", "", nil)
	slotIndexBig := new(big.Int).SetUint64(slotIndex)
	packedArgs, err := abi.Arguments{
		{Type: addressTy},
		{Type: uint256Ty},
	}.Pack(rCtx.Acc.Address, slotIndexBig)
	if err != nil {
		log.Fatalf("Failed to pack arguments: %v", err)
	}

	return crypto.Keccak256Hash(packedArgs)
}

// isZeroBytes checks if a byte slice consists only of zero bytes
func IsZeroBytes(b []byte) bool {
	for _, v := range b {
		if v != 0 {
			return false
		}
	}
	return true
}

func MustBeautifyLogs(logs []types.Log) string {
	receiptsJSON, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal receipts: %v", err)
	}
	return string(receiptsJSON)
}

// MustBeautifyLog returns a formatted string representing the details of an Ethereum log
func MustBeautifyLog(l types.Log) string {
	jsonData, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal log to JSON: %v", err)
	}

	return string(jsonData)
}

func toFilterArg(q ethereum.FilterQuery) (interface{}, error) {
	arg := map[string]interface{}{
		"address": q.Addresses,
		"topics":  q.Topics,
	}
	if q.BlockHash != nil {
		arg["blockHash"] = *q.BlockHash
		if q.FromBlock != nil || q.ToBlock != nil {
			return nil, errors.New("cannot specify both BlockHash and FromBlock/ToBlock")
		}
	} else {
		if q.FromBlock == nil {
			arg["fromBlock"] = "0x0"
		} else {
			arg["fromBlock"] = toBlockNumArg(q.FromBlock)
		}
		arg["toBlock"] = toBlockNumArg(q.ToBlock)
	}
	return arg, nil
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	if number.Sign() >= 0 {
		return hexutil.EncodeBig(number)
	}
	// It's negative.
	if number.IsInt64() {
		return rpc.BlockNumber(number.Int64()).String()
	}
	// It's negative and large, which is invalid.
	return fmt.Sprintf("<invalid %d>", number)
}
