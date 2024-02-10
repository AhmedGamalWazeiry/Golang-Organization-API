package pkg

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"org.com/org/pkg/api/routes"
	"org.com/org/pkg/database/mongodb"
	"org.com/org/pkg/utils"
)

// Application represents the application.
type Application struct {
	Router *gin.Engine
}

// NewApp initializes and returns a new Application instance.
func NewApp() *Application {
	app := &Application{
		Router: gin.Default(),
	}

	viper.SetConfigType("yaml")
    viper.SetConfigFile("./config/database-config.yaml") 
	
	err := viper.ReadInConfig()
	if err != nil {
		panic("Failed to read the configuration file")
	}

	dbName := viper.GetString("mongodb.database")
	uri := viper.GetString("mongodb.uri")

	// Initialize database
	mongodb.InitDB(dbName,uri)

	utils.InitRedis()

	// Initialize user routes
	routes.InitUserRoutes(app.Router)
	routes.InitOrganizationRoutes(app.Router)

	return app
}

// Run starts the application.
func (app *Application) Run(addr string) {
	app.Router.Run(addr)
}