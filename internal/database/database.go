package database

import (
	"fmt"
	"time"

	"app/internal/config"

	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100" json:"name"`
	Username  string    `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email     string    `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	Role      string    `gorm:"size:20;default:member" json:"role"`
	Bio       string    `gorm:"type:text" json:"bio"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Post struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	UserID         uint       `gorm:"index;not null" json:"user_id"`
	User           User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Title          string     `gorm:"size:255;not null" json:"title"`
	Slug           string     `gorm:"size:255;uniqueIndex;not null" json:"slug"`
	Category       string     `gorm:"size:100" json:"category"`
	Excerpt        string     `gorm:"type:text" json:"excerpt"`
	Content        string     `gorm:"type:text" json:"content"`
	CoverImagePath string     `gorm:"type:text" json:"cover_image_path"`
	IsPublished    bool       `gorm:"default:false" json:"is_published"`
	PublishedAt    *time.Time `json:"published_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func Connect(cfg *config.Config) (*gorm.DB, error) {
	gcfg := &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)}
	if cfg.DatabaseURL != "" {
		return gorm.Open(postgres.Open(cfg.DatabaseURL), gcfg)
	}
	if cfg.DBConnection == "sqlite" {
		return gorm.Open(sqlite.Open(cfg.DBDatabase), gcfg)
	}
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBPort, cfg.DBUsername, cfg.DBPassword, cfg.DBDatabase, cfg.DBSSLMode)
	return gorm.Open(postgres.Open(dsn), gcfg)
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Post{})
}

func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func Seed(db *gorm.DB, cfg *config.Config) error {
	var admin User
	if err := db.Where("username = ?", cfg.AdminUsername).First(&admin).Error; err != nil {
		hash, herr := HashPassword(cfg.AdminPassword)
		if herr != nil {
			return herr
		}
		admin = User{Name: "Administrator", Username: cfg.AdminUsername, Email: cfg.AdminUsername + "@delcom.org", Password: hash, Role: "admin"}
		if err := db.Create(&admin).Error; err != nil {
			return err
		}
	} else if admin.Role != "admin" {
		db.Model(&admin).Update("role", "admin")
	}
	if !cfg.SeedSample {
		return nil
	}
	var count int64
	db.Model(&Post{}).Count(&count)
	if count > 0 {
		return nil
	}
	now := time.Now()
	for i, s := range samplePosts {
		pub := now.Add(-time.Duration(i+1) * time.Hour)
		p := Post{UserID: admin.ID, Title: s[0], Slug: s[1], Category: s[2], Excerpt: s[3], Content: s[4], IsPublished: true, PublishedAt: &pub}
		if err := db.Omit("User").Create(&p).Error; err != nil {
			return err
		}
	}
	return nil
}

