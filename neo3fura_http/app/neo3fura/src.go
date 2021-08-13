package main

import (
	"context"
	"fmt"
	"github.com/robfig/cron"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v2"
	"log"
	"neo3fura_http/biz/api"
	"neo3fura_http/biz/job"
	"neo3fura_http/lib/cli"
	"neo3fura_http/lib/joh"
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
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Pass     string `yaml:"pass"`
		Database string `yaml:"database"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`
	Proxy struct {
		Uri []string `yaml:"uri"`
	} `yaml:"proxy"`
}

func main() {
	fmt.Println("YOUR ENV IS " + os.ExpandEnv("${RUNTIME}"))
	cfg, err := OpenConfigFile()
	if err != nil {
		log.Fatalln(err)
	}
	ctx := context.TODO()

	co := intializeMongoOnlineClient(ctx)
	cl := intializeMongoLocalClient(ctx)

	client := &cli.T{
		C_online: co,
		C_local:  cl,
		Ctx:      ctx,
		RpcCli:   neoRpc.NewClient(""), // placeholder
		RpcPorts: cfg.Proxy.Uri,
	}
	rpc.Register(&api.T{
		Client: client,
	})

	j := &job.T{
		Client: client,
	}

	c := cron.New()
	spec := "0 0/10 * * * *"
	err = c.AddFunc(spec, func() {
		go j.GetPopularTokens()
		go j.GetHoldersByContractHash()
		go j.GetNewAddresses()
		go j.GetActiveAddresses()
		go j.GetAddressCount()
		go j.GetContractCount()
		go j.GetCandidateCount()
		go j.GetTransactionList()
		go j.GetTransactionCount()
		go j.GetBlockCount()
		go j.GetBlockInfoList()
		go j.GetAddressCount()
		go j.GetAssetCount()
	})
	if err != nil {
		log.Fatal(err)
	}
	c.Start()

	listen := os.ExpandEnv("0.0.0.0:1926")
	log.Println("[LISTEN]", listen)
	err = http.ListenAndServe(listen, &joh.T{})
	if err != nil {
		log.Fatal(err)
	}
}

func intializeMongoOnlineClient(ctx context.Context) *mongo.Client {
	rt := os.ExpandEnv("${RUNTIME}")
	var clientOptions *options.ClientOptions
	switch rt {
	case "DEV":
		clientOptions = options.Client().ApplyURI("mongodb://" + cfg.Database_Dev.User + ":" + cfg.Database_Dev.Pass + "@" + cfg.Database_Dev.Host + ":" + cfg.Database_Dev.Port + "/" + cfg.Database_Dev.Database)
	case "TEST":
		clientOptions = options.Client().ApplyURI("mongodb://" + cfg.Database_Test.User + ":" + cfg.Database_Test.Pass + "@" + cfg.Database_Test.Host + ":" + cfg.Database_Test.Port + "/" + cfg.Database_Test.Database)
	case "STAGING":
		clientOptions = options.Client().ApplyURI("mongodb://" + cfg.Database_Staging.User + ":" + cfg.Database_Staging.Pass + "@" + cfg.Database_Staging.Host + ":" + cfg.Database_Staging.Port + "/" + cfg.Database_Staging.Database)
	default:
		log.Fatal("err")
	}
	co, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = co.Ping(ctx, nil)
	if err != nil {
		log.Fatal("err")
	}
	return co
}

func intializeMongoLocalClient(ctx context.Context) *mongo.Client {
	var clientOptions *options.ClientOptions
	clientOptions = options.Client().ApplyURI("mongodb://" + cfg.Database_Local.Host + ":" + cfg.Database_Local.Port + "/" + cfg.Database_Local.Database)
	cl, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = cl.Ping(ctx, nil)
	if err != nil {
		log.Fatal("err")
	}
	return cl
}