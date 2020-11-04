package main

import (
	"flag"
	"imagine2/config"
	"imagine2/controllers"
	"imagine2/storage"
	"imagine2/utils"
	"os"

	"github.com/fasthttp/router"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func main() {
	configFile := flag.String("c", getDefaultConfigFilepath(), "configuration file path")

	flag.Parse()

	if !utils.IsFileExists(*configFile) {
		logrus.Error("config not found: ", *configFile)
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
	logrus.Info("init storage")
	storage.Initialize()
}

func initializeConfig(configFile string) {
	config.InitializeConfig(configFile)
}

func initializeServer() {
	logrus.Info("init http router")
	router := router.New()

	router.GET("/", controllers.StatsController)
	router.POST("/upload", controllers.UploadController)
	router.POST("/save_base64", controllers.SaveBase64Controller)

	router.GET("/file", controllers.FileController)
	router.GET("/show", controllers.ShowController)
	router.GET("/file/{filepath:*}", controllers.FileByPathController)

	router.GET("/render/{filepath:*}", controllers.RenderController)

	logrus.Info("bind service to address ", config.Context.Service.ServerAddress)

	logrus.Fatal(fasthttp.ListenAndServe(config.Context.Service.ServerAddress, router.Handler))
}

func getDefaultConfigFilepath() string {
	path, err := os.Getwd()
	if err != nil {
		logrus.Error(err.Error())
	}

	return path + "/config.json"
}