var samplePosts = [][5]string{
	{
		"Mengenal Clean Architecture pada Backend Golang",
		"mengenal-clean-architecture-pada-backend-golang",
		"Teknologi",
		"Clean Architecture memisahkan logika bisnis dari detail teknis sehingga backend Golang lebih mudah diuji, dirawat, dan dikembangkan.",
		`## Mengapa Arsitektur Itu Penting?

Seiring bertambahnya fitur, kode backend yang tidak terstruktur akan sulit dipahami dan diuji. **Clean Architecture** hadir untuk memisahkan logika bisnis dari detail teknis seperti framework, database, dan protokol HTTP.

## Lapisan Utama

1. **Handler** menerima request HTTP dan memvalidasi input.
2. **Service** menampung aturan bisnis, misalnya membuat slug unik atau menghitung waktu baca.
3. **Repository** bertanggung jawab atas query ke database melalui GORM.

## Pendekatan Feature-First

Alih-alih mengelompokkan file berdasarkan jenis (semua handler di satu folder), pendekatan *feature-first* mengelompokkan kode per fitur seperti auth, post, dan ai. Setiap fitur memiliki dto, repository, service, dan handler sendiri sehingga perubahan pada satu fitur tidak merembet ke fitur lain.

## Kesimpulan

Dengan memisahkan tanggung jawab, setiap lapisan dapat diuji secara terisolasi dan target **100% code coverage** menjadi jauh lebih realistis.`,
	},
	{
		"Strategi Belajar Efektif untuk Mahasiswa Informatika",
		"strategi-belajar-efektif-untuk-mahasiswa-informatika",
		"Pendidikan",
		"Kombinasi active recall, proyek nyata, dan jadwal belajar konsisten membantu mahasiswa informatika menguasai materi lebih cepat.",
		`## Belajar dengan Membangun

Cara terbaik memahami pemrograman adalah membangun proyek nyata. Mulailah dari aplikasi kecil seperti todolist, lalu kembangkan menjadi aplikasi blog lengkap dengan autentikasi dan REST API.

## Teknik yang Terbukti Efektif

- **Active recall**: tutup catatan, lalu coba jelaskan kembali konsep dengan kata-kata sendiri.
- **Spaced repetition**: ulangi materi dalam interval yang semakin panjang.
- **Pomodoro**: fokus 25 menit, istirahat 5 menit, agar konsentrasi tetap terjaga.

## Manfaatkan Komunitas

Diskusi dengan teman sekelas dan membaca kode orang lain mempercepat pemahaman. Biasakan pula menulis dokumentasi agar ilmu yang dipelajari tidak mudah hilang.

> Konsistensi kecil setiap hari mengalahkan belajar maraton semalam sebelum ujian.`,
	},
	{
		"Kebiasaan Sehat bagi Pekerja Remote dan Mahasiswa",
		"kebiasaan-sehat-bagi-pekerja-remote-dan-mahasiswa",
		"Gaya Hidup",
		"Duduk berjam-jam di depan laptop berisiko bagi kesehatan. Berikut kebiasaan sederhana agar tetap bugar dan produktif.",
		`## Atur Ruang Kerja yang Ergonomis

Posisikan layar sejajar mata, gunakan kursi yang menopang punggung, dan pastikan pencahayaan cukup agar mata tidak cepat lelah.

## Gerak Setiap Jam

Bangun dan regangkan tubuh minimal lima menit setiap jam. Aktivitas ringan seperti berjalan atau peregangan leher membantu sirkulasi darah.

## Jaga Pola Tidur dan Hidrasi

- Tidur tujuh hingga delapan jam setiap malam.
- Minum air putih secara berkala.
- Batasi konsumsi kafein pada sore hari.

## Kesimpulan

Kesehatan adalah fondasi produktivitas. Kebiasaan kecil yang dilakukan secara konsisten akan memberikan dampak besar dalam jangka panjang.`,
	},
	{
		"Peluang Bisnis Digital untuk UMKM di Era Kecerdasan Buatan",
		"peluang-bisnis-digital-untuk-umkm-di-era-kecerdasan-buatan",
		"Bisnis",
		"Kecerdasan buatan membuka peluang baru bagi UMKM untuk memasarkan produk, melayani pelanggan, dan mengambil keputusan berbasis data.",
		`## AI Bukan Hanya untuk Perusahaan Besar

Kini banyak layanan kecerdasan buatan yang terjangkau sehingga UMKM dapat memanfaatkannya untuk membuat konten pemasaran, menjawab pertanyaan pelanggan, hingga menganalisis penjualan.

## Contoh Pemanfaatan

1. **Konten otomatis**: membuat deskripsi produk dan artikel blog dengan bantuan AI.
2. **Chatbot layanan pelanggan** yang aktif 24 jam.
3. **Analisis data penjualan** untuk menentukan produk terlaris dan waktu promosi terbaik.

## Tips Memulai

Mulailah dari satu masalah yang paling menyita waktu, uji coba solusi AI sederhana, lalu ukur hasilnya sebelum memperluas penggunaan.`,
	},
}
