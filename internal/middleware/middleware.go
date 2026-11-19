package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"app/internal/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

var revoked sync.Map

func GenerateToken(secret string, userID uint, username, role string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

func RevokeToken(token string) {
	if token != "" {
		revoked.Store(token, time.Now())
	}
}

func BearerToken(c *gin.Context) string {
	h := strings.TrimSpace(c.GetHeader("Authorization"))
	if len(h) < 7 || !strings.EqualFold(h[:7], "Bearer ") {
		return ""
	}
	return strings.TrimSpace(h[7:])
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		h := c.Writer.Header()
		if origin != "" {
			h.Set("Access-Control-Allow-Origin", origin)
			h.Set("Access-Control-Allow-Credentials", "true")
			h.Add("Vary", "Origin")
		} else {
			h.Set("Access-Control-Allow-Origin", "*")
		}
		h.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		h.Set("Access-Control-Max-Age", "86400")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := BearerToken(c)
		if tokenStr == "" {
			response.Fail(c, http.StatusUnauthorized, "Token autentikasi tidak ditemukan")
			return
		}
		if _, ok := revoked.Load(tokenStr); ok {
			response.Fail(c, http.StatusUnauthorized, "Token sudah tidak berlaku, silakan login kembali")
			return
		}
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
			return []byte(secret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		if err != nil || !token.Valid || claims.UserID == 0 {
			response.Fail(c, http.StatusUnauthorized, "Token tidak valid atau sudah kedaluwarsa")
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("token", tokenStr)
		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		for _, r := range roles {
			if r == role {
				c.Next()
				return
			}
		}
		response.Fail(c, http.StatusForbidden, "Akses ditolak: Anda tidak memiliki hak akses")
	}
}

func RateLimiter(limit int, window time.Duration) gin.HandlerFunc {
	type entry struct {
		count int
		reset time.Time
	}
	var mu sync.Mutex
	hits := map[string]*entry{}
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()
		mu.Lock()
		e, ok := hits[ip]
		if !ok || now.After(e.reset) {
			e = &entry{reset: now.Add(window)}
			hits[ip] = e
		}
		e.count++
		over := e.count > limit
		mu.Unlock()
		if over {
			response.Fail(c, http.StatusTooManyRequests, "Terlalu banyak permintaan, coba lagi nanti")
			return
		}
		c.Next()
	}
}
