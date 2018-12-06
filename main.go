package main

import (
	"./parsers"
	"./restapi"
)

func main() {
	//Parsing Games.log to JSON format
	parsedJSON, err := parsers.ParsetoJSON()
	//If there is some error during the parsing process
	if err != nil {
		panic("Error parsing the game log file")
	} else {
		restapi.StartServer(parsedJSON)
	}

}
