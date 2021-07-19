package main

import (
	"context"
	"gopkg.in/yaml.v2"
	"log"
	"neo3fura/biz/api"
	"neo3fura/biz/data"
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
		Data: &data.T{
			Client: client,
		},
	})
	listen := os.ExpandEnv("0.0.0.0:1926")
	log.Println("[LISTEN]", listen)
	http.ListenAndServe(listen, &joh.T{})
}
