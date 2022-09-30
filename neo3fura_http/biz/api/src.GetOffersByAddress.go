package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"io/ioutil"
	"log"
	"neo3fura_http/lib/mapsort"
	"neo3fura_http/lib/type/OfferState"
	_ "neo3fura_http/lib/type/OfferState"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/type/strval"
	"neo3fura_http/var/stderr"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func (me *T) GetOffersByAddress(args struct {
	Address    h160.T
	OfferState strval.T //state:vaild  received
	Limit      int64
	Skip       int64
	Filter     map[string]interface{}
}, ret *json.RawMessage) error {
	currentTime := time.Now().UnixNano() / 1e6
	if args.Address.Valid() == false {
		return stderr.ErrInvalidArgs
	}

	pipeline := []bson.M{}

	if args.OfferState.Val() == OfferState.Valid.Val() { //拍卖中  accont >0 && auctionType =2 &&  owner=market && runtime <deadline
		pipeline = []bson.M{
			bson.M{"$match": bson.M{"user": args.Address, "eventname": "Offer"}},
			bson.M{"$lookup": bson.M{
				"from": "Nep11Properties",
				"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid"},
				"pipeline": []bson.M{
					bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
						bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
						bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
					}}}},
					bson.M{"$project": bson.M{"asset": 1, "tokenid": 1, "properties": 1}},
				},
				"as": "properties"},
			},
		}
	} else if args.OfferState.Val() == OfferState.Received.Val() {
		pipeline = []bson.M{
			bson.M{"$match": bson.M{"eventname": "Offer", "$or": []interface{}{
				bson.M{"extendData": bson.M{"$regex": "originOwner\":\"" + args.Address}},
				bson.M{"extendData": bson.M{"$regex": "originOwner\": \"" + args.Address}},
			}}},
			bson.M{"$lookup": bson.M{
				"from": "Nep11Properties",
				"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid"},
				"pipeline": []bson.M{
					bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
						bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
						bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
					}}}},
					bson.M{"$project": bson.M{"asset": 1, "tokenid": 1, "properties": 1}},
				},
				"as": "properties"},
			},
		}
	} else {
		pipeline = []bson.M{
			bson.M{"$match": bson.M{"eventname": "Offer", "$or": []interface{}{
				bson.M{"extendData": bson.M{"$regex": "originOwner\":\"" + args.Address}},
				bson.M{"extendData": bson.M{"$regex": "originOwner\": \"" + args.Address}},
				bson.M{"user": args.Address},
			}}},
			bson.M{"$lookup": bson.M{
				"from": "Nep11Properties",
				"let":  bson.M{"asset": "$asset", "tokenid": "$tokenid"},
				"pipeline": []bson.M{
					bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
						bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
						bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
					}}}},
					bson.M{"$project": bson.M{"properties": 1}},
				},
				"as": "properties"},
			},
		}
	}

	var r1, err = me.Client.QueryAggregate(
		struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Pipeline   []bson.M
			Query      []string
		}{
			Collection: "MarketNotification",
			Index:      "GetOffersByAddress",
			Sort:       bson.M{"timestamp": -1},
			Filter:     bson.M{},
			Pipeline:   pipeline,
			Query:      []string{},
		}, ret)
	if err != nil {
		return err
	}

	result := make([]map[string]interface{}, 0)
	for _, item := range r1 {

		//查看offer 当前状态
		offer_nonce := item["nonce"]
		offer, _ := me.Client.QueryOne(struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Query      []string
		}{
			Collection: "MarketNotification",
			Index:      "getOfferSate",
			Sort:       bson.M{},
			Filter: bson.M{
				"nonce":   offer_nonce,
				"asset":   item["asset"],
				"tokenid": item["tokenid"],
				//"eventname":"CancelOffer",
				"$or": []interface{}{
					bson.M{"eventname": "CompleteOffer"},
					bson.M{"eventname": "CancelOffer"},
				},
			},
			Query: []string{},
		}, ret)

		if len(offer) > 0 {
			continue
		}

		if item["extendData"] != nil {
			extendData := item["extendData"].(string)
			if extendData != "" {
				var data map[string]interface{}
				if err2 := json.Unmarshal([]byte(extendData), &data); err2 == nil {
					item["originOwner"] = data["originOwner"]
					item["offerAsset"] = data["offerAsset"]
					oa := data["offerAmount"].(string)
					offerAmount, err := strconv.ParseInt(oa, 10, 64)
					if err != nil {
						return err
					}
					item["offerAmount"] = offerAmount
					dl := data["deadline"].(string)
					deadline, err := strconv.ParseInt(dl, 10, 64)
					if err != nil {
						return err
					}
					item["deadline"] = deadline

				} else {
					return err2
				}

			}
		}

		nftproperties := item["properties"]
		if nftproperties != nil && nftproperties != "" {
			pp := nftproperties.(primitive.A)
			if len(pp) > 0 {
				it := pp[0].(map[string]interface{})
				extendData := it["properties"].(string)
				asset := it["asset"].(string)
				tokenid := it["tokenid"].(string)
				if extendData != "" {
					properties := make(map[string]interface{})
					var data map[string]interface{}
					if err1 := json.Unmarshal([]byte(extendData), &data); err1 == nil {
						image, ok := data["image"]
						if ok {
							properties["image"] = image
							item["image"] = image
						} else {
							item["image"] = ""
						}
						tokenurl, ok := data["tokenURI"]
						if ok {
							if image == "" {
								image, err = GetImgFromTokenURL(tokenurl.(string), asset, tokenid)
								item["image"] = image
								if err != nil {
									return err
								}
							}
						}

					} else {
						return err
					}

				} else {
					item["image"] = ""
				}

			}

		}

		if item["deadline"].(int64) > currentTime {
			result = append(result, item)
		}
		delete(item, "extendData")
		delete(item, "properties")
	}

	if args.OfferState.Val() == OfferState.Received.Val() {
		result = mapsort.MapSort(result, "offerAmount")
	}
	pageResult := make([]map[string]interface{}, 0)
	for i, item := range result {
		if int64(i) < args.Skip {
			continue
		} else if int64(i) > args.Skip+args.Limit-1 {
			continue
		} else {
			pageResult = append(pageResult, item)
		}
	}

	var count = int64(len(pageResult))

	if err != nil {
		return err
	}
	r2, err := me.FilterArrayAndAppendCount(result, count, args.Filter)
	if err != nil {
		return err
	}
	r, err := json.Marshal(r2)
	if err != nil {
		return err
	}

	*ret = json.RawMessage(r)
	return nil
}

