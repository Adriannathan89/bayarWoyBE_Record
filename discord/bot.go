package discord

import (
	"bayar-woy-project/config"
	"errors"
	"log"

	"github.com/bwmarrin/discordgo"
)

// Init opens the discordgo Session, registers the /verify slash command,
// and attaches the InteractionCreate handler.
// Called once from main.go at startup.
func Init() error {
	token := config.GetEnv("DISCORD_BOT_TOKEN")
	if token == "" {
		return errors.New("DISCORD_BOT_TOKEN is empty")
	}

	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		return err
	}

	sess.Identify.Intents = discordgo.IntentsDirectMessages
	sess.AddHandler(handleInteraction)

	if err := sess.Open(); err != nil {
		return err
	}

	Session = sess

	if err := registerSlashCommands(sess); err != nil {
		log.Printf("[discord] failed to register slash commands: %v", err)
		// Do not fail Init — bot still works for DMs, just commands missing
	}

	log.Println("[discord] bot ready")
	return nil
}

// Shutdown closes the discordgo session cleanly.
// Called via defer in main.go.
func Shutdown() {
	if Session != nil {
		if err := Session.Close(); err != nil {
			log.Printf("[discord] error closing session: %v", err)
		}
	}
}

// registerSlashCommands registers /verify as a global, user-installable application command.
// Idempotent: Discord deduplicates by command name on the same application.
func registerSlashCommands(sess *discordgo.Session) error {
	cmd := &discordgo.ApplicationCommand{
		Name:        "verify",
		Description: "Hubungkan akun Discord kamu dengan BayarWoy",
		IntegrationTypes: &[]discordgo.ApplicationIntegrationType{
			discordgo.ApplicationIntegrationUserInstall,
			discordgo.ApplicationIntegrationGuildInstall,
		},
		Contexts: &[]discordgo.InteractionContextType{
			discordgo.InteractionContextBotDM,
			discordgo.InteractionContextPrivateChannel,
			discordgo.InteractionContextGuild,
		},
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "code",
				Description: "6-digit code dari halaman Profile BayarWoy",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
				MinLength:   intPtr(6),
				MaxLength:   6,
			},
		},
	}

	_, err := sess.ApplicationCommandCreate(sess.State.User.ID, "", cmd)
	return err
}

// handleInteraction dispatches incoming interactions to the right handler.
// Currently only /verify is supported.
func handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	defer recoverPanic("handleInteraction")

	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	switch i.ApplicationCommandData().Name {
	case "verify":
		handleVerifyCommand(s, i)
	}
}

func intPtr(v int) *int {
	return &v
}
