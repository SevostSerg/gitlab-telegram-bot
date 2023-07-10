package botTypes

import (
	botDB "GitlabTgBot/db"
	"fmt"
	"log"

	"GitlabTgBot/extensions"
)

const (
	DefaultUserRole          = "Employe"
	AdminUserRole            = "Admin"
	MasterUserRole           = "Master"
	DefaultActivityTimeLimit = 1
)

type User struct {
	ChatID             int64
	UserName           string
	Token              string
	UserRole           string
	GitlabUsername     string
	PushOption         int
	MergeRequestOption int
	CommentsOption     int
	PipelineOption     int
	TagOption          int
	ActivityTimeLimit  int
}

// Key: tg bot chat ID
var usersChatIDKey map[int64]*User

// Key: gitlab username
var usersGitlabUsernameKey map[string]*User

func GetUserListFromDB() (map[int64]*User, map[string]*User) {
	userListChatIDKey := make(map[int64]*User)
	userListGitlabUserNameKey := make(map[string]*User)

	result, err := botDB.GetBotDB().Query(fmt.Sprintf("select * from %s", botDB.UsersTableName))
	if err != nil {
		log.Panic(err)
	}

	for result.Next() {
		u := User{}
		Token := []byte{}
		err := result.Scan(&u.ChatID, &u.UserName, &Token, &u.UserRole, &u.GitlabUsername,
			&u.PushOption, &u.MergeRequestOption, &u.CommentsOption, &u.PipelineOption, &u.TagOption, &u.ActivityTimeLimit)
		if err != nil {
			log.Panic(err)
		}

		u.Token = extensions.Decrypt(Token)
		userListChatIDKey[u.ChatID] = &u
		userListGitlabUserNameKey[u.GitlabUsername] = &u
	}

	usersChatIDKey = userListChatIDKey
	usersGitlabUsernameKey = userListGitlabUserNameKey
	log.Print("Userlist loaded successfully")
	return userListChatIDKey, userListGitlabUserNameKey
}

func GetUserWithChatID(chatID int64) (*User, bool) {
	if user, contains := usersChatIDKey[chatID]; contains {
		return user, contains
	}

	return nil, false
}

func GetUserWithGitlabName(gitlabName string) (*User, bool) {
	if user, contains := usersGitlabUsernameKey[gitlabName]; contains {
		return user, contains
	}

	return nil, false
}
