package userProperties

import (
	"fmt"
	"strconv"

	botTypes "GitlabTgBot/botTypes"
	botDB "GitlabTgBot/db"
)

const (
	maxLimit int = 15
	minLimit int = 1
)

func ChangeTimeLimit(newValue string, user *botTypes.User) (string, bool) {
	intVal, err := strconv.Atoi(newValue)
	if err != nil {
		return "Please enter an INTEGER value!", false
	}

	if intVal == user.ActivityTimeLimit {
		return "This value is already set", true
	}

	if intVal > maxLimit || intVal < minLimit {
		return fmt.Sprintf("Value must be greater or equal %s and less or equal %s!", strconv.Itoa(minLimit), strconv.Itoa(maxLimit)), false
	}

	user.ActivityTimeLimit = intVal
	botDB.UpdateValInDB(user.ChatID, botDB.ActivityTimeLimitColumn, strconv.Itoa(intVal), botDB.UsersTableName)
	return fmt.Sprintf("Great! Now the hooks will accumulate for %s minutes!", newValue), true
}
