package post

import (
	"strings"

	"app/internal/database"

	"gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) PublicList(f ListFilter) ([]database.Post, error) {
	var posts []database.Post
	tx := r.db.Preload("User").Where("is_published = ?", true)
	if q := strings.TrimSpace(f.Query); q != "" {
		like := "%" + strings.ToLower(q) + "%"
		tx = tx.Where("(LOWER(title) LIKE ? OR LOWER(excerpt) LIKE ? OR LOWER(content) LIKE ? OR LOWER(category) LIKE ?)", like, like, like, like)
	}
	if f.Category != "" {
		tx = tx.Where("category = ?", f.Category)
	}
	err := tx.Order("published_at DESC").Order("id DESC").Find(&posts).Error
	return posts, err
}

func (r *Repository) FindPublishedBySlug(slug string) (*database.Post, error) {
	var p database.Post
	if err := r.db.Preload("User").Where("slug = ? AND is_published = ?", slug, true).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) ListByUser(userID uint) ([]database.Post, error) {
	var posts []database.Post
	err := r.db.Preload("User").Where("user_id = ?", userID).Order("updated_at DESC").Order("id DESC").Find(&posts).Error
	return posts, err
}

func (r *Repository) ListAll() ([]database.Post, error) {
	var posts []database.Post
	err := r.db.Preload("User").Order("created_at DESC").Order("id DESC").Find(&posts).Error
	return posts, err
}

func (r *Repository) FindByIDForUser(id, userID uint) (*database.Post, error) {
	var p database.Post
	if err := r.db.Preload("User").Where("id = ? AND user_id = ?", id, userID).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) SlugExists(slug string, excludeID uint) bool {
	var n int64
	r.db.Model(&database.Post{}).Where("slug = ? AND id <> ?", slug, excludeID).Count(&n)
	return n > 0
}

func (r *Repository) Create(p *database.Post) error { return r.db.Omit("User").Create(p).Error }

func (r *Repository) Save(p *database.Post) error { return r.db.Omit("User").Save(p).Error }

func (r *Repository) Delete(p *database.Post) error { return r.db.Delete(&database.Post{}, p.ID).Error }

func (r *Repository) CountByUser(userID uint, published *bool) int64 {
	var n int64
	tx := r.db.Model(&database.Post{}).Where("user_id = ?", userID)
	if published != nil {
		tx = tx.Where("is_published = ?", *published)
	}
	tx.Count(&n)
	return n
}

func (r *Repository) Users() ([]database.User, error) {
	var users []database.User
	err := r.db.Order("id ASC").Find(&users).Error
	return users, err
}

func (r *Repository) FindUser(id uint) (*database.User, error) {
	var u database.User
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
