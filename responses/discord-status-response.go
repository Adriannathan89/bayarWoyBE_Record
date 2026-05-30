package responses

type DiscordStatusResponse struct {
	Verified        bool   `json:"verified"`
	DiscordUsername string `json:"discordUsername"`
}
