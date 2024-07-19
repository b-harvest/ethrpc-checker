package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// Copy compiled contract bytecode
const ContractByteCode = "60806040526040518060400160405280600881526020017f4d7920546f6b656e000000000000000000000000000000000000000000000000815250600090816100489190610382565b506040518060400160405280600381526020017f4d544b00000000000000000000000000000000000000000000000000000000008152506001908161008d9190610382565b506012600260006101000a81548160ff021916908360ff160217905550600260009054906101000a900460ff1660ff16600a6100c991906105b6565b620f42406100d79190610601565b6003553480156100e657600080fd5b50600354600460003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550610643565b600081519050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b600060028204905060018216806101b357607f821691505b6020821081036101c6576101c561016c565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b60006008830261022e7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff826101f1565b61023886836101f1565b95508019841693508086168417925050509392505050565b6000819050919050565b6000819050919050565b600061027f61027a61027584610250565b61025a565b610250565b9050919050565b6000819050919050565b61029983610264565b6102ad6102a582610286565b8484546101fe565b825550505050565b600090565b6102c26102b5565b6102cd818484610290565b505050565b5b818110156102f1576102e66000826102ba565b6001810190506102d3565b5050565b601f82111561033657610307816101cc565b610310846101e1565b8101602085101561031f578190505b61033361032b856101e1565b8301826102d2565b50505b505050565b600082821c905092915050565b60006103596000198460080261033b565b1980831691505092915050565b60006103728383610348565b9150826002028217905092915050565b61038b82610132565b67ffffffffffffffff8111156103a4576103a361013d565b5b6103ae825461019b565b6103b98282856102f5565b600060209050601f8311600181146103ec57600084156103da578287015190505b6103e48582610366565b86555061044c565b601f1984166103fa866101cc565b60005b82811015610422578489015182556001820191506020850194506020810190506103fd565b8683101561043f578489015161043b601f891682610348565b8355505b6001600288020188555050505b505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60008160011c9050919050565b6000808291508390505b60018511156104da578086048111156104b6576104b5610454565b5b60018516156104c55780820291505b80810290506104d385610483565b945061049a565b94509492505050565b6000826104f357600190506105af565b8161050157600090506105af565b8160018114610517576002811461052157610550565b60019150506105af565b60ff84111561053357610532610454565b5b8360020a91508482111561054a57610549610454565b5b506105af565b5060208310610133831016604e8410600b84101617156105855782820a9050838111156105805761057f610454565b5b6105af565b6105928484846001610490565b925090508184048111156105a9576105a8610454565b5b81810290505b9392505050565b60006105c182610250565b91506105cc83610250565b92506105f97fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff84846104e3565b905092915050565b600061060c82610250565b915061061783610250565b925082820261062581610250565b9150828204841483151761063c5761063b610454565b5b5092915050565b610ddc806106526000396000f3fe608060405234801561001057600080fd5b50600436106100935760003560e01c8063313ce56711610066578063313ce5671461013457806370a082311461015257806395d89b4114610182578063a9059cbb146101a0578063dd62ed3e146101d057610093565b806306fdde0314610098578063095ea7b3146100b657806318160ddd146100e657806323b872dd14610104575b600080fd5b6100a0610200565b6040516100ad9190610985565b60405180910390f35b6100d060048036038101906100cb9190610a40565b61028e565b6040516100dd9190610a9b565b60405180910390f35b6100ee610380565b6040516100fb9190610ac5565b60405180910390f35b61011e60048036038101906101199190610ae0565b610386565b60405161012b9190610a9b565b60405180910390f35b61013c610678565b6040516101499190610b4f565b60405180910390f35b61016c60048036038101906101679190610b6a565b61068b565b6040516101799190610ac5565b60405180910390f35b61018a6106a3565b6040516101979190610985565b60405180910390f35b6101ba60048036038101906101b59190610a40565b610731565b6040516101c79190610a9b565b60405180910390f35b6101ea60048036038101906101e59190610b97565b6108d0565b6040516101f79190610ac5565b60405180910390f35b6000805461020d90610c06565b80601f016020809104026020016040519081016040528092919081815260200182805461023990610c06565b80156102865780601f1061025b57610100808354040283529160200191610286565b820191906000526020600020905b81548152906001019060200180831161026957829003601f168201915b505050505081565b600081600560003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9258460405161036e9190610ac5565b60405180910390a36001905092915050565b60035481565b600081600460008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054101561040a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161040190610c83565b60405180910390fd5b81600560008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205410156104c9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104c090610cef565b60405180910390fd5b81600460008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546105189190610d3e565b9250508190555081600460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825461056e9190610d72565b9250508190555081600560008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546106019190610d3e565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516106659190610ac5565b60405180910390a3600190509392505050565b600260009054906101000a900460ff1681565b60046020528060005260406000206000915090505481565b600180546106b090610c06565b80601f01602080910402602001604051908101604052809291908181526020018280546106dc90610c06565b80156107295780601f106106fe57610100808354040283529160200191610729565b820191906000526020600020905b81548152906001019060200180831161070c57829003601f168201915b505050505081565b600081600460003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205410156107b5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107ac90610c83565b60405180910390fd5b81600460003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546108049190610d3e565b9250508190555081600460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825461085a9190610d72565b925050819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516108be9190610ac5565b60405180910390a36001905092915050565b6005602052816000526040600020602052806000526040600020600091509150505481565b600081519050919050565b600082825260208201905092915050565b60005b8381101561092f578082015181840152602081019050610914565b60008484015250505050565b6000601f19601f8301169050919050565b6000610957826108f5565b6109618185610900565b9350610971818560208601610911565b61097a8161093b565b840191505092915050565b6000602082019050818103600083015261099f818461094c565b905092915050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006109d7826109ac565b9050919050565b6109e7816109cc565b81146109f257600080fd5b50565b600081359050610a04816109de565b92915050565b6000819050919050565b610a1d81610a0a565b8114610a2857600080fd5b50565b600081359050610a3a81610a14565b92915050565b60008060408385031215610a5757610a566109a7565b5b6000610a65858286016109f5565b9250506020610a7685828601610a2b565b9150509250929050565b60008115159050919050565b610a9581610a80565b82525050565b6000602082019050610ab06000830184610a8c565b92915050565b610abf81610a0a565b82525050565b6000602082019050610ada6000830184610ab6565b92915050565b600080600060608486031215610af957610af86109a7565b5b6000610b07868287016109f5565b9350506020610b18868287016109f5565b9250506040610b2986828701610a2b565b9150509250925092565b600060ff82169050919050565b610b4981610b33565b82525050565b6000602082019050610b646000830184610b40565b92915050565b600060208284031215610b8057610b7f6109a7565b5b6000610b8e848285016109f5565b91505092915050565b60008060408385031215610bae57610bad6109a7565b5b6000610bbc858286016109f5565b9250506020610bcd858286016109f5565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680610c1e57607f821691505b602082108103610c3157610c30610bd7565b5b50919050565b7f496e73756666696369656e742062616c616e6365000000000000000000000000600082015250565b6000610c6d601483610900565b9150610c7882610c37565b602082019050919050565b60006020820190508181036000830152610c9c81610c60565b9050919050565b7f416c6c6f77616e63652065786365656465640000000000000000000000000000600082015250565b6000610cd9601283610900565b9150610ce482610ca3565b602082019050919050565b60006020820190508181036000830152610d0881610ccc565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610d4982610a0a565b9150610d5483610a0a565b9250828203905081811115610d6c57610d6b610d0f565b5b92915050565b6000610d7d82610a0a565b9150610d8883610a0a565b9250828201905080821115610da057610d9f610d0f565b5b9291505056fea2646970667358221220b10024f1dc62ba720dc6c2479d79b18cd10c8432e7bb7d00a7003947039bee0164736f6c63430008190033"

