package restapi

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"../structs"
)

//parsedJSON is a map that is going to read the JSON representation of the Game log
var parsedJSON []map[string]structs.GameDetail

var (
	//flagPort is the open port the application listens on
	flagPort = flag.String("port", "9000", "Port to listen on")
)

//bodyStruct is the representation of what we are expecting when the user sends a POST to our server eg.: {"game":"game_1"}
type bodyStruct struct {
	Game string
}

//StartServer is responsible to start and respond all REST requests
func StartServer(mappedGameLog []map[string]structs.GameDetail) {
	//initializing parsedJSON receiving the content from the parser
	parsedJSON = mappedGameLog

	//godoc: ServeMux is an HTTP request multiplexer.
	// It matches the URL of each incoming request against a list of registered patterns and calls the handler for the pattern that most closely matches the URL.
	mux := http.NewServeMux()
	//Pointing the GET method to func getHandler
	mux.HandleFunc("/get", getHandler)
	//Pointing the GET method to func postHandler
	mux.HandleFunc("/post", postHandler)
	//prints port that the app is listening on
	log.Printf("listening on port %s", *flagPort)
	//Start and prepare for a Fatal error for the port and Mux initialized before
	log.Fatal(http.ListenAndServe(":"+*flagPort, mux))
}

//getHandler handles the localhost:9000/get route
func getHandler(w http.ResponseWriter, r *http.Request) {
	//Check if the request method is a GET one
	if r.Method == "GET" {
		json.NewEncoder(w).Encode(parsedJSON)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

//PostHandler handles the localhost:9000/post route
func postHandler(w http.ResponseWriter, r *http.Request) {
	//Check if the request method is a POST one
	if r.Method == "POST" {
		//Reads all request body and checks for any error (on reading process)
		body, err := ioutil.ReadAll(r.Body)
		//In case some error occurs during the reading process
		if err != nil {
			//In case some error is explicited call HTTP ERROR with code 500
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
		} else { //In case there is no error on the reading body process
			//Initialize a body struct receiving some JSON similar to -> {"game":"game_1"}
			var bd bodyStruct
			//Tries to unmarshal (convert) the byte request to JSON format
			err = json.Unmarshal(body, &bd)
			//In case some error occurs during the unmarshaling process
			if err != nil {
				//In case some error is explicited call HTTP ERROR with code 500
				http.Error(w, "Error reading request body", http.StatusInternalServerError)
			}
			//Regex that gets only numeric chars
			gameRegex, err := regexp.Compile(`\d`)
			//If occurs some error during the regex compile process stop the followings actions
			if err != nil {
				return
			}
			//In case on the post method the sended body contais "game_" on the game key as value
			if strings.Contains(bd.Game, "game_") {
				//Extract the game number
				stringJSONPosition := gameRegex.FindStringSubmatch(bd.Game)
				//Subtract the game number by 1 because the first game starts with 1
				intPosition, err := strconv.Atoi(stringJSONPosition[0])
				//if there is some error when converting string to int
				if err != nil {
					//If some error occurs runs a Bad Request
					http.Error(w, "Bad request", http.StatusBadRequest)
				} else {
					//In case everything runs well print the JSON that reffers to the requested game
					json.NewEncoder(w).Encode(parsedJSON[intPosition-1][bd.Game])
				}

			} else { //In case doesn't find "game_" in the request body
				http.Error(w, "Bad request", http.StatusBadRequest)
			}
		}
	} else { //In case the request Method wasn't POST
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

//OLD FUNCTION THAT WAS USED TO FIND THE GAME KEY SUBSTITUTED TO REGEX SOLUTION TO SPEED UP THE SEARCHING PERFORMANCE
//Runs every position on the
// for i := range parsedJSON {
// 	value, ok = parsedJSON[i][bd.Game]
// 	//IF there is no key related to a playerStats name
// 	if ok {
// 		json.NewEncoder(w).Encode(value)
// 		break
// 	}
// }
// if !ok {
// 	http.Error(w, "Game not found", http.StatusNotFound)
// }
