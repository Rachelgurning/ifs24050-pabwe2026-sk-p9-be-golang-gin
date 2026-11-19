package post

import (
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"app/internal/database"
	"app/internal/response"
)

var (
	nonSlug = regexp.MustCompile(`[^a-z0-9]+`)
	mdStrip = regexp.MustCompile("[#*_`>\\[\\]()!~|]+")
)

func Slugify(s string) string {
	s = nonSlug.ReplaceAllString(strings.ToLower(strings.TrimSpace(s)), "-")
	s = strings.Trim(s, "-")
	if len(s) > 200 {
		s = strings.Trim(s[:200], "-")
	}
	if s == "" {
		s = "artikel"
	}
	return s
}

func ReadingTime(content string) (int, int) {
	words := len(strings.Fields(content))
	minutes := int(math.Ceil(float64(words) / 200.0))
	if minutes < 1 {
		minutes = 1
	}
	return minutes, words
}

func NormalizeCategory(c string) (string, bool) {
	c = strings.TrimSpace(c)
	if c == "" {
		return "Umum", true
	}
	for _, cat := range Categories {
		if strings.EqualFold(cat, c) {
			return cat, true
		}
	}
	return "", false
}

func MakeExcerpt(content string) string {
	plain := strings.Join(strings.Fields(mdStrip.ReplaceAllString(content, " ")), " ")
	if utf8.RuneCountInString(plain) <= 160 {
		return plain
	}
	r := []rune(plain)
	return strings.TrimSpace(string(r[:157])) + "..."
}

func ToResponse(p *database.Post) PostResponse {
	minutes, words := ReadingTime(p.Content)
	status := "draft"
	if p.IsPublished {
		status = "published"
	}
	res := PostResponse{
		ID: p.ID, UserID: p.UserID, Title: p.Title, Slug: p.Slug, Category: p.Category,
		Excerpt: p.Excerpt, Content: p.Content, CoverImagePath: p.CoverImagePath,
		IsPublished: p.IsPublished, Status: status, PublishedAt: p.PublishedAt,
		ReadingTime: minutes, ReadingTimeText: fmt.Sprintf("%d menit baca", minutes), WordCount: words,
		CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
	res.ReadingTimeMinutes = minutes
	if p.User.ID != 0 {
		res.Author = &AuthorInfo{ID: p.User.ID, Name: p.User.Name, Username: p.User.Username}
		res.AuthorName = p.User.Name
		res.AuthorUsername = p.User.Username
	}
	return res
}

func ToResponses(posts []database.Post) []PostResponse {
	out := make([]PostResponse, 0, len(posts))
	for i := range posts {
		out = append(out, ToResponse(&posts[i]))
	}
	return out
}

type Service struct{ repo *Repository }

func NewService(repo *Repository) *Service { return &Service{repo: repo} }

func notFound() error { return response.NewError(http.StatusNotFound, "Artikel tidak ditemukan") }

func bad(msg string) error { return response.NewError(http.StatusBadRequest, msg) }

func (s *Service) owned(userID, id uint) (*database.Post, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return nil, notFound()
	}
	if p.UserID != userID {
		return nil, response.NewError(http.StatusForbidden, "Anda tidak memiliki akses ke artikel ini")
	}
	return p, nil
}

func (s *Service) uniqueSlug(title string, excludeID uint) string {
	base := Slugify(title)
	slug := base
	for i := 2; s.repo.SlugExists(slug, excludeID); i++ {
		slug = fmt.Sprintf("%s-%d", base, i)
	}
	return slug
}

func setPublished(p *database.Post, v bool) {
	if v && !p.IsPublished {
		now := time.Now()
		p.PublishedAt = &now
	}
	if !v {
		p.PublishedAt = nil
	}
	p.IsPublished = v
}

func (s *Service) apply(p *database.Post, req PostRequest, creating bool) error {
	if req.Title != nil {
		t := strings.TrimSpace(*req.Title)
		if t == "" {
			return bad("Judul artikel wajib diisi")
		}
		if utf8.RuneCountInString(t) > 255 {
			return bad("Judul artikel maksimal 255 karakter")
		}
		p.Title = t
	}
	if req.Category != nil {
		cat, ok := NormalizeCategory(*req.Category)
		if !ok {
			return bad("Kategori tidak valid. Pilihan: " + strings.Join(Categories, ", "))
		}
		p.Category = cat
	}
	if req.Content != nil {
		c := strings.TrimSpace(*req.Content)
		if c == "" {
			return bad("Konten artikel wajib diisi")
		}
		p.Content = c
	}
	if req.Excerpt != nil {
		p.Excerpt = strings.TrimSpace(*req.Excerpt)
	}
	if req.CoverImagePath != nil {
		p.CoverImagePath = strings.TrimSpace(*req.CoverImagePath)
	}
	if creating {
		if p.Title == "" {
			return bad("Judul artikel wajib diisi")
		}
		if p.Content == "" {
			return bad("Konten artikel wajib diisi")
		}
		if p.Category == "" {
			p.Category = DefaultCategory
		}
	}
	if p.Excerpt == "" {
		p.Excerpt = MakeExcerpt(p.Content)
	}
	if req.IsPublished != nil {
		setPublished(p, *req.IsPublished)
	}
	return nil
}

func (s *Service) PublicList(f ListFilter) ([]PostResponse, error) {
	if strings.EqualFold(f.Category, "semua") || strings.EqualFold(f.Category, "all") {
		f.Category = ""
	} else if cat, ok := NormalizeCategory(f.Category); ok && f.Category != "" {
		f.Category = cat
	}
	posts, err := s.repo.PublicList(f)
	if err != nil {
		return nil, err
	}
	out := ToResponses(posts)
	for i := range out {
		out[i].UserID = 0 // jangan bocorkan data internal ke endpoint publik
	}
	return out, nil
}

func (s *Service) PublicDetail(slug string) (*PostResponse, error) {
	p, err := s.repo.FindPublishedBySlug(slug)
	if err != nil {
		return nil, notFound()
	}
	res := ToResponse(p)
	res.UserID = 0 // jangan bocorkan data internal ke endpoint publik
	return &res, nil
}

func (s *Service) ListMine(userID uint) ([]PostResponse, error) {
	posts, err := s.repo.ListByUser(userID)
	if err != nil {
		return nil, err
	}
	return ToResponses(posts), nil
}

func (s *Service) GetMine(userID, id uint) (*PostResponse, error) {
	p, err := s.owned(userID, id)
	if err != nil {
		return nil, err
	}
	res := ToResponse(p)
	return &res, nil
}

func (s *Service) Create(userID uint, req PostRequest) (*PostResponse, error) {
	p := &database.Post{UserID: userID}
	if err := s.apply(p, req, true); err != nil {
		return nil, err
	}
	p.Slug = s.uniqueSlug(p.Title, 0)
	if err := s.repo.Create(p); err != nil {
		return nil, err
	}
	return s.GetMine(userID, p.ID)
}

func (s *Service) Update(userID, id uint, req PostRequest) (*PostResponse, error) {
	p, err := s.owned(userID, id)
	if err != nil {
		return nil, err
	}
	oldTitle := p.Title
	if err := s.apply(p, req, false); err != nil {
		return nil, err
	}
	if p.Title != oldTitle {
		p.Slug = s.uniqueSlug(p.Title, p.ID)
	}
	if err := s.repo.Save(p); err != nil {
		return nil, err
	}
	res := ToResponse(p)
	return &res, nil
}

func (s *Service) TogglePublish(userID, id uint, value *bool) (*PostResponse, error) {
	p, err := s.owned(userID, id)
	if err != nil {
		return nil, err
	}
	next := !p.IsPublished
	if value != nil {
		next = *value
	}
	setPublished(p, next)
	if err := s.repo.Save(p); err != nil {
		return nil, err
	}
	res := ToResponse(p)
	return &res, nil
}

func (s *Service) Delete(userID, id uint) error {
	p, err := s.owned(userID, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(p)
}

func (s *Service) Stats(userID uint) Stats {
	t := true
	return NewStats(s.repo.CountByUser(userID, nil), s.repo.CountByUser(userID, &t))
}

func (s *Service) summary(u *database.User) UserSummary {
	t := true
	total := s.repo.CountByUser(u.ID, nil)
	pub := s.repo.CountByUser(u.ID, &t)
	return UserSummary{
		ID: u.ID, Name: u.Name, Username: u.Username, Email: u.Email, Role: u.Role,
		TotalArticles: total, TotalSelesai: pub, BelumSelesai: total - pub,
		TotalPosts: total, PublishedPosts: pub, TotalPublished: pub, DraftPosts: total - pub,
		CreatedAt: u.CreatedAt,
	}
}

func (s *Service) AdminUsers() ([]UserSummary, error) {
	users, err := s.repo.Users()
	if err != nil {
		return nil, err
	}
	out := make([]UserSummary, 0, len(users))
	for i := range users {
		out = append(out, s.summary(&users[i]))
	}
	return out, nil
}

func (s *Service) AdminPosts() ([]PostResponse, error) {
	posts, err := s.repo.ListAll()
	if err != nil {
		return nil, err
	}
	return ToResponses(posts), nil
}

func (s *Service) AdminUserPosts(id uint) ([]PostResponse, error) {
	posts, err := s.repo.ListByUser(id)
	if err != nil {
		return nil, err
	}
	return ToResponses(posts), nil
}