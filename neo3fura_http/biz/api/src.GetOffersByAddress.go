package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"neo3fura_http/lib/mapsort"
	"neo3fura_http/lib/type/OfferState"
	_ "neo3fura_http/lib/type/OfferState"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/type/strval"
	"neo3fura_http/var/stderr"
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