package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// Config will be used to store the configuration values
type Config struct {
	DBLocation string     `json:"database_location"`
	Port       int        `json:"port"`
	Currencies []Currency `json:"currencies"`
}

func main() {

	configFile := os.Args[1:]

	fmt.Println(configFile)
	// Load Configuration
	values, err := LoadConfig(configFile[0])
	if err != nil {
		log.Fatalf("Problem loading configuration with Error: %v", err)
	}

	a := App{}

	a.Initialize(values)

	path := fmt.Sprintf(":%d", values.Port)
	a.Run(path)
}

// LoadConfig will load configuation variables
func LoadConfig(fileName string) (Config, error) {
	var c Config
	var m Currency

	if fileName == "" || fileName == "default" {

		c.Port = 8080
		c.DBLocation = "./bank.db"

		m.CurrencyCode = "USD"
		m.DecimalPlaces = 2
		m.ActiveSaturday = false
		m.CurrencyTimeZone = "NYC"
		m.ReconTime = "16:00"
		c.Currencies = append(c.Currencies, m)

		m.CurrencyCode = "EUR"
		m.DecimalPlaces = 2
		m.ActiveSaturday = false
		m.CurrencyTimeZone = "NYC"
		m.ReconTime = "21:00"
		c.Currencies = append(c.Currencies, m)

		return c, nil
	}

	// Open, Read File and Unmarshall it into json
	cFile, err := os.Open(fileName)
	if err != nil {
		return c, err
	}
	defer cFile.Close()
	jsonParser := json.NewDecoder(cFile)
	if err = jsonParser.Decode(&c); err != nil {
		return c, err
	}

	// TODO:   Log config file

	return c, nil
}
