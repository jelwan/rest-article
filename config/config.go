package config

import (
	"github.com/spf13/viper"
	"rest-article/log"
)

const (
	appConfigPath = "data/config/app.yaml"
)

var (
	logger    = log.NewLogger().WithField("module", "config")
	configApp *AppConfig
)

func init() {
	configApp = loadAppConfig()
}

func App() *AppConfig {
	if configApp == nil {
		configApp = loadAppConfig()
	}
	return configApp
}

func loadAppConfig() *AppConfig {
	configMap := map[string]*AppConfig{}

	_, err := readConfigFromFile("app", appConfigPath, &configMap)
	if err != nil {
		logger.Fatalf("Unable to read app config with error %v", err)
	}

	app, ok := configMap["defaults"]
	if app == nil || !ok {
		logger.Fatalf("unable to get environment [defaults] in config [%s] with error: %v", appConfigPath, err)
	}

	return app
}

func readConfigFromFile(filename, path string, target interface{}) (*viper.Viper, error) {
	vip := viper.New()

	vip.SetConfigType("yaml")
	vip.SetConfigFile(path)
	vip.AutomaticEnv()

	err := vip.ReadInConfig()
	if err != nil {
		logger.WithFields(log.ErrorFields("config", "ReadConfig")).
			Errorf("Unable to read config [%s] with error %v", filename, err)
		return nil, err
	}

	err = vip.Unmarshal(&target)
	if err != nil {
		logger.WithFields(log.ErrorFields("config", "Unmarshal")).
			Errorf("Unable to Unmarshall config [%s] with error: %v", filename, err)
		return nil, err
	}

	return vip, nil
}
