package transport

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userID"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "error: empty authorization header")
		return
	}

	headerIntoParts := strings.Split(header, " ")
	if len(headerIntoParts) != 2 || headerIntoParts[0] != "Bearer" {
		newErrorResponse(c, http.StatusUnauthorized, "error: invalid authorization header")
		return
	}

	if len(headerIntoParts[1]) == 0 {
		newErrorResponse(c, http.StatusUnauthorized, "error: empty token")
		return
	}

	userID, err := h.services.Authorization.ParseToken(headerIntoParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userCtx, userID)
}

func getUserID(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		return 0, errors.New("user id is of invalid type")
	}

	return idInt, nil
}
