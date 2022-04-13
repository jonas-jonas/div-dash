package identity

import (
	"div-dash/internal/db"
	"div-dash/internal/httputil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) getUser(c *gin.Context) {
	id := c.Param("id")

	user, err := h.Queries().GetUser(c, id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			httputil.AbortNotFound(c)
			return
		}
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, userResponseFromUser(user))
}

func userResponseFromUser(user db.User) UserResponse {
	return UserResponse{
		ID:    user.ID,
		Email: user.Email,
	}
}
