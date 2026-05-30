package discord

import (
	"bayar-woy-project/config"
	"bayar-woy-project/models"
	"errors"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	ErrCodeNotFound         = errors.New("verification code not found")
	ErrCodeExpired          = errors.New("verification code expired")
	ErrDiscordAlreadyLinked = errors.New("discord account already linked to another BayarWoy user")
	ErrUserAlreadyLinked    = errors.New("user already linked to a different Discord account")
)

// ProcessVerifyCode performs the core verification logic.
// Separated from Discord interaction handling for testability.
//
// On success: updates User row, deletes verification row, returns confirmation message.
// On error: returns one of the Err* sentinel errors.
func ProcessVerifyCode(code, discordID, discordUsername string) (string, error) {
	var verification models.DiscordVerification
	if err := config.DB.Where("code = ?", code).First(&verification).Error; err != nil {
		return "", ErrCodeNotFound
	}

	if time.Now().After(verification.ExpiresAt) {
		return "", ErrCodeExpired
	}

	var existingByDiscord models.User
	res := config.DB.Where("discord_id = ?", discordID).First(&existingByDiscord)
	if res.Error == nil && existingByDiscord.ID != verification.UserID {
		return "", ErrDiscordAlreadyLinked
	}

	var targetUser models.User
	if err := config.DB.First(&targetUser, "id = ?", verification.UserID).Error; err != nil {
		return "", ErrCodeNotFound
	}

	if targetUser.DiscordID != nil && *targetUser.DiscordID != discordID {
		return "", ErrUserAlreadyLinked
	}

	tx := config.DB.Begin()
	if err := tx.Model(&targetUser).Updates(map[string]interface{}{
		"discord_id":       discordID,
		"discord_username": discordUsername,
		"is_validated":     true,
	}).Error; err != nil {
		tx.Rollback()
		return "", err
	}
	if err := tx.Delete(&verification).Error; err != nil {
		tx.Rollback()
		return "", err
	}
	if err := tx.Commit().Error; err != nil {
		return "", err
	}

	return "✅ Berhasil terhubung dengan akun BayarWoy `" + targetUser.Username + "`. Kamu akan terima notifikasi transaksi.", nil
}

// handleVerifyCommand is called from bot.go when /verify slash command arrives.
// Extracts the code option, identifies the Discord user, processes verification,
// and replies in Discord.
func handleVerifyCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	defer recoverPanic("handleVerifyCommand")

	opts := i.ApplicationCommandData().Options
	if len(opts) == 0 {
		respondInteraction(s, i, "❌ Code tidak ditemukan dalam command.")
		return
	}
	code := opts[0].StringValue()

	var discordID, discordUsername string
	if i.Member != nil && i.Member.User != nil {
		discordID = i.Member.User.ID
		discordUsername = i.Member.User.Username
	} else if i.User != nil {
		discordID = i.User.ID
		discordUsername = i.User.Username
	} else {
		respondInteraction(s, i, "❌ Tidak bisa mengidentifikasi akun Discord kamu.")
		return
	}

	msg, err := ProcessVerifyCode(code, discordID, discordUsername)
	if err != nil {
		switch err {
		case ErrCodeNotFound:
			respondInteraction(s, i, "❌ Code tidak ditemukan atau sudah expired.")
		case ErrCodeExpired:
			respondInteraction(s, i, "⏱️ Code expired. Generate code baru di BayarWoy.")
		case ErrDiscordAlreadyLinked:
			respondInteraction(s, i, "⚠️ Akun Discord ini sudah terhubung dengan akun BayarWoy lain. Disconnect dari sana dulu.")
		case ErrUserAlreadyLinked:
			respondInteraction(s, i, "ℹ️ Akun BayarWoy ini sudah terhubung dengan Discord lain. Disconnect dulu sebelum link ulang.")
		default:
			log.Printf("[discord:verify] unexpected error: %v", err)
			respondInteraction(s, i, "❌ Terjadi kesalahan. Coba lagi nanti.")
		}
		return
	}

	respondInteraction(s, i, msg)
}

func respondInteraction(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Printf("[discord:verify] failed to respond: %v", err)
	}
}
