package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/google/go-cmp/cmp"
)

// GethVersion is the version of the Geth client used in the tests
// Update it when go-ethereum of go.mod is updated
const GethVersion = "1.14.7"

type RpcName string
type RpcCall func(rCtx *RpcContext) (*RpcResult, error)

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

func RpcSendRawTransaction(rCtx *RpcContext) (*RpcResult, error) {
	if result := rCtx.AlreadyTested(SendRawTransaction); result != nil {
		return result, nil
	}

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

	err = rCtx.EthCli.SendTransaction(context.Background(), signedTx)
	if err != nil {
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
				return nil
			}
		}
	}
}
