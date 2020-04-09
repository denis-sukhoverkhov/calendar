package calendar

import (
	"flag"
	"github.com/spf13/viper"
	"log"
	"path"
	"strings"
)

type httpListen struct {
	Ip   string
	Port int
}

type logs struct {
	PathToLogFile string
	Level         string // it can be (error/warn/info/debug)
}

type Configuration struct {
	HttpListen httpListen
	Logs       logs
}

func InitConfig() Configuration {
	var pathToConfigFile = flag.String("config", "", "path to configuration file")
	flag.Parse()

	if *pathToConfigFile != "" {
		dir := path.Dir(*pathToConfigFile)
		ext := path.Ext(*pathToConfigFile)
		base := strings.Replace(path.Base(*pathToConfigFile), ext, "", 1)
		viper.SetConfigName(base)
		viper.AddConfigPath(dir)
		viper.SetConfigType(ext[1:])

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file, %s", err)
		}

	} else {
		viper.SetDefault("HttpListen.Ip", "0.0.0.0")
		viper.SetDefault("HttpListen.Port", "8080")
		viper.SetDefault("Logs.PathToLogFile", "all.log")
		viper.SetDefault("Logs.Level", "debug")
	}

	var configuration Configuration

	if err := viper.Unmarshal(&configuration); err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	return configuration
}
