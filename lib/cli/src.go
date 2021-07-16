package cli

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joeqian10/neo3-gogogo/rpc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
)

// T ...
type T struct {
	C        *mongo.Client
	Ctx      context.Context
	RpcCli   *rpc.RpcClient
	RpcPorts []string
}

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Pass     string `yaml:"pass"`
		Database string `yaml:"database"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`
}

func (me *T) OpenConfigFile() (Config, error) {
	absPath, _ := filepath.Abs("./config.yml")
	f, err := os.Open(absPath)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()
	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, err
}

func (me *T) ListDatabaseNames() error {
	databases, err := me.C.ListDatabaseNames(me.Ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)
	return nil
}

func (me *T) ListCollections() error {
	cfg, err := me.OpenConfigFile()
	if err != nil {
		return err
	}
	collections, err := me.C.Database(cfg.Database.DBName).ListCollectionNames(me.Ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(collections)
	return nil
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
	str := hex.EncodeToString(h.Sum(nil))
	fmt.Println(str)

	cfg, err := me.OpenConfigFile()
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	convert := make(map[string]interface{})
	collection := me.C.Database(cfg.Database.DBName).Collection(args.Collection)
	opts := options.FindOne().SetSort(args.Sort)
	err = collection.FindOne(me.Ctx, args.Filter, opts).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("NOT FOUND")
	} else if err != nil {
		log.Fatal(err)
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
		return nil, err
	}
	*ret = json.RawMessage(r)
	return convert, err
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
	cfg, err := me.OpenConfigFile()
	if err != nil {
		return nil, 0, err
	}
	var results []map[string]interface{}
	convert := make([]map[string]interface{}, 0)
	collection := me.C.Database(cfg.Database.DBName).Collection(args.Collection)
	op := options.Find()
	op.SetSort(args.Sort)
	op.SetLimit(args.Limit)
	op.SetSkip(args.Skip)
	co := options.CountOptions{}
	count, err := collection.CountDocuments(me.Ctx, args.Filter, &co)
	cursor, err := collection.Find(me.Ctx, args.Filter, op)
	if err == mongo.ErrNoDocuments {
		return nil, 0, errors.New("NOT FOUNT")
	}
	if err != nil {
		return nil, 0, err
	}
	if err = cursor.All(me.Ctx, &results); err != nil {
		return nil, 0, err
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
		return nil, 0, err
	}
	*ret = json.RawMessage(r)
	return convert, count, nil
}
