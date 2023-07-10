package userProperties

import (
	"log"
	"strings"

	botDB "GitlabTgBot/db"
	extensions "GitlabTgBot/extensions"
)

const commandForDB string = "UPDATE TelebotDB SET Token = \"%s\" WHERE chatID = %s"

// returns result and new token to update user map
func UpdateToken(userCommand string, userTgID int64) (string, string, bool) {
	splittedCommand := strings.Split(userCommand, "=")
	if extensions.IsTokenCorrect(splittedCommand[len(splittedCommand)-1]) {
		err := UpdateTokenInDB(splittedCommand[len(splittedCommand)-1], userTgID)
		if err != nil {
			log.Printf("Error while token updating: %s", err)
			return err.Error(), "", false
		}

		return "✅Token successfully updated!", splittedCommand[len(splittedCommand)-1], true
	}

	return "❌Wrong token!", "", false
}

func UpdateTokenInDB(newToken string, userTgID int64) error {
	cryptoToken, err := extensions.Encrypt([]byte(newToken))
	if err != nil {
		return err
	}

	botDB.UpdateValInDB(userTgID, botDB.TokenColumn, string(cryptoToken), botDB.UsersTableName)
	return nil
}
