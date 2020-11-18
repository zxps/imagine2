package config

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/tkanos/gonfig"
)

// ServiceConfig ...
type ServiceConfig struct {
	ServerAddress string

	FilesPath string
	CachePath string

	DatabaseDSN        string
	DatabaseFilesTable string
	DatabaseDialect    string

	GeneratorPathTimePattern string
	GeneratorFilenameLength  int

	GeneratorMinFileDirIndex int
	GeneratorMaxFileDirIndex int

	ImageDefaultQuality int

	EnableOptimizator bool
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
		logrus.Error("io error. unable to retrieve root path ", err.Error())
		return
	}

	Context.RootPath = rootPath

	logrus.Info("working path: ", Context.RootPath)
}
