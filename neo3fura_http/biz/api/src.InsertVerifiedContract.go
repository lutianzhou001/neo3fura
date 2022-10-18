package api

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"neo3fura_http/lib/cli"
	log2 "neo3fura_http/lib/log"
	"neo3fura_http/lib/type/h160"
	"neo3fura_http/var/stderr"
)

func (me *T) InsertVerifiedContract(args struct {
	ContractHash  h160.T
	UpdateCounter int32
	Id            int32
}, ret *json.RawMessage) error {
	if args.ContractHash.Valid() == false {
		return stderr.ErrInvalidArgs
	}
	var clientOptions *options.ClientOptions
	var dbOnline string
	dbname := me.Client.Db_online
	if dbname == "neofura" {
		clientOptions = options.Client().ApplyURI("mongodb://Mindy:QMRhLk9m8rqXWC3X9pMJ@20.106.201.244:27019/bakN3")
		dbOnline = "bakN3"
	} else if dbname == "bakN3" {
		clientOptions = options.Client().ApplyURI("mongodb://Mindy:QMRhLk9m8rqXWC3X9pMJ@20.106.201.244:27018/neofura")
		dbOnline = "neofura"
	}

	clientOptions.SetMaxPoolSize(50)
	co, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log2.Fatalf("mongo connect error:%s", err)
	}

	otherclient := &cli.T{
		Redis:     me.Client.Redis,
		Db_online: dbOnline,
		C_online:  co,
		C_local:   me.Client.C_local,
		Ctx:       me.Client.Ctx,
		RpcCli:    me.Client.RpcCli, // placeholder
		RpcPorts:  me.Client.RpcPorts,
		NeoFs:     me.Client.NeoFs,
	}

	rr1, err := otherclient.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: "VerifyContractModel",
		Index:      "GetVerifiedContract",
		Sort:       bson.M{},
		Filter:     bson.M{"hash": args.ContractHash.Val()},
		Query:      []string{},
	}, ret)

	if len(rr1) > 0 {
		return stderr.ErrExistsDocument
	}

	rr2, err := me.Client.QueryOne(struct {
		Collection string
		Index      string
		Sort       bson.M
		Filter     bson.M
		Query      []string
	}{
		Collection: "VerifyContractModel",
		Index:      "GetVerifiedContract",
		Sort:       bson.M{},
		Filter:     bson.M{"hash": args.ContractHash.Val()},
		Query:      []string{},
	}, ret)

	if len(rr2) > 0 {
		return stderr.ErrExistsDocument
	}

	_, err = otherclient.InsertDocument(struct {
		Collection string
		Index      string
		Insert     *cli.Insert
	}{
		Collection: "VerifyContractModel",
		Index:      "InsertVerifiedContract",
		Insert: &cli.Insert{
			Hash:          args.ContractHash,
			Id:            args.Id,
			UpdateCounter: args.UpdateCounter,
		},
	}, ret)
	if err != nil {
		return err
	}

	r2, err := me.Client.InsertDocument(struct {
		Collection string
		Index      string
		Insert     *cli.Insert
	}{
		Collection: "VerifyContractModel",
		Index:      "InsertVerifiedContract",
		Insert: &cli.Insert{
			Hash:          args.ContractHash,
			Id:            args.Id,
			UpdateCounter: args.UpdateCounter,
		},
	}, ret)

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
