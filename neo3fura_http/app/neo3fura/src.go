package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/robfig/cron"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v2"
	"neo3fura_http/biz/api"
	"neo3fura_http/biz/job"
	"neo3fura_http/lib/cli"
	"neo3fura_http/lib/joh"
	log2 "neo3fura_http/lib/log"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"

	neoRpc "github.com/joeqian10/neo3-gogogo/rpc"
)

func OpenConfigFile() (Config, error) {
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
	Proxy struct {
		Uri []string `yaml:"uri"`
	} `yaml:"proxy"`
}

func main() {
	log2.Infof("YOUR ENV IS %s", os.ExpandEnv("${RUNTIME}"))
	cfg, err := OpenConfigFile()
	if err != nil {
		log2.Fatalf("open file error:%s", err)
	}
	ctx := context.TODO()
	co, dbOnline := intializeMongoOnlineClient(cfg, ctx)
	cl := intializeMongoLocalClient(cfg, ctx)
	rds := initializeRedisLocalClient(cfg, ctx)

	client := &cli.T{
		Redis:     rds,
		Db_online: dbOnline,
		C_online:  co,
		C_local:   cl,
		Ctx:       ctx,
		RpcCli:    neoRpc.NewClient(""), // placeholder
		RpcPorts:  cfg.Proxy.Uri,
	}
	rpc.Register(&api.T{
		Client: client,
	})

	j := &job.T{
		Client: client,
	}

	c1 := cron.New()
	spec1 := "* 11 * * * *"

	c2 := cron.New()
	spec2 := "* 21 * * * *"

	c3 := cron.New()
	spec3 := "* 31 * * * *"

	c4 := cron.New()
	spec4 := "* 41 * * * *"

	c5 := cron.New()
	spec5 := "* 51 * * * *"

	c6 := cron.New()
	spec6 := "* 01 * * * *"


	err = c1.AddFunc(spec1, func() {
		go j.GetPopularTokens()
		//go j.GetHoldersByContractHash()
		//go j.GetNewAddresses()
		//go j.GetActiveAddresses()
		//go j.GetTransactionList()
		//go j.GetBlockInfoList()
	})

	err = c2.AddFunc(spec2, func() {
		//go j.GetPopularTokens()
		go j.GetHoldersByContractHash()
		//go j.GetNewAddresses()
		//go j.GetActiveAddresses()
		//go j.GetTransactionList()
		//go j.GetBlockInfoList()
	})

	err = c3.AddFunc(spec3, func() {
		//go j.GetPopularTokens()
		//go j.GetHoldersByContractHash()
		go j.GetNewAddresses()
		//go j.GetActiveAddresses()
		//go j.GetTransactionList()
		//go j.GetBlockInfoList()
	})

	err = c4.AddFunc(spec4, func() {
		//go j.GetPopularTokens()
		//go j.GetHoldersByContractHash()
		//go j.GetNewAddresses()
		go j.GetActiveAddresses()
		//go j.GetTransactionList()
		//go j.GetBlockInfoList()
	})

	err = c5.AddFunc(spec5, func() {
		//go j.GetPopularTokens()
		//go j.GetHoldersByContractHash()
		//go j.GetNewAddresses()
		//go j.GetActiveAddresses()
		go j.GetTransactionList()
		//go j.GetBlockInfoList()
	})

	err = c6.AddFunc(spec6, func() {
		//go j.GetPopularTokens()
		//go j.GetHoldersByContractHash()
		//go j.GetNewAddresses()
		//go j.GetActiveAddresses()
		//go j.GetTransactionList()
		go j.GetBlockInfoList()
	})

	if err != nil {
		log2.Fatal("add job function error:%s", err)
	}
	c1.Start()
	c2.Start()
	c3.Start()
	c4.Start()
	c5.Start()
	c6.Start()

	listen := os.ExpandEnv("0.0.0.0:1926")
	log2.Infof("NOW LISTEN ON: %s", listen)
	err = http.ListenAndServe(listen, &joh.T{})
	if err != nil {
		log2.Fatalf("linsten and server error:%s", err)
	}
}

func intializeMongoOnlineClient(cfg Config, ctx context.Context) (*mongo.Client, string) {
	rt := os.ExpandEnv("${RUNTIME}")
	var clientOptions *options.ClientOptions
	var dbOnline string
	switch rt {
	case "dev":
		clientOptions = options.Client().ApplyURI("mongodb://" + cfg.Database_Dev.User + ":" + cfg.Database_Dev.Pass + "@" + cfg.Database_Dev.Host + ":" + cfg.Database_Dev.Port + "/" + cfg.Database_Dev.Database)
		dbOnline = cfg.Database_Dev.Database
	case "test":
		clientOptions = options.Client().ApplyURI("mongodb://" + cfg.Database_Test.User + ":" + cfg.Database_Test.Pass + "@" + cfg.Database_Test.Host + ":" + cfg.Database_Test.Port + "/" + cfg.Database_Test.Database)
		dbOnline = cfg.Database_Test.Database
	case "staging":
		clientOptions = options.Client().ApplyURI("mongodb://" + cfg.Database_Staging.User + ":" + cfg.Database_Staging.Pass + "@" + cfg.Database_Staging.Host + ":" + cfg.Database_Staging.Port + "/" + cfg.Database_Staging.Database)
		dbOnline = cfg.Database_Staging.Database
	default:
		log2.Fatalf("runtime environment mismatch")
	}
	co, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log2.Fatalf("mongo connect error:%s", err)
	}
	err = co.Ping(ctx, nil)
	if err != nil {
		log2.Fatalf("ping mongo error:%s", err)
	}
	return co, dbOnline
}

func intializeMongoLocalClient(cfg Config, ctx context.Context) *mongo.Client {
	var clientOptions *options.ClientOptions
	clientOptions = options.Client().ApplyURI("mongodb://" + cfg.Database_Local.Host + ":" + cfg.Database_Local.Port + "/" + cfg.Database_Local.Database)
	cl, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log2.Fatalf("connect to mongo error:%s", err)
	}
	err = cl.Ping(ctx, nil)
	if err != nil {
		log2.Fatalf("ping mongo error:%s", err)
	}
	return cl
}

func initializeRedisLocalClient(cfg Config, ctx context.Context) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return rdb
}
