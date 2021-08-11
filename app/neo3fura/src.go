package main

import (
	"context"
	"fmt"
	"github.com/robfig/cron"
	"gopkg.in/yaml.v2"
	"log"
	"neo3fura/biz/api"
	"neo3fura/biz/job"
	"neo3fura/lib/cli"
	"neo3fura/lib/joh"
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
	client := &cli.T{
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
		err = j.GetPopularTokens()
		err = j.GetHoldersByContractHash()
		err = j.GetNewAddresses()
		err = j.GetActiveAddresses()
		err = j.GetAddressCount()
		err = j.GetContractCount()
		err = j.GetCandidateCount()
		err = j.GetTransactionList()
		err = j.GetTransactionCount()
		err = j.GetBlockCount()
		err = j.GetBlockInfoList()
		err = j.GetAddressCount()
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
