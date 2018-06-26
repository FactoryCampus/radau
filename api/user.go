package api

import (
	"os"
	"time"

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
	LastQuery       time.Time         `sql:"lastquery"`
}

type UserCreation struct {
	Username        string            `json:"username"`
	ExtraProperties map[string]string `json:"extraProperties"`
}

type UserOutput struct {
	Username        string            `json:"username"`
	ExtraProperties map[string]string `json:"extraProperties"`
	LastQuery       string            `json:"lastQuery"`
}

type UserOutputWithToken struct {
	UserOutput
	Token string `json:"token"`
}

func InitUserHandler(router gin.IRouter, db *pg.DB) {
	h := &userHandler{&baseHandler{
		db: db,
	}}

	authManagementKey := os.Getenv("API_KEY_MANAGEMENT")
	adminRoutes := router.Group("", HandleApiKeyAuth(authManagementKey))

	adminRoutes.GET("/users", h.getUsers)
	adminUserGroup := adminRoutes.Group("/user")
	{
		adminUserGroup.POST("", h.createUser)
		adminUserGroup.PUT("/:username", h.updateUser)
		adminUserGroup.DELETE("/:username", h.deleteUser)
	}

	router.GET("/user/:username", h.EnsureTokenOrKey(authManagementKey), h.getUser)
}

func serializeUser(user *User) UserOutput {
	// Make a copy because go-pg will always give us the identical object
	// and modify it in subsequent invocations, instead of yielding new rows.
	extraProps := map[string]string{}
	for k, v := range user.ExtraProperties {
		extraProps[k] = v
	}

	o := UserOutput{
		Username:        user.Username,
		ExtraProperties: extraProps,
		LastQuery:       user.LastQuery.String(),
	}
	if user.LastQuery.IsZero() {
		o.LastQuery = ""
	}
	return o
}

func (h *userHandler) getUsers(c *gin.Context) {
	limit := c.GetInt("limit")
	offset := c.GetInt("offset")

	userList := []UserOutputWithToken{}
	h.db.Model(&User{}).
		Limit(limit).
		Offset(offset).
		ForEach(func(u *User) error {
			o := UserOutputWithToken{
				UserOutput: serializeUser(u),
				Token:      u.Token,
			}
			userList = append(userList, o)
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
