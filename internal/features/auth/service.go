package auth

import (
	"net/http"
	"net/mail"
	"regexp"
	"strings"

	"app/internal/database"
	"app/internal/middleware"
	"app/internal/response"

	"golang.org/x/crypto/bcrypt"
)

var usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9_.]{3,50}$`)

type Service struct {
	repo   *Repository
	secret string
}

func NewService(repo *Repository, secret string) *Service {
	return &Service{repo: repo, secret: secret}
}

func validEmail(e string) bool {
	addr, err := mail.ParseAddress(e)
	return err == nil && addr.Address == e && strings.Contains(e, ".")
}

func bad(msg string) error { return response.NewError(http.StatusBadRequest, msg) }

func (s *Service) Register(req RegisterRequest) (*database.User, error) {
	username := strings.TrimSpace(req.Username)
	email := strings.ToLower(strings.TrimSpace(req.Email))
	name := strings.TrimSpace(req.Name)
	if username == "" || email == "" || req.Password == "" {
		return nil, bad("Username, email, dan password wajib diisi")
	}
	if !usernamePattern.MatchString(username) {
		return nil, bad("Username 3-50 karakter dan hanya boleh berisi huruf, angka, titik, atau underscore")
	}
	if !validEmail(email) {
		return nil, bad("Format email tidak valid")
	}
	if len(req.Password) < 6 {
		return nil, bad("Password minimal 6 karakter")
	}
	if s.repo.UsernameTaken(username, 0) {
		return nil, response.NewError(http.StatusBadRequest, "Username sudah digunakan")
	}
	if s.repo.EmailTaken(email, 0) {
		return nil, response.NewError(http.StatusBadRequest, "Email sudah terdaftar")
	}
	if name == "" {
		name = username
	}
	hash, err := database.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	role := "member"
	if strings.EqualFold(strings.TrimSpace(req.Role), "admin") {
		role = "admin"
	}
	u := &database.User{Name: name, Username: username, Email: email, Password: hash, Role: role}
	if err := s.repo.Create(u); err != nil {
		return nil, response.NewError(http.StatusBadRequest, "Username atau email sudah terdaftar")
	}
	return u, nil
}

func (s *Service) IssueToken(u *database.User) (string, error) {
	return middleware.GenerateToken(s.secret, u.ID, u.Username, u.Role)
}

func (s *Service) Login(req LoginRequest) (*LoginResult, error) {
	identifier := strings.TrimSpace(req.Username)
	if identifier == "" {
		identifier = strings.TrimSpace(req.Email)
	}
	if identifier == "" || req.Password == "" {
		return nil, bad("Username dan password wajib diisi")
	}
	u, err := s.repo.FindByLogin(identifier)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)) != nil {
		return nil, response.NewError(http.StatusUnauthorized, "Username atau password salah")
	}
	token, err := s.IssueToken(u)
	if err != nil {
		return nil, err
	}
	return &LoginResult{Token: token, AccessToken: token, TokenType: "Bearer", User: u}, nil
}

func (s *Service) Profile(id uint) (*database.User, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return nil, response.NewError(http.StatusNotFound, "Pengguna tidak ditemukan")
	}
	return u, nil
}

func (s *Service) UpdateProfile(id uint, req UpdateProfileRequest) (*database.User, error) {
	u, err := s.Profile(id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		n := strings.TrimSpace(*req.Name)
		if n == "" {
			return nil, bad("Nama tidak boleh kosong")
		}
		u.Name = n
	}
	if req.Username != nil {
		n := strings.TrimSpace(*req.Username)
		if !usernamePattern.MatchString(n) {
			return nil, bad("Username 3-50 karakter dan hanya boleh berisi huruf, angka, titik, atau underscore")
		}
		if s.repo.UsernameTaken(n, u.ID) {
			return nil, response.NewError(http.StatusBadRequest, "Username sudah digunakan")
		}
		u.Username = n
	}
	if req.Email != nil {
		e := strings.ToLower(strings.TrimSpace(*req.Email))
		if !validEmail(e) {
			return nil, bad("Format email tidak valid")
		}
		if s.repo.EmailTaken(e, u.ID) {
			return nil, response.NewError(http.StatusBadRequest, "Email sudah terdaftar")
		}
		u.Email = e
	}
	if req.Bio != nil {
		u.Bio = strings.TrimSpace(*req.Bio)
	}
	if req.Password != nil && *req.Password != "" {
		if len(*req.Password) < 6 {
			return nil, bad("Password minimal 6 karakter")
		}
		hash, err := database.HashPassword(*req.Password)
		if err != nil {
			return nil, err
		}
		u.Password = hash
	}
	if err := s.repo.Save(u); err != nil {
		return nil, err
	}
	return u, nil
}
