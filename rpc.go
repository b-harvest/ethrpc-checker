package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/google/go-cmp/cmp"
	"github.com/status-im/keycard-go/hexutils"
)

// GethVersion is the version of the Geth client used in the tests
// Update it when go-ethereum of go.mod is updated
const GethVersion = "1.14.7"

type RpcName string
type CallRPC func(rCtx *RpcContext) (*RpcResult, error)

const (
	SendRawTransaction                  RpcName = "eth_sendRawTransaction"
	GetBlockNumber                      RpcName = "eth_blockNumber"
	GetGasPrice                         RpcName = "eth_gasPrice"
	GetMaxPriorityFeePerGas             RpcName = "eth_maxPriorityFeePerGas"
	GetChainId                          RpcName = "eth_chainId"
	GetBalance                          RpcName = "eth_getBalance"
	GetBlockByHash                      RpcName = "eth_getBlockByHash"
	GetBlockByNumber                    RpcName = "eth_getBlockByNumber"
	GetBlockReceipts                    RpcName = "eth_getBlockReceipts"
	GetTransactionByHash                RpcName = "eth_getTransactionByHash"
	GetTransactionByBlockHashAndIndex   RpcName = "eth_getTransactionByBlockHashAndIndex"
	GetTransactionByBlockNumberAndIndex RpcName = "eth_getTransactionByBlockNumberAndIndex"
	GetTransactionReceipt               RpcName = "eth_getTransactionReceipt"
	GetTransactionCount                 RpcName = "eth_getTransactionCount"
	GetTransactionCountByHash           RpcName = "eth_getTransactionCountByHash"
	GetBlockTransactionCountByHash      RpcName = "eth_getBlockTransactionCountByHash"
	GetCode                             RpcName = "eth_getCode"
	GetStorageAt                        RpcName = "eth_getStorageAt"
	NewFilter                           RpcName = "eth_newFilter"
	GetFilterLogs                       RpcName = "eth_getFilterLogs"
	NewBlockFilter                      RpcName = "eth_newBlockFilter"
	GetFilterChanges                    RpcName = "eth_getFilterChanges"
	UninstallFilter                     RpcName = "eth_uninstallFilter"
	GetLogs                             RpcName = "eth_getLogs"
	EstimateGas                         RpcName = "eth_estimateGas"
	Call                                RpcName = "eth_call"
)

type Account struct {
	Address common.Address
	PrivKey *ecdsa.PrivateKey
}

type RpcContext struct {
	Conf                  *Config
	EthCli                *ethclient.Client
	Acc                   *Account
	ChainId               *big.Int
	MaxPriorityFeePerGas  *big.Int
	GasPrice              *big.Int
	ProcessedTransactions []common.Hash
	BlockNumsIncludingTx  []uint64
	AlreadyTestedRPCs     []*RpcResult
	ERC20Abi              *abi.ABI
	ERC20ByteCode         []byte
	ERC20Addr             common.Address
	FilterQuery           ethereum.FilterQuery
	FilterId              string
	BlockFilterId         string
}

func NewContext(conf *Config) (*RpcContext, error) {
	// Connect to the Ethereum client
	ethCli, err := ethclient.Dial(conf.RpcEndpoint)
	if err != nil {
		return nil, err
	}

	ecdsaPrivKey, err := crypto.HexToECDSA(conf.RichPrivKey)
	if err != nil {
		return nil, err
	}

	return &RpcContext{
		Conf:   conf,
		EthCli: ethCli,
		Acc: &Account{
			Address: crypto.PubkeyToAddress(ecdsaPrivKey.PublicKey),
			PrivKey: ecdsaPrivKey,
		},
	}, nil
}

func (rCtx *RpcContext) AlreadyTested(rpc RpcName) *RpcResult {
	for _, testedRPC := range rCtx.AlreadyTestedRPCs {
		if rpc == testedRPC.Method {
			return testedRPC
		}
	}
	return nil

}

func RpcGetBlockNumber(rCtx *RpcContext) (*RpcResult, error) {
	if result := rCtx.AlreadyTested(GetBlockNumber); result != nil {
		return result, nil
	}
	blockNumber, err := rCtx.EthCli.BlockNumber(context.Background())
	if err != nil {
		return nil, err
	}

	var warnings []string
	if blockNumber == 0 {
		warnings = append(warnings, "blockNumber is zero")
	}

	status := Ok
	if len(warnings) > 0 {
		status = Warning
	}

	result := &RpcResult{
		Method:   GetBlockNumber,
		Status:   status,
		Value:    blockNumber,
		Warnings: warnings,
	}
	rCtx.AlreadyTestedRPCs = append(rCtx.AlreadyTestedRPCs, result)

	return result, nil
}

