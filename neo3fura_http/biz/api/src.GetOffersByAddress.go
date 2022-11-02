package api

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"io/ioutil"
	"log"
	log2 "neo3fura_http/lib/log"
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
	"strings"
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
							//item["image"] = image
							item["image"] = ImagUrl(asset, image.(string), "images")
						} else {
							item["image"] = ""
						}
						thumbnail, ok := data["thumbnail"]
						if ok {
							tb, err22 := base64.URLEncoding.DecodeString(thumbnail.(string))
							if err22 != nil {
								return err22
							}
							//item["image"] = string(tb[:])
							item["thumbnail"] = ImagUrl(asset, string(tb[:]), "thumbnail")
						} else {
							if item["thumbnail"] == nil {
								if image != nil && image != "" {
									if image == nil {
										item["thumbnail"] = item["image"]
									} else {
										item["thumbnail"] = ImagUrl(asset, image.(string), "thumbnail")
									}
								}
							}
						}
						tokenuri, ok := data["tokenURI"]
						if ok {
							ppjson, err := GetImgFromTokenURL(tokenurl(tokenuri.(string)), asset, tokenid)
							if err != nil {
								return err
							}
							for key, value := range ppjson {
								item[key] = value
								properties[key] = value
								if key == "image" {
									img := value.(string)
									item["thumbnail"] = ImagUrl(asset, img, "thumbnail")
									item["image"] = ImagUrl(asset, img, "images")
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

	//	//OfferCollection
	rt := os.ExpandEnv("${RUNTIME}")

	var secondMarketHash string
	if rt == "staging" {
		secondMarketHash = "0xd2e7cf18ee0d9b509fac02457f54b63e47b25e29"
	} else if rt == "test2" {
		secondMarketHash = "0xc198d687cc67e244662c3b9c1325f095f8e663b1"
	} else {
		secondMarketHash = "0xc198d687cc67e244662c3b9c1325f095f8e663b1"

	}

	raw := make(map[string]interface{}, 0)
	err = me.GetMarketAssetOwnedByAddress(struct {
		Address    h160.T
		MarketHash h160.T
		NFTState   string
		Limit      int64
		Skip       int64
		Filter     map[string]interface{}
		Raw        *map[string]interface{}
	}{Address: args.Address, MarketHash: h160.T(secondMarketHash), NFTState: "", Raw: &raw}, ret)
	if err != nil {
		return err
	}
	if raw != nil {
		asset := raw["assetlist"]
		offerCollection, err := me.Client.QueryAggregate(struct {
			Collection string
			Index      string
			Sort       bson.M
			Filter     bson.M
			Pipeline   []bson.M
			Query      []string
		}{Collection: "MarketNotification",
			Index:  "getOfferCollection",
			Sort:   bson.M{},
			Filter: bson.M{},
			Pipeline: []bson.M{
				bson.M{"$match": bson.M{"eventname": "OfferCollection", "asset": bson.M{"$in": asset}}},
				bson.M{"$lookup": bson.M{
					"from": "Nep11Properties",
					"let":  bson.M{"asset": "$asset"},
					"pipeline": []bson.M{
						bson.M{"$match": bson.M{"$expr": bson.M{"$and": []interface{}{
							//bson.M{"$eq": []interface{}{"$tokenid", "$$tokenid"}},
							bson.M{"$eq": []interface{}{"$asset", "$$asset"}},
						}}}},
						bson.M{"$project": bson.M{"properties": 1, "asset": 1, "tokenid": 1}},
					},
					"as": "properties"},
				},
			},
			Query: []string{},
		}, ret)
		if err != nil {
			return err
		}

		if offerCollection != nil {
			for _, item := range offerCollection {
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
							bson.M{"eventname": "CompleteOfferCollection"},
							bson.M{"eventname": "CancelOfferCollection"},
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
							//item["originOwner"] = data["originOwner"]
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

							//properties

							if item["properties"] != nil {
								properties := item["properties"].(primitive.A)
								for _, ppit := range properties {
									copyItem := make(map[string]interface{})
									copyItem = item

									ppitem := ppit.(map[string]interface{})
									ppinfo := ppitem["properties"]
									ppAsset := ppitem["asset"].(string)
									ppTokenid := ppitem["tokenid"].(string)
									var ppdata map[string]interface{}
									copyItem["tokenid"] = ppTokenid
									copyItem["asset"] = ppAsset
									copyItem["properties"] = ppinfo

									//
									// OWNER
									marketInfo, err := me.Client.QueryOne(struct {
										Collection string
										Index      string
										Sort       bson.M
										Filter     bson.M
										Query      []string
									}{Collection: "Market",
										Index:  "GetMarketInfo",
										Sort:   bson.M{},
										Filter: bson.M{"amount": bson.M{"$gt": 0}, "asset": item["asset"], "tokenid": ppTokenid},
										Query:  []string{},
									}, ret)

									if err != nil {
										return stderr.ErrGetNFTInfo
									}
									market := marketInfo["market"]
									owner := marketInfo["owner"]
									bidder := marketInfo["bidder"]
									bidAmount := marketInfo["bidAmount"].(primitive.Decimal128).String()
									ddl := marketInfo["deadline"].(int64)

									if market == owner && ddl > currentTime { // 上架未过期
										copyItem["originOwner"] = marketInfo["auctor"]
									} else if market == owner && ddl < currentTime { //上架过期
										if bidAmount == "0" {
											copyItem["originOwner"] = marketInfo["auctor"]
										} else {
											copyItem["originOwner"] = bidder
										}
									} else { //未上架
										copyItem["originOwner"] = owner
									}

									if copyItem["originOwner"].(string) != args.Address.Val() {
										continue
									}
									//properties
									if ppinfo != nil {
										if err1 := json.Unmarshal([]byte(ppinfo.(string)), &ppdata); err1 == nil {
											image, ok := ppdata["image"]
											if ok {

												//item["image"] = image
												copyItem["image"] = ImagUrl(ppAsset, image.(string), "images")
											} else {
												copyItem["image"] = ""
											}
											thumbnail, ok := ppdata["thumbnail"]
											if ok {
												tb, err22 := base64.URLEncoding.DecodeString(thumbnail.(string))
												if err22 != nil {
													return err22
												}
												//item["image"] = string(tb[:])
												copyItem["thumbnail"] = ImagUrl(ppAsset, string(tb[:]), "thumbnail")
											} else {
												if item["thumbnail"] == nil {
													if image != nil && image != "" {
														if image == nil {
															copyItem["thumbnail"] = item["image"]
														} else {
															copyItem["thumbnail"] = ImagUrl(ppAsset, image.(string), "thumbnail")
														}
													}
												}
											}
											tokenuri, ok := ppdata["tokenURI"]
											if ok {
												ppjson, err := GetImgFromTokenURL(tokenurl(tokenuri.(string)), ppAsset, ppTokenid)
												if err != nil {
													return err
												}
												for key, value := range ppjson {
													item[key] = value

													if key == "image" {
														img := value.(string)
														copyItem["thumbnail"] = ImagUrl(ppAsset, img, "thumbnail")
														copyItem["image"] = ImagUrl(ppAsset, img, "images")
													}
												}
											}
										} else {
											return err
										}
									}

									delete(copyItem, "extendData")
									delete(copyItem, "properties")
									result = append(result, copyItem)

								}
							}

						} else {
							return err2
						}

					}
				}

			}

		}

		//append(result, )
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
	if tokenurl == "" {
		return nil, nil
	}
	currentPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	//
	num, _ := base64.StdEncoding.DecodeString(tokenid)
	if err != nil {
		return nil, err
	}
	fname, err := bytesToIntS(num)
	if err != nil {
		return nil, err
	}
	filename := strconv.FormatInt(int64(fname), 10) + ".json"
	path := currentPath + "/tokenURI/" + asset + "/" + filename
	isExit, _ := PathExists(path)
	jsonData := make(map[string]interface{})
	if !isExit { //读取数据并保存到本地
		filepath := CreateDateDir(currentPath+"/tokenURI/", asset)
		response, err := http.Get(tokenurl)
		if err != nil {
			log.Println("http get error: ", err)
			return nil, err
		}

		raw := response.Body
		defer raw.Close()

		out, err := os.Create(filepath + "/" + filename)
		if err != nil {
			panic(err)
			return nil, err
		}

		wt := bufio.NewWriter(out)
		defer out.Close()

		n, err := io.Copy(wt, response.Body)
		fmt.Println("write", n)
		if err != nil {
			return nil, err
		}
		wt.Flush()

	}
	//从文件读数据
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println("error opening json file")
		return nil, err
	}
	defer jsonFile.Close()

	body, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("error reading json file")
		return nil, err
	}

	if len(body) > 0 {
		err := json.Unmarshal([]byte(string(body)), &jsonData)
		if err != nil {
			log.Println("imag from json error :", err, tokenurl)
			return nil, err
		}

		attributes, ok := jsonData["attributes"]

		if ok {
			attribute := attributes.([]interface{})
			for _, item := range attribute {
				it := item.(map[string]interface{})
				jsonData[it["trait_type"].(string)] = it["value"]
			}
			delete(jsonData, "attributes")
		}
		delete(jsonData, "number")

	}

	return jsonData, nil
}
func Imgname(asset string, url string) string {
	imgname := strings.ReplaceAll(url, "/", "")

	rt := os.ExpandEnv("${RUNTIME}")
	const test = "0xaecbad96ccc77c8b147a52e45723a6b5886454e0"
	const main = "0x9f344fe24c963d70f5dcf0cfdeb536dc9c0acb3a"
	split := strings.Split(url, ".")
	suf := split[len(split)-1]
	pre := "ipfs://bafybeiapiufkjejfj2mdvjyigrga5vt3o2sd6xf35372tnptiah7kygm7m/1.gif"
	if rt == "staging" && asset == main && suf == "gif" {
		imgname = strings.ReplaceAll(pre, "/", "")
	} else if rt == "test2" && asset == test && suf == "gif" {
		imgname = strings.ReplaceAll(pre, "/", "")
	}

	return imgname
}
func ImagUrl(asset string, imgurl string, pre string) string {
	rt := os.ExpandEnv("${RUNTIME}")
	name := Imgname(asset, imgurl)
	url := ""
	switch rt {
	case "test":
		url = "https://img.megaoasis.io/testnet/" + pre + "/" + asset + "/" + name
	case "test2":
		url = "https://img.megaoasis.io/testnet/" + pre + "/" + asset + "/" + name
	case "staging":
		url = "https://img.megaoasis.io/" + pre + "/" + asset + "/" + name
	default:
		log2.Fatalf("runtime environment mismatch")
	}
	return url
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

func bytesToIntS(b []byte) (int, error) {
	if len(b) == 3 {
		b = append([]byte{0}, b...)
	}
	bytesBuffer := bytes.NewBuffer(b)
	switch len(b) {
	case 1:
		var tmp int8
		err := binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
		return int(tmp), err
	case 2:
		var tmp int16
		err := binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
		return int(tmp), err
	case 4:
		var tmp int32
		err := binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
		return int(tmp), err
	default:
		return 0, fmt.Errorf("%s", "BytesToInt bytes lenth is invaild!")
	}
}
