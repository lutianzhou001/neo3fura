package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"math"
	"math/big"
	log2 "neo3fura_http/lib/log"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"
	"net/http"
	"path/filepath"
	"strconv"
)

func (me *T) GetMarketIndexByAsset(args struct {
	AssetHash h160.T
	Filter    map[string]interface{}
	Raw       *map[string]interface{}
}, ret *json.RawMessage) error {

	if args.AssetHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}

	result, err := me.Client.QueryOneJob(struct {
		Collection string
		Filter     bson.M
	}{Collection: "MarketIndex", Filter: bson.M{"asset": args.AssetHash}})

	if err != nil && err.Error() == "mongo: no documents in result" {
		return err
	}

	r, err := json.Marshal(result)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)

	return nil
}

func GetPrice(asset string) (float64, error) {

	client := &http.Client{}
	reqBody := []byte(`["` + asset + `"]`)
	url := "https://onegate.space/api/quote?convert=usd"
	//str :=[]string{asset}
	req, _ :=
		http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log2.Fatal("request price err :", err)
		return 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log2.Fatal("readall price err :", err)
		return 0, err
	}
	response := string(body)
	re := response[1 : len(response)-1]
	price, err1 := strconv.ParseFloat(re, 64)
	if err1 != nil {
		log2.Fatal("price parsefloat err :", err)
		return 0, err
	}
	price, err = strconv.ParseFloat(fmt.Sprintf("%.8f", price), 64)
	if err != nil {
		log2.Fatal("price parsefloat decimal err :", err)
		return 0, err
	}
	return price, nil
}
func GetPrice2(asset string, amount primitive.Decimal128) (*big.Float, error) {

	client := &http.Client{}
	reqBody := []byte(`["` + asset + `"]`)
	url := "https://onegate.space/api/quote?convert=usd"
	//str :=[]string{asset}
	req, _ :=
		http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return big.NewFloat(float64(0)), stderr.ErrPrice
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return big.NewFloat(float64(0)), stderr.ErrPrice
	}
	response := string(body)
	re := response[1 : len(response)-1]
	price, err1 := strconv.ParseFloat(re, 64)

	//获取decimal
	decimal := int64(1)
	if asset != "" {
		dd, _ := OpenAssetHashFile()
		decimal = dd[asset] //获取精度
		if decimal == int64(0) {
			decimal = int64(1)
		}
	}

	var usdAuctionAmount *big.Float
	//计算价格
	bamount, _, err := amount.BigInt()
	bfauctionAmount := new(big.Float).SetInt(bamount)
	flag := bamount.Cmp(big.NewInt(0))

	if flag == 1 {
		bfprice := big.NewFloat(price)
		ffprice := big.NewFloat(1).Mul(bfprice, bfauctionAmount)
		de := math.Pow(10, float64(decimal))
		usdAuctionAmount = new(big.Float).Quo(ffprice, big.NewFloat(de))

	} else {
		usdAuctionAmount = big.NewFloat(float64(0))
	}

	if err1 != nil {
		return big.NewFloat(0), stderr.ErrPrice
	}
	return usdAuctionAmount, nil
}

func OpenAssetHashFile() (map[string]int64, error) {
	absPath, _ := filepath.Abs("./assethash.json")

	b, err := ioutil.ReadFile(absPath)
	if err != nil {
		fmt.Print(err)
	}
	whitelist := map[string]int64{}
	err = json.Unmarshal([]byte(string(b)), &whitelist)
	if err != nil {
		panic(err)
	}

	return whitelist, err
}