func RpcGetGasPrice(rCtx *RpcContext) (*RpcResult, error) {
	if result := rCtx.AlreadyTested(GetGasPrice); result != nil {
		return result, nil
	}

	gasPrice, err := rCtx.EthCli.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	var warnings []string
	if gasPrice == nil || gasPrice.Cmp(big.NewInt(0)) == 0 {
		warnings = append(warnings, "gasPrice is nil or zero")
	}

	status := Ok
	if len(warnings) > 0 {
		status = Warning
	}

	result := &RpcResult{
		Method:   GetGasPrice,
		Status:   status,
		Value:    gasPrice.String(),
		Warnings: warnings,
	}
	rCtx.AlreadyTestedRPCs = append(rCtx.AlreadyTestedRPCs, result)

	return result, nil
}

func RpcGetMaxPriorityFeePerGas(rCtx *RpcContext) (*RpcResult, error) {
	if result := rCtx.AlreadyTested(GetMaxPriorityFeePerGas); result != nil {
		return result, nil
	}

	maxPriorityFeePerGas, err := rCtx.EthCli.SuggestGasTipCap(context.Background())
	if err != nil {
		return nil, err
	}

	var warnings []string
	if maxPriorityFeePerGas.Cmp(big.NewInt(0)) == 0 {
		warnings = append(warnings, "maxPriorityFeePerGas is nil or zero")
	}

	status := Ok
	if len(warnings) > 0 {
		status = Warning
	}

	result := &RpcResult{
		Method:   GetMaxPriorityFeePerGas,
		Status:   status,
		Value:    maxPriorityFeePerGas.String(),
		Warnings: warnings,
	}
	rCtx.AlreadyTestedRPCs = append(rCtx.AlreadyTestedRPCs, result)

	return result, nil
}

func RpcGetChainId(rCtx *RpcContext) (*RpcResult, error) {
	if result := rCtx.AlreadyTested(GetChainId); result != nil {
		return result, nil
	}

	chainId, err := rCtx.EthCli.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	var warnings []string
	if chainId.Cmp(big.NewInt(0)) == 0 {
		warnings = append(warnings, "chainId is nil")
	}

	status := Ok
	if len(warnings) > 0 {
		status = Warning
	}

	result := &RpcResult{
		Method:   GetChainId,
		Status:   status,
		Value:    chainId.String(),
		Warnings: warnings,
	}
	rCtx.AlreadyTestedRPCs = append(rCtx.AlreadyTestedRPCs, result)

	return result, nil
}

func RpcGetBalance(rCtx *RpcContext) (*RpcResult, error) {
	if result := rCtx.AlreadyTested(GetBalance); result != nil {
		return result, nil
	}

	balance, err := rCtx.EthCli.BalanceAt(context.Background(), rCtx.Acc.Address, nil)
	if err != nil {
		return nil, err
	}

	var warnings []string
	if balance.Cmp(big.NewInt(0)) == 0 {
		warnings = append(warnings, "balance is zero")
	}

	status := Ok
	if len(warnings) > 0 {
		status = Warning
	}

	result := &RpcResult{
		Method:   GetBalance,
		Status:   status,
		Value:    balance.String(),
		Warnings: warnings,
	}
	rCtx.AlreadyTestedRPCs = append(rCtx.AlreadyTestedRPCs, result)

	return result, nil
}

func RpcGetTransactionCount(rCtx *RpcContext) (*RpcResult, error) {
	if result := rCtx.AlreadyTested(GetTransactionCount); result != nil {
		return result, nil
	}

	nonce, err := rCtx.EthCli.PendingNonceAt(context.Background(), rCtx.Acc.Address)
	if err != nil {
		return nil, err
	}

	var warnings []string
	if nonce == 0 {
		warnings = append(warnings, "nonce is zero")
	}

	status := Ok
	if len(warnings) > 0 {
		status = Warning
	}

	return &RpcResult{
		Method:   GetTransactionCount,
		Status:   status,
		Value:    nonce,
		Warnings: warnings,
	}, nil
}

func RpcGetBlockByHash(rCtx *RpcContext) (*RpcResult, error) {
	if result := rCtx.AlreadyTested(GetBlockByHash); result != nil {
		return result, nil
	}

	blkNum, err := rCtx.EthCli.BlockNumber(context.Background())
	if err != nil {
		return nil, err
	}

	blk, err := rCtx.EthCli.BlockByNumber(context.Background(), new(big.Int).SetUint64(blkNum))
	if err != nil {
		return nil, err
	}

	block, err := rCtx.EthCli.BlockByHash(context.Background(), blk.Hash())
	if err != nil {
		return nil, err
	}

	if !cmp.Equal(blk, block) {
		return nil, errors.New("implementation error: blockByNumber and blockByHash return different blocks")
	}

	result := &RpcResult{
		Method: GetBlockByHash,
		Status: Ok,
		Value:  MustBeautifyBlock(NewRpcBlock(block)),
	}
	rCtx.AlreadyTestedRPCs = append(rCtx.AlreadyTestedRPCs, result)

	return result, nil
}

