package main

import (
	"strings"
	"testing"
)

func TestJsonRepository(t *testing.T) {
	const stream = `
			{"id": "foo", "srid": 3857, "table": "foo", "column": "geom"}
			{"id": "bar", "srid": 4326, "table": "bar", "column": "the_geom"}`

	t.Run("Factory", func(t *testing.T) {
		r := NewJsonRepo(strings.NewReader(stream))
		if len(r.layers) != 2 {
			t.Error("Expected 2, got", len(r.layers))
		}
	})

	t.Run("Find", func(t *testing.T) {
		r := NewJsonRepo(strings.NewReader(stream))
		l := r.Find("bar")
		if l.Column != "the_geom" {
			t.Error("Expected the_geom, got", l.Column)
		}
	})

	t.Run("Save", func(t *testing.T) {
		r := NewJsonRepo(strings.NewReader(stream))
		r.Save(Layer{"baz", 1111, "baz", "geom"})
		l := r.Find("baz")
		if l.Srid != 1111 {
			t.Error("Expected 1111, got", l.Srid)
		}
	})
}
