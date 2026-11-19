package auth

import (
	"net/http"

	"app/internal/middleware"
	"app/internal/response"

	"github.com/gin-gonic/gin"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Format request tidak valid")
		return
	}
	user, err := h.svc.Register(req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	token, err := h.svc.IssueToken(user)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusCreated, "Registrasi berhasil", gin.H{
		"id": user.ID, "name": user.Name, "username": user.Username, "email": user.Email, "role": user.Role,
		"bio": user.Bio, "created_at": user.CreatedAt, "updated_at": user.UpdatedAt,
		"user": user, "token": token,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Format request tidak valid")
		return
	}
	result, err := h.svc.Login(req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Login berhasil", result)
}

func (h *Handler) Logout(c *gin.Context) {
	middleware.RevokeToken(c.GetString("token"))
	response.Success(c, http.StatusOK, "Logout berhasil", nil)
}

func (h *Handler) Profile(c *gin.Context) {
	user, err := h.svc.Profile(c.GetUint("user_id"))
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Profil berhasil diambil", user)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Format request tidak valid")
		return
	}
	user, err := h.svc.UpdateProfile(c.GetUint("user_id"), req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Profil berhasil diperbarui", user)
}
