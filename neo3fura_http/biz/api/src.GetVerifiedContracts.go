package api

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"neo3fura_http/lib/cli"
	log2 "neo3fura_http/lib/log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
)

func (me *T) GetVerifiedContracts(args struct {
	Filter map[string]interface{}
	Limit  int64
	Skip   int64
}, ret *json.RawMessage) error {

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

	r1, _, err := client.QueryAll(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
		Limit      int64
		Skip       int64
	}{
		Collection: getDocumentByEnv("VerifyContractModel"),
		Index:      "GetVerifiedContracts",
		Sort:       bson.M{},
		Filter:     bson.M{},
		Query:      []string{},
		Limit:      args.Limit,
		Skip:       args.Skip,
	}, ret)
	if err != nil {
		return err
	}
	r2, err := json.Marshal(r1)
	if err != nil {
		return err
	}
	*ret = json.RawMessage(r2)
	return nil
}

func getDocumentByEnv(docname string) string {
	rt := os.ExpandEnv("${RUNTIME}")
	if rt != "staging" && rt != "test" && rt != "test2" {
		rt = "mainnet"
	}
	switch rt {
	case "staging":
		docname = "main_" + docname
	case "test":
		docname = "main_" + docname
	case "test2":
		docname = "magnet_" + docname
	}
	return docname
}
