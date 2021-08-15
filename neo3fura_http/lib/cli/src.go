package cli

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/joeqian10/neo3-gogogo/rpc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// T ...
type T struct {
	Redis     *redis.Client
	Db_online string
	C_online  *mongo.Client
	C_local   *mongo.Client
	Ctx       context.Context
	RpcCli    *rpc.RpcClient
	RpcPorts  []string
}

func (me *T) QueryOne(args struct {
	Collection string
	Index      string
	Sort       bson.M
	Filter     bson.M
	Query      []string
}, ret *json.RawMessage) (map[string]interface{}, error) {
	var kvs string
	kvs = kvs + args.Collection
	kvs = kvs + args.Index
	for k, v := range args.Sort {
		kvs = kvs + k + fmt.Sprintf("%v", v)
	}
	for k, v := range args.Filter {
		kvs = kvs + k + fmt.Sprintf("%v", v)
	}
	for _, v := range args.Query {
		kvs = kvs + v
	}
	h := sha1.New()
	h.Write([]byte(kvs))
	hash := hex.EncodeToString(h.Sum(nil))
	val, err := me.Redis.Get(me.Ctx, hash).Result()
	// if sort != nil, it may have several results, we have to pick the sorted one
	if err == redis.Nil || args.Sort != nil {
		var result map[string]interface{}
		convert := make(map[string]interface{})
		collection := me.C_online.Database(me.Db_online).Collection(args.Collection)
		opts := options.FindOne().SetSort(args.Sort)
		err = collection.FindOne(me.Ctx, args.Filter, opts).Decode(&result)
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("document not found:%s", errors.New("NOT FOUND"))
		} else if err != nil {
			return nil, fmt.Errorf("find document error:%s", err)
		}
		if len(args.Query) == 0 {
			convert = result
		} else {
			for _, v := range args.Query {
				convert[v] = result[v]
			}
		}
		r, err := json.Marshal(convert)
		if err != nil {
			return nil, fmt.Errorf("json marshal error:%s", err)
		}
		err = me.Redis.Set(me.Ctx, hash, hex.EncodeToString(r), 0).Err()
		if err != nil {
			return nil, fmt.Errorf("write to redis error:%s", err)
		}
		*ret = json.RawMessage(r)
		return convert, nil
	} else {
		r, err := hex.DecodeString(val)
		if err != nil {
			return nil, fmt.Errorf("decoding to hexstring error:%s", err)
		}

		*ret = json.RawMessage(r)
		convert := make(map[string]interface{})
		err = json.Unmarshal(r, &convert)
		if convert["_id"] != nil {
			convert["_id"], err = primitive.ObjectIDFromHex(convert["_id"].(string))
		}
		if err != nil {
			return nil, fmt.Errorf("convert to string error:%s", err)
		}
		return convert, nil
	}
	return nil, nil
}

func (me *T) QueryAll(args struct {
	Collection string
	Index      string
	Sort       bson.M
	Filter     bson.M
	Query      []string
	Limit      int64
	Skip       int64
}, ret *json.RawMessage) ([]map[string]interface{}, int64, error) {
	var results []map[string]interface{}
	convert := make([]map[string]interface{}, 0)
	collection := me.C_online.Database(me.Db_online).Collection(args.Collection)
	op := options.Find()
	op.SetSort(args.Sort)
	op.SetLimit(args.Limit)
	op.SetSkip(args.Skip)
	co := options.CountOptions{}
	count, err := collection.CountDocuments(me.Ctx, args.Filter, &co)
	if err != nil {
		return nil, 0, fmt.Errorf("count documents error:%s", err)
	}
	cursor, err := collection.Find(me.Ctx, args.Filter, op)
	defer cursor.Close(me.Ctx)
	if err == mongo.ErrNoDocuments {
		return nil, 0, fmt.Errorf("document not found:%s", errors.New("NOT FOUND"))
	}
	if err != nil {
		return nil, 0, fmt.Errorf("get cursor error:%s", err)
	}
	if err = cursor.All(me.Ctx, &results); err != nil {
		return nil, 0, fmt.Errorf("find documents error:%s", err)
	}
	for _, item := range results {
		if len(args.Query) == 0 {
			convert = append(convert, item)
		} else {
			temp := make(map[string]interface{})
			for _, v := range args.Query {
				temp[v] = item[v]
			}
			convert = append(convert, temp)
		}
	}
	r, err := json.Marshal(convert)
	if err != nil {
		return nil, 0, fmt.Errorf("json marshal error:%s", err)
	}
	*ret = json.RawMessage(r)
	return convert, count, nil
}

func (me *T) SaveJob(args struct {
	Collection string
	Data       bson.M
}) (bool, error) {
	collection := me.C_local.Database("job").Collection(args.Collection)
	_, err := collection.InsertOne(me.Ctx, args.Data)
	if err != nil {
		return false, fmt.Errorf("insert job error:%s", err)
	}
	return true, nil
}

func (me *T) QueryLastJob(args struct {
	Collection string
}) (map[string]interface{}, error) {
	collection := me.C_local.Database("job").Collection(args.Collection)
	var result map[string]interface{}
	opts := options.FindOne().SetSort(bson.M{"_id": -1})
	err := collection.FindOne(me.Ctx, bson.M{}, opts).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("find last job error:%s", err)
	}
	return result, nil
}

func (me *T) QueryAggregate(args struct {
	Collection string
	Index      string
	Sort       bson.M
	Filter     bson.M
	Pipeline   []bson.M
	Query      []string
}, ret *json.RawMessage) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	convert := make([]map[string]interface{}, 0)
	collection := me.C_online.Database(me.Db_online).Collection(args.Collection)
	op := options.AggregateOptions{}
	cursor, err := collection.Aggregate(me.Ctx, args.Pipeline, &op)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("document not found:%s", errors.New("NOT FOUND"))
	}
	if err != nil {
		return nil, fmt.Errorf("get cursor error:%s", err)
	}
	if err = cursor.All(me.Ctx, &results); err != nil {
		return nil, fmt.Errorf("find documents error:%s", err)
	}
	for _, item := range results {
		if len(args.Query) == 0 {
			convert = append(convert, item)
		} else {
			temp := make(map[string]interface{})
			for _, v := range args.Query {
				temp[v] = item[v]
			}
			convert = append(convert, temp)
		}
	}
	r, err := json.Marshal(convert)
	if err != nil {
		return nil, fmt.Errorf("json marshal error:%s", err)
	}
	*ret = json.RawMessage(r)
	return convert, nil
}

func (me *T) QueryDocument(args struct {
	Collection string
	Index      string
	Sort       bson.M
	Filter     bson.M
}, ret *json.RawMessage) (map[string]interface{}, error) {
	co := options.CountOptions{}
	collection := me.C_online.Database(me.Db_online).Collection(args.Collection)
	count, err := collection.CountDocuments(me.Ctx, args.Filter, &co)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("NOT FOUNT")
	}
	convert := make(map[string]interface{})
	convert["total counts"] = count
	r, err := json.Marshal(convert)
	if err != nil {
		return nil, fmt.Errorf("json marshal error:%s", err)
	}
	*ret = json.RawMessage(r)
	return convert, nil
}
