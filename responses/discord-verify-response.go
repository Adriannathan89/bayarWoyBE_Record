package responses

import "time"

type DiscordVerifyResponse struct {
	Code        string    `json:"code"`
	BotUsername string    `json:"botUsername"`
	ExpiresAt   time.Time `json:"expiresAt"`
	InstallURL  string    `json:"installUrl"`
}
