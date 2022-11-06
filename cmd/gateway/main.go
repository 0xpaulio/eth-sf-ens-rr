package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
)

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

type GatewayResponse struct {
	Data string `json:"data"`
}

type GatewayResponseData struct {
	StateRootProof
	StateTrieWitness   string
	StorageTrieWitness string
}

type StateRootProofInput struct {
	ResolverAddr string `json:"resolver_addr"`
	AddrSlot     string `json:"addr_slot"`
}

type StateRootProof struct {
	StateRoot            string `json:"stateRoot"`
	StateRootBatchHeader struct {
		BatchIndex struct {
			Type string `json:"type"`
			Hex  string `json:"hex"`
		} `json:"batchIndex"`
		BatchRoot string `json:"batchRoot"`
		BatchSize struct {
			Type string `json:"type"`
			Hex  string `json:"hex"`
		} `json:"batchSize"`
		PrevTotalElements struct {
			Type string `json:"type"`
			Hex  string `json:"hex"`
		} `json:"prevTotalElements"`
		ExtraData string `json:"extraData"`
	} `json:"stateRootBatchHeader"`
	StateRootProof struct {
		Index    int `json:"index"`
		Siblings []struct {
			Type string `json:"type"`
			Data []byte `json:"data"`
		} `json:"siblings"`
	} `json:"stateRootProof"`
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
		l2ResolverAddress: common.HexToAddress(GetOrDefault("L2_RESOLVER_ADDR", "0xE933897412cc2164331e542B2a2Be491612C233F")),
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
	r.Get("/gateway/{sender}/{data}.json", gateway.getGateway)
	http.ListenAndServe(":41234", r)
}

func getoSLOForResolver(node [32]byte) []byte {
	resp := make([]byte, 32)
	new(big.Int).Add(new(big.Int).SetBytes(getSLOForRecords(node)), big.NewInt(1)).FillBytes(resp)
	return resp
}

func getoSLOForOwner(node [32]byte) []byte {
	return getSLOForRecords(node)
}

func getSLOForRecords(node [32]byte) []byte {
	return crypto.Keccak256(
		append(node[:], make([]byte, 32)...),
	)
}

func getSLO(methodName string, decodedCalldata []interface{}) ([]byte, error) {
	switch methodName {
	case "owner":
		node, ok := decodedCalldata[0].([32]byte)
		if !ok {
			return nil, errors.New("get SLO for owner")
		}
		return getoSLOForOwner(node), nil
	case "resolver":
		node, ok := decodedCalldata[0].([32]byte)
		if !ok {
			return nil, errors.New("get SLO for address")
		}
		return getoSLOForResolver(node), nil
	default:
		return nil, fmt.Errorf("get SLO for unknown method %s", methodName)
	}
}

func (g *Gateway) decode(hexCalldata string) (methodName string, decoded []interface{}, err error) {
	calldata, err := DecodeHex(hexCalldata)
	if err != nil {
		return "", nil, fmt.Errorf("decoding hex calldata: %w", err)
	}
	functionSignature := calldata[:4]
	functionParameters := calldata[4:]
	var sig [4]byte
	copy(sig[:], functionSignature)
	methodName, ok := g.selectors[sig]
	if !ok {
		return "", nil, fmt.Errorf("unknown function signature: %s", string(functionSignature))
	}
	unpacked, err := abi.Methods[methodName].Inputs.Unpack(functionParameters)
	return methodName, unpacked, err
}

func (g *Gateway) getGateway(w http.ResponseWriter, r *http.Request) {
	sender := chi.URLParam(r, "sender")
	hexCalldata := chi.URLParam(r, "data")
	log.Printf("sender: %s, hexCalldata: %s", sender, hexCalldata)
	methodName, decoded, err := g.decode(hexCalldata)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	log.Printf("method: %s decoded: %+v", methodName, decoded)
	slo, err := getSLO(methodName, decoded)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	log.Printf("slo: %+v", slo)
	blockNum, err := g.l2EthClient.BlockNumber(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
	}
	addressSlot := fmt.Sprintf("0x%s", hex.EncodeToString(slo))
	res, err := g.l2GethClient.GetProof(r.Context(), g.l2ResolverAddress, []string{addressSlot}, new(big.Int).SetUint64(blockNum))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
	}
	log.Printf("getProof result: %+v", res)
	jsonBody, err := json.Marshal(StateRootProofInput{
		ResolverAddr: g.l2ResolverAddress.Hex(),
		AddrSlot:     addressSlot,
	})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
	}
	bodyReader := bytes.NewReader(jsonBody)
	req, err := http.NewRequest(http.MethodPost, "http://localhost:41235/storage_proof", bodyReader)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("reading stateRootProof body: %s", err)
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
	}
	stateRootProof := &StateRootProof{}
	err = json.Unmarshal(body, stateRootProof)
	if err != nil {
		fmt.Printf("unmarshaling stateRootProof body: %s", err)
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
	}
	log.Printf("stateRootProof: %+v", stateRootProof)

	w.Write(body)
}

func encodeResponse() (resp []byte, err error) {
	resp, err = abi.Methods["ownerWithProof"].Inputs.Pack()
	if err != nil {
		return
	}

}

// calldata
// 0x02571be36e4409edbffcd311c18c7dbfadfc428daf0edd143693491dd2b4782bc494dc61
