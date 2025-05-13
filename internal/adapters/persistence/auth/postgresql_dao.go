package persistence

import (
	"time"

	"github.com/google/uuid"

	user "github.com/sopial42/cleanic/internal/domains/user"
	uPersistence "github.com/sopial42/cleanic/internal/adapters/persistence/user"
	"github.com/sopial42/cleanic/internal/adapters/rest/utils/jwt"
	"github.com/uptrace/bun"
)

type tokenDAO struct {
	bun.BaseModel `bun:"table:refresh_token"`

	ID        uuid.UUID             `bun:"id,pk,autoincrement"`
	UserID    int64                 `bun:"user_id,notnull"`
	User      *uPersistence.UserDAO `bun:"rel:belongs-to,join:user_id=id"`
	ExpiresAt time.Time             `bun:"expires_at,notnull"`
	IssuedAt  time.Time             `bun:"issued_at,notnull"`
}

func fromTokenClaimsToTokenDAO(claims jwt.RefreshTokenClaims) tokenDAO {
	issueTime := time.Unix(claims.IssuedAt, 0)
	expireTime := time.Unix(claims.ExpiresAt, 0)

	return tokenDAO{
		ID:        claims.ID,
		UserID:    int64(claims.Subject),
		IssuedAt:  issueTime,
		ExpiresAt: expireTime,
	}
}

func fromTokenDAOToTokenClaims(tokenDAO *tokenDAO) jwt.RefreshTokenClaims {
	return jwt.RefreshTokenClaims{
		ID:        tokenDAO.ID,
		Subject:   user.ID(tokenDAO.UserID),
		IssuedAt:  tokenDAO.IssuedAt.Unix(),
		ExpiresAt: tokenDAO.ExpiresAt.Unix(),
	}
}
