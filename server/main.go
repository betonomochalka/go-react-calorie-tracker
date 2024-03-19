package main

import(
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin" //its web app framwork
	"github.com/betonomochalka/go-react-calorie-tracker/routes" //routes
)

func main(){

	port := os.Getenv("PORT") //"PORT" if I have an .env file 
	if port == ""{ //"if struct" if I dont have any
		port = "8000" 
	}
	

	router := gin.New() //creates a router without any middleware by default
	
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	router.Use(gin.Logger())
	router.Use(cors.Default())

	router.POST("/entry/create", routes.CreateEntries)
	router.GET("/entries", routes.GetEntries)
	router.GET("/entries/:id", routes.GetEntriesByID)
	router.GET("/ingredient/:ingredient", routes.GetEntriesByIngredients)

	router.PUT("/entry/update/:id", routes.UpdateEntry)
	router.PUT("/ingredient/update/:id", routes.UpdateIngredient)
	router.DELETE("/entry/delete/:id", routes.DeleteEntry)
	router.Run(":" + port) //runs the server on port 8000
}