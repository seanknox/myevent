package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/seanknox/myevent/pkg/persistence/dblayer"
	"github.com/seanknox/myevent/pkg/rest"

	"github.com/seanknox/myevent/config"
)

func main() {
	confPath := flag.String("conf", `.\config\config.json`, "flag to set path of configuration file")
	flag.Parse()

	// extract config
	config, _ := config.ExtractConfiguration(*confPath)

	fmt.Println("Connecting to database...")
	dbhandler, err := dblayer.NewPersistenceLayer(config.DatabaseType, config.DBConnection)
	if err != nil {
		log.Fatalf("couldn't connect to database: %+v", err)
	}
	fmt.Println("Connected to DB.")

	// API start
	log.Fatal(rest.ServeAPI(config.RestfulEndpoint, dbhandler))
}