func RpcGetBlockByNumber(rCtx *RpcContext) (*RpcResult, error) {
	if result := rCtx.AlreadyTested(GetBlockByNumber); result != nil {
		return result, nil
	}

	blkNum, err := rCtx.EthCli.BlockNumber(context.Background())
	if err != nil {
		return nil, err
	}

	blk, err := rCtx.EthCli.BlockByNumber(context.Background(), new(big.Int).SetUint64(blkNum))
	if err != nil {
		return nil, err
	}

	result := &RpcResult{
		Method: GetBlockByNumber,
		Status: Ok,
		Value:  MustBeautifyBlock(NewRpcBlock(blk)),
	}
	rCtx.AlreadyTestedRPCs = append(rCtx.AlreadyTestedRPCs, result)

	return result, nil
}

func RpcSendRawTransactionTransferValue(rCtx *RpcContext) (*RpcResult, error) {
	// testedRPCs is a slice of RpcResult that will be appended to rCtx.AlreadyTestedRPCs
	// if the transaction is successfully sent
	var testedRPCs []*RpcResult
	var err error
	// Create a new transaction
	if rCtx.ChainId, err = rCtx.EthCli.ChainID(context.Background()); err != nil {
		return nil, err
	}
	testedRPCs = append(testedRPCs, &RpcResult{
		Method: GetChainId,
		Status: Ok,
		Value:  rCtx.ChainId.String(),
	})

	nonce, err := rCtx.EthCli.PendingNonceAt(context.Background(), rCtx.Acc.Address)
	if err != nil {
		return nil, err
	}
	testedRPCs = append(testedRPCs, &RpcResult{
		Method: GetTransactionCount,
		Status: Ok,
		Value:  nonce,
	})

	if rCtx.MaxPriorityFeePerGas, err = rCtx.EthCli.SuggestGasTipCap(context.Background()); err != nil {
		return nil, err
	}
	testedRPCs = append(testedRPCs, &RpcResult{
		Method: GetMaxPriorityFeePerGas,
		Status: Ok,
		Value:  rCtx.MaxPriorityFeePerGas.String(),
	})
	if rCtx.GasPrice, err = rCtx.EthCli.SuggestGasPrice(context.Background()); err != nil {
		return nil, err
	}
	testedRPCs = append(testedRPCs, &RpcResult{
		Method: GetGasPrice,
		Status: Ok,
		Value:  rCtx.GasPrice.String(),
	})

	randomRecipient := MustCreateRandomAccount().Address
	value := new(big.Int).SetUint64(1)
	balanceBeforeSend, err := rCtx.EthCli.BalanceAt(context.Background(), rCtx.Acc.Address, nil)
	if err != nil {
		return nil, err
	}
	testedRPCs = append(testedRPCs, &RpcResult{
		Method: GetBalance,
		Status: Ok,
		Value:  balanceBeforeSend.String(),
	})

	if balanceBeforeSend.Cmp(value) < 0 {
		return nil, errors.New("insufficient balanceBeforeSend")
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   rCtx.ChainId,
		Nonce:     nonce,
		GasTipCap: rCtx.MaxPriorityFeePerGas,
		GasFeeCap: new(big.Int).Add(rCtx.GasPrice, big.NewInt(1000000000)),
		Gas:       21000, // fixed gas limit for transfer
		To:        &randomRecipient,
		Value:     value,
	})

	// TODO: Make signer using types.MakeSigner with chain params
	signer := types.NewLondonSigner(rCtx.ChainId)
	signedTx, err := types.SignTx(tx, signer, rCtx.Acc.PrivKey)
	if err != nil {
		return nil, err
	}

	if err = rCtx.EthCli.SendTransaction(context.Background(), signedTx); err != nil {
		return nil, err
	}
	result := &RpcResult{
		Method: SendRawTransaction,
		Status: Ok,
		Value:  signedTx.Hash().Hex(),
	}
	testedRPCs = append(testedRPCs, result)

	// wait for the transaction to be mined
	tout, _ := time.ParseDuration(rCtx.Conf.Timeout)
	if err = WaitForTx(rCtx, signedTx.Hash(), tout); err != nil {
		return nil, err
	}

	balance, err := rCtx.EthCli.BalanceAt(context.Background(), rCtx.Acc.Address, nil)
	if err != nil {
		return nil, err
	}
	// check if the balance decreased by the value of the transaction (+ gas fee)
	if new(big.Int).Sub(balanceBeforeSend, balance).Cmp(value) < 0 {
		return nil, errors.New("balanceBeforeSend mismatch, maybe the transaction was not mined or implementation is incorrect")
	}
	rCtx.AlreadyTestedRPCs = append(rCtx.AlreadyTestedRPCs, testedRPCs...)

	return result, nil
}

