package main

import (
	"fmt"
	"log"
	"net/http"

	middlewares "web-forum/internal/Middlewares"
	"web-forum/internal/app"
	"web-forum/internal/routers"
	"web-forum/pkg/config"
	"web-forum/pkg/logger"
)

func main() {
	logger, err := logger.Create_Logger()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	configurations := config.LoadConfig()
	application := app.NewApp(configurations)
	mux := routers.SetupRoutes(application)
	defer application.DB.Close()

	addr := fmt.Sprintf(":%d", configurations.Port)
	fmt.Printf("Server is running on http://localhost%s\n", addr)

	// lets wrap our mux within the cors middleware to enable cors safty acces
	log.Fatal(http.ListenAndServe(addr, middlewares.CORSMiddleware(mux)))

	/// dont forgot to add garefuly shutdown server
}
