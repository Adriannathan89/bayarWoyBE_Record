package loader

import (
	"bayar-woy-project/models"

	"gorm.io/gorm"
)

func SeedCategories(db *gorm.DB) {
	categories := []models.Category{
		// Primary categories
		{Name: "makanan", Type: "primary"},
		{Name: "minuman", Type: "primary"},
		{Name: "transport", Type: "primary"},
		{Name: "belanja", Type: "primary"},
		{Name: "hiburan", Type: "primary"},
		{Name: "tagihan", Type: "primary"},
		{Name: "kesehatan", Type: "primary"},
		{Name: "gaji", Type: "primary"},
		{Name: "hadiah", Type: "primary"},

		// Secondary categories
		{Name: "jajanan", Type: "secondary"},
		{Name: "makanan", Type: "secondary"},
		{Name: "minuman", Type: "secondary"},
		{Name: "elektronik", Type: "secondary"},
		{Name: "fashion", Type: "secondary"},
		{Name: "harian", Type: "secondary"},
		{Name: "kecantikan", Type: "secondary"},
		{Name: "online_shop", Type: "secondary"},
		{Name: "bonus", Type: "secondary"},
		{Name: "gaji", Type: "secondary"},
		{Name: "komisi", Type: "secondary"},
		{Name: "kado", Type: "secondary"},
		{Name: "hadiah", Type: "secondary"},
		{Name: "reward", Type: "secondary"},
		{Name: "sumbangan", Type: "secondary"},
		{Name: "undian", Type: "secondary"},
		{Name: "pemberian_masuk", Type: "secondary"},
		{Name: "aktivitas", Type: "secondary"},
		{Name: "game", Type: "secondary"},
		{Name: "hiburan", Type: "secondary"},
		{Name: "liburan", Type: "secondary"},
		{Name: "streaming", Type: "secondary"},
		{Name: "tontonan", Type: "secondary"},
		{Name: "apotek", Type: "secondary"},
		{Name: "kesehatan", Type: "secondary"},
		{Name: "konsul", Type: "secondary"},
		{Name: "obat", Type: "secondary"},
		{Name: "prosedur", Type: "secondary"},
		{Name: "asuransi", Type: "secondary"},
		{Name: "gaji_pihak3", Type: "secondary"},
		{Name: "internet_pulsa", Type: "secondary"},
		{Name: "iuran", Type: "secondary"},
		{Name: "kredit", Type: "secondary"},
		{Name: "sewa", Type: "secondary"},
		{Name: "tagihan", Type: "secondary"},
		{Name: "utilitas", Type: "secondary"},
		{Name: "online", Type: "secondary"},
		{Name: "pesawat", Type: "secondary"},
		{Name: "pribadi", Type: "secondary"},
		{Name: "transport", Type: "secondary"},
		{Name: "umum", Type: "secondary"},
		{Name: "belanja", Type: "secondary"},
	}

	for _, cat := range categories {
		db.Where(models.Category{Name: cat.Name, Type: cat.Type}).
			FirstOrCreate(&cat)
	}
}
