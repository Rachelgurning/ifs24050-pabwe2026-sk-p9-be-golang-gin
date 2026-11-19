package post

import "time"

var Categories = []string{"Teknologi", "Pendidikan", "Gaya Hidup", "Bisnis", "Opini", "Umum"}

type PostRequest struct {
	Title          *string `json:"title"`
	Category       *string `json:"category"`
	Excerpt        *string `json:"excerpt"`
	Content        *string `json:"content"`
	CoverImagePath *string `json:"cover_image_path"`
	IsPublished    *bool   `json:"is_published"`
}

type ListFilter struct {
	Query    string
	Category string
}

type AuthorInfo struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type PostResponse struct {
	ID              uint        `json:"id"`
	UserID          uint        `json:"user_id"`
	Title           string      `json:"title"`
	Slug            string      `json:"slug"`
	Category        string      `json:"category"`
	Excerpt         string      `json:"excerpt"`
	Content         string      `json:"content"`
	CoverImagePath  string      `json:"cover_image_path"`
	IsPublished     bool        `json:"is_published"`
	Status          string      `json:"status"`
	PublishedAt     *time.Time  `json:"published_at"`
	ReadingTime     int         `json:"reading_time"`
	ReadingTimeText string      `json:"reading_time_text"`
	WordCount       int         `json:"word_count"`
	Author          *AuthorInfo `json:"author,omitempty"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

type Stats struct {
	TotalPosts     int64 `json:"total_posts"`
	PublishedPosts int64 `json:"published_posts"`
	DraftPosts     int64 `json:"draft_posts"`
	Total          int64 `json:"total"`
	Published      int64 `json:"published"`
	Draft          int64 `json:"draft"`
}

type UserSummary struct {
	ID             uint      `json:"id"`
	Name           string    `json:"name"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	Role           string    `json:"role"`
	TotalPosts     int64     `json:"total_posts"`
	PublishedPosts int64     `json:"published_posts"`
	DraftPosts     int64     `json:"draft_posts"`
	CreatedAt      time.Time `json:"created_at"`
}
