package bot

import (
	"bayar-woy-project/handler"

	"github.com/bwmarrin/discordgo"
)

var BotSession *discordgo.Session

func InitDiscordBot(token string) error  {
	BotSession, err := discordgo.New("Bot " + token)
	
	if err != nil {
		return err
	}
	
	handler.DiscordHandlerRegister(BotSession)

	return BotSession.Open()
}