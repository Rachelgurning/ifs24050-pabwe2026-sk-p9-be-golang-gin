package auth

import (
	"app/internal/database"

	"gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) FindByLogin(identifier string) (*database.User, error) {
	var u database.User
	if err := r.db.Where("LOWER(username) = LOWER(?) OR LOWER(email) = LOWER(?)", identifier, identifier).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) FindByID(id uint) (*database.User, error) {
	var u database.User
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) UsernameTaken(username string, excludeID uint) bool {
	var n int64
	r.db.Model(&database.User{}).Where("LOWER(username) = LOWER(?) AND id <> ?", username, excludeID).Count(&n)
	return n > 0
}

func (r *Repository) EmailTaken(email string, excludeID uint) bool {
	var n int64
	r.db.Model(&database.User{}).Where("LOWER(email) = LOWER(?) AND id <> ?", email, excludeID).Count(&n)
	return n > 0
}

func (r *Repository) Create(u *database.User) error { return r.db.Create(u).Error }

func (r *Repository) Save(u *database.User) error { return r.db.Save(u).Error }
