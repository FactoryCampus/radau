package api

import (
	"errors"
	"os"
	"regexp"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	errorsPkg "github.com/pkg/errors"
)

var bearerRegex = regexp.MustCompile(`^(?:B|b)earer (\S+$)`)
var jwtParser = jwt.Parser{
	ValidMethods: []string{jwt.SigningMethodHS256.Name},
}

type Claims struct {
	jwt.StandardClaims
	Email string `json:"email"`
}

func IsValidApiKey(ctx *gin.Context, validKey string) error {
	apiKey := ctx.GetHeader("Authorization")
	if apiKey != validKey {
		return errors.New("Invalid API key")
	}

	ctx.Set("apiKey", validKey)

	return nil
}

func HandleApiKeyAuth(validKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := IsValidApiKey(c, validKey); err != nil {
			c.AbortWithStatus(401)
		}
	}
}

func extractToken(ctx *gin.Context) (string, error) {
	authorization := ctx.GetHeader("Authorization")
	if authorization == "" {
		return "", errors.New("Header not found")
	}

	matches := bearerRegex.FindStringSubmatch(authorization)
	if len(matches) != 2 {
		return "", errors.New("Could not split Bearer")
	}

	return matches[1], nil
}

func parseJWT(ctx *gin.Context) (*jwt.Token, error) {
	tokenValue, err := extractToken(ctx)
	if err != nil {
		return nil, err
	}

	token, parseErr := jwtParser.ParseWithClaims(
		tokenValue,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		},
	)

	if parseErr != nil {
		return nil, parseErr
	}

	if !token.Valid {
		return nil, errors.New("Invalid token")
	}

	return token, nil
}

func (h *baseHandler) CanEditUser(ctx *gin.Context) error {
	user := h.findUserOrHandleMissing(ctx)
	if user == nil {
		return errors.New("User not found")
	}

	token, err := parseJWT(ctx)
	if err != nil {
		return err
	}

	claims := token.Claims.(*Claims)

	if user.Username != claims.Email {
		return errors.New("Action not allowed")
	}

	return nil
}

func (h *baseHandler) EnsureTokenOrKey(validKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		keyErr := IsValidApiKey(ctx, validKey)
		if keyErr == nil {
			return
		}

		tokenErr := h.CanEditUser(ctx)
		if tokenErr == nil {
			return
		}

		ctx.AbortWithError(401, errorsPkg.Wrap(tokenErr, "Token auth failed"))
	}
}
