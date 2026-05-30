package discord

import (
	"bayar-woy-project/config"
	"bayar-woy-project/models"
	"log"
	"time"
)

// lastReportSent tracks when the weekly report job was most recently run.
// In-memory: trade-off accepted for MVP. Restart in trigger window may cause duplicates.
var lastReportSent time.Time

// StartScheduler starts a single background goroutine that:
//   - Fires the weekly report job every Saturday between 20:00 and 20:04 WIB (with 1-hour reentry guard)
//   - Runs cleanup of expired DiscordVerification rows daily at 01:00 WIB
//
// Called once from main.go after Init.
func StartScheduler() {
	go schedulerLoop()
}

func schedulerLoop() {
	defer recoverPanic("scheduler")

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Printf("[discord:scheduler] failed to load WIB timezone: %v", err)
		return
	}

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	log.Println("[discord:scheduler] started")

	for {
		now := time.Now().In(loc)

		// Weekly report trigger: Saturday 20:00–20:04 WIB, with 1-hour reentry guard
		if now.Weekday() == time.Saturday &&
			now.Hour() == 20 &&
			now.Minute() < 5 &&
			now.Sub(lastReportSent) > 1*time.Hour {
			runWeeklyReportJob(loc)
			lastReportSent = now
		}

		// Daily cleanup: 01:00 WIB
		if now.Hour() == 1 && now.Minute() == 0 {
			cleanupExpiredCodes()
		}

		<-ticker.C
	}
}

func runWeeklyReportJob(loc *time.Location) {
	defer recoverPanic("runWeeklyReportJob")

	var users []models.User
	if err := config.DB.Where("is_validated = ? AND discord_id IS NOT NULL AND discord_weekly_notif_enabled = ?",
		true, true).Find(&users).Error; err != nil {
		log.Printf("[discord:weekly-job] query failed: %v", err)
		return
	}

	log.Printf("[discord:weekly-job] sending to %d users", len(users))

	now := time.Now()
	for _, user := range users {
		report := BuildWeeklyReport(user.ID, now, loc)
		NotifyWeeklyReport(user, report)
		time.Sleep(500 * time.Millisecond)
	}
}

func cleanupExpiredCodes() {
	defer recoverPanic("cleanupExpiredCodes")

	cutoff := time.Now().Add(-1 * time.Hour)
	res := config.DB.Where("expires_at < ?", cutoff).Delete(&models.DiscordVerification{})
	if res.Error != nil {
		log.Printf("[discord:cleanup] failed: %v", res.Error)
		return
	}
	if res.RowsAffected > 0 {
		log.Printf("[discord:cleanup] deleted %d expired verification rows", res.RowsAffected)
	}
}
