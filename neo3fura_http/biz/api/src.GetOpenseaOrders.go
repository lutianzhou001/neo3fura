package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"neo3fura_http/lib/type/h160"
	"net/http"
	"strconv"
)

func (me *T) GetOpenseaOrders(args struct {
	AssetContractAddress h160.T
	PaymentTokenAddress  h160.T
	Marker               string
	Taker                string
	Owner                string
	IsEnlish             string
	Bundled              bool
	IncludeBundled       bool
	ListedAfter          string
	ListedBefore         string
	TokenId              string
	TokenIds             []string
	Side                 string
	SaleKind             string
	Limit                int64
	Offset               int64
	OrderBy              string
	OrderDirection       string
	ApiKey               string
	Filter               map[string]interface{}
}, ret *json.RawMessage) error {

	params := ""

	params = "offset=" + fmt.Sprintf("%d", args.Offset) + "&bundled=" + strconv.FormatBool(args.Bundled) + "&include_bundled=" + strconv.FormatBool(args.IncludeBundled) + "&"

	tokenids := ""
	if len(args.TokenIds) > 0 {
		for _, item := range args.TokenIds {
			if len(item) > 0 {
				tokenids += "token_ids=" + item + "&"
			}
		}
	}

	if args.IsEnlish != "" {
		if args.IsEnlish != "true" && args.IsEnlish != "false" {
			args.IsEnlish = ""
		}
	}
	if args.AssetContractAddress != "" {
		params = params + "asset_contract_address=" + args.AssetContractAddress.Val() + "&"
	}
	if args.PaymentTokenAddress != "" {
		params = params + "payment_token_address=" + args.PaymentTokenAddress.Val() + "&"
	}
	if args.Marker != "" {
		params = params + "maker=" + args.Marker + "&"
	}
	if args.Taker != "" {
		params = params + "taker=" + args.Taker + "&"
	}
	if args.Owner != "" {
		params = params + "owner=" + args.Owner + "&"
	}
	if args.IsEnlish != "" {
		params = params + "is_english=" + args.IsEnlish + "&"
	}
	if args.ListedAfter != "" {
		params = params + "listed_after=" + args.ListedAfter + "&"
	}
	if args.ListedBefore != "" {
		params = params + "listed_before=" + args.ListedBefore + "&"
	}
	if len(args.TokenIds) > 0 {
		params = params + tokenids
	}
	if args.Side != "" {
		params = params + "side=" + fmt.Sprintf("%d", args.Side) + "&"
	} else {
		params = params + "side=1&"
	}
	if args.SaleKind != "" {
		params = params + "sale_kind=" + args.SaleKind + "&"
	}
	if args.Limit > 0 {
		params = params + "limit=" + fmt.Sprintf("%d", args.Limit) + "&"
	}
	if args.Limit == 0 {
		params = params + "limit=20&"
	}
	if args.OrderBy != "" {
		params = params + "order_by=" + args.OrderBy + "&"
	}
	if args.OrderBy == "" {
		params = params + "order_by=created_date&"
	}
	if args.OrderDirection != "" {
		params = params + "order_direction=" + args.OrderDirection + "&"
	}
	if args.OrderDirection == "" {
		params = params + "order_direction=desc&"
	}

	var requestGetURLNoParams = "https://api.opensea.io/wyvern/v1/orders?" + params
	//fmt.Println(requestGetURLNoParams1)
	//consts requestGetURLNoParams = "https://api.opensea.io/wyvern/v1/orders?bundled=false&include_bundled=false&side=1&limit=20&offset=0&order_by=created_date&order_direction=desc"
	client := &http.Client{}
	req, err := http.NewRequest("GET", requestGetURLNoParams, nil)
	if err != nil {
		panic(err)
		return err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Set("X-API-KEY", args.ApiKey)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	resbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	re := make(map[string]interface{})
	err = json.Unmarshal(resbody, &re)
	if err != nil {
		return err
	}

	r, err := json.Marshal(re)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}
