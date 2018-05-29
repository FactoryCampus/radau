package pkg

import (
	"crypto/rand"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
)

type tokenHandler struct {
	*baseHandler
}

func InitTokenHandler(router gin.IRoutes, db *pg.DB) {
	h := &tokenHandler{&baseHandler{
		db: db,
	}}

	router.POST("/:username", h.createToken)
	router.DELETE("/:username", h.deleteToken)
}

const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const charMask = 1<<6 - 1 // 64-bit mask

func randString(n int) string {
	bufferSize := int(float64(n) * 1.3)
	result := make([]byte, bufferSize)
	rand.Read(result)
	for i := range result {
		result[i] = characters[int(result[i]&charMask)%len(characters)]
	}
	return string(result)
}

func (h *tokenHandler) createToken(c *gin.Context) {
	user := h.findUserOrHandleMissing(c)
	if user == nil {
		return
	}

	user.Token = randString(32)

	err := h.db.Update(user)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"token": user.Token,
	})
}

func (h *tokenHandler) deleteToken(c *gin.Context) {
	user := h.findUserOrHandleMissing(c)
	if user == nil {
		return
	}

	user.Token = ""
	err := h.db.Update(user)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.Status(204)
}
