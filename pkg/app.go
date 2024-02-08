package pkg

import (
	"github.com/gin-gonic/gin"
	"org.com/org/pkg/api/routes"
	"org.com/org/pkg/database/mongodb"
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

	// Initialize database
	mongodb.InitDB()

	// Initialize user routes
	routes.InitUserRoutes(app.Router)

	return app
}

// Run starts the application.
func (app *Application) Run(addr string) {
	app.Router.Run(addr)
}