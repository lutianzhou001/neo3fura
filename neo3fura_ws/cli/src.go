package cli

import (
	"context"
	"fmt"
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
type T struct{}

type Config struct {
	Database_Local struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Pass     string `yaml:"pass"`
		Database string `yaml:"database"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database_local"`
}

func (me *T) getConnection(database string) (uc *mongo.Client, err error) {
	cfg, err := me.OpenConfigFile()
	if err != nil {
		log.Fatalln(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	co := options.Client().ApplyURI("mongodb://" + cfg.Database_Local.Host + ":" + cfg.Database_Local.Port + "/" + cfg.Database_Local.Database)
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

func (me *T) GetCollection(args struct {
	Collection string
}) (*mongo.Collection, error) {
	cfg, err := me.OpenConfigFile()
	if err != nil {
		return nil, err
	}
	uc, err := me.getConnection("LOCAL")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	dbName := cfg.Database_Local.DBName
	collection := uc.Database(dbName).Collection(args.Collection)
	return collection, nil
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
	dbName := cfg.Database_Local.DBName
	collection := uc.Database(dbName).Collection(args.Collection)
	var result map[string]interface{}
	opts := options.FindOne().SetSort(bson.M{"_id": -1})
	err = collection.FindOne(context.TODO(), bson.M{}, opts).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
