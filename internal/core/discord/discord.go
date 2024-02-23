package discord

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// SlashCommand is a custom struct definition for all slash command based interactions.
type SlashCommand struct {
	Name        string
	Description string
	Options     []*discordgo.ApplicationCommandOption
	Handler     func(d *Discord, i *discordgo.InteractionCreate) error
}

// Discord represents a thin wrapper around a Discord client.
type Discord struct {
	Session       *discordgo.Session
	guildID       string
	applicationID string
	slashCommands map[string]*SlashCommand
}

// NewDiscord initializes the Discord client with a given token and guild identifier. It validates
// that the bot has the access it needs to function within the guild.
func NewDiscord(token string, guildID string, applicationID string) (*Discord, error) {
	s, err := discordgo.New(token)
	if err != nil {
		return nil, err
	}

	if err = s.Open(); err != nil {
		return nil, err
	}

	return &Discord{
		Session:       s,
		guildID:       guildID,
		applicationID: applicationID,
		slashCommands: make(map[string]*SlashCommand, 0),
	}, nil
}

// RegisterSlashCommand associates a command with the server, ensuring it is created and
// subsequently deleted whenever the application stops.
func (d *Discord) RegisterSlashCommand(cmd *SlashCommand) error {
	log.Printf("discord_slash_cmd(%v): registered", cmd.Name)

	// register the command definition with Discord
	createdCmd, err := d.Session.ApplicationCommandCreate(
		d.applicationID,
		d.guildID,
		&discordgo.ApplicationCommand{
			Name:        cmd.Name,
			Description: cmd.Description,
		},
	)
	if err != nil {
		return fmt.Errorf("discord_slash_cmd(%v): create command err: %w", cmd.Name, err)
	} else {
		log.Printf("discord_slash_cmd(%v): created id=%v", cmd.Name, createdCmd.ID)
	}

	d.slashCommands[cmd.Name] = cmd

	return nil
}

// Listen executes the core components of the Discord bot. It waits for interactions and
// subsequently acts on them.
func (d *Discord) Listen() error {
	// add slash command handler
	d.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := d.slashCommands[i.ApplicationCommandData().Name]; ok {
			log.Printf("discord_slash_cmd(%v): executing", h.Name)
			h.Handler(d, i)
			log.Printf("discord_slash_cmd(%v): ending", h.Name)
		}
	})

	return nil
}
