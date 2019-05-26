package main

type TileRepository interface {
	Tile(w, s, e, n float64, layer Layer) ([]byte, error)
}

type Tiler interface {
	GetTile(x, y, z int, layer Layer) ([]byte, error)
}

type MvtTiler struct {
	repo TileRepository
}

func (t MvtTiler) GetTile(x, y, z int, layer Layer) ([]byte, error) {
	w, s, e, n := boundsFromTile(x, y, z)
	return t.repo.Tile(w, s, e, n, layer)
}

type LayerRepository interface {
	Find(id string) Layer
	Save(layer Layer)
}

type Layer struct {
	Id     string
	Srid   int
	Table  string
	Column string
}