func RpcSendRawTransactionDeployContract(rCtx *RpcContext) (*RpcResult, error) {
	// testedRPCs is a slice of RpcResult that will be appended to rCtx.AlreadyTestedRPCs
	// if the transaction is successfully sent
	var testedRPCs []*RpcResult
	var err error
	// Create a new transaction
	if rCtx.ChainId, err = rCtx.EthCli.ChainID(context.Background()); err != nil {
		return nil, err
	}
	testedRPCs = append(testedRPCs, &RpcResult{
		Method: GetChainId,
		Status: Ok,
		Value:  rCtx.ChainId.String(),
	})

	nonce, err := rCtx.EthCli.PendingNonceAt(context.Background(), rCtx.Acc.Address)
	if err != nil {
		return nil, err
	}
	testedRPCs = append(testedRPCs, &RpcResult{
		Method: GetTransactionCount,
		Status: Ok,
		Value:  nonce,
	})

	if rCtx.MaxPriorityFeePerGas, err = rCtx.EthCli.SuggestGasTipCap(context.Background()); err != nil {
		return nil, err
	}
	testedRPCs = append(testedRPCs, &RpcResult{
		Method: GetMaxPriorityFeePerGas,
		Status: Ok,
		Value:  rCtx.MaxPriorityFeePerGas.String(),
	})
	if rCtx.GasPrice, err = rCtx.EthCli.SuggestGasPrice(context.Background()); err != nil {
		return nil, err
	}
	testedRPCs = append(testedRPCs, &RpcResult{
		Method: GetGasPrice,
		Status: Ok,
		Value:  rCtx.GasPrice.String(),
	})

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   rCtx.ChainId,
		Nonce:     nonce,
		GasTipCap: rCtx.MaxPriorityFeePerGas,
		GasFeeCap: new(big.Int).Add(rCtx.GasPrice, big.NewInt(1000000000)),
		Gas:       10000000,
		Data:      common.FromHex(ContractByteCode),
	})

	// TODO: Make signer using types.MakeSigner with chain params
	signer := types.NewLondonSigner(rCtx.ChainId)
	signedTx, err := types.SignTx(tx, signer, rCtx.Acc.PrivKey)
	if err != nil {
		return nil, err
	}

	if err = rCtx.EthCli.SendTransaction(context.Background(), signedTx); err != nil {
		return nil, err
	}
	result := &RpcResult{
		Method: SendRawTransaction,
		Status: Ok,
		Value:  signedTx.Hash().Hex(),
	}
	testedRPCs = append(testedRPCs, result)

	// wait for the transaction to be mined
	tout, _ := time.ParseDuration(rCtx.Conf.Timeout)
	if err = WaitForTx(rCtx, signedTx.Hash(), tout); err != nil {
		return nil, err
	}

	if rCtx.ERC20Addr == (common.Address{}) {
		return nil, errors.New("contract address is empty, failed to deploy")
	}

	rCtx.AlreadyTestedRPCs = append(rCtx.AlreadyTestedRPCs, testedRPCs...)

	return result, nil
}

