package pkg

import (
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/imdario/mergo"
)

type userHandler struct {
	*baseHandler
}

type User struct {
	ID              int
	Username        string            `sql:"username,notnull,unique"`
	Token           string            `sql:"token"`
	ExtraProperties map[string]string `sql:"extraProperties,hstore"`
}

type UserCreation struct {
	Username        string            `json:"username"`
	ExtraProperties map[string]string `json:"extraProperties"`
}

func InitUserHandler(router gin.IRouter, db *pg.DB) {
	h := &userHandler{&baseHandler{
		db: db,
	}}

	router.GET("/users", h.getUsers)

	group := router.Group("/user")
	{
		group.POST("", h.createUser)
		group.GET("/:username", h.getUser)
		group.PUT("/:username", h.updateUser)
		group.DELETE("/:username", h.deleteUser)
	}
}

func serializeUser(user *User) gin.H {
	o := gin.H{
		"username":        user.Username,
		"extraProperties": user.ExtraProperties,
	}
	if user.ExtraProperties == nil {
		o["extraProperties"] = map[string]string{}
	}
	return o
}

func (h *userHandler) getUsers(c *gin.Context) {
	limit := c.GetInt("limit")
	offset := c.GetInt("offset")

	userList := []gin.H{}
	h.db.Model(&User{}).
		Limit(limit).
		Offset(offset).
		ForEach(func(u *User) error {
			userList = append(userList, gin.H{
				"username":        u.Username,
				"extraProperties": u.ExtraProperties,
				"token":           u.Token,
			})
			return nil
		})
	c.JSON(200, gin.H{
		"users":     &userList,
		"userCount": len(userList),
	})
}

func (h *userHandler) createUser(c *gin.Context) {
	var body UserCreation
	bodyMissing := c.BindJSON(&body)
	if bodyMissing != nil {
		c.JSON(400, gin.H{
			"error": "Request body missing",
		})
		return
	}

	numUsers, err := h.db.Model(new(User)).Where("username = ?", body.Username).Count()
	if err != nil {
		c.Status(500)
		return
	}
	if numUsers > 0 {
		c.JSON(422, gin.H{
			"error": "User already exists",
		})
		return
	}

	userInsert := &User{
		Username:        body.Username,
		ExtraProperties: body.ExtraProperties,
	}
	insertFailed := h.db.Insert(userInsert)
	if insertFailed != nil {
		c.Status(500)
		return
	}

	c.JSON(200, serializeUser(userInsert))
}

func (h *userHandler) getUser(c *gin.Context) {
	user := h.findUserOrHandleMissing(c)
	if user == nil {
		return
	}

	c.JSON(200, serializeUser(user))
}

func (h *userHandler) updateUser(c *gin.Context) {
	user := h.findUserOrHandleMissing(c)
	if user == nil {
		return
	}

	payload := &User{}
	bodyParseError := c.BindJSON(payload)
	if bodyParseError != nil {
		c.String(400, bodyParseError.Error())
		return
	}

	mergo.Merge(user, payload)
	user.ExtraProperties = payload.ExtraProperties

	err := h.db.Update(user)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.JSON(200, serializeUser(user))
}

func (h *userHandler) deleteUser(c *gin.Context) {
	user := h.findUserOrHandleMissing(c)
	if user == nil {
		return
	}

	err := h.db.Delete(user)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.Status(204)
}
