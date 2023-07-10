package extensions

import (
	"fmt"
	"time"

	config "GitlabTgBot/configuration"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	unbanMessage string = "You are unbanned! Dont be like that anymore❤️❤️❤️"
	banMessage   string = "You are banned for %v minutes for spam!💩"
)

type ConnectionInfo struct {
	ConnectionAttempt  int
	LastConnectionTime time.Time
}

var activeUsers map[int64]*ConnectionInfo = make(map[int64]*ConnectionInfo)
var bannedUsers map[int64]int64 = make(map[int64]int64)

func IsUserBannedForSpam(userId int64, tgBot *tgbotapi.BotAPI) bool {
	if _, isUserBanned := bannedUsers[userId]; isUserBanned {
		return true
	}

	config := config.GetConfigInstance()
	now := time.Now()
	var lastConnectionAttempt *ConnectionInfo
	if val, contains := activeUsers[userId]; !contains {
		lastConnectionAttempt = &ConnectionInfo{
			ConnectionAttempt:  1,
			LastConnectionTime: now,
		}

		activeUsers[userId] = lastConnectionAttempt
	} else {
		lastConnectionAttempt = val
	}

	if now.Sub(lastConnectionAttempt.LastConnectionTime) > time.Second*time.Duration(config.LimitOfSecondsBetweenAttempts) {
		delete(activeUsers, userId)
		return false
	}

	lastConnectionAttempt.ConnectionAttempt++
	if lastConnectionAttempt.ConnectionAttempt > config.MaxSpamMessages {
		go BanUser(userId, tgBot)
		tgBot.Send(tgbotapi.NewMessage(userId, fmt.Sprintf(banMessage, config.BanMinutes)))
		return true
	}

	return false
}

func BanUser(userId int64, tgBot *tgbotapi.BotAPI) {
	bannedUsers[userId] = userId
	timer := *time.NewTimer(time.Minute * time.Duration(config.GetConfigInstance().BanMinutes))
	<-timer.C
	tgBot.Send(tgbotapi.NewMessage(userId, unbanMessage))
	delete(bannedUsers, userId)
	delete(activeUsers, userId)
}
