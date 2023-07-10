package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	config "GitlabTgBot/configuration"
	botDB "GitlabTgBot/db"
	messageReceiver "GitlabTgBot/messageReceiver"
	webhookReceiver "GitlabTgBot/webhookReceiver"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	autoUpdateTime int    = 60
	logFileName    string = "Session %d_%d_%d'%dh%dm%ds.txt"
)

var configInfo *config.Configuration

func main() {
	configInfo = config.GetInfoFromConfigFile()
	CreateLogFile()
	db := botDB.StartDB()
	defer db.Close()
	bot, updates := StartBot()
	go webhookReceiver.StartWebhookReceiving()
	messageReceiver.Start(bot, updates)
}

func CreateLogFile() {
	logsFolder := configInfo.LogsFolderName
	logsPath := path.Join(logsFolder, logFileName)
	timeNow := time.Now()
	if _, err := os.Stat(logsFolder); os.IsNotExist(err) {
		os.Mkdir(logsFolder, 0755)
	}

	logsPath = fmt.Sprintf(logsPath,
		timeNow.Day(), timeNow.Month(), timeNow.Year(), timeNow.Hour(), timeNow.Minute(), timeNow.Second())
	file, err := os.OpenFile(logsPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModeAppend)
	if err != nil {
		log.Fatal("error opening log file &s", err)
	}

	log.SetOutput(file)
}

func StartBot() (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel) {
	bot, err := tgbotapi.NewBotAPI(configInfo.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(10)
	u.Timeout = autoUpdateTime
	updates, _ := bot.GetUpdatesChan(u)
	return bot, updates
}
