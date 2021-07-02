package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"sort"
)

func (me *T) GetNetFeeRange(ret *json.RawMessage) error {
	r1, err := me.Data.Client.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: "Transaction",
		Index:      "someIndex",
		Sort:       bson.M{"_id": -1},
		Filter:     bson.M{},
		Query:      []string{},
	}, ret)
	if err != nil {
		return err
	}
	r2, _, err := me.Data.Client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: "Transaction",
		Index:      "someIndex",
		Sort:       bson.M{},
		Filter:     bson.M{"blockhash": r1["hash"]},
		Query:      []string{},
		Limit:      500,
		Skip:       0,
	}, ret)
	if err != nil {
		return err
	}
	// get net fee
	netFees := make([]float64, len(r2))
	for i, tx := range r2 {
		netFees[i] = tx["netfee"].(float64) // todo, type doesn't match may throw error
	}
	// sort, split to three parts
	sort.Float64s(netFees)
	r3 := splitAndGetEachMedian(netFees)
	r4 := map[string]interface{}{
		"slow":    r3[0],
		"fast":    r3[1],
		"fastest": r3[2],
	}
	r, err := json.Marshal(r4)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r)
	return nil
}

func splitAndGetEachMedian(values []float64) []float64 {
	subL := len(values) / 3
	if len(values)%3 == 2 {
		subL++
	}
	r := make([][]float64, 3)
	r[0] = values[:subL]         // lowest
	r[1] = values[subL : 2*subL] // medium
	r[2] = values[2*subL:]       // highest
	results := make([]float64, 3)
	for i, v := range r {
		results[i] = getMedian(v)
	}
	return results
}

func getMedian(values []float64) float64 {
	n := len(values)
	if n == 0 {
		return 0
	}
	if n%2 == 1 {
		return values[n/2]
	} else {
		return (values[n/2-1] + values[n/2]) / 2 // 8 decimals, tested
	}
}
