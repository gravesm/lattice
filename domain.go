package main

type TileRepository interface {
	Tile(w, s, e, n float64) ([]byte, error)
}

type Tiler interface {
	GetTile(x, y, z int) ([]byte, error)
}

type MvtTiler struct {
	repo TileRepository
}

func (t MvtTiler) GetTile(x, y, z int) ([]byte, error) {
	w, s, e, n := boundsFromTile(x, y, z)
	return t.repo.Tile(w, s, e, n)
}
