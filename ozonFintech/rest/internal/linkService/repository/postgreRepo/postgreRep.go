package postgreRepo

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"ozonFintech/internal/linkService"
	"ozonFintech/pkg/logger"
)

type Postgres struct {
	logger zerolog.Logger
	pool   *pgxpool.Pool
}

func (p Postgres) SelectLink(ctx context.Context, abbreviatedLink string) (string, error) {
	conn, err := p.pool.Acquire(ctx)
	if err != nil {
		p.logger.Warn().Err(err).Msg("unable to acquire conn while select link")
		return "", err
	}
	defer conn.Release()
	q := `select original from links where abbreviated = $1`
	var originalLink string
	err = conn.QueryRow(ctx, q, abbreviatedLink).Scan(&originalLink)
	if err != nil {
		p.logger.Warn().Err(err).Msg("unable to select link while scan")
		return "", err
	}
	p.logger.Info().Msg("originalLink selected successfully.")
	return originalLink, nil
}

func (p Postgres) InsertLink(ctx context.Context, abbreviatedLink, originalLink string) error {
	conn, err := p.pool.Acquire(ctx)
	if err != nil {
		p.logger.Warn().Err(err).Msg("unable to acquire conn while insert link")
		return err
	}
	defer conn.Release()
	q := `insert into links values ($1, $2)`
	_, err = conn.Exec(ctx, q, abbreviatedLink, originalLink)
	if err != nil {
		p.logger.Warn().Err(err).Msg("unable to insert link into table")
		return err
	}
	return err
}

func NewPostgreRep(pool *pgxpool.Pool) linkService.Repo {
	return &Postgres{
		logger: logger.GetLogger(),
		pool:   pool,
	}
}
