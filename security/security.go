package security

import (
	"encoding/base64"
	"errors"

	"log"
	"strings"
	"time"

	"user_service/types"
	"user_service/util"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var SIGNING_KEY []byte

func GenerateJWT(txid uuid.UUID, user_id uuid.UUID, config types.Config) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)
	claims["iat"] = time.Now().UTC().Unix()
	claims["exp"] = time.Now().Add(time.Duration(config.App.LoginExpirationMs) * time.Millisecond).UTC().Unix()
	// TODO [drd] update this to be auth.jfl.com (or whatever the url ends up being) and jti will take over the txid functionality
	claims["iss"] = txid
	claims["jti"] = txid
	// TODO tie to user agent as well
	claims["user_id"] = user_id.String()
	signed_token, err := token.SignedString(SIGNING_KEY)
	if err != nil {
		return "", err
	}

	return signed_token, nil
}

func GenerateSigningKey() {
	SIGNING_KEY = []byte(util.RandomString(64))
}

func GetBasicAuth(auth string, config types.Config) (string, string, bool, error) {
	// Basically copied from gofiber/basicauth/main.go
	// Check if header is valid
	if len(auth) > 6 && strings.ToLower(auth[:5]) == "basic" {
		// Try to decode
		raw, err := base64.StdEncoding.DecodeString(auth[6:])
		if err != nil {
			return "", "", false, err
		}
		credentials := string(raw)
		// Find semicolumn
		for i := 0; i < len(credentials); i++ {
			if credentials[i] == ':' {
				// Split into user & pass
				username := credentials[:i]
				password := credentials[i+1:]
				return username, password, true, nil
			}
		}
	}
	return "", "", false, errors.New("invalid header")
}

func parseToken(token string) (jwt.MapClaims, error) {
	token = strings.TrimPrefix(token, "Bearer ")
	parsed_token, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return "", errors.New("invalid signing method")
		}
		return SIGNING_KEY, nil
	})
	if err != nil || !parsed_token.Valid {
		log.Println(err.Error())
		return nil, errors.New("invalid jwt")
	}

	claims, ok := parsed_token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("missing claims")
	}
	return claims, nil
}

func ValidateJWT(c *fiber.Ctx, config types.Config) (jwt.MapClaims, error) {
	token := c.Get(fiber.HeaderAuthorization)
	if strings.HasPrefix(token, "Bearer ") {
		passed_claims, err := parseToken(token)
		if err != nil {
			return nil, errors.New("failed to parse current token")
		}
		passed_claims["user_id"], err = uuid.Parse(passed_claims["user_id"].(string))
		if err != nil {
			return nil, errors.New("invalid user")
		}

		// Make sure the token is valid
		if !passed_claims.VerifyExpiresAt(time.Now().UTC().Unix(), true) {
			return nil, errors.New("token expired or not set")
		}
		if !passed_claims.VerifyIssuedAt(time.Now().UTC().Unix(), true) {
			return nil, errors.New("issued_at not set")
		}
		if !passed_claims.VerifyIssuer(config.App.Host.Issuer, true) {
			return nil, errors.New("issuer not set")
		}
		return passed_claims, nil
	}
	return nil, errors.New("invalid credentials")
}
