package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// Config will be used to store the configuration values
type Config struct {
	DBLocation string     `json:"database_location"`
	Port       int        `json:"port"`
	TLS        bool       `json:"tls"`
	Currencies []Currency `json:"currencies"`
}

const version = "v0.4-beta"

func main() {

	var conf string
	flag.StringVar(&conf, "c", "default", "(short-hand) configuration file containing setup information")
	flag.StringVar(&conf, "config", "default", "configuration file containing setup information")
	var v bool
	flag.BoolVar(&v, "v", false, "(short-hand) application version")
	flag.BoolVar(&v, "version", false, "application version")
	flag.Parse()

	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("    fc2-mock-bank -c [<file name>|default] \n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if v {
		fmt.Printf("%s version %s\n\n", os.Args[0], version)
		os.Exit(1)
	}

	log.Printf("Configuration File: %s", conf)

	// Load Configuration
	values, err := LoadConfig(conf)
	if err != nil {
		log.Fatalf("Problem loading configuration with Error: %v", err)
	}

	a := App{}

	a.Initialize(values)

	path := fmt.Sprintf(":%d", values.Port)
	a.Run(path, values.TLS)
}
