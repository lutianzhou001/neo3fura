package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	for {
		task()
	}
}

func init() {
	addressesNEORPCPOPPER = strings.Split(os.ExpandEnv("${NEORPC_POPPERADDRESSES}"), " ")
	addressesNEOCLI = strings.Split(os.ExpandEnv("${NEOCLI_ADDRESSES}"), " ")
}

var addressesNEORPCPOPPER []string
var addressesNEOCLI []string

func task() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("[!!!!][ERROR]", r)
			time.Sleep(time.Second)
		}
	}()
	addressPOPPER := addressesNEORPCPOPPER[rand.Intn(len(addressesNEORPCPOPPER))]
	resp, err := http.Get(addressPOPPER)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	payload, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      rand.Uint32(),
		"method":  "sendrawtransaction",
		"params": []interface{}{
			hex.EncodeToString(data),
		},
	})
	if err != nil {
		return
	}

	for i := time.Millisecond; i < time.Second; i = i * 2 {
		addressNEOCLI := addressesNEOCLI[rand.Intn(len(addressesNEOCLI))]
		resp, err := http.Post(addressNEOCLI, "application/json", bytes.NewReader(payload))
		if err != nil {
			log.Println("[????][REQ]", err)
			continue
		}
		defer resp.Body.Close()
	}
}
