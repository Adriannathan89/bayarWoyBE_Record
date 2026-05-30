package discord

import (
	"errors"
	"log"

	"github.com/bwmarrin/discordgo"
)

// Session holds the active discordgo session. Set by Init, read by other handlers.
// Public so notifier and verify_handler can access it.
var Session *discordgo.Session

// SendDM opens a DM channel with the given Discord user ID and sends plain text content.
// Best-effort: returns error so callers can log, but typical use is fire-and-forget.
func SendDM(discordID string, content string) error {
	if Session == nil {
		return errors.New("discord session not initialized")
	}

	channel, err := Session.UserChannelCreate(discordID)
	if err != nil {
		return err
	}

	_, err = Session.ChannelMessageSend(channel.ID, content)
	return err
}

// SendDMEmbed opens a DM channel with the given Discord user ID and sends an embed.
// Best-effort: errors logged by caller, not returned for fire-and-forget paths.
func SendDMEmbed(discordID string, embed *discordgo.MessageEmbed) error {
	if Session == nil {
		return errors.New("discord session not initialized")
	}

	channel, err := Session.UserChannelCreate(discordID)
	if err != nil {
		return err
	}

	_, err = Session.ChannelMessageSendEmbed(channel.ID, embed)
	return err
}

// recoverPanic is used by goroutines (notifier, scheduler) to prevent panics
// from crashing the backend. Logs the error and lets the goroutine continue.
func recoverPanic(label string) {
	if r := recover(); r != nil {
		log.Printf("[discord:%s] recovered from panic: %v", label, r)
	}
}
