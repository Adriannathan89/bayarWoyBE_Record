package handler

import (
	"bayar-woy-project/bot-model"
	"bayar-woy-project/bot-service"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func MessageHandler(m *discordgo.MessageCreate) {
	var bot botmodel.DiscordBotSession

	if m.Author.Bot {
		return
	}

	args := strings.Split(m.Content, " ")
	command := args[0]

	switch command {
		case "!register":
			userId := m.Author.ID
			res := botservice.RegisterUserToBot(args[1], userId)

			bot.Sesion.ChannelMessageSend(m.ChannelID, res)
	}	

}