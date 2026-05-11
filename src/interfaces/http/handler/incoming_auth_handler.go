package handler

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/infraLinkit/mediaplatform-datasource-v2/src/domain/entity"
)

// audienceMatches: aud claim bisa string ATAU []interface{}. Cek match dgn expected.
func audienceMatches(audClaim interface{}, expected string) bool {
	switch v := audClaim.(type) {
	case string:
		return v == expected
	case []interface{}:
		for _, a := range v {
			if s, ok := a.(string); ok && s == expected {
				return true
			}
		}
	}
	return false
}

// RevokeJWT: tambah jti ke blacklist Redis dgn TTL = sisa exp token.
// Pakai saat logout. No-op kalau Redis nil (degraded mode).
func (h *IncomingHandler) RevokeJWT(jti string, ttl time.Duration) error {
	if h.RCP == nil || jti == "" {
		return nil
	}
	return h.RCP.Set("jwt:blacklist:"+jti, "1", ttl).Err()
}

func (h *IncomingHandler) AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		h.Logs.Warn("AuthMiddleware: Missing or invalid Authorization header")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid token"})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	secret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if secret == "" {
		h.Logs.Error("AuthMiddleware: JWT_SECRET environment variable not set")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "JWT secret not configured"})
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		h.Logs.WithError(err).Warn("AuthMiddleware: Token parse error")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}
	if !token.Valid {
		h.Logs.Warn("AuthMiddleware: Token is invalid")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		h.Logs.Warn("AuthMiddleware: Invalid token claims")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid claims"})
	}

	now := time.Now().Unix()

	if exp, ok := claims["exp"].(float64); ok {
		if int64(exp) < now {
			h.Logs.Warn("AuthMiddleware: Token expired")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token expired"})
		}
	}

	// nbf: tolak token belum aktif (with 30s leeway)
	if nbf, ok := claims["nbf"].(float64); ok {
		if int64(nbf)-30 > now {
			h.Logs.Warn("AuthMiddleware: Token not yet valid (nbf)")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token not yet active"})
		}
	}

	// aud: opt-in via env. Kalau JWT_AUDIENCE diset, token wajib match.
	if expectedAud := strings.TrimSpace(os.Getenv("JWT_AUDIENCE")); expectedAud != "" {
		if !audienceMatches(claims["aud"], expectedAud) {
			h.Logs.Warn("AuthMiddleware: Invalid audience")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid audience"})
		}
	}

	// iss: opt-in via env.
	if expectedIss := strings.TrimSpace(os.Getenv("JWT_ISSUER")); expectedIss != "" {
		if iss, _ := claims["iss"].(string); iss != expectedIss {
			h.Logs.Warn("AuthMiddleware: Invalid issuer")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid issuer"})
		}
	}

	if tokenType, ok := claims["type"].(string); !ok || tokenType != "access" {
		h.Logs.Warn("AuthMiddleware: Invalid token type")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token type"})
	}

	jti, ok := claims["jti"].(string)
	if !ok || jti == "" {
		h.Logs.Warn("AuthMiddleware: Missing or invalid jti in token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token (jti)"})
	}

	// jti blacklist check via Redis. Skip kalau Redis nil (degraded mode #4).
	if h.RCP != nil {
		blKey := "jwt:blacklist:" + jti
		if v := h.RCP.Get(blKey); v.Err() == nil && v.Val() != "" {
			h.Logs.Warnf("AuthMiddleware: Token revoked (jti=%s)", jti)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token revoked"})
		}
	}

	var userID int
	switch v := claims["sub"].(type) {
	case float64:
		userID = int(v)
	case string:
		_, err := fmt.Sscanf(v, "%d", &userID)
		if err != nil {
			h.Logs.WithError(err).Warn("AuthMiddleware: Invalid user ID format in token sub claim")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID in token"})
		}
	default:
		h.Logs.Warn("AuthMiddleware: Unknown type for user ID in token sub claim")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID in token"})
	}

	var companies []entity.Company
	err = h.DB.Raw(`
		SELECT c.id, c.name
		FROM user_companies uc
		JOIN companies c ON uc.company_id = c.id
		WHERE uc.user_id = ? AND uc.status = true
	`, userID).Scan(&companies).Error
	if err != nil {
		h.Logs.WithError(err).Errorf("AuthMiddleware: Database error fetching companies for user %d", userID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	var companyNames []string
	for _, cpy := range companies {
		companyNames = append(companyNames, cpy.Name)
	}

	c.Locals("companies", companyNames)

	var adnets []entity.AdnetList
	err = h.DB.Raw(`
		SELECT a.id, a.code
		FROM user_adnets ua
		JOIN adnet_lists a ON ua.adnet_id = a.id
		WHERE ua.user_id = ? AND ua.status = true
	`, userID).Scan(&adnets).Error
	if err != nil {
		h.Logs.WithError(err).Errorf("AuthMiddleware: Database error fetching adnets for user %d", userID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	var adnetCodes []string
	for _, adn := range adnets {
		adnetCodes = append(adnetCodes, adn.Code)
	}

	c.Locals("adnets", adnetCodes)
	c.Locals("user_id", userID)

	var agencies []entity.Agency
	err = h.DB.Raw(`
		SELECT a.id, upper(a.name) as name
		FROM user_agencies ua
		JOIN agencies a ON ua.agency_id = a.id
		WHERE ua.user_id = ? AND ua.status = true
	`, userID).Scan(&agencies).Error
	if err != nil {
		h.Logs.WithError(err).Errorf("AuthMiddleware: Database error fetching agencies for user %d", userID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	var agencyNames []string
	for _, agency := range agencies {
		agencyNames = append(agencyNames, agency.Name)
	}

	c.Locals("agencies", agencyNames)

	return c.Next()
}
