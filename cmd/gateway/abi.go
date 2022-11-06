package main

import (
	"log"
	"strings"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
)

var abiJSON string = `
[
    {
      "inputs": [
        {
          "internalType": "uint64",
          "name": "_chainId",
          "type": "uint64"
        },
        {
          "internalType": "address",
          "name": "_l2Registrar",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "_ovmAddressManager",
          "type": "address"
        },
        {
          "internalType": "string[]",
          "name": "_gatewayUrls",
          "type": "string[]"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "constructor"
    },
    {
      "inputs": [],
      "name": "AccountDNE",
      "type": "error"
    },
    {
      "inputs": [],
      "name": "InvalidStateRoot",
      "type": "error"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "sender",
          "type": "address"
        },
        {
          "internalType": "string[]",
          "name": "urls",
          "type": "string[]"
        },
        {
          "internalType": "bytes",
          "name": "callData",
          "type": "bytes"
        },
        {
          "internalType": "bytes4",
          "name": "callbackFunction",
          "type": "bytes4"
        },
        {
          "internalType": "bytes",
          "name": "extraData",
          "type": "bytes"
        }
      ],
      "name": "OffchainLookup",
      "type": "error"
    },
    {
      "inputs": [],
      "name": "StorageDNE",
      "type": "error"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "chainId",
          "type": "uint256"
        },
        {
          "internalType": "address",
          "name": "contractAddress",
          "type": "address"
        }
      ],
      "name": "StorageHandledByL2",
      "type": "error"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "owner",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "operator",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "bool",
          "name": "approved",
          "type": "bool"
        }
      ],
      "name": "ApprovalForAll",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "bytes32",
          "name": "node",
          "type": "bytes32"
        },
        {
          "indexed": true,
          "internalType": "bytes32",
          "name": "label",
          "type": "bytes32"
        },
        {
          "indexed": false,
          "internalType": "address",
          "name": "owner",
          "type": "address"
        }
      ],
      "name": "NewOwner",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "bytes32",
          "name": "node",
          "type": "bytes32"
        },
        {
          "indexed": false,
          "internalType": "address",
          "name": "resolver",
          "type": "address"
        }
      ],
      "name": "NewResolver",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "bytes32",
          "name": "node",
          "type": "bytes32"
        },
        {
          "indexed": false,
          "internalType": "uint64",
          "name": "ttl",
          "type": "uint64"
        }
      ],
      "name": "NewTTL",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "bytes32",
          "name": "node",
          "type": "bytes32"
        },
        {
          "indexed": false,
          "internalType": "address",
          "name": "owner",
          "type": "address"
        }
      ],
      "name": "Transfer",
      "type": "event"
    },
    {
      "inputs": [],
      "name": "L2_REGISTRY_CHAIN_ID",
      "outputs": [
        {
          "internalType": "uint64",
          "name": "",
          "type": "uint64"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "L2_REGISTRY_CONTRACT_ADDRESS",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "SLO__L2_REGISTRY__OPERATORS",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "SLO__L2_REGISTRY__RECORDS",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "name": "gatewayUrls",
      "outputs": [
        {
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_owner",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "_operator",
          "type": "address"
        }
      ],
      "name": "isApprovedForAll",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_owner",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "_operator",
          "type": "address"
        },
        {
          "components": [
            {
              "internalType": "bytes32",
              "name": "stateRoot",
              "type": "bytes32"
            },
            {
              "components": [
                {
                  "internalType": "uint256",
                  "name": "batchIndex",
                  "type": "uint256"
                },
                {
                  "internalType": "bytes32",
                  "name": "batchRoot",
                  "type": "bytes32"
                },
                {
                  "internalType": "uint256",
                  "name": "batchSize",
                  "type": "uint256"
                },
                {
                  "internalType": "uint256",
                  "name": "prevTotalElements",
                  "type": "uint256"
                },
                {
                  "internalType": "bytes",
                  "name": "extraData",
                  "type": "bytes"
                }
              ],
              "internalType": "struct Lib_OVMCodec.ChainBatchHeader",
              "name": "stateRootBatchHeader",
              "type": "tuple"
            },
            {
              "components": [
                {
                  "internalType": "uint256",
                  "name": "index",
                  "type": "uint256"
                },
                {
                  "internalType": "bytes32[]",
                  "name": "siblings",
                  "type": "bytes32[]"
                }
              ],
              "internalType": "struct Lib_OVMCodec.ChainInclusionProof",
              "name": "stateRootProof",
              "type": "tuple"
            },
            {
              "internalType": "bytes",
              "name": "stateTrieWitness",
              "type": "bytes"
            },
            {
              "internalType": "bytes",
              "name": "storageTrieWitness",
              "type": "bytes"
            }
          ],
          "internalType": "struct OptimismHelper.L2StateProof",
          "name": "_proof",
          "type": "tuple"
        }
      ],
      "name": "isApprovedForAllWithProof",
      "outputs": [
        {
          "internalType": "bool",
          "name": "approved_",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "libAddressManager",
      "outputs": [
        {
          "internalType": "contract Lib_AddressManager",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "bytes32",
          "name": "_node",
          "type": "bytes32"
        }
      ],
      "name": "owner",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "bytes32",
          "name": "_node",
          "type": "bytes32"
        },
        {
          "components": [
            {
              "internalType": "bytes32",
              "name": "stateRoot",
              "type": "bytes32"
            },
            {
              "components": [
                {
                  "internalType": "uint256",
                  "name": "batchIndex",
                  "type": "uint256"
                },
                {
                  "internalType": "bytes32",
                  "name": "batchRoot",
                  "type": "bytes32"
                },
                {
                  "internalType": "uint256",
                  "name": "batchSize",
                  "type": "uint256"
                },
                {
                  "internalType": "uint256",
                  "name": "prevTotalElements",
                  "type": "uint256"
                },
                {
                  "internalType": "bytes",
                  "name": "extraData",
                  "type": "bytes"
                }
              ],
              "internalType": "struct Lib_OVMCodec.ChainBatchHeader",
              "name": "stateRootBatchHeader",
              "type": "tuple"
            },
            {
              "components": [
                {
                  "internalType": "uint256",
                  "name": "index",
                  "type": "uint256"
                },
                {
                  "internalType": "bytes32[]",
                  "name": "siblings",
                  "type": "bytes32[]"
                }
              ],
              "internalType": "struct Lib_OVMCodec.ChainInclusionProof",
              "name": "stateRootProof",
              "type": "tuple"
            },
            {
              "internalType": "bytes",
              "name": "stateTrieWitness",
              "type": "bytes"
            },
            {
              "internalType": "bytes",
              "name": "storageTrieWitness",
              "type": "bytes"
            }
          ],
          "internalType": "struct OptimismHelper.L2StateProof",
          "name": "_proof",
          "type": "tuple"
        }
      ],
      "name": "ownerWithProof",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "bytes32",
          "name": "_node",
          "type": "bytes32"
        }
      ],
      "name": "recordExists",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "bytes32",
          "name": "_node",
          "type": "bytes32"
        },
        {
          "components": [
            {
              "internalType": "bytes32",
              "name": "stateRoot",
              "type": "bytes32"
            },
            {
              "components": [
                {
                  "internalType": "uint256",
                  "name": "batchIndex",
                  "type": "uint256"
                },
                {
                  "internalType": "bytes32",
                  "name": "batchRoot",
                  "type": "bytes32"
                },
                {
                  "internalType": "uint256",
                  "name": "batchSize",
                  "type": "uint256"
                },
                {
                  "internalType": "uint256",
                  "name": "prevTotalElements",
                  "type": "uint256"
                },
                {
                  "internalType": "bytes",
                  "name": "extraData",
                  "type": "bytes"
                }
              ],
              "internalType": "struct Lib_OVMCodec.ChainBatchHeader",
              "name": "stateRootBatchHeader",
              "type": "tuple"
            },
            {
              "components": [
                {
                  "internalType": "uint256",
                  "name": "index",
                  "type": "uint256"
                },
                {
                  "internalType": "bytes32[]",
                  "name": "siblings",
                  "type": "bytes32[]"
                }
              ],
              "internalType": "struct Lib_OVMCodec.ChainInclusionProof",
              "name": "stateRootProof",
              "type": "tuple"
            },
            {
              "internalType": "bytes",
              "name": "stateTrieWitness",
              "type": "bytes"
            },
            {
              "internalType": "bytes",
              "name": "storageTrieWitness",
              "type": "bytes"
            }
          ],
          "internalType": "struct OptimismHelper.L2StateProof",
          "name": "_proof",
          "type": "tuple"
        }
      ],
      "name": "recordExistsWithProof",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "_name",
          "type": "string"
        }
      ],
      "name": "resolve",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "bytes32",
          "name": "_node",
          "type": "bytes32"
        }
      ],
      "name": "resolver",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "bytes32",
          "name": "_node",
          "type": "bytes32"
        },
        {
          "components": [
            {
              "internalType": "bytes32",
              "name": "stateRoot",
              "type": "bytes32"
            },
            {
              "components": [
                {
                  "internalType": "uint256",
                  "name": "batchIndex",
                  "type": "uint256"
                },
                {
                  "internalType": "bytes32",
                  "name": "batchRoot",
                  "type": "bytes32"
                },
                {
                  "internalType": "uint256",
                  "name": "batchSize",
                  "type": "uint256"
                },
                {
                  "internalType": "uint256",
                  "name": "prevTotalElements",
                  "type": "uint256"
                },
                {
                  "internalType": "bytes",
                  "name": "extraData",
                  "type": "bytes"
                }
              ],
              "internalType": "struct Lib_OVMCodec.ChainBatchHeader",
              "name": "stateRootBatchHeader",
              "type": "tuple"
            },
            {
              "components": [
                {
                  "internalType": "uint256",
                  "name": "index",
                  "type": "uint256"
                },
                {
                  "internalType": "bytes32[]",
                  "name": "siblings",
                  "type": "bytes32[]"
                }
              ],
              "internalType": "struct Lib_OVMCodec.ChainInclusionProof",
              "name": "stateRootProof",
              "type": "tuple"
            },
            {
              "internalType": "bytes",
              "name": "stateTrieWitness",
              "type": "bytes"
            },
            {
              "internalType": "bytes",
              "name": "storageTrieWitness",
              "type": "bytes"
            }
          ],
          "internalType": "struct OptimismHelper.L2StateProof",
          "name": "_proof",
          "type": "tuple"
        }
      ],
      "name": "resolverWithProof",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "operator",
          "type": "address"
        },
        {
          "internalType": "bool",
          "name": "approved",
          "type": "bool"
        }
      ],
      "name": "setApprovalForAll",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "bytes32",
          "name": "node",
          "type": "bytes32"
        },
        {
          "internalType": "address",
          "name": "owner",
          "type": "address"
        }
      ],
      "name": "setOwner",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "bytes32",
          "name": "node",
          "type": "bytes32"
        },
        {
          "internalType": "address",
          "name": "owner",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "resolver",
          "type": "address"
        },
        {
          "internalType": "uint64",
          "name": "ttl",
          "type": "uint64"
        }
      ],
      "name": "setRecord",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "bytes32",
          "name": "node",
          "type": "bytes32"
        },
        {
          "internalType": "address",
          "name": "resolver",
          "type": "address"
        }
      ],
      "name": "setResolver",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "bytes32",
          "name": "node",
          "type": "bytes32"
        },
        {
          "internalType": "bytes32",
          "name": "label",
          "type": "bytes32"
        },
        {
          "internalType": "address",
          "name": "owner",
          "type": "address"
        }
      ],
      "name": "setSubnodeOwner",
      "outputs": [
        {
          "internalType": "bytes32",
          "name": "",
          "type": "bytes32"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "bytes32",
          "name": "node",
          "type": "bytes32"
        },
        {
          "internalType": "bytes32",
          "name": "label",
          "type": "bytes32"
        },
        {
          "internalType": "address",
          "name": "owner",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "resolver",
          "type": "address"
        },
        {
          "internalType": "uint64",
          "name": "ttl",
          "type": "uint64"
        }
      ],
      "name": "setSubnodeRecord",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "bytes32",
          "name": "node",
          "type": "bytes32"
        },
        {
          "internalType": "uint64",
          "name": "ttl",
          "type": "uint64"
        }
      ],
      "name": "setTTL",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "bytes32",
          "name": "_node",
          "type": "bytes32"
        }
      ],
      "name": "ttl",
      "outputs": [
        {
          "internalType": "uint64",
          "name": "",
          "type": "uint64"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "bytes32",
          "name": "_node",
          "type": "bytes32"
        },
        {
          "components": [
            {
              "internalType": "bytes32",
              "name": "stateRoot",
              "type": "bytes32"
            },
            {
              "components": [
                {
                  "internalType": "uint256",
                  "name": "batchIndex",
                  "type": "uint256"
                },
                {
                  "internalType": "bytes32",
                  "name": "batchRoot",
                  "type": "bytes32"
                },
                {
                  "internalType": "uint256",
                  "name": "batchSize",
                  "type": "uint256"
                },
                {
                  "internalType": "uint256",
                  "name": "prevTotalElements",
                  "type": "uint256"
                },
                {
                  "internalType": "bytes",
                  "name": "extraData",
                  "type": "bytes"
                }
              ],
              "internalType": "struct Lib_OVMCodec.ChainBatchHeader",
              "name": "stateRootBatchHeader",
              "type": "tuple"
            },
            {
              "components": [
                {
                  "internalType": "uint256",
                  "name": "index",
                  "type": "uint256"
                },
                {
                  "internalType": "bytes32[]",
                  "name": "siblings",
                  "type": "bytes32[]"
                }
              ],
              "internalType": "struct Lib_OVMCodec.ChainInclusionProof",
              "name": "stateRootProof",
              "type": "tuple"
            },
            {
              "internalType": "bytes",
              "name": "stateTrieWitness",
              "type": "bytes"
            },
            {
              "internalType": "bytes",
              "name": "storageTrieWitness",
              "type": "bytes"
            }
          ],
          "internalType": "struct OptimismHelper.L2StateProof",
          "name": "_proof",
          "type": "tuple"
        }
      ],
      "name": "ttlWithProof",
      "outputs": [
        {
          "internalType": "uint64",
          "name": "",
          "type": "uint64"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    }
  ]
`

var abi *ethabi.ABI = mustParseABI(abiJSON)

var methodNames = []string{
	"owner",
	"resolver",
	"ttl",
	"recordExists",
	"isApprovedForAll",
}

func mustGetSelector(parsedABI *ethabi.ABI, methodName string) [4]byte {
	method, ok := parsedABI.Methods[methodName]
	if !ok {
		log.Fatal("failed to find method", methodName)
	}
	var selector [4]byte
	copy(selector[:], method.ID)
	return selector
}

func mustParseABI(json string) *ethabi.ABI {
	a, err := ethabi.JSON(strings.NewReader(json))
	if err != nil {
		log.Fatal("failed to parse ABI", err)
	}
	return &a
}