func GetImgFromTokenURL(tokenurl string, asset string, tokenid string) (map[string]interface{}, error) {
	//检查该tokenurl 文件是否本地存在
	path := "./tokenURL/" + asset + "/" + tokenid
	isExit, _ := PathExists(path)
	jsonData := make(map[string]interface{})
	var body []byte
	if isExit { //从本地读数据
		jsonFile, err := os.Open(path)
		if err != nil {
			fmt.Println("error opening json file")
			return nil, err
		}
		defer jsonFile.Close()

		body, err = ioutil.ReadAll(jsonFile)
		if err != nil {
			fmt.Println("error reading json file")
			return nil, err
		}

	} else { //读取数据并保存到本地
		currentPath, err := os.Getwd()
		fmt.Println(currentPath)
		filepath := CreateDateDir(currentPath+"/tokenURI/", asset)
		response, err := http.Get(tokenurl)
		if err != nil {
			log.Println("http get error: ", err)
			return nil, err
		}
		defer response.Body.Close()
		out, err := os.Create(filepath + "/" + tokenid)
		if err != nil {
			panic(err)
		}
		wt := bufio.NewWriter(out)
		defer out.Close()
		n, err := io.Copy(wt, response.Body)
		fmt.Println("write", n)
		if err != nil {
			panic(err)
		}
		wt.Flush()
		body, err = ioutil.ReadAll(response.Body)
		if err != nil {
			log.Println("ioutil read error: ", err)
		}
	}

	err := json.Unmarshal([]byte(string(body)), &jsonData)
	if err != nil {
		log.Println("imag from json error :", err, tokenurl)
		return nil, err
	}

	return jsonData, nil
}

func CreateDateDir(basepath string, folderName string) string {

	folderPath := filepath.Join(basepath, folderName)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.MkdirAll(folderPath, 0777)
		if err != nil {
			fmt.Println("Create dir error: %v", err)
		}
		err = os.Chmod(folderPath, 0777)
		if err != nil {
			fmt.Println("Chmod error: %v", err)
		}
	}
	return folderPath
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	//当为空文件或文件夹存在
	if err == nil {
		return true, nil
	}
	//os.IsNotExist(err)为true，文件或文件夹不存在
	if os.IsNotExist(err) {
		return false, nil
	}
	//其它类型，不确定是否存在
	return false, err
}
