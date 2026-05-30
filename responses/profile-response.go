package responses

type DiscordProfileInfo struct {
	Connected             bool   `json:"connected"`
	Username              string `json:"username"`
	CommitNotifEnabled    bool   `json:"commitNotifEnabled"`
	WeeklyNotifEnabled    bool   `json:"weeklyNotifEnabled"`
}

type ProfileResponse struct {
	ID       string             `json:"id"`
	Username string             `json:"username"`
	Discord  DiscordProfileInfo `json:"discord"`
}
