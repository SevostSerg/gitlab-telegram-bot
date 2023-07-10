package configuration

import (
	"encoding/json"
	"io/fs"
	"log"
	"os"
	"path"
)

type Configuration struct {
	BotToken                      string
	GitlabAccessToken             string
	Port                          string
	LogsFolderName                string
	DBFolderName                  string
	GitlabURL                     string
	WebhookURL                    string
	WebhookToken                  string
	EncryptionKey                 string
	CommentsTimeLimit             string
	DontSendMessagesSince         int
	DontSendMessagesUntil         int
	LimitOfSecondsBetweenAttempts int
	BanMinutes                    int
	MaxSpamMessages               int
}

const (
	configFileName string      = "config.json"
	configFolder   string      = "ConfigurationData"
	fileMode       fs.FileMode = 0755
)

var configurationInfo *Configuration = nil
var errorMessage string = "Config: must call ReadConf method first, configInfo is empty"

func GetInfoFromConfigFile() *Configuration {
	if _, err := os.Stat(configFolder); os.IsNotExist(err) {
		err := os.Mkdir(configFolder, fileMode)
		if err != nil {
			log.Panic(err)
		}
	}

	configPath := path.Join(configFolder, configFileName)
	file, err := os.Open(configPath)
	if err != nil {
		log.Panic("Unable to open configuration file!")
	}

	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configurationInfo)
	if err != nil {
		log.Println("Config: ", err)
	}

	return configurationInfo
}

func GetConfigInstance() *Configuration {
	if configurationInfo != nil {
		return configurationInfo
	}

	log.Panic(errorMessage)
	panic(errorMessage)
}

func GetConfigFolder() string {
	return configFolder
}

func GetConfigFilename() string {
	return configFileName
}
