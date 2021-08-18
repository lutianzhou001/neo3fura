package main

import (
	"context"
	"encoding/json"
	"flag"
	"neo3fura_ws/home"
	"neo3fura_ws/lib/cli"
	log2 "neo3fura_ws/lib/log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v2"
)

var add = flag.String("addr", "0.0.0.0:2026", "http service address")
var upgrader = websocket.Upgrader{
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

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


func mainpage(w http.ResponseWriter, r *http.Request) {
	log2.Infof("DETECT CONNECTION")
	cfg, err := OpenConfigFile()
	if err != nil {
		log2.Fatalf("open file error:%s", err)
	}
	ctx := context.TODO()
	co, dbOnline := intializeMongoOnlineClient(cfg, ctx)
	client := &cli.T{
		Db_online: dbOnline,
		C_online:  co,
	}
	c := &home.T{
		Client: client,
	}
	wsc, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log2.Fatalf("upgrade error:%s", err)
	}
	mt, _, err := wsc.ReadMessage()
	if err != nil {
		log2.Fatalf("read message error:%s", err)
	}

	var responseChannel = make(chan map[string]interface{}, 20)

	go c.GetAddressCount(&responseChannel)
	go c.GetAssetCount(&responseChannel)
	go c.GetBlockCount(&responseChannel)
	// go c.GetBlockInfoList(&responseChannel)
	go c.GetCandidateCount(&responseChannel)
	go c.GetContractCount(&responseChannel)
	go c.GetTransactionCount(&responseChannel)
	// go c.GetTransactionList(&responseChannel)
	go ResponseController(mt, wsc, &responseChannel)
}

func ResponseController(mt int, wsc *websocket.Conn, ch *chan map[string]interface{}) {
	str := "hello neo3fura"
	err := wsc.WriteMessage(mt, []byte(str))
	if err != nil {
		log2.Fatalf("write hello message error:%s", err)
	}
	for {
		b := <-*ch
		sent, err := json.Marshal(b)
		if err != nil {
			log2.Fatalf("json marshal error:%s", err)
		}
		err = wsc.WriteMessage(mt, sent)
		if err != nil {
			log2.Fatalf("write message error:%s", err)
		}
	}
}

func main() {
	http.HandleFunc("/home", mainpage)
	log2.Fatal(http.ListenAndServe(*add, nil))
}
