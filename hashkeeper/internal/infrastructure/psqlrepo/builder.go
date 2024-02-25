package psqlrepo

import (
	"hashkeeper/internal/domain/datahash"

	sq "github.com/Masterminds/squirrel"
)

const _table = "string_hashes"

func prepareInsertHashes(hashes []datahash.HashContent) (string, []interface{}, error) {
	psqlSq := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	rawQuery := psqlSq.Insert(_table).
		Columns("hash")

	for _, hash := range hashes {
		rawQuery = rawQuery.Values(hash)
	}

	rawQuery = rawQuery.Suffix(`
		ON CONFLICT DO NOTHING
	`)

	return rawQuery.ToSql()
}

func prepareFindHashesByContent(hashes []datahash.HashContent) (string, []interface{}, error) {
	psqlSq := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	rawQuery := psqlSq.Select("id, hash").
		From(_table).
		Where(sq.Eq{"hash": hashes})

	return rawQuery.ToSql()
}

func prepareFindHashesByID(ids []datahash.HashID) (string, []interface{}, error) {
	psqlSq := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	rawQuery := psqlSq.Select("id, hash").
		From(_table).
		Where(sq.Eq{"id": ids})

	return rawQuery.ToSql()
}
