package main

import (
	"context"
	"fmt"
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

	// "strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
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
}

func main() {
	cfg, err := OpenConfigFile()
	if err != nil {
		log.Fatalln(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1000000*time.Hour)
	c, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+cfg.Database.User+":"+cfg.Database.Pass+"@"+cfg.Database.Host+":"+cfg.Database.Port+"/"+cfg.Database.Database))
	fmt.Println("connected")
	defer cancel()
	//address := os.ExpandEnv("${NEODB_ADDRESS}")
	//poolsize, err := strconv.Atoi(os.ExpandEnv("${NEODB_POOLSIZE}"))
	if err != nil {
		log.Fatalln(err)
	}
	client := &cli.T{
		C:   c,
		Ctx: ctx,
	}
	defer cancel()
	rpc.Register(&api.T{
		Data: &data.T{
			Client: client,
		},
	})
	listen := os.ExpandEnv("0.0.0.0:1926")
	log.Println("[LISTEN]", listen)
	http.ListenAndServe(listen, &joh.T{})
}
