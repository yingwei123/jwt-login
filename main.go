package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"mongoTest.io/auth"
	"mongoTest.io/env"
	"mongoTest.io/mongodb"
	"mongoTest.io/server"
)

func main() {
	cfg, err := env.LoadEnvironment()
	if err != nil {
		log.Fatalf("could not load environment variables: %v", err)
	}

	authenticator := auth.CreateAuthenticator(cfg.JWTKEY)

	serverURL := cfg.ServerBaseURL
	if serverURL == "http://localhost" {
		serverURL = fmt.Sprintf("%s:%d", cfg.ServerBaseURL, cfg.ServerPort)
	}

	resourcesPath := filepath.Join(cfg.ApplicationRootPath, "resources")

	mongodbClient, err := mongodb.CreateMongoClient(cfg.AtlasURI)
	if err != nil {
		log.Fatalf("could not connect to mongodb with uri")
	}

	router := server.Router{
		ResourcesPath: resourcesPath,
		ServerURL:     serverURL,
		MongoDBClient: mongodbClient,
		Authenticator: authenticator,
	}

	s := &http.Server{
		Addr:    fmt.Sprint(":", cfg.ServerPort),
		Handler: router,
	}

	defer mongodbClient.Disconnect()
	log.Printf("starting server on port %d\n", cfg.ServerPort)
	log.Fatal(s.ListenAndServe())
}
