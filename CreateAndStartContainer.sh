#!/bin/bash
echo "Enter port like in the config"
read port
docker stop GitlabTgBot
docker rm GitlabTgBot
docker build -t gitlab-tg-bot .
docker volume create GitlabTgBotLogs
docker volume create GitlabTgBotConfig
docker volume create GitlabTgBotDB
docker create --name=GitlabTgBot --restart=always -p $port:$port -v GitlabTgBotLogs:/app/Logs -v GitlabTgBotConfig:/app/ConfigurationData -v GitlabTgBotDB:/app/BotDB gitlab-tg-bot
docker cp config.json GitlabTgBot:/app/ConfigurationData
docker start GitlabTgBot