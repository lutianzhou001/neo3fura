package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// T ...
type T struct {
	Db_online string
	C_online  *mongo.Client
}

type Config struct {
	Database_Local struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Pass     string `yaml:"pass"`
		Database string `yaml:"database"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database_local"`
}

func (me *T) GetCollection(args struct {
	Collection string
}) (*mongo.Collection, error) {
	collection := me.C_online.Database(me.Db_online).Collection(args.Collection)
	return collection, nil
}

func (me *T) QueryDocument(args struct {
	Collection string
	Index      string
	Sort       bson.M
	Filter     bson.M
}, ret *json.RawMessage) (map[string]interface{}, error) {
	co := options.CountOptions{}
	collection := me.C_online.Database(me.Db_online).Collection(args.Collection)
	count, err := collection.CountDocuments(context.TODO(), args.Filter, &co)
	if err == mongo.ErrNoDocuments {
		return nil, err
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
