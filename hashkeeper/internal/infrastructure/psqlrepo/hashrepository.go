package psqlrepo

import (
	"context"
	"hashkeeper/internal/domain/datahash"
	"hashkeeper/pkg/hashlog"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PSQLHashRepository interface {
	datahash.HashRepository

	Dial(context.Context) error
	Close() error
}

type hashRepositoryImpl struct {
	dsn    string
	client *sqlx.DB
}

type dbRow struct {
	ID   int64  `db:"id"`
	Hash string `db:"hash"`
}

var _ PSQLHashRepository = (*hashRepositoryImpl)(nil)

type HashRepositoryCfg struct {
	// Data Source Name
	DSN string `arg:"--psql-dsn,env:PSQL_DSN" default:"host=localhost port=5432 user=postgres password=postgres dbname=stringhashes sslmode=disable"`
}

func NewHashRepository(cfg HashRepositoryCfg) (PSQLHashRepository, error) {
	res := hashRepositoryImpl{
		dsn: cfg.DSN,
	}
	return &res, nil
}

func (hr hashRepositoryImpl) Store(ctx context.Context, hashes []datahash.HashContent) error {
	hashlog.LogReqID(ctx).Debug("store hashes staring")
	hashlog.LogReqID(ctx).WithField("hashes", hashes).Trace("data to store")

	query, args, err := prepareInsertHashes(hashes)
	if err != nil {
		return hashlog.WithStackErrorf("failed to prepare insert hashes query: %w", err)
	}

	_, err = hr.client.ExecContext(ctx, query, args...)
	if err != nil {
		return hashlog.WithStackErrorf("failed to execute insert hashes query: %w", err)
	}

	hashlog.LogReqID(ctx).Debug("store hashes finished successfully")

	return nil
}

func (hr hashRepositoryImpl) FindByContent(ctx context.Context, hashes []datahash.HashContent) ([]datahash.Hash, error) {
	hashlog.LogReqID(ctx).Debug("find hashes by content staring")
	hashlog.LogReqID(ctx).WithField("hashes", hashes).Trace("content to search")

	query, args, err := prepareFindHashesByContent(hashes)
	if err != nil {
		return nil, hashlog.WithStackErrorf("failed to prepare find hashes by content query: %w", err)
	}

	var rows []dbRow
	err = hr.client.SelectContext(ctx, &rows, query, args...)
	if err != nil {
		return nil, hashlog.WithStackErrorf("failed to execute find hashes by content query: %w", err)
	}

	result := make([]datahash.Hash, len(rows))
	for i, row := range rows {
		result[i] = datahash.NewHash(datahash.HashID(row.ID), datahash.HashContent(row.Hash))
	}

	hashlog.LogReqID(ctx).Debug("find hashes by content finished successfully")

	return result, nil
}

func (hr hashRepositoryImpl) FindByID(ctx context.Context, ids []datahash.HashID) ([]datahash.Hash, error) {
	hashlog.LogReqID(ctx).Debug("find hashes by id staring")
	hashlog.LogReqID(ctx).WithField("ids", ids).Trace("ids to search")

	query, args, err := prepareFindHashesByID(ids)
	if err != nil {
		return nil, hashlog.WithStackErrorf("failed to prepare find hashes by id query: %w", err)
	}

	var rows []dbRow
	err = hr.client.SelectContext(ctx, &rows, query, args...)
	if err != nil {
		return nil, hashlog.WithStackErrorf("failed to execute find hashes by id query: %w", err)
	}

	result := make([]datahash.Hash, len(rows))
	for i, row := range rows {
		result[i] = datahash.NewHash(datahash.HashID(row.ID), datahash.HashContent(row.Hash))
	}

	hashlog.LogReqID(ctx).Debug("find hashes by id finished successfully")

	return result, nil
}

func (hr *hashRepositoryImpl) Dial(ctx context.Context) error {
	db, err := sqlx.Open("postgres", hr.dsn)
	if err != nil {
		return hashlog.WithStackErrorf("can`t connect to database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return hashlog.WithStackErrorf("can`t ping database: %w", err)
	}

	hr.client = db

	return nil
}

func (hr *hashRepositoryImpl) Close() error {
	err := hr.client.Close()

	hr.client = nil

	return err
}
