package main

import (
	"./parsers"
	"./restapi"
)

func main() {
	//Parsing Games.log to JSON format
	parsedJSON, err := parsers.ParsetoJSON()
	if err != nil {
		println("Error parsing the game log file")
	} else {
		restapi.StartServer(parsedJSON)
	}

}
