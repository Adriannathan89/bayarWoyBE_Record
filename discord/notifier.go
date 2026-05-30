package discord

import (
	"bayar-woy-project/models"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	colorIncome  = 0x10b981 // green
	colorExpense = 0xef4444 // red
	colorReport  = 0x3b82f6 // blue
)

// NotifyCommit sends a DM embed describing a newly committed transaction.
// Called as fire-and-forget goroutine from commit-record-service.
// Errors are logged, not returned.
func NotifyCommit(user models.User, record models.Record) {
	defer recoverPanic("NotifyCommit")

	if !user.IsValidated || !user.DiscordCommitNotifEnabled || user.DiscordID == nil {
		return
	}

	embed := buildCommitEmbed(record)
	if err := SendDMEmbed(*user.DiscordID, embed); err != nil {
		log.Printf("[discord:commit-notif] failed for user %s: %v", user.ID, err)
	}
}

// NotifyWeeklyReport sends a DM embed with the user's weekly financial summary.
// Called as fire-and-forget from scheduler.
func NotifyWeeklyReport(user models.User, report WeeklyReport) {
	defer recoverPanic("NotifyWeeklyReport")

	if !user.IsValidated || !user.DiscordWeeklyNotifEnabled || user.DiscordID == nil {
		return
	}

	embed := buildWeeklyReportEmbed(report)
	if err := SendDMEmbed(*user.DiscordID, embed); err != nil {
		log.Printf("[discord:weekly-notif] failed for user %s: %v", user.ID, err)
	}
}

func buildCommitEmbed(record models.Record) *discordgo.MessageEmbed {
	title := "💸 Pengeluaran"
	color := colorExpense
	if record.Type == "income" {
		title = "💰 Pemasukan"
		color = colorIncome
	}

	primaryName, secondaryName := "-", "-"
	for _, cat := range record.Categories {
		if cat.Type == "primary" {
			primaryName = cat.Name
		} else if cat.Type == "secondary" {
			secondaryName = cat.Name
		}
	}

	fields := []*discordgo.MessageEmbedField{
		{Name: "Jumlah", Value: "Rp " + formatRupiah(record.Amount), Inline: true},
		{Name: "Kategori", Value: primaryName + " / " + secondaryName, Inline: true},
	}
	if strings.TrimSpace(record.Description) != "" {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Deskripsi",
			Value:  record.Description,
			Inline: false,
		})
	}

	return &discordgo.MessageEmbed{
		Title:       title,
		Description: "**" + record.Title + "**",
		Color:       color,
		Fields:      fields,
		Footer:      &discordgo.MessageEmbedFooter{Text: "BayarWoy"},
		Timestamp:   record.CreatedAt.Format(time.RFC3339),
	}
}

func buildWeeklyReportEmbed(report WeeklyReport) *discordgo.MessageEmbed {
	dateRange := fmt.Sprintf("%s – %s",
		report.WindowStart.Format("2 Jan 2006"),
		report.WindowEnd.Format("2 Jan 2006"))

	if report.TotalIncome == 0 && report.TotalExpense == 0 {
		return &discordgo.MessageEmbed{
			Title:       "📊 Laporan Mingguan",
			Description: dateRange + "\n\nBelum ada transaksi minggu ini.",
			Color:       colorReport,
			Footer:      &discordgo.MessageEmbedFooter{Text: "BayarWoy Weekly Report"},
		}
	}

	net := report.TotalIncome - report.TotalExpense
	netStr := "+Rp " + formatRupiah(net)
	if net < 0 {
		netStr = "-Rp " + formatRupiah(-net)
	}

	fields := []*discordgo.MessageEmbedField{
		{Name: "Pemasukan", Value: "Rp " + formatRupiah(report.TotalIncome), Inline: true},
		{Name: "Pengeluaran", Value: "Rp " + formatRupiah(report.TotalExpense), Inline: true},
		{Name: "Net minggu ini", Value: netStr, Inline: false},
	}

	if len(report.TopExpenseCategories) > 0 {
		var sb strings.Builder
		for _, cat := range report.TopExpenseCategories {
			sb.WriteString(fmt.Sprintf("%s %s — Rp %s\n", categoryEmoji(cat.Name), cat.Name, formatRupiah(cat.Total)))
		}
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Top kategori pengeluaran",
			Value:  sb.String(),
			Inline: false,
		})
	}

	return &discordgo.MessageEmbed{
		Title:       "📊 Laporan Mingguan",
		Description: dateRange,
		Color:       colorReport,
		Fields:      fields,
		Footer:      &discordgo.MessageEmbedFooter{Text: "BayarWoy Weekly Report"},
	}
}

func formatRupiah(amount float32) string {
	intAmt := int(amount)
	if intAmt < 0 {
		intAmt = -intAmt
	}
	s := fmt.Sprintf("%d", intAmt)
	n := len(s)
	if n <= 3 {
		return s
	}
	var out strings.Builder
	pre := n % 3
	if pre > 0 {
		out.WriteString(s[:pre])
		if n > pre {
			out.WriteString(".")
		}
	}
	for i := pre; i < n; i += 3 {
		out.WriteString(s[i : i+3])
		if i+3 < n {
			out.WriteString(".")
		}
	}
	return out.String()
}

func categoryEmoji(name string) string {
	emojis := map[string]string{
		"makanan":   "🍜",
		"minuman":   "☕",
		"transport": "🚗",
		"belanja":   "🛍️",
		"hiburan":   "🎮",
		"tagihan":   "📋",
		"kesehatan": "💊",
		"gaji":      "💰",
		"hadiah":    "🎁",
	}
	if e, ok := emojis[name]; ok {
		return e
	}
	return "📦"
}