func RpcSendRawTransactionTransferERC20(rCtx *RpcContext) (*RpcResult, error) {
	// testedRPCs is a slice of RpcResult that will be appended to rCtx.AlreadyTestedRPCs
	// if the transaction is successfully sent
	var testedRPCs []*RpcResult
	var err error
	// Create a new transaction
	if rCtx.ChainId, err = rCtx.EthCli.ChainID(context.Background()); err != nil {
		return nil, err
	}
	testedRPCs = append(testedRPCs, &RpcResult{
		Method: GetChainId,
		Status: Ok,
		Value:  rCtx.ChainId.String(),
	})

	nonce, err := rCtx.EthCli.PendingNonceAt(context.Background(), rCtx.Acc.Address)
	if err != nil {
		return nil, err
	}
	testedRPCs = append(testedRPCs, &RpcResult{
		Method: GetTransactionCount,
		Status: Ok,
		Value:  nonce,
	})

	if rCtx.MaxPriorityFeePerGas, err = rCtx.EthCli.SuggestGasTipCap(context.Background()); err != nil {
		return nil, err
	}
	testedRPCs = append(testedRPCs, &RpcResult{
		Method: GetMaxPriorityFeePerGas,
		Status: Ok,
		Value:  rCtx.MaxPriorityFeePerGas.String(),
	})
	if rCtx.GasPrice, err = rCtx.EthCli.SuggestGasPrice(context.Background()); err != nil {
		return nil, err
	}
	testedRPCs = append(testedRPCs, &RpcResult{
		Method: GetGasPrice,
		Status: Ok,
		Value:  rCtx.GasPrice.String(),
	})

	randomRecipient := MustCreateRandomAccount().Address
	data, err := rCtx.ERC20Abi.Pack("transfer", randomRecipient, new(big.Int).SetUint64(1))
	if err != nil {
		log.Fatalf("Failed to pack transaction data: %v", err)
	}

	// Erc20 transfer
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   rCtx.ChainId,
		Nonce:     nonce,
		GasTipCap: rCtx.MaxPriorityFeePerGas,
		GasFeeCap: new(big.Int).Add(rCtx.GasPrice, big.NewInt(1000000000)),
		Gas:       10000000,
		To:        &rCtx.ERC20Addr,
		Data:      data,
	})

	// TODO: Make signer using types.MakeSigner with chain params
	signer := types.NewLondonSigner(rCtx.ChainId)
	signedTx, err := types.SignTx(tx, signer, rCtx.Acc.PrivKey)
	if err != nil {
		return nil, err
	}

	if err = rCtx.EthCli.SendTransaction(context.Background(), signedTx); err != nil {
		return nil, err
	}

	result := &RpcResult{
		Method: SendRawTransaction,
		Status: Ok,
		Value:  signedTx.Hash().Hex(),
	}
	testedRPCs = append(testedRPCs, result)

	// wait for the transaction to be mined
	tout, _ := time.ParseDuration(rCtx.Conf.Timeout)
	if err = WaitForTx(rCtx, signedTx.Hash(), tout); err != nil {
		return nil, err
	}

	rCtx.AlreadyTestedRPCs = append(rCtx.AlreadyTestedRPCs, testedRPCs...)

	return result, nil
}

func RpcGetBlockReceipts(rCtx *RpcContext) (*RpcResult, error) {
	if result := rCtx.AlreadyTested(GetBlockReceipts); result != nil {
		return result, nil
	}

	if len(rCtx.BlockNumsIncludingTx) == 0 {
		return nil, errors.New("no blocks with transactions")

	}

	// TODO: Random pick
	// pick a block with transactions
	blkNum := rCtx.BlockNumsIncludingTx[0]
	rpcBlockNum := rpc.BlockNumber(blkNum)
	receipts, err := rCtx.EthCli.BlockReceipts(context.Background(), rpc.BlockNumberOrHash{BlockNumber: &rpcBlockNum})
	if err != nil {
		return nil, err
	}

	result := &RpcResult{
		Method: GetBlockReceipts,
		Status: Ok,
		Value:  MustBeautifyReceipts(receipts),
	}
	rCtx.AlreadyTestedRPCs = append(rCtx.AlreadyTestedRPCs, result)

	return result, nil
}

func RpcGetTransactionByHash(rCtx *RpcContext) (*RpcResult, error) {
	if result := rCtx.AlreadyTested(GetTransactionByHash); result != nil {
		return result, nil
	}

	if len(rCtx.ProcessedTransactions) == 0 {
		return nil, errors.New("no transactions")
	}

	// TODO: Random pick
	txHash := rCtx.ProcessedTransactions[0]
	tx, _, err := rCtx.EthCli.TransactionByHash(context.Background(), txHash)
	if err != nil {
		return nil, err
	}

	result := &RpcResult{
		Method: GetTransactionByHash,
		Status: Ok,
		Value:  MustBeautifyTransaction(tx),
	}
	rCtx.AlreadyTestedRPCs = append(rCtx.AlreadyTestedRPCs, result)

	return result, nil
}

func RpcGetTransactionByBlockHashAndIndex(rCtx *RpcContext) (*RpcResult, error) {
	if result := rCtx.AlreadyTested(GetTransactionByBlockHashAndIndex); result != nil {
		return result, nil
	}

	if len(rCtx.BlockNumsIncludingTx) == 0 {
		return nil, errors.New("no blocks with transactions")
	}

	// TODO: Random pick
	blkNum := rCtx.BlockNumsIncludingTx[0]
	blk, err := rCtx.EthCli.BlockByNumber(context.Background(), new(big.Int).SetUint64(blkNum))
	if err != nil {
		return nil, err
	}

	if len(blk.Transactions()) == 0 {
		return nil, errors.New("no transactions in the block")
	}

	tx, err := rCtx.EthCli.TransactionInBlock(context.Background(), blk.Hash(), 0)
	if err != nil {
		return nil, err
	}

	result := &RpcResult{
		Method: GetTransactionByBlockHashAndIndex,
		Status: Ok,
		Value:  MustBeautifyTransaction(tx),
	}
	rCtx.AlreadyTestedRPCs = append(rCtx.AlreadyTestedRPCs, result)

	return result, nil
}

