package app

import (
	"context"
	"time"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/joshsoftware/golang-boilerplate/config"
	"go.uber.org/zap"
)

var (
	db     *mongo.Database
	client  *mongo.Client
	logger *zap.SugaredLogger
	ctx    context.Context

)

func Init() {
	InitLogger()

	err := initDB()
	if err != nil {
		panic(err)
	}
}

func InitLogger() {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	logger = zapLogger.Sugar()
}

func initDB() (err error) {
	dbConfig := config.Database()
	client, err := mongo.NewClient(dbConfig.ConnectionURL())
	if err != nil { return err }
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return
	}
	db      =  client.Database("function_junction")
	// db, err = sqlx.Open(dbConfig.Driver(), dbConfig.ConnectionURL())
	
	if err = client.Ping(ctx, nil); err != nil {
		return
	}

	// db.SetMaxIdleConns(dbConfig.MaxPoolSize())
	// db.SetMaxOpenConns(dbConfig.MaxOpenConns())
	// db.SetConnMaxLifetime(time.Duration(dbConfig.MaxLifeTimeMins()) * time.Minute)

	return
}
func GetCollection(name string) *mongo.Collection{
	collsection := db.Collection(name)
	return collsection
}
func GetDB() *mongo.Database {
	return db
}

func GetLogger() *zap.SugaredLogger {
	return logger
}

func Close() {
	logger.Sync()
	client.Disconnect(nil)
}
