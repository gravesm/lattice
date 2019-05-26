package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DbRepository struct {
	conn *sqlx.DB
}

type DbMvtRepository DbRepository

func (db DbMvtRepository) Tile(w, s, e, n float64, layer Layer) ([]byte, error) {
	var tile []byte
	sql := `
		SELECT ST_AsMVT(q, NULL, 4096, 'geom')
		FROM (
			SELECT
				ST_AsMVTGeom(
					ST_Transform(%s, 3857),
					ST_MakeEnvelope($1, $2, $3, $4, 3857),
					4096,
					256,
					true
				) AS geom
			FROM %s
		) AS q`
	err := db.conn.Get(&tile, fmt.Sprintf(sql, layer.Column, layer.Table), w, s, e, n)
	return tile, err
}

func DB(dsName string) *sqlx.DB {
	db, _ := sqlx.Open("postgres", dsName)
	return db
}

type MemLayerRepository struct {
	sync.RWMutex
	layers map[string]Layer
}

func (repo *MemLayerRepository) Find(id string) Layer {
	repo.RLock()
	l := repo.layers[id]
	repo.RUnlock()
	return l
}

func (repo *MemLayerRepository) Save(layer Layer) {
	repo.Lock()
	repo.layers[layer.Id] = layer
	repo.Unlock()
}

func NewJsonRepo(reader io.Reader) *MemLayerRepository {
	repo := &MemLayerRepository{
		layers: make(map[string]Layer),
	}
	dec := json.NewDecoder(reader)
	for {
		var l Layer
		if err := dec.Decode(&l); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		repo.Save(l)
	}
	return repo
}
