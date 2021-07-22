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
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"time"
)

// T ...
type T struct {
	Ctx      context.Context
	RpcCli   *rpc.RpcClient
	RpcPorts []string
}

type Config struct {
	Database_Dev struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Pass     string `yaml:"pass"`
		Database string `yaml:"database"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database_dev"`
	Database_Test struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Pass     string `yaml:"pass"`
		Database string `yaml:"database"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database_test"`
	Database_Staging struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Pass     string `yaml:"pass"`
		Database string `yaml:"database"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database_staging"`
	Database_Local struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Pass     string `yaml:"pass"`
		Database string `yaml:"database"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database_local"`
	Redis struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"redis"`
}

func (me *T) chooseDatabase(database string, cfg Config) (co *options.ClientOptions, err error) {
	switch database {
	case "DEV":
		clientOptions := options.Client().ApplyURI("mongodb://" + cfg.Database_Dev.User + ":" + cfg.Database_Dev.Pass + "@" + cfg.Database_Dev.Host + ":" + cfg.Database_Dev.Port + "/" + cfg.Database_Dev.Database)
		return clientOptions, nil
	case "TEST":
		clientOptions := options.Client().ApplyURI("mongodb://" + cfg.Database_Test.User + ":" + cfg.Database_Test.Pass + "@" + cfg.Database_Test.Host + ":" + cfg.Database_Test.Port + "/" + cfg.Database_Test.Database)
		return clientOptions, nil
	case "STAGING":
		clientOptions := options.Client().ApplyURI("mongodb://" + cfg.Database_Staging.User + ":" + cfg.Database_Staging.Pass + "@" + cfg.Database_Staging.Host + ":" + cfg.Database_Staging.Port + "/" + cfg.Database_Staging.Database)
		return clientOptions, nil
	case "LOCAL":
		clientOptions := options.Client().ApplyURI("mongodb://" + cfg.Database_Local.Host + ":" + cfg.Database_Local.Port + "/" + cfg.Database_Local.Database)
		return clientOptions, nil
	default:
		return nil, err
	}
}

func (me *T) getDbName(cfg Config, database string) string {
	switch database {
	case "DEV":
		dbName := cfg.Database_Dev.DBName
		return dbName
	case "TEST":
		dbName := cfg.Database_Test.DBName
		return dbName
	case "STAGING":
		dbName := cfg.Database_Staging.DBName
		return dbName
	case "LOCAL":
		dbName := cfg.Database_Local.DBName
		return dbName
	default:
		return ""
	}
}

func (me *T) getConnection(database string) (uc *mongo.Client, err error) {
	cfg, err := me.OpenConfigFile()
	if err != nil {
		log.Fatalln(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	co, err := me.chooseDatabase(database, cfg)
	if err != nil {
		return nil, err
	}
	co = co.SetMaxPoolSize(50)
	userClient, err := mongo.Connect(ctx, co)
	if err != nil {
		log.Fatal(err)
	}
	err = userClient.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	return userClient, nil
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
	uc, err := me.getConnection(os.Getenv("RUNTIME"))
	if err != nil {
		return err
	}
	databases, err := uc.ListDatabaseNames(me.Ctx, bson.M{})
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
	uc, err := me.getConnection(os.Getenv("RUNTIME"))
	if err != nil {
		return err
	}
	dbName := me.getDbName(cfg, os.Getenv("RUNTIME"))
	collections, err := uc.Database(dbName).ListCollectionNames(me.Ctx, bson.M{})
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

	// connect to redis
	// if found return conver,ret
	// if not found
	cfg, err := me.OpenConfigFile()
	if err != nil {
		return nil, err
	}

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

	var ctx = context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	val, err := rdb.Get(ctx, hash).Result()
	// if sort != nil, it may have several results, we have to pick the sorted one
	if err == redis.Nil || args.Sort != nil {
		var result map[string]interface{}
		convert := make(map[string]interface{})
		uc, err := me.getConnection(os.Getenv("RUNTIME"))
		if err != nil {
			return nil, err
		}
		dbName := me.getDbName(cfg, os.Getenv("RUNTIME"))
		collection := uc.Database(dbName).Collection(args.Collection)

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
		err = rdb.Set(ctx, hash, hex.EncodeToString(r), 0).Err()
		if err != nil {
			return nil, err
		}
		*ret = json.RawMessage(r)
		return convert, err
	} else {
		r, err := hex.DecodeString(val)
		if err != nil {
			return nil, err
		}

		*ret = json.RawMessage(r)
		convert := make(map[string]interface{})
		err = json.Unmarshal(r, &convert)
		if convert["_id"] != nil {
			convert["_id"], err = primitive.ObjectIDFromHex(convert["_id"].(string))
		}
		if err != nil {
			return nil, err
		}
		return convert, err
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
	cfg, err := me.OpenConfigFile()
	if err != nil {
		return nil, 0, err
	}
	var results []map[string]interface{}
	convert := make([]map[string]interface{}, 0)
	uc, err := me.getConnection(os.Getenv("RUNTIME"))
	if err != nil {
		return nil, 0, err
	}
	dbName := me.getDbName(cfg, os.Getenv("RUNTIME"))
	collection := uc.Database(dbName).Collection(args.Collection)
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

func (me *T) SaveJob(args struct {
	Collection string
	Data       bson.M
}) (bool, error) {
	cfg, err := me.OpenConfigFile()
	if err != nil {
		return false, err
	}

	uc, err := me.getConnection("LOCAL")
	if err != nil {
		return false, err
	}
	dbName := me.getDbName(cfg, "LOCAL")
	collection := uc.Database(dbName).Collection(args.Collection)
	_, err = collection.InsertOne(me.Ctx, args.Data)
	return true, nil
}

func (me *T) QueryLastJob(args struct {
	Collection string
}) (map[string]interface{}, error) {
	cfg, err := me.OpenConfigFile()
	if err != nil {
		return nil, err
	}
	uc, err := me.getConnection("LOCAL")
	if err != nil {
		return nil, err
	}
	dbName := me.getDbName(cfg, "LOCAL")
	collection := uc.Database(dbName).Collection("PopularTokens")
	var result map[string]interface{}
	opts := options.FindOne().SetSort(bson.M{"_id": -1})
	err = collection.FindOne(me.Ctx, bson.M{}, opts).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
