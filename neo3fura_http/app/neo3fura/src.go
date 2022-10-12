package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v2"
	"neo3fura_http/biz/api"
	"neo3fura_http/biz/job"
	"neo3fura_http/biz/watch"
	"neo3fura_http/lib/cli"
	"neo3fura_http/lib/joh"
	log2 "neo3fura_http/lib/log"
	"neo3fura_http/lib/monitor"
	"neo3fura_http/lib/verify"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"
	"time"

	"github.com/go-redis/redis/v8"
	neoRpc "github.com/joeqian10/neo3-gogogo/rpc"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron"
	"github.com/rs/cors"
)

func OpenConfigFile() (Config, error) {
	absPath, _ := filepath.Abs("./config.yml")
	f, err := os.Open(absPath)
	if err != nil {
		return Config{}, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log2.Fatalf("Closing file error: %v", err)
		}
	}(f)
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
	Replica    string `yaml:"replica"`
	NeoFs_Main struct {
		Host        string `yaml:"host"`
		Port        string `yaml:"port"`
		ContainerId string `yaml:"containerid"`
	} `yaml:"neofs_main"`
	NeoFs_Test struct {
		Host        string `yaml:"host"`
		Port        string `yaml:"port"`
		ContainerId string `yaml:"containerid"`
	} `yaml:"neofs_test"`
}

func main() {
	log2.InitLog(1, "./Logs/", os.Stdout)
	log2.InitLog(12) //test
	log2.Infof("YOUR ENV IS %s", os.ExpandEnv("${RUNTIME}"))
	cfg, err := OpenConfigFile()
	if err != nil {
		log2.Fatalf("open file error:%s", err)
	}
	ctx := context.TODO()
	co, dbOnline := initializeMongoOnlineClient(cfg, ctx)
	cl := initializeMongoLocalClient(cfg, ctx)
	rds := initializeRedisLocalClient(cfg, ctx)
	fs := initializeNeoFsHost(cfg)

	client := &cli.T{
		Redis:     rds,
		Db_online: dbOnline,
		C_online:  co,
		C_local:   cl,
		Ctx:       ctx,
		RpcCli:    neoRpc.NewClient(""), // placeholder
		RpcPorts:  cfg.Proxy.Uri,
		NeoFs:     fs,
	}

	rpc.Register(&api.T{
		Client: client,
	})

	j := &job.T{
		Client: client,
	}

	w := &watch.T{
		Client: client,
	}

	h := &joh.T{}
	v := &verify.T{
		Client: client,
	}

	// reset qps
	go func() {
		for {
			monitor.Http_request_qps.Set(0)
			time.Sleep(1 * time.Second)
		}
	}()

	if cfg.Replica == "master" {
		go func() {
			err := w.GetFirstEventByTransactionHash()
			if err != nil {
				log2.Fatalf("run watching error:%v", err)
			}
		}()

		c1 := cron.New()
		c2 := cron.New()
		c3 := cron.New()

		err = c1.AddFunc("@daily", func() {
			log2.Infof("Start daily job")
			go j.GetPopularTokens()
			go j.GetDailyTransactions()
			go j.GetNewAddresses()
			go j.GetActiveAddresses()
			go j.GetMarketDailyVolume() //获取market 前一天的交易数据
		})
		err = c2.AddFunc("@hourly", func() { //@hourly
			log2.Infof("Start hourly job")
			go j.GetHoldersByContractHash()
			go j.GetTransactionList()
			go j.GetBlockInfoList()
			go j.GetHourlyTransactions()
			go j.GetMarketHourlyVolume() //获取market当天的交易数据
		})

		err = c3.AddFunc("@every 10m", func() {
			log2.Infof("Start mintnue job")
			go j.GetMarketSupply()
			go j.GetMarketTxAmount()
			go j.GetMarketOwnerCount()
			go j.GetNFTFloorPrice()
			//go j.GetNFTIndex()
		})
		if err != nil {
			log2.Fatal("add job function error:%s", err)
		}
		c1.Start()
		c2.Start()
		c3.Start()
	}

	listen := os.ExpandEnv("0.0.0.0:1926")
	log2.Infof("NOW LISTEN ON: %s", listen)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		monitor.Http_request_qps.Inc()
		monitor.Http_request_total.Inc()
		monitor.Http_request_in_flight.Inc()
		defer monitor.Http_request_in_flight.Dec()
		monitor.Http_request_duration_seconds.Observe(time.Since(time.Now()).Seconds())
		h.ServeHTTP(writer, request)
	})
	mux.HandleFunc("/upload", func(writer http.ResponseWriter, request *http.Request) {
		v.MultipleFile(writer, request)
	})
	mux.Handle("/metrics", promhttp.Handler())
	handler := cors.Default().Handler(mux)
	err = http.ListenAndServe(listen, handler)
	if err != nil {
		log2.Fatalf("listen and server error:%s", err)
	}
}

func initializeMongoOnlineClient(cfg Config, ctx context.Context) (*mongo.Client, string) {
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
	case "test2":
		clientOptions = options.Client().ApplyURI("mongodb://" + cfg.Database_Test.User + ":" + cfg.Database_Test.Pass + "@" + cfg.Database_Test.Host + ":" + cfg.Database_Test.Port + "/" + cfg.Database_Test.Database)
		dbOnline = cfg.Database_Test.Database
	case "staging":
		clientOptions = options.Client().ApplyURI("mongodb://" + cfg.Database_Staging.User + ":" + cfg.Database_Staging.Pass + "@" + cfg.Database_Staging.Host + ":" + cfg.Database_Staging.Port + "/" + cfg.Database_Staging.Database)
		dbOnline = cfg.Database_Staging.Database
	default:
		log2.Fatalf("runtime environment mismatch")
	}

	clientOptions.SetMaxPoolSize(50)
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
func initializeNeoFsHost(cfg Config) string {
	rt := os.ExpandEnv("${RUNTIME}")
	var neoFsHost string
	switch rt {
	case "test":
		neoFsHost = cfg.NeoFs_Test.Host + ":" + cfg.NeoFs_Test.Port + "/gate" + "/get/" + cfg.NeoFs_Test.ContainerId + "/"
	case "staging":
		neoFsHost = cfg.NeoFs_Main.Host + ":" + cfg.NeoFs_Main.Port + "/gate" + "/get/" + cfg.NeoFs_Main.ContainerId + "/"
	default:
		log2.Fatalf("runtime environment mismatch")
	}
	return neoFsHost
}
func initializeMongoLocalClient(cfg Config, ctx context.Context) *mongo.Client {
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
