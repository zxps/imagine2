package main

import (
	"flag"
	"imagine2/config"
	"imagine2/controllers"
	"imagine2/storage"
	"imagine2/tasks"
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
	initializeTasks()
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

	router.GET("/", controllers.Stats)
	router.GET("/file", controllers.File)
	router.GET("/files", controllers.FilesController)
	router.GET("/delete", controllers.Delete)
	router.GET("/show", controllers.Show)
	router.GET("/render/{filepath:*}", controllers.Render)

	router.POST("/upload", controllers.Upload)
	router.POST("/save_base64", controllers.SaveBase64)

	logrus.Info("bind service to address ", config.Context.Service.ServerAddress)

	server := &fasthttp.Server{
		Handler:         router.Handler,
		WriteBufferSize: 2048,
		ReadBufferSize:  2048,
	}

	err := server.ListenAndServe(config.Context.Service.ServerAddress)

	logrus.Fatal(err)
}

func initializeTasks() {
	if config.Context.Service.EnableOptimizator {
		logrus.Info("start optimizator")
		tasks.StartOptimizator()
	}
}

func getDefaultConfigFilepath() string {
	path, err := os.Getwd()
	if err != nil {
		logrus.Error(err.Error())
	}

	return path + "/config.json"
}
