package botservice

import (
	"bayar-woy-project/config"
	"bayar-woy-project/models"
)

func formatFriendsList(friends []models.Friendship) string {
	var formattedList string
	for _, friend := range friends {
		formattedList += "- " + friend.Friend.Username + "\n"
	}
	return formattedList
}

func GetFriendsList(discordId string) string {
	var user models.User
	var friends []models.Friendship

	if err := config.DB.Where("discord_id = ?", discordId).First(&user); err != nil {
		return "User not found"
	}

	if err := config.DB.Where("user_id = ?", user.ID).Find(&friends); err != nil {
		return "Error fetching friends"
	}
	
	if len(friends) == 0 {
		return "You have no friends yet."
	}
	return "Your friends:\n" + formatFriendsList(friends)

}