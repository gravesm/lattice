package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/urfave/cli"
)

type Env struct {
	tiler  Tiler
	layers LayerRepository
}

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name: "run",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "address",
					Value: ":6000",
					Usage: "Address to listen on",
				},
				cli.StringFlag{
					Name:  "db-hostname",
					Value: "localhost",
					Usage: "Database hostname",
				},
				cli.StringFlag{
					Name:  "db-name",
					Value: "postgres",
					Usage: "Database name",
				},
				cli.StringFlag{
					Name:  "db-password",
					Value: "",
					Usage: "Database password",
				},
				cli.IntFlag{
					Name:  "db-port",
					Value: 5432,
					Usage: "Database port",
				},
				cli.StringFlag{
					Name:  "db-username",
					Value: "postgres",
					Usage: "Database user",
				},
				cli.StringFlag{
					Name:  "layers",
					Usage: "Filename of JSON layer list",
				},
			},
			Action: func(c *cli.Context) error {
				var layers LayerRepository
				conn := fmt.Sprintf("host=%s port=%d user=%s password='%s' dbname=%s "+
					"sslmode=disable", c.String("db-hostname"), c.Int("db-port"),
					c.String("db-username"), c.String("db-password"), c.String("db-name"))
				db := DB(conn)
				defer db.Close()
				repo := DbMvtRepository{conn: db}
				tiler := MvtTiler{repo: repo}
				if c.String("layers") != "" {
					f, err := os.Open(c.String("layers"))
					if err != nil {
						log.Fatal(err)
					}
					layers = NewJsonRepo(f)
					f.Close()
				} else {
					layers = &MemLayerRepository{}
				}
				env := Env{tiler, layers}
				r := mux.NewRouter()
				r.HandleFunc("/{layer:[a-z]+}/{x:[0-9]+}/{y:[0-9]+}/{z:[0-9]+}", env.getTile)
				http.Handle("/", r)
				srv := &http.Server{
					Handler:      r,
					Addr:         c.String("address"),
					WriteTimeout: 15 * time.Second,
					ReadTimeout:  15 * time.Second,
				}
				log.Fatal(srv.ListenAndServe())

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func (env *Env) getTile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	x, _ := strconv.Atoi(vars["x"])
	y, _ := strconv.Atoi(vars["y"])
	z, _ := strconv.Atoi(vars["z"])
	l := env.layers.Find(vars["layer"])
	mvt, err := env.tiler.GetTile(x, y, z, l)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", err)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/x-protobuf")
	w.Write(mvt)
}
