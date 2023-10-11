package bootstrap

import (
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/mongo"
	"gorm.io/gorm"
)

type Application struct {
	Env   *Env
	Mongo mongo.Client
	Mysql *gorm.DB
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	app.Mongo = NewMongoDatabase(app.Env)
	// app.Mysql = NewMysqlDatabase(app.Env)
	return *app
}

func (app *Application) CloseDBConnection() {
	CloseMongoDBConnection(app.Mongo)
}
