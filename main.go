package main

import (
	"flag"
	"imagine2/config"
	"imagine2/controllers"
	"imagine2/storage"
	"imagine2/utils"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func main() {
	configFile := flag.String("c", getDefaultConfigPath(), "configuration file path")
	flag.Parse()

	if !utils.IsFileExists(*configFile) {
		log.Error("config file not found", *configFile)
		os.Exit(1)
	}

	initialize(*configFile)
}

func initialize(configFile string) {
	initializeConfig(configFile)
	initializeStorage()
	initializeServer()
}

func initializeStorage() {
	log.Info("initializing storage")
	storage.Initialize()
}

func initializeConfig(configFile string) {
	config.InitializeConfig(configFile)
}

func initializeServer() {
	log.Info("initializing service http router")
	router := router.New()

	router.GET("/", controllers.StatsController)
	router.POST("/upload", controllers.UploadController)

	router.GET("/file", controllers.FileController)
	router.GET("/show", controllers.ShowController)
	router.GET("/file/{filepath:*}", controllers.FileByPathController)

	router.GET("/render/{filepath:*}", controllers.RenderController)

	log.Info("bind service to address ", config.Context.Service.Address)
	log.Fatal(fasthttp.ListenAndServe(config.Context.Service.Address, router.Handler))
}

func getDefaultConfigPath() string {
	path, err := os.Getwd()
	if err != nil {
		log.Error(err.Error())
	}

	return path + "/config.json"
}
