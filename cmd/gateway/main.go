package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

func mustParseABI(json string) *abi.ABI {
	a, err := abi.JSON(strings.NewReader(json))
	if err != nil {
		log.Fatal("failed to parse ABI", err)
	}
	return &a
}

func mustGetSelector(parsedABI *abi.ABI, methodName string) []byte {
	method, ok := parsedABI.Methods[methodName]
	if !ok {
		log.Fatal("failed to find method", methodName)
	}
	return method.ID
}

/*
decode calldata
calculate storage slot
generate proof
encode response
*/

type Gateway struct {
	l2Client          *ethclient.Client
	l2ResolverAddress common.Address
}

func main() {
	l2RPCURL := GetOrDefault("L2_RPC_URL", "https://goerli.optimism.io")
	l2Client, err := ethclient.Dial(l2RPCURL)
	if err != nil {
		log.Fatal("dialing L2 RPC", err)
	}

	gateway := Gateway{
		l2Client:          l2Client,
		l2ResolverAddress: common.HexToAddress(GetOrDefault("L2_RESOLVER_ADDR", "0xbeefbabe")),
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/block_number", gateway.getBlockNumber)
	r.Post("/query", gateway.query)
	http.ListenAndServe(":3000", r)
}

func (g *Gateway) getBlockNumber(w http.ResponseWriter, r *http.Request) {
	blockNum, err := g.l2Client.BlockNumber(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
	}
	w.Write([]byte(fmt.Sprintf("block_number: %d", blockNum)))
}

func decode(calldata []byte) {
	// var mockSelector []byte
	// if calldata[:4] == mockSelector {

	// }
}

func (g *Gateway) query(w http.ResponseWriter, r *http.Request) {
	hexCalldata, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	calldata, err := DecodeHex(string(hexCalldata))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	decode(calldata)
	w.Write([]byte(fmt.Sprintf("calldata: %s", hexCalldata)))
	// g.l2Client.StorageAt()
}
