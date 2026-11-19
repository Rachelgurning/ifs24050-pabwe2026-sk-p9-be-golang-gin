package post

import (
	"net/http"
	"strconv"

	"app/internal/response"

	"github.com/gin-gonic/gin"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func parseID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, http.StatusBadRequest, "ID tidak valid")
		return 0, false
	}
	return uint(id), true
}

func bindPost(c *gin.Context) (PostRequest, bool) {
	var req PostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Format request tidak valid")
		return req, false
	}
	return req, true
}

func (h *Handler) Categories(c *gin.Context) {
	response.Success(c, http.StatusOK, "Daftar kategori berhasil diambil", gin.H{"categories": Categories})
}

func (h *Handler) PublicList(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		q = c.Query("search")
	}
	posts, err := h.svc.PublicList(ListFilter{Query: q, Category: c.Query("category")})
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Daftar artikel berhasil diambil", gin.H{"posts": posts, "total": len(posts)})
}

func (h *Handler) PublicDetail(c *gin.Context) {
	p, err := h.svc.PublicDetail(c.Param("slug"))
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Detail artikel berhasil diambil", gin.H{"post": p})
}

func (h *Handler) MyPosts(c *gin.Context) {
	posts, err := h.svc.ListMine(c.GetUint("user_id"))
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Daftar artikel berhasil diambil", gin.H{"posts": posts, "total": len(posts)})
}

func (h *Handler) Create(c *gin.Context) {
	req, ok := bindPost(c)
	if !ok {
		return
	}
	p, err := h.svc.Create(c.GetUint("user_id"), req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusCreated, "Artikel berhasil dibuat", gin.H{"post": p})
}

func (h *Handler) Stats(c *gin.Context) {
	st := h.svc.Stats(c.GetUint("user_id"))
	response.Success(c, http.StatusOK, "Statistik artikel berhasil diambil", gin.H{
		"stats": st, "total_posts": st.TotalPosts, "published_posts": st.PublishedPosts, "draft_posts": st.DraftPosts,
		"total": st.Total, "published": st.Published, "draft": st.Draft,
	})
}

func (h *Handler) Detail(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	p, err := h.svc.GetMine(c.GetUint("user_id"), id)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Detail artikel berhasil diambil", gin.H{"post": p})
}

func (h *Handler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	req, ok := bindPost(c)
	if !ok {
		return
	}
	p, err := h.svc.Update(c.GetUint("user_id"), id, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Artikel berhasil diperbarui", gin.H{"post": p})
}

func (h *Handler) TogglePublish(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var body struct {
		IsPublished *bool `json:"is_published"`
	}
	_ = c.ShouldBindJSON(&body)
	p, err := h.svc.TogglePublish(c.GetUint("user_id"), id, body.IsPublished)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	msg := "Artikel dikembalikan ke draf"
	if p.IsPublished {
		msg = "Artikel berhasil dipublikasikan"
	}
	response.Success(c, http.StatusOK, msg, gin.H{"post": p})
}

func (h *Handler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.svc.Delete(c.GetUint("user_id"), id); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Artikel berhasil dihapus", nil)
}

func (h *Handler) AdminUsers(c *gin.Context) {
	users, err := h.svc.AdminUsers()
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Daftar pengguna berhasil diambil", gin.H{"users": users, "total": len(users)})
}

func (h *Handler) AdminPosts(c *gin.Context) {
	posts, err := h.svc.AdminPosts()
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Daftar seluruh artikel berhasil diambil", gin.H{"posts": posts, "total": len(posts)})
}

func (h *Handler) AdminUserPosts(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	user, posts, err := h.svc.AdminUserPosts(id)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Daftar artikel pengguna berhasil diambil", gin.H{"user": user, "posts": posts, "total": len(posts)})
}
