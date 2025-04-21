package main

import (
	"github.com/azr4e1/discord-bot-test/bot"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error")
	}
	openWeatherToken := os.Getenv("OPEN_WEATHER_TOKEN")
	botToken := os.Getenv("BOT_TOKEN")
	bot.BotToken = botToken
	bot.OpenWeatherToken = openWeatherToken
	bot.Run()
}