func RpcGetTransactionByBlockNumberAndIndex(rCtx *RpcContext) (*RpcResult, error) {
	if result := rCtx.AlreadyTested(GetTransactionByBlockNumberAndIndex); result != nil {
		return result, nil
	}

	if len(rCtx.BlockNumsIncludingTx) == 0 {
		return nil, errors.New("no blocks with transactions")
	}

	// TODO: Random pick
	blkNum := rCtx.BlockNumsIncludingTx[0]
	var tx types.Transaction
	if err := rCtx.EthCli.Client().CallContext(context.Background(), &tx, string(GetTransactionByBlockNumberAndIndex), blkNum, "0x0"); err != nil {
		return nil, err
	}

	result := &RpcResult{
		Method: GetTransactionByBlockNumberAndIndex,
		Status: Ok,
		Value:  MustBeautifyTransaction(&tx),
	}
	rCtx.AlreadyTestedRPCs = append(rCtx.AlreadyTestedRPCs, result)

	return result, nil
}

func RpcGetTransactionCountByHash(rCtx *RpcContext) (*RpcResult, error) {
	if result := rCtx.AlreadyTested(GetTransactionCountByHash); result != nil {
		return result, nil
	}

	if len(rCtx.BlockNumsIncludingTx) == 0 {
		return nil, errors.New("no transactions")
	}

	// get block
	blkNum := rCtx.BlockNumsIncludingTx[0]
	blk, err := rCtx.EthCli.BlockByNumber(context.Background(), new(big.Int).SetUint64(blkNum))
	if err != nil {
		return nil, err
	}

	var count uint64
	if err = rCtx.EthCli.Client().CallContext(context.Background(), &count, string(GetTransactionCountByHash), blk.Hash()); err != nil {
		return nil, err
	}

	result := &RpcResult{
		Method: GetTransactionCountByHash,
		Status: Ok,
		Value:  count,
	}
	rCtx.AlreadyTestedRPCs = append(rCtx.AlreadyTestedRPCs, result)

	return result, nil
}

func RpcGetTransactionReceipt(rCtx *RpcContext) (*RpcResult, error) {
	if result := rCtx.AlreadyTested(GetTransactionReceipt); result != nil {
		return result, nil
	}

	if len(rCtx.ProcessedTransactions) == 0 {
		return nil, errors.New("no transactions")
	}

	txHash := rCtx.ProcessedTransactions[0]
	receipt, err := rCtx.EthCli.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		return nil, err
	}

	result := &RpcResult{
		Method: GetTransactionReceipt,
		Status: Ok,
		Value:  MustBeautifyReceipt(receipt),
	}
	rCtx.AlreadyTestedRPCs = append(rCtx.AlreadyTestedRPCs, result)

	return result, nil
}

func RpcGetBlockTransactionCountByHash(rCtx *RpcContext) (*RpcResult, error) {
	if result := rCtx.AlreadyTested(GetBlockTransactionCountByHash); result != nil {
		return result, nil
	}

	if len(rCtx.BlockNumsIncludingTx) == 0 {
		return nil, errors.New("no blocks with transactions")
	}

	blkNum := rCtx.BlockNumsIncludingTx[0]
	blk, err := rCtx.EthCli.BlockByNumber(context.Background(), new(big.Int).SetUint64(blkNum))
	if err != nil {
		return nil, err
	}

	count, err := rCtx.EthCli.TransactionCount(context.Background(), blk.Hash())
	if err != nil {
		return nil, err
	}

	result := &RpcResult{
		Method: GetBlockTransactionCountByHash,
		Status: Ok,
		Value:  count,
	}
	rCtx.AlreadyTestedRPCs = append(rCtx.AlreadyTestedRPCs, result)

	return result, nil
}

func RpcGetCode(rCtx *RpcContext) (*RpcResult, error) {
	if result := rCtx.AlreadyTested(GetCode); result != nil {
		return result, nil
	}

	if rCtx.ERC20Addr == (common.Address{}) {
		return nil, errors.New("no contract address, must be deployed first")
	}

	code, err := rCtx.EthCli.CodeAt(context.Background(), rCtx.ERC20Addr, nil)
	if err != nil {
		return nil, err
	}

	result := &RpcResult{
		Method: GetCode,
		Status: Ok,
		Value:  hexutils.BytesToHex(code),
	}
	rCtx.AlreadyTestedRPCs = append(rCtx.AlreadyTestedRPCs, result)

	return result, nil
}

func RpcGetStorageAt(rCtx *RpcContext) (*RpcResult, error) {
	if result := rCtx.AlreadyTested(GetStorageAt); result != nil {
		return result, nil
	}

	if rCtx.ERC20Addr == (common.Address{}) {
		return nil, errors.New("no contract address, must be deployed first")
	}

	key := MustCalculateSlotKey(rCtx, 4)
	storage, err := rCtx.EthCli.StorageAt(context.Background(), rCtx.ERC20Addr, key, nil)
	if err != nil {
		return nil, err
	}

	var warnings []string
	status := Ok
	// check storage is zero
	if IsZeroBytes(storage) {
		warnings = append(warnings, "storage is zero bytes, should try another slot")
		status = Warning
	}

	result := &RpcResult{
		Method:   GetStorageAt,
		Status:   status,
		Value:    hexutils.BytesToHex(storage),
		Warnings: warnings,
	}
	rCtx.AlreadyTestedRPCs = append(rCtx.AlreadyTestedRPCs, result)

	return result, nil
}

