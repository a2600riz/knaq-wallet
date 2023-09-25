package security

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"knaq-wallet/config"
	"strconv"
	"strings"
)

const ClaimDataKey = "jwtClaimData"

type JwtCustomClaims struct {
	jwt.StandardClaims
	Data JwtCustomData `json:"data,omitempty"`
}
type JwtCustomData struct {
	UserId   string `json:"userId,omitempty"`
	DeviceId string `json:"deviceId,omitempty"`
	Role     string `json:"role,omitempty"`
	Side     string `json:"side,omitempty"`
}

type UserClaimData struct {
	Email  string
	UserID uint64
}

func ParseClaimUserData(ctx echo.Context) (userClaimData UserClaimData, err error) {
	tokenString := ctx.Request().Header.Get("Authorization")
	tokenString = strings.Replace(tokenString, "Bearer ", "", -1)
	claims := jwt.MapClaims{}
	//dec, err := base64.URLEncoding.DecodeString(config.Config.GetJWTSecret())
	//if err != nil {
	//	return userClaimData, err
	//}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return config.Config.GetJWTSecret(), nil
	})
	// ... error handling
	//if err != nil {
	//	if err.Error() != "Token is expired" {
	//		return userClaimData, err
	//	}
	//}

	claims, ok := token.Claims.(jwt.MapClaims)
	userId, ok := claims["sub"].(string) // get email from JWT
	email, ok := claims["email"].(string)
	if !ok {
		return userClaimData, fmt.Errorf("couldn't parse claims")
	}

	uid, err := strconv.Atoi(userId)
	if err != nil {
		return userClaimData, err
	}

	userClaimData = UserClaimData{
		Email:  email,
		UserID: uint64(uid),
	}

	return userClaimData, nil
}

func ProcessUserClaim(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var (
			err       error
			claimData UserClaimData
		)
		claimData, err = ParseClaimUserData(ctx)
		if err != nil {
			return err
		}
		ctx.Set(ClaimDataKey, claimData)

		return next(ctx)
	}
}
