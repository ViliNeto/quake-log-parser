package restapi

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"../structs"
)

var parsedJSON []map[string]structs.GameDetail

var (
	// flagPort is the open port the application listens on
	flagPort = flag.String("port", "9000", "Port to listen on")
)

type bodyStruct struct {
	Game string
}

//StartServer is responsible to start and respond all REST requests
func StartServer(intf []map[string]structs.GameDetail) {
	//Creating REST API's
	parsedJSON = intf
	mux := http.NewServeMux()
	mux.HandleFunc("/get", getHandler)
	mux.HandleFunc("/post", postHandler)

	log.Printf("listening on port %s", *flagPort)
	log.Fatal(http.ListenAndServe(":"+*flagPort, mux))
}

// getHandler handles the index route
func getHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		json.NewEncoder(w).Encode(parsedJSON)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// PostHandler converts post request body to string
func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		var ok bool
		var value structs.GameDetail

		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		} else {
			var bd bodyStruct
			err = json.Unmarshal(body, &bd)

			for i := range parsedJSON {

				value, ok = parsedJSON[i][bd.Game]
				//IF there is no key related to a playerStats name
				if ok {
					json.NewEncoder(w).Encode(value)
					break
				}
			}
			if !ok {
				http.Error(w, "Game not found", http.StatusNotFound)
			}

		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
