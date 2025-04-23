package main

import (
	"context"
	"log"

	_ "memory-book/docs"
	"memory-book/internal/app"
)

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server Petstore server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		petstore.swagger.io
// @BasePath	/v2
func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app, err := app.New()
	if err != nil {
		panic(err)
	}

	err = app.Run(ctx)
	if err != nil {
		panic(err)
	}
	log.Print("Server started")
	//TODO log and db вынести сюда или возвращать из прилжения?

}
