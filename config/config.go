package config

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/tkanos/gonfig"
)

// ServiceConfig ...
type ServiceConfig struct {
	Address    string
	FilesPath  string
	CachePath  string
	Database   string
	FilesTable string
}

// ContextConfig ...
type ContextConfig struct {
	// Main Serivce configuration
	Service ServiceConfig

	// Service root path
	RootPath string
}

// Context ...
var Context ContextConfig = ContextConfig{}

// InitializeConfig - initialize service configuration
func InitializeConfig(configFile string) {
	Context.Service = ServiceConfig{}
	gonfig.GetConf(configFile, &Context.Service)

	rootPath, err := os.Getwd()
	if err != nil {
		log.Error("io error. Unable to retrieve root path ", err.Error())
		return
	}

	Context.RootPath = rootPath

	log.Info("current path ", Context.RootPath)
}
