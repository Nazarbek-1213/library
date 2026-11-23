package middleware

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims — JWT token payload qismida saqlanadigan ma'lumotlar.
// Talab bo'yicha: id va role bo'lishi shart.
type Claims struct {
	ID   int64  `json:"id"`   // Foydalanuvchining ma'lumotlar bazasidagi ID raqami
	Role string `json:"role"` // Foydalanuvchi roli: "member" yoki "admin"
	jwt.RegisteredClaims
}

// GenerateToken — yangi a'zo yoki admin uchun JWT token yaratadi.
// Payload ichiga foydalanuvchi id'si va roli yoziladi.
func GenerateToken(id int64, role string) (string, error) {
	if id <= 0 || (role != "admin" && role != "member") {
		return "", errors.New("noto'g'ri id yoki role")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET topilmadi")
	}

	claims := Claims{
		ID:   id,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token 24 soat amal qiladi
			IssuedAt:  jwt.NewNumericDate(time.Now()),                     // Token berilgan vaqt
		},
	}

	// HS256 algoritmi va maxfiy kalit (secret) bilan imzolaymiz
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// NewToken — GenerateToken uchun alias (mavjud testlar buzilmasligi uchun).
func NewToken(id int64, role string) (string, error) {
	return GenerateToken(id, role)
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Authorization sarlavhasidagi Bearer tokenini olamiz.
		header := strings.TrimSpace(r.Header.Get("Authorization"))
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		claims := &Claims{}
		// JWT imzosini va ichidagi id, role qiymatlarini tekshiramiz.
		token, err := jwt.ParseWithClaims(parts[1], claims, func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid || claims.ID <= 0 || (claims.Role != "admin" && claims.Role != "member") {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		// Tekshirilgan claimsni keyingi middleware olishi uchun contextga qo'yamiz.
		next.ServeHTTP(w, r.WithContext(WithClaims(r.Context(), claims)))
	})
}

// AuthMiddleware — Authenticate uchun alias (PDF hujjatidagi nom bilan mos bo'lishi uchun).
func AuthMiddleware(next http.Handler) http.Handler {
	return Authenticate(next)
}

type contextKey struct{}

// WithClaims — Claims ma'lumotlarini context ichiga saqlaydi.
func WithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, contextKey{}, claims)
}

// GetClaims — context ichidan Claims ma'lumotlarini oladi.
func GetClaims(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(contextKey{}).(*Claims)
	return claims, ok
}

// GetUserID — context ichidan foydalanuvchi ID raqamini oladi (PDF: a'zo endpointlari ID ni tokendan oladi).
func GetUserID(ctx context.Context) (int64, bool) {
	if claims, ok := GetClaims(ctx); ok && claims != nil {
		return claims.ID, true
	}
	return 0, false
}

// RequireRole — endpointga faqat ko'rsatilgan rolga ega foydalanuvchini kiritadi (masalan, RequireRole("admin")).
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetClaims(r.Context())
			if !ok || claims == nil || claims.Role != role {
				http.Error(w, "forbidden: ruxsat etilmagan", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RequireAdmin — faqat admin roli uchun ruxsat beradi (RequireRole("admin") ga asoslangan).
func RequireAdmin(next http.Handler) http.Handler {
	return RequireRole("admin")(next)
}
