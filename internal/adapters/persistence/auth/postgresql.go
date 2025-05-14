package persistence

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"

	utils "github.com/sopial42/cleanic/internal/adapters/rest/utils/jwt"
	user "github.com/sopial42/cleanic/internal/domains/user"
	"github.com/sopial42/cleanic/internal/services/auth"
)

type pgPersistence struct {
	clientDB *bun.DB
}

func NewPGClient(client *bun.DB) auth.Persistence {
	return &pgPersistence{clientDB: client}
}

func (p *pgPersistence) StoreRefreshTokenClaims(ctx context.Context, claims utils.RefreshTokenClaims) error {
	tokenDAO := fromTokenClaimsToTokenDAO(claims)
	_, err := p.clientDB.NewInsert().
		Model(&tokenDAO).
		On("CONFLICT (user_id) DO UPDATE").
		Set("id = EXCLUDED.id, issued_at = EXCLUDED.issued_at, expires_at = EXCLUDED.expires_at").
		Returning("*").
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("unable to update token claims: %w", err)
	}

	// ID == 0 means that the insert failed
	if len(tokenDAO.ID) == 0 {
		return fmt.Errorf("unable to insert refresh token claims: %+v", tokenDAO)
	}

	return nil
}

func (p *pgPersistence) GetRefreshTokenClaimsByUserID(ctx context.Context, userID user.ID) (utils.RefreshTokenClaims, error) {
	tokenDAO := &tokenDAO{}
	err := p.clientDB.NewSelect().
		Model(tokenDAO).
		Where("user_id = ?", userID).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return utils.RefreshTokenClaims{}, fmt.Errorf("unable to get refresh token: %w", err)
	}

	if tokenDAO == nil || len(tokenDAO.ID) == 0 {
		return utils.RefreshTokenClaims{}, fmt.Errorf("no refresh token found for userID: %v", userID)
	}

	return fromTokenDAOToTokenClaims(tokenDAO), nil
}

func (p *pgPersistence) DeleteRefreshTokenClaims(ctx context.Context, userID user.ID) error {
	_, err := p.clientDB.NewDelete().
		Model((*tokenDAO)(nil)).
		Where("user_id = ?", userID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("unable to delete refresh token: %w", err)
	}

	return nil
}
