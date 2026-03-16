package bot

import (
	botmodel "bayar-woy-project/bot-model"
	"bayar-woy-project/handler"

	"github.com/bwmarrin/discordgo"
)

func InitDiscordBot(token string) (*botmodel.DiscordBotSession, error) {
	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	bot.AddHandler(handler.MessageHandler)

	return &botmodel.DiscordBotSession{Sesion: bot}, nil
}

func Start(b *botmodel.DiscordBotSession) error {
	return b.Sesion.Open()
}

func Stop(b *botmodel.DiscordBotSession) {
	b.Sesion.Close()
}