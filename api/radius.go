package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/imdario/mergo"
)

type radiusHandler struct {
	*baseHandler
}

func InitRadiusHandler(router gin.IRoutes, db *pg.DB) {
	h := &radiusHandler{&baseHandler{
		db: db,
	}}

	router.GET("/:username", h.getRadiusInfo)
}

func (h *radiusHandler) getRadiusInfo(c *gin.Context) {
	user := h.findUserOrHandleMissing(c)
	if user == nil {
		return
	}

	if user.Token == "" {
		c.Status(401)
		return
	}

	radiusResponse := &map[string]string{}
	mergo.Merge(radiusResponse, user.ExtraProperties, mergo.WithOverride)
	mergo.Merge(radiusResponse, map[string]string{
		"control:Cleartext-Password": user.Token,
	}, mergo.WithOverride)
	c.JSON(200, radiusResponse)
}