func RpcNewFilter(rCtx *RpcContext) (*RpcResult, error) {
	if result := rCtx.AlreadyTested(NewFilter); result != nil {
		return result, nil
	}

	fErc20Transfer := ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(rCtx.BlockNumsIncludingTx[0] - 1),
		Addresses: []common.Address{rCtx.ERC20Addr},
		Topics: [][]common.Hash{
			{rCtx.ERC20Abi.Events["Transfer"].ID}, // Filter for Transfer event
		},
	}
	args, err := toFilterArg(fErc20Transfer)
	if err != nil {
		return nil, err
	}
	var rpcId string
	if err = rCtx.EthCli.Client().CallContext(context.Background(), &rpcId, string(NewFilter), args); err != nil {
		return nil, err
	}

	result := &RpcResult{
		Method: NewFilter,
		Status: Ok,
		Value:  rpcId,
	}
	rCtx.AlreadyTestedRPCs = append(rCtx.AlreadyTestedRPCs, result)
	rCtx.FilterId = rpcId
	rCtx.FilterQuery = fErc20Transfer

	return result, nil
}

func RpcGetFilterLogs(rCtx *RpcContext) (*RpcResult, error) {
	if result := rCtx.AlreadyTested(GetFilterLogs); result != nil {
		return result, nil
	}

	if rCtx.FilterId == "" {
		return nil, errors.New("no filter id, must create a filter first")
	}

	if _, err := RpcSendRawTransactionTransferERC20(rCtx); err != nil {
		return nil, errors.New("transfer ERC20 must be succeeded before checking filter logs")
	}

	var logs []types.Log
	if err := rCtx.EthCli.Client().CallContext(context.Background(), &logs, string(GetFilterLogs), rCtx.FilterId); err != nil {
		return nil, err
	}

	result := &RpcResult{
		Method: GetFilterLogs,
		Status: Ok,
		Value:  MustBeautifyLogs(logs),
	}
	rCtx.AlreadyTestedRPCs = append(rCtx.AlreadyTestedRPCs, result)

	return result, nil
}

func RpcNewBlockFilter(rCtx *RpcContext) (*RpcResult, error) {
	if result := rCtx.AlreadyTested(NewBlockFilter); result != nil {
		return result, nil
	}

	var rpcId string
	if err := rCtx.EthCli.Client().CallContext(context.Background(), &rpcId, string(NewBlockFilter)); err != nil {
		return nil, err
	}

	result := &RpcResult{
		Method: NewBlockFilter,
		Status: Ok,
		Value:  rpcId,
	}
	rCtx.AlreadyTestedRPCs = append(rCtx.AlreadyTestedRPCs, result)
	rCtx.BlockFilterId = rpcId

	return result, nil
}

func RpcGetFilterChanges(rCtx *RpcContext) (*RpcResult, error) {
	if result := rCtx.AlreadyTested(GetFilterChanges); result != nil {
		return result, nil
	}

	if rCtx.BlockFilterId == "" {
		return nil, errors.New("no block filter id, must create a block filter first")
	}

	// TODO: Make it configurable
	time.Sleep(3 * time.Second) // wait for a new block to be mined

	var changes []interface{}
	if err := rCtx.EthCli.Client().CallContext(context.Background(), &changes, string(GetFilterChanges), rCtx.BlockFilterId); err != nil {
		return nil, err
	}

	status := Ok
	warnings := []string{}
	if len(changes) == 0 {
		status = Warning
		warnings = append(warnings, "no new blocks")
	}

	result := &RpcResult{
		Method:   GetFilterChanges,
		Status:   status,
		Value:    changes,
		Warnings: warnings,
	}
	rCtx.AlreadyTestedRPCs = append(rCtx.AlreadyTestedRPCs, result)

	return result, nil
}

