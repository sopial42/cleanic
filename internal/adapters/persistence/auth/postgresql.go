package persistence

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"

	utils "github.com/sopial42/cleanic/internal/adapters/rest/utils/jwt"
	"github.com/sopial42/cleanic/internal/services/auth"
)

type pgPersistence struct {
	clientDB *bun.DB
}

func NewPGClient(client *bun.DB) auth.Persistence {
	return &pgPersistence{clientDB: client}
}

func (p *pgPersistence) RegisterRefreshToken(ctx context.Context, token utils.RefreshToken) (utils.RefreshToken, error) {
	tokenDAO := fromTokenToTokenDAO(token)
	_, err := p.clientDB.NewInsert().
		Model(&tokenDAO).
		Returning("*").
		Exec(ctx)
	if err != nil {
		return utils.RefreshToken{}, fmt.Errorf("err insert: %w", err)
	}

	// ID == 0 means that the insert failed
	if len(tokenDAO.ID) == 0 {
		return utils.RefreshToken{}, fmt.Errorf("unable to create a new token: %+v", tokenDAO)
	}

	return fromTokenDAOToToken(tokenDAO)
}
