package bot

import (
	"bayar-woy-project/models"

	"github.com/bwmarrin/discordgo"
)

func InitDiscordBot(token string) (*models.DiscordBotSession, error) {
	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	return &models.DiscordBotSession{Sesion: bot}, nil
}

func Start(b *models.DiscordBotSession) error {
	return b.Sesion.Open()
}

func Stop(b *models.DiscordBotSession) {
	b.Sesion.Close()
}