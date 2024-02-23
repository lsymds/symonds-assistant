package betterebaywatchlist

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/lsymds/symonds-assistant/internal/core/discord"
	"github.com/lsymds/symonds-assistant/internal/core/sqlite"
	"github.com/lsymds/symonds-assistant/internal/core/utils"
)

// Register registers all neccessary functionality for the Better eBay Watchlist package to work.
func Register(channelID string, db *sqlite.Database, d *discord.Discord) error {
	b := wrapper{
		db:      db,
		discord: d,
	}

	// register all Discord slash commands
	err := b.discord.RegisterSlashCommand(
		&discord.SlashCommand{
			Name:        "bew-add-watchlist-item",
			Description: "Adds an eBay auction to the watchlist.",
			Handler:     b.addNewAuctionEndingNotification,
		},
	)
	if err != nil {
		log.Printf("package registration error: %v", err)
	}

	// register any background jobs
	utils.RunInBackground("bew:send-due-notifications", b.sendDueNotifications, 30*time.Second)
	utils.RunInBackground("bew:update-listing-details", b.updateListingDetails, 30*time.Minute)

	// log that the package has been registered
	log.Println("package registered: bew - Better eBay Watchlist")

	return nil
}

type wrapper struct {
	db      *sqlite.Database
	discord *discord.Discord
}

func (w *wrapper) addNewAuctionEndingNotification(d *discord.Discord, i *discordgo.InteractionCreate) error {
	d.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "It worked!",
		},
	})

	return nil
}

func (w *wrapper) sendDueNotifications() error {
	return nil
}

func (w *wrapper) updateListingDetails() error {
	return nil
}
