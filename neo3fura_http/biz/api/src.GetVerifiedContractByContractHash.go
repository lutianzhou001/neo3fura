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

func (me *T) GetVerifiedContractByContractHash(args struct {
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

	clientOptions := options.Client().ApplyURI("mongodb://Mindy:QMRhLk9m8rqXWC3X9pMJ@10.0.7.38:27018/ContractSource")
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

	r1, err := client.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: getDocumentByEnv("VerifyContractModel"),
		Index:      "GetVerifiedContractByContractHash",
		Sort:       bson.M{},
		Filter:     bson.M{"hash": args.ContractHash, "updatecounter": args.UpdateCounter},
		Query:      []string{},
	}, ret)

	if err != nil {
		return err
	}

	r2, err := me.Filter(r1, args.Filter)
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
