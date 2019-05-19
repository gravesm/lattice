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
	tiler Tiler
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
			},
			Action: func(c *cli.Context) error {
				conn := fmt.Sprintf("host=%s port=%d user=%s password='%s' dbname=%s "+
					"sslmode=disable", c.String("db-hostname"), c.Int("db-port"),
					c.String("db-username"), c.String("db-password"), c.String("db-name"))
				db := DB(conn)
				defer db.Close()
				repo := DbMvtRepository{conn: db}
				tiler := MvtTiler{repo: repo}
				env := Env{tiler}
				r := mux.NewRouter()
				r.HandleFunc("/{x:[0-9]+}/{y:[0-9]+}/{z:[0-9]+}", env.getTile)
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
	mvt, err := env.tiler.GetTile(x, y, z)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "%v", err)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/x-protobuf")
	w.Write(mvt)
}
