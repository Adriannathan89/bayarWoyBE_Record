package discord

import (
	"bayar-woy-project/config"
	"bayar-woy-project/models"
	"sort"
	"time"
)

// WeeklyReport contains aggregated financial data for a single user
// over a single week window (Sunday 00:00 WIB to Saturday 19:59:59 WIB).
type WeeklyReport struct {
	WindowStart          time.Time
	WindowEnd            time.Time
	TotalIncome          float32
	TotalExpense         float32
	TopExpenseCategories []CategoryTotal
}

// CategoryTotal is a single primary category with the sum of its expenses
// during the report window.
type CategoryTotal struct {
	Name  string
	Total float32
}

// BuildWeeklyReport computes the user's report for the week containing `now`.
// Window: Sunday 00:00 WIB → Saturday 19:59:59 WIB (so report.CreatedAt < windowEnd).
// Only committed records are included.
func BuildWeeklyReport(userID string, now time.Time, loc *time.Location) WeeklyReport {
	nowLocal := now.In(loc)

	// Find this week's Saturday at 20:00 (the report time)
	// Then window is from previous Sunday 00:00 up to (not including) Saturday 20:00.
	daysSinceSunday := int(nowLocal.Weekday()) // Sunday=0, Saturday=6
	sunday := time.Date(nowLocal.Year(), nowLocal.Month(), nowLocal.Day(), 0, 0, 0, 0, loc).AddDate(0, 0, -daysSinceSunday)
	saturday2000 := sunday.AddDate(0, 0, 6).Add(20 * time.Hour)

	windowStart := sunday
	windowEnd := saturday2000

	var records []models.Record
	config.DB.Preload("Categories").
		Where("owner_id = ? AND is_committed = ? AND created_at >= ? AND created_at < ?",
			userID, true, windowStart, windowEnd).
		Find(&records)

	var totalIncome, totalExpense float32
	expenseByCategory := make(map[string]float32)

	for _, r := range records {
		if r.Type == "income" {
			totalIncome += r.Amount
		} else if r.Type == "expense" {
			totalExpense += r.Amount
			for _, cat := range r.Categories {
				if cat.Type == "primary" {
					expenseByCategory[cat.Name] += r.Amount
					break
				}
			}
		}
	}

	var topCats []CategoryTotal
	for name, total := range expenseByCategory {
		topCats = append(topCats, CategoryTotal{Name: name, Total: total})
	}
	sort.Slice(topCats, func(i, j int) bool {
		return topCats[i].Total > topCats[j].Total
	})
	if len(topCats) > 3 {
		topCats = topCats[:3]
	}

	return WeeklyReport{
		WindowStart:          windowStart,
		WindowEnd:            windowEnd.Add(-1 * time.Second),
		TotalIncome:          totalIncome,
		TotalExpense:         totalExpense,
		TopExpenseCategories: topCats,
	}
}
