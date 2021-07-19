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
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Pass     string `yaml:"pass"`
		Database string `yaml:"database"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`
	Redis struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
	} `yaml:"redis"`
}

func (me *T) getConnection() (uc *mongo.Client, err error) {
	cfg, err := me.OpenConfigFile()
	if err != nil {
		log.Fatalln(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://" + cfg.Database.User + ":" + cfg.Database.Pass + "@" + cfg.Database.Host + ":" + cfg.Database.Port + "/" + cfg.Database.Database)
	clientOptions = clientOptions.SetMaxPoolSize(50)
	userClient, err := mongo.Connect(ctx, clientOptions)
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
	uc, err := me.getConnection()
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
	uc, err := me.getConnection()
	if err != nil {
		return err
	}
	collections, err := uc.Database(cfg.Database.DBName).ListCollectionNames(me.Ctx, bson.M{})
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
		Addr:   cfg.Redis.Host + ":" + cfg.Redis.Port,
		// Addr:     "docker.for.mac.host.internal:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	val, err := rdb.Get(ctx, hash).Result()
	if err == redis.Nil {
		var result map[string]interface{}
		convert := make(map[string]interface{})
		uc, err := me.getConnection()
		if err != nil {
			return nil, err
		}
		collection := uc.Database(cfg.Database.DBName).Collection(args.Collection)
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
	} else if err != nil {
		return nil, err
	} else {
		// return the data
		r, err := hex.DecodeString(val)
		if err != nil {
			return nil, err
		}
		*ret = json.RawMessage(r)
		convert := make(map[string]interface{})
		err = json.Unmarshal(r, &convert)
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
	uc, err := me.getConnection()
	if err != nil {
		return nil, 0, err
	}
	collection := uc.Database(cfg.Database.DBName).Collection(args.Collection)
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
