package api

import (
	"crypto/rand"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
)

type tokenHandler struct {
	*baseHandler
}

func InitTokenHandler(router gin.IRouter, db *pg.DB) {
	h := &tokenHandler{&baseHandler{
		db: db,
	}}

	group := router.Group("/token")
	{
		group.POST("/:username", h.createToken)
		group.DELETE("/:username", h.deleteToken)
	}
}

const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const charMask = 1<<6 - 1 // 64-bit mask

func randString(n int) string {
	result := make([]byte, n)
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

	tokenLength := 32

	tokenLengthStr, tokenLengthExists := os.LookupEnv("TOKEN_LENGTH")
	if tokenLengthExists {
		tokenLengthParsed, err := strconv.Atoi(tokenLengthStr)
		if err == nil {
			tokenLength = tokenLengthParsed
		}
	}

	user.Token = randString(tokenLength)

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
