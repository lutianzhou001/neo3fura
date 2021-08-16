package cli

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// T ...
type T struct{
	C_local  *mongo.Client
}

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

func (me *T) GetCollection(args struct {
	Collection string
}) (*mongo.Collection, error) {
	collection := me.C_local.Database("job").Collection(args.Collection)
	return collection, nil
}