package persistence

import (
	"time"

	uPersistence "github.com/sopial42/cleanic/internal/adapters/persistence/user"
	"github.com/sopial42/cleanic/internal/adapters/rest/utils/jwt"
	"github.com/sopial42/cleanic/internal/domains/user"
	"github.com/uptrace/bun"
)

type tokenDAO struct {
	bun.BaseModel `bun:"table:refresh_token"`

	ID        string                `bun:"id,pk,autoincrement"`
	UserID    int64                 `bun:"user_id,notnull"`
	User      *uPersistence.UserDAO `bun:"rel:belongs-to,join:user_id=id"`
	Token     string                `bun:"token,notnull"`
	ExpiresAt time.Time             `bun:"expires_at,notnull"`
	CreatedAt time.Time             `bun:"created_at,timestamp"`
	UpdatedAt time.Time             `bun:"updated_at,timestamp"`
}

func fromTokenToTokenDAO(token jwt.RefreshToken) tokenDAO {
	return tokenDAO{
		UserID:    int64(token.Claims.Subject),
		Token:     string(token.Token),
		ExpiresAt: token.Claims.ExpiresAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func fromTokenDAOToToken(dao tokenDAO) (jwt.RefreshToken, error) {
	return jwt.RefreshToken{
		Token: jwt.SignedRefreshToken(dao.Token),
		Claims: jwt.RefreshTokenClaims{
			Subject:   user.ID(dao.UserID),
			ExpiresAt: dao.ExpiresAt,
			IssuedAt:  dao.CreatedAt,
		},
	}, nil
}
