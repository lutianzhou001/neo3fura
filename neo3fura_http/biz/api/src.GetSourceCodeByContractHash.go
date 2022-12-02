package api

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"neo3fura_http/lib/cli"
	log2 "neo3fura_http/lib/log"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/lib/type/uintval"
	"neo3fura_http/var/stderr"

	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetSourceCodeByContractHash(args struct {
	Limit         int64
	Skip          int64
	ContractHash  h160.T
	UpdateCounter uintval.T
	Filter        map[string]interface{}
}, ret *json.RawMessage) error {
	if args.ContractHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	if args.UpdateCounter.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	if args.Limit <= 0 {
		args.Limit = 20
	}
	clientOptions := options.Client().ApplyURI("mongodb://Mindy:QMRhLk9m8rqXWC3X9pMJ@20.106.201.244:27019/ContractSource")
	dbOnline := "ContractSource"
	clientOptions.SetMaxPoolSize(50)
	co, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log2.Fatalf("mongo connect error:%s", err)
	}

	client := &cli.T{
		Redis:     me.Client.Redis,
		Db_online: dbOnline,
		C_online:  co,
		C_local:   me.Client.C_local,
		Ctx:       me.Client.Ctx,
		RpcCli:    me.Client.RpcCli, // placeholder
		RpcPorts:  me.Client.RpcPorts,
		NeoFs:     me.Client.NeoFs,
	}

	r1, count, err := client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: getDocumentByEnv("ContractSourceCode"),
		Index:      "someIndex",
		Sort:       bson.M{},
		Filter:     bson.M{"hash": args.ContractHash, "updatecounter": args.UpdateCounter},
		Query:      []string{},
		Limit:      args.Limit,
		Skip:       args.Skip,
	}, ret)
	if err != nil {
		return err
	}
	r2, err := me.FilterArrayAndAppendCount(r1, count, args.Filter)
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
