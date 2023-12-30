package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/easonchen147/foundation/cfg"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Mgo struct {
	Client *mongo.Client
	Db     *mongo.Database
}

var (
	mgo *Mgo
)

func init() {
	err := InitMongo(cfg.AppConf)
	if err != nil {
		panic(fmt.Sprintf("init mongo failed: %s", err))
	}
}

func InitMongo(cfg *cfg.AppConfig) error {
	if cfg.MongoConfig == nil {
		return nil
	}
	var err error
	mgo, err = connectMongo(cfg)
	if err != nil {
		return err
	}
	return nil
}

func Mongo() *Mgo {
	if mgo == nil {
		panic(errors.New("mongodb is not ready"))
	}
	return mgo
}

func connectMongo(cfg *cfg.AppConfig) (*Mgo, error) {
	option := options.Client().ApplyURI(cfg.MongoConfig.Uri).
		SetConnectTimeout(time.Duration(cfg.MongoConfig.ConnectTimeout) * time.Second).
		SetMaxConnecting(cfg.MongoConfig.MaxOpenConn).
		SetMaxPoolSize(cfg.MongoConfig.MaxPoolSize).SetMinPoolSize(cfg.MongoConfig.MinPoolSize)
	client, err := mongo.NewClient(option)

	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return &Mgo{
		Client: client,
		Db:     client.Database(cfg.MongoConfig.Db),
	}, nil
}
