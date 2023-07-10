package webhookReceiver

import (
	"GitlabTgBot/configuration"
	extensions "GitlabTgBot/extensions"
	"time"
)

func handleHookSilent(payload interface{}) {
	now := time.Now()
	currMinute := int(now.Minute())
	currHour := int(now.Hour())
	remainingTime := extensions.Abs(configuration.GetConfigInstance().DontSendMessagesUntil-(currHour-1))*60 + (60 - currMinute)
	timer := *time.NewTimer(time.Minute * time.Duration(remainingTime+1))
	<-timer.C
	ParseWebhook(payload)
}

// OverBicycle64
func isItSilentTime() bool {
	config := configuration.GetConfigInstance()
	since := config.DontSendMessagesSince
	until := config.DontSendMessagesUntil
	currHour := time.Now().Hour()
	if since > until {
		// 0||||||||||||||||||(until)==========o=========(since)|||||||||||||||24
		// or
		// 0||||||||o|||||||||(until)====================(since)|||||||||||||||24
		// 0||||||||||||||||||(until)====================(since)|||||||o|||||||24
		return currHour < since && currHour >= until
	}

	if since < until {
		// 0===========(since)||||||||o||||||||(until)==============24
		// or
		// 0======o====(since)|||||||||||||||||(until)==============24
		// 0===========(since)|||||||||||||||||(until)=======o======24
		return currHour >= since && currHour < until
	}

	//if since == until
	return false
}
