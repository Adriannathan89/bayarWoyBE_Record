package handler

import (
	botservice "bayar-woy-project/bot-service"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func DiscordHandlerRegister(s *discordgo.Session) {
	s.AddHandler(MessageHandler)
}

func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	bot := s


	if m.Author.Bot {
		return
	}

	args := strings.Split(m.Content, " ")
	command := args[0]
	userId := m.Author.ID

	switch command {
		case "!register":
			if(len(args) < 2) {
				bot.ChannelMessageSend(m.ChannelID, "Invailid command format. Use !register <username>")
				return
			}
			res := botservice.RegisterUserToBot(args[1], userId)

			bot.ChannelMessageSend(m.ChannelID, res)
			
		case "!friendList":
			res := botservice.GetFriendsList(userId)
			bot.ChannelMessageSend(m.ChannelID, res)

		case "!ping":
			bot.ChannelMessageSend(m.ChannelID, "Pong!")
	}	

}