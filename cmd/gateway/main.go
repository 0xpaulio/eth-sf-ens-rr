package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
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

func Must(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("%s must be set", key)
	}
	return val
}

func GetOrDefault(key, def string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	return val
}

func DecodeHex(h string) ([]byte, error) {
	return hex.DecodeString(strings.TrimPrefix(h, "0x"))
}

func mustParseABI(json string) *ethabi.ABI {
	a, err := ethabi.JSON(strings.NewReader(json))
	if err != nil {
		log.Fatal("failed to parse ABI", err)
	}
	return &a
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

/*
decode calldata
calculate storage slot
generate proof
encode response
*/

type Gateway struct {
	l2GethClient      *gethclient.Client
	l2EthClient       *ethclient.Client
	l2ResolverAddress common.Address
	selectors         map[[4]byte]string
}

func main() {
	l2RPCURL := GetOrDefault("L2_RPC_URL", "https://goerli.optimism.io")
	l2RPCClient, err := rpc.Dial(l2RPCURL)
	if err != nil {
		log.Fatal("dialing L2 RPC", err)
	}
	selectors := make(map[[4]byte]string, len(methodNames))
	for _, name := range methodNames {
		selectors[mustGetSelector(abi, name)] = name
	}

	gateway := Gateway{
		l2GethClient:      gethclient.New(l2RPCClient),
		l2EthClient:       ethclient.NewClient(l2RPCClient),
		l2ResolverAddress: common.HexToAddress(GetOrDefault("L2_RESOLVER_ADDR", "0xbeefbabe")),
		selectors:         selectors,
	}

	logger := httplog.NewLogger("httplog-example", httplog.Options{
		LogLevel: "Debug",
		Concise:  true,
	})

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(httplog.RequestLogger(logger))
	r.Get("/block_number", gateway.getBlockNumber)
	r.Get("/gateway/{sender}/{data}.json", gateway.getGateway)
	http.ListenAndServe(":3000", r)
}

func (g *Gateway) getBlockNumber(w http.ResponseWriter, r *http.Request) {
	blockNum, err := g.l2EthClient.BlockNumber(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
	}
	w.Write([]byte(fmt.Sprintf("block_number: %d", blockNum)))
}

func (g *Gateway) decode(hexCalldata string) (decoded []interface{}, err error) {
	calldata, err := DecodeHex(hexCalldata)
	if err != nil {
		return nil, fmt.Errorf("decoding hex calldata: %w", err)
	}
	functionSignature := calldata[:4]
	functionParameters := calldata[4:]
	var sig [4]byte
	copy(sig[:], functionSignature)
	methodName, ok := g.selectors[sig]
	if !ok {
		return nil, fmt.Errorf("unknown function signature: %s", string(functionSignature))
	}
	return abi.Methods[methodName].Inputs.Unpack(functionParameters)
}

func (g *Gateway) getGateway(w http.ResponseWriter, r *http.Request) {
	sender := chi.URLParam(r, "sender")
	hexCalldata := chi.URLParam(r, "data")
	log.Printf("sender: %s, hexCalldata: %s", sender, hexCalldata)
	decoded, err := g.decode(hexCalldata)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	log.Printf("decoded: %+v", decoded...)
	// g.l2GethClient.GetProof()
	w.Write([]byte(fmt.Sprintf("decoded: %+v", decoded)))
}

// calldata
// 0x02571be36e4409edbffcd311c18c7dbfadfc428daf0edd143693491dd2b4782bc494dc61
