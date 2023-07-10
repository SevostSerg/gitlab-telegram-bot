package db

import (
	config "GitlabTgBot/configuration"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"

	_ "github.com/glebarez/go-sqlite"
)

const (
	UsersTableName           = "TelebotDB"
	StatsTableName           = "Stats"
	ChatIDColumn             = "ChatID"
	UserNameColumn           = "UserName"
	TokenColumn              = "Token"
	UserRoleColumn           = "UserRole"
	GitlabUsernameColumn     = "GitlabUsername"
	PushOptionColumn         = "PushOption"
	MergeRequestOptionColumn = "MergeRequestOption"
	CommentsOptionColumn     = "CommentsOption"
	PipelineOptionColumn     = "PipelineOption"
	TagOptionColumn          = "TagOption"
	ActivityTimeLimitColumn  = "ActivityTimeLimit"
	PushColumnName           = "PushActions"
	MRColumnName             = "MRActions"
	TagColumnName            = "ColumnName"
)

const (
	usersDBName             string = "GitlabBotUsers.db"
	sqlDriver               string = "sqlite"
	updateDBCommand         string = "UPDATE \"%s\" SET %s = %s WHERE chatID = %s"
	createUsersTableCommand string = "CREATE TABLE \"%s\" (\"%s\" INTEGER, \"%s\" TEXT," +
		"\"%s\" TEXT, \"%s\" TEXT, \"%s\" TEXT, " +
		"\"%s\" INTEGER, \"%s\" INTEGER, \"%s\" INTEGER, \"%s\" INTEGER, " +
		"\"%s\" INTEGER, \"%s\" INTEGER)"
	createStatsTableCommand string = "CREATE TABLE \"%s\" (\"%s\" TEXT, \"%s\" INTEGER, \"%s\" INTEGER, \"%s\" INTEGER)"
)

var botDB *sql.DB

func StartDB() *sql.DB {
	dbFolderName := config.GetConfigInstance().DBFolderName
	botDB = CheckDB(path.Join(dbFolderName, usersDBName))
	return botDB
}

func CheckDB(path string) *sql.DB {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		db, err := sql.Open(sqlDriver, path)
		if err != nil {
			log.Fatal(err)
		}

		return db
	}

	dbFolder := config.GetConfigInstance().DBFolderName
	if _, err := os.Stat(dbFolder); os.IsNotExist(err) {
		err := os.Mkdir(dbFolder, 0755)
		if err != nil {
			log.Panic(err)
		}
	}

	os.Create(path)
	db, err := sql.Open(sqlDriver, path)
	if err != nil {
		log.Fatal(err)
	}

	err = CreateTables(db)
	if err != nil {
		log.Panic(err)
	}

	return db
}

func CreateTables(db *sql.DB) error {
	_, err := db.Exec(fmt.Sprintf(
		createUsersTableCommand,
		UsersTableName,
		ChatIDColumn,
		UserNameColumn,
		TokenColumn,
		UserRoleColumn,
		GitlabUsernameColumn,
		PushOptionColumn,
		MergeRequestOptionColumn,
		CommentsOptionColumn,
		PipelineOptionColumn,
		TagOptionColumn,
		ActivityTimeLimitColumn))
	if err != nil {
		return err
	}

	_, err = db.Exec(fmt.Sprintf(createStatsTableCommand, GitlabUsernameColumn, StatsTableName, PushColumnName, MRColumnName, TagColumnName))
	if err != nil {
		return err
	}

	return nil
}

func UpdateValInDB(userID int64, parameter string, newValue string, tableName string) {
	result, err := botDB.Exec(fmt.Sprintf(updateDBCommand, tableName, parameter, newValue, strconv.FormatInt(userID, 10)))
	if err != nil {
		log.Print("DB error:" + err.Error())
	}

	log.Print(result)
}

func GetBotDB() *sql.DB {
	return botDB
}
