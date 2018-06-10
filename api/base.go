package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
)

type baseHandler struct {
	db *pg.DB
}

func (h *baseHandler) findUserOrHandleMissing(c *gin.Context) *User {
	var users []User
	err := h.db.Model(&users).Where("username = ?", c.Param("username")).Limit(1).Select()
	if err != nil {
		c.String(500, err.Error())
		return nil
	}
	if len(users) == 0 {
		c.Status(404)
		return nil
	}

	return &users[0]
}
