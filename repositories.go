package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DbRepository struct {
	conn *sqlx.DB
}

type DbMvtRepository DbRepository

func (db DbMvtRepository) Tile(w, s, e, n float64) ([]byte, error) {
	var tile []byte
	sql := `
		SELECT ST_AsMVT(q, 'states', 4096, 'geom') mvt
		FROM (
			SELECT
				ST_AsMVTGeom(
					geom,
					ST_MakeEnvelope($1, $2, $3, $4, 4326),
					4096,
					80,
					false
				) geom
			FROM states
		) q`
	err := db.conn.Get(&tile, sql, w, s, e, n)
	return tile, err
}

func DB(dsName string) *sqlx.DB {
	db, _ := sqlx.Open("postgres", dsName)
	return db
}
