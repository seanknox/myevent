package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/seanknox/myevent/lib/persistence/dblayer"
)

var (
	DBTypeDefault       = dblayer.DBTYPE("mongodb")
	DBConnectionDefault = "mongodb://127.0.0.1:27017"
	RestfulEPDefault    = "localhost:8182"
	RestfulTLSEPDefault = "localhost:8183"
)

type ServiceConfig struct {
	DatabaseType       dblayer.DBTYPE `json:"databasetype"`
	DBConnection       string         `json:"dbconnection"`
	RestfulEndpoint    string         `json:"restfulapi_endpoint"`
	RestfulTLSEndpoint string         `json:"restfultlsapi_endpoint"`
}

func ExtractConfiguration(filename string) (ServiceConfig, error) {
	conf := ServiceConfig{
		DBTypeDefault,
		DBConnectionDefault,
		RestfulEPDefault,
		RestfulTLSEPDefault,
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Configuration file not found. Continuing with default values.")
		return conf, err
	}

	err = json.NewDecoder(file).Decode(&conf)
	return conf, err
}