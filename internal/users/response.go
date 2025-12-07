package users

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type successResponse struct {
	Success bool        `json:"success"`
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// problemDetails follows RFC 7807 to provide standardized error bodies.
type problemDetails struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	// Status is the HTTP status code.
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

func responseSuccess(c *gin.Context, status int, message string, data interface{}) {
	if message == "" {
		message = http.StatusText(status)
	}

	res := successResponse{
		Success: true,
		Status:  status,
		Message: message,
		Data:    data,
	}

	c.JSON(status, res)
}

func responseError(c *gin.Context, status int, detail string) {
	if status < 400 {
		status = http.StatusInternalServerError
	}

	problem := problemDetails{
		Type:     fmt.Sprintf("https://httpstatuses.com/%d", status),
		Title:    http.StatusText(status),
		Status:   status,
		Detail:   detail,
		Instance: c.Request.URL.Path,
	}

	c.Header("Content-Type", "application/problem+json")
	c.JSON(status, problem)
}
