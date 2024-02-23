package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	betterebaywatchlist "github.com/lsymds/symonds-assistant/internal/better_ebay_watchlist"
	"github.com/lsymds/symonds-assistant/internal/core/discord"
	"github.com/lsymds/symonds-assistant/internal/core/sqlite"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Panicf("godotenv: %v", err.Error())
	}

	db, err := sqlite.NewDatabase(os.Getenv("SYMONDSASSISTANT_DB_PATH"))
	if err != nil {
		log.Panicf("sqlite: %v", err.Error())
	}

	dcord, err := discord.NewDiscord(
		fmt.Sprintf("Bot %s", os.Getenv("SYMONDSASSISTANT_DISCORD_TOKEN")),
		os.Getenv("SYMONDSASSISTANT_DISCORD_GUILDID"),
		os.Getenv("SYMONDSASSISTANT_DISCORD_APPLICATIONID"),
	)
	if err != nil {
		log.Panicf("discord: %v", err.Error())
	}

	// register any packages
	betterebaywatchlist.Register("", db, dcord)

	// start handling Discord interactions
	if err = dcord.Listen(); err != nil {
		log.Panicf("discord: %v", err.Error())
	}

	// block
	done := make(chan bool)
	<-done
}
