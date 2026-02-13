package security

import (
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"os"

	"log"
	"strings"
	"time"

	"github.com/thedanisaur/jfl_platform/types"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func GenerateJWT(txid uuid.UUID, user_claims types.UserClaims, config types.Config, private_key *rsa.PrivateKey) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iat"] = time.Now().UTC().Unix()
	claims["exp"] = time.Now().Add(time.Duration(config.App.LoginExpirationMs) * time.Millisecond).UTC().Unix()
	claims["iss"] = config.App.Host.Issuer
	claims["jti"] = txid.String()
	claims["user_id"] = user_claims.UserID.String()
	claims["issuing_unit"] = user_claims.IssuingUnit
	claims["role_name"] = user_claims.RoleName
	signed_token, err := token.SignedString(private_key)
	if err != nil {
		return "", err
	}

	return signed_token, nil
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

func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(bytes)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func LoadPublicKey(path string) (*rsa.PublicKey, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(bytes)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func mapToUserClaims(txid uuid.UUID, claims map[string]interface{}) (types.UserClaims, error) {
	user_claims := types.UserClaims{}

	user_id, err := uuid.Parse(claims["user_id"].(string))
	if err != nil {
		log.Printf("%s | missing user id\n", txid.String())
		return user_claims, errors.New("invalid user claims")
	}
	user_claims.UserID = user_id
	issuing_unit, ok := claims["issuing_unit"].(string)
	if !ok {
		log.Printf("%s | missing issuing unit\n", txid.String())
		return user_claims, errors.New("invalid user claims")
	}
	user_claims.IssuingUnit = issuing_unit
	role_name, ok := claims["role_name"].(string)
	if !ok {
		log.Printf("%s | missing role name\n", txid.String())
		return user_claims, errors.New("invalid user claims")
	}
	user_claims.RoleName = role_name

	return user_claims, nil
}

func parseToken(token string, public_key *rsa.PublicKey) (jwt.MapClaims, error) {
	token = strings.TrimPrefix(token, "Bearer ")
	parsed_token, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodRSA)
		if !ok {
			return nil, errors.New("invalid signing method")
		}
		return public_key, nil
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

func ValidateJWT(txid uuid.UUID, c *fiber.Ctx, config types.Config, public_key *rsa.PublicKey) (types.UserClaims, error) {
	token := c.Get(fiber.HeaderAuthorization)
	if !strings.HasPrefix(token, "Bearer ") {
		return types.UserClaims{}, errors.New("invalid credentials")
	}
	passed_claims, err := parseToken(token, public_key)
	if err != nil {
		return types.UserClaims{}, errors.New("failed to parse current token")
	}
	// Make sure the token is valid
	if !passed_claims.VerifyExpiresAt(time.Now().UTC().Unix(), true) {
		return types.UserClaims{}, errors.New("token expired or not set")
	}
	if !passed_claims.VerifyIssuedAt(time.Now().UTC().Unix(), true) {
		return types.UserClaims{}, errors.New("issued_at not set or invalid")
	}
	if !passed_claims.VerifyIssuer(config.App.Host.Issuer, true) {
		return types.UserClaims{}, errors.New("issuer not set or invalid")
	}

	// Make sure the user is valid
	user_claims, err := mapToUserClaims(txid, passed_claims)
	return user_claims, err
}
