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
		SELECT ST_AsMVT(q, NULL, 4096, 'geom')
		FROM (
			SELECT
				ST_AsMVTGeom(
					ST_Transform(geom, 3857),
					ST_MakeEnvelope($1, $2, $3, $4, 3857),
					4096,
					256,
					true
				) AS geom
			FROM states
		) AS q`
	err := db.conn.Get(&tile, sql, w, s, e, n)
	return tile, err
}

func DB(dsName string) *sqlx.DB {
	db, _ := sqlx.Open("postgres", dsName)
	return db
}