func main() {
	verbose := flag.Bool("v", false, "Enable verbose output")
	outputExcel := flag.Bool("xlsx", false, "Save output as xlsx")
	flag.Parse()

	// Load configuration from config.yaml
	config := MustLoadConfig("config.yaml")

	rCtx, err := NewContext(config)
	if err != nil {
		log.Fatalf("Failed to create context: %v", err)
	}

	rCtx = MustLoadContractInfo(rCtx)

	// Collect json rpc results
	var results []*RpcResult

	rpcs := []struct {
		name RpcName
		test RpcCall
	}{
		{SendRawTransaction, RpcSendRawTransactionTransferValue},
		{SendRawTransaction, RpcSendRawTransactionDeployContract},
		{SendRawTransaction, RpcSendRawTransactionTransferERC20},
		{GetBlockNumber, RpcGetBlockNumber},
		{GetGasPrice, RpcGetGasPrice},
		{GetMaxPriorityFeePerGas, RpcGetMaxPriorityFeePerGas},
		{GetChainId, RpcGetChainId},
		{GetBalance, RpcGetBalance},
		{GetTransactionCount, RpcGetTransactionCount},
		{GetBlockByHash, RpcGetBlockByHash},
		{GetBlockByNumber, RpcGetBlockByNumber},
		{GetBlockReceipts, RpcGetBlockReceipts},
		{GetTransactionByHash, RpcGetTransactionByHash},
		{GetTransactionByBlockHashAndIndex, RpcGetTransactionByBlockHashAndIndex},
		{GetTransactionByBlockNumberAndIndex, RpcGetTransactionByBlockNumberAndIndex},
		{GetTransactionReceipt, RpcGetTransactionReceipt},
		{GetTransactionCountByHash, RpcGetTransactionCountByHash},
		{GetBlockTransactionCountByHash, RpcGetBlockTransactionCountByHash},
		{GetCode, RpcGetCode},
		{GetStorageAt, RpcGetStorageAt},
		{NewFilter, RpcNewFilter},
		{GetFilterLogs, RpcGetFilterLogs},
		{NewBlockFilter, RpcNewBlockFilter},
		{GetFilterChanges, RpcGetFilterChanges},
		{UninstallFilter, RpcUninstallFilter},
		{GetLogs, RpcGetLogs},
	}

	for _, r := range rpcs {
		_, err := r.test(rCtx)
		if err != nil {
			// add error to results
			results = append(results, &RpcResult{
				Method: r.name,
				Status: Error,
				ErrMsg: err.Error(),
			})
			continue
		}
	}
	results = append(results, rCtx.AlreadyTestedRPCs...)

	ReportResults(results, *verbose, *outputExcel)
}

func MustLoadContractInfo(rCtx *RpcContext) *RpcContext {
	// Read the ABI file
	abiFile, err := os.ReadFile("ERC20Token.abi")
	if err != nil {
		log.Fatalf("Failed to read ABI file: %v", err)
	}
	// Parse the ABI
	parsedABI, err := abi.JSON(strings.NewReader(string(abiFile)))
	if err != nil {
		log.Fatalf("Failed to parse ERC20 ABI: %v", err)
	}
	rCtx.ERC20Abi = &parsedABI
	// Read the compiled contract bytecode
	bytecode, err := os.ReadFile("ERC20Token.bin")
	if err != nil {
		log.Fatalf("Failed to read contract bytecode: %v", err)
	}
	// Decode the hex string to bytes
	contractBytecode := common.FromHex(string(bytecode))
	rCtx.ERC20ByteCode = contractBytecode

	return rCtx
}
