package bootstrap

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/amitshekhariitbhu/go-backend-clean-architecture/api/db"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/mongo"
	"gorm.io/gorm"
)

func NewMongoDatabase(env *Env) mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	dbHost := env.DBHost
	dbPort := env.DBPort
	dbUser := env.DBUser
	dbPass := env.DBPass

	mongodbURI := fmt.Sprintf("mongodb://%s:%s@%s:%s", dbUser, dbPass, dbHost, dbPort)

	if dbUser == "" || dbPass == "" {
		mongodbURI = fmt.Sprintf("mongodb://%s:%s", dbHost, dbPort)
	}

	client, err := mongo.NewClient(mongodbURI)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(" err 1 := ", err)
	}

	err = client.Ping(ctx)
	if err != nil {
		log.Fatal(" err 2 := ", err)
	}

	return client
}

func CloseMongoDBConnection(client mongo.Client) {
	if client == nil {
		return
	}

	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connection to MongoDB closed.")
}

func NewMysqlDatabase(env *Env) *gorm.DB {

	return db.NewDatabase().DB
}
