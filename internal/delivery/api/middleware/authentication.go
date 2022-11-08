package middleware

import (
	"errors"
	"fmt"
	"strings"

	"gitlab.com/stoqu/stoqu-be/internal/config"

	baseModel "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	ctxval "gitlab.com/stoqu/stoqu-be/pkg/util/ctxval"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	var (
		jwtKey = config.Config.JWT.Secret
	)

	return func(c echo.Context) error {
		authToken := c.Request().Header.Get("Authorization")
		if authToken == "" {
			return res.ErrorBuilder(res.Constant.Error.Unauthorized, errors.New("auth token not found")).Send(c)
		}

		splitToken := strings.Split(authToken, "Bearer ")
		token, err := jwt.Parse(splitToken[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method :%v", token.Header["alg"])
			}

			return []byte(jwtKey), nil
		})

		if !token.Valid || err != nil {
			return res.ErrorBuilder(res.Constant.Error.Unauthorized, err).Send(c)
		}

		var id string
		destructID := token.Claims.(jwt.MapClaims)["id"]
		if destructID != nil {
			id = destructID.(string)
		}

		var name string
		destructName := token.Claims.(jwt.MapClaims)["name"]
		if destructName != nil {
			name = destructName.(string)
		}

		var email string
		destructEmail := token.Claims.(jwt.MapClaims)["email"]
		if destructEmail != nil {
			email = destructEmail.(string)
		}

		authCtx := &baseModel.AuthContext{
			ID:    id,
			Name:  name,
			Email: email,
		}
		ctx := ctxval.SetAuthValue(c.Request().Context(), authCtx)
		newRequest := c.Request().WithContext(ctx)
		c.SetRequest(newRequest)

		return next(c)
	}
}
