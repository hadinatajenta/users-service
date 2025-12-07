package users

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Health(c *gin.Context) {
	responseSuccess(c, http.StatusOK, "Service is healthy", gin.H{"status": "ok"})
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserInput
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.CreateUser(c.Request.Context(), req)
	if err != nil {
		responseError(c, http.StatusBadRequest, err.Error())
		return
	}

	responseSuccess(c, http.StatusCreated, "User created", user)
}

func (h *Handler) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		responseError(c, http.StatusBadRequest, "invalid id")
		return
	}

	user, err := h.service.GetUser(c.Request.Context(), id)
	if err != nil {
		responseError(c, http.StatusNotFound, err.Error())
		return
	}

	responseSuccess(c, http.StatusOK, "User retrieved", user)
}

func (h *Handler) ListUsers(c *gin.Context) {
	users, err := h.service.ListUsers(c.Request.Context())
	if err != nil {
		responseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	responseSuccess(c, http.StatusOK, "Users retrieved", users)
}