func RpcUninstallFilter(rCtx *RpcContext) (*RpcResult, error) {
	if result := rCtx.AlreadyTested(UninstallFilter); result != nil {
		return result, nil
	}

	if rCtx.FilterId == "" {
		return nil, errors.New("no filter id, must create a filter first")
	}

	var res bool
	if err := rCtx.EthCli.Client().CallContext(context.Background(), &res, string(UninstallFilter), rCtx.FilterId); err != nil {
		return nil, err
	}
	if !res {
		return nil, errors.New("uninstall filter failed")
	}

	if err := rCtx.EthCli.Client().CallContext(context.Background(), &res, string(UninstallFilter), rCtx.FilterId); err != nil {
		return nil, err
	}
	if res {
		return nil, errors.New("uninstall filter should be failed because it was already uninstalled")
	}

	result := &RpcResult{
		Method: UninstallFilter,
		Status: Ok,
		Value:  rCtx.FilterId,
	}
	rCtx.AlreadyTestedRPCs = append(rCtx.AlreadyTestedRPCs, result)

	return result, nil
}

func RpcGetLogs(rCtx *RpcContext) (*RpcResult, error) {
	if result := rCtx.AlreadyTested(GetLogs); result != nil {
		return result, nil
	}

	if _, err := RpcNewFilter(rCtx); err != nil {
		return nil, errors.New("failed to create a filter")
	}

	if _, err := RpcSendRawTransactionTransferERC20(rCtx); err != nil {
		return nil, errors.New("transfer ERC20 must be succeeded before checking filter logs")
	}

	// set from block because of limit
	logs, err := rCtx.EthCli.FilterLogs(context.Background(), rCtx.FilterQuery)
	if err != nil {
		return nil, err
	}

	status := Ok
	warnings := []string{}
	if len(logs) == 0 {
		status = Warning
		warnings = append(warnings, "no logs")
	}

	result := &RpcResult{
		Method:   GetLogs,
		Status:   status,
		Value:    MustBeautifyLogs(logs),
		Warnings: warnings,
	}
	rCtx.AlreadyTestedRPCs = append(rCtx.AlreadyTestedRPCs, result)

	return result, nil
}

func RpcEstimateGas(rCtx *RpcContext) (*RpcResult, error) {
	if result := rCtx.AlreadyTested(EstimateGas); result != nil {
		return result, nil
	}

	if rCtx.ERC20Addr == (common.Address{}) {
		return nil, errors.New("no contract address, must be deployed first")
	}

	data, err := rCtx.ERC20Abi.Pack("transfer", rCtx.Acc.Address, new(big.Int).SetUint64(1))
	if err != nil {
		log.Fatalf("Failed to pack transaction data: %v", err)
	}

	msg := ethereum.CallMsg{
		From: rCtx.Acc.Address,
		To:   &rCtx.ERC20Addr,
		Data: data,
	}
	gas, err := rCtx.EthCli.EstimateGas(context.Background(), msg)
	if err != nil {
		return nil, err
	}

	result := &RpcResult{
		Method: EstimateGas,
		Status: Ok,
		Value:  gas,
	}
	rCtx.AlreadyTestedRPCs = append(rCtx.AlreadyTestedRPCs, result)

	return result, nil
}

func RPCCall(rCtx *RpcContext) (*RpcResult, error) {
	if result := rCtx.AlreadyTested(Call); result != nil {
		return result, nil
	}

	if rCtx.ERC20Addr == (common.Address{}) {
		return nil, errors.New("no contract address, must be deployed first")
	}

	data, err := rCtx.ERC20Abi.Pack("balanceOf", rCtx.Acc.Address)
	if err != nil {
		log.Fatalf("Failed to pack transaction data: %v", err)
	}

	msg := ethereum.CallMsg{
		To:   &rCtx.ERC20Addr,
		Data: data,
	}
	res, err := rCtx.EthCli.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}

	result := &RpcResult{
		Method: Call,
		Status: Ok,
		Value:  hexutils.BytesToHex(res),
	}
	rCtx.AlreadyTestedRPCs = append(rCtx.AlreadyTestedRPCs, result)

	return result, nil
}

func WaitForTx(rCtx *RpcContext, txHash common.Hash, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(500 * time.Millisecond) // Check every 500ms
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout exceeded while waiting for transaction %s", txHash.Hex())
		case <-ticker.C:
			receipt, err := rCtx.EthCli.TransactionReceipt(context.Background(), txHash)
			if err != nil && !errors.Is(err, ethereum.NotFound) {
				return err
			}
			if err == nil {
				rCtx.ProcessedTransactions = append(rCtx.ProcessedTransactions, txHash)
				rCtx.BlockNumsIncludingTx = append(rCtx.BlockNumsIncludingTx, receipt.BlockNumber.Uint64())
				rCtx.AlreadyTestedRPCs = append(rCtx.AlreadyTestedRPCs, &RpcResult{
					Method: GetTransactionReceipt,
					Status: Ok,
					Value:  MustBeautifyReceipt(receipt),
				})
				if receipt.ContractAddress != (common.Address{}) {
					rCtx.ERC20Addr = receipt.ContractAddress
				}
				if receipt.Status == 0 {
					return fmt.Errorf("transaction %s failed", txHash.Hex())
				}
				return nil
			}
		}
	}
}
