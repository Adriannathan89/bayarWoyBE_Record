package handler

import (
	"bayar-woy-project/models"
)

func SendDM(s *models.DiscordBotSession, userID string, message string) error {
	channel, err := s.Sesion.UserChannelCreate(userID)
	if err != nil {
		return err
	}

	_, err = s.Sesion.ChannelMessageSend(channel.ID, message)
	return err
}