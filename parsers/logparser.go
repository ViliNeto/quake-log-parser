package parsers

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"

	"../structs"
)

//Initialize all Regex's
var userRegex, killRegex, initRegex, shutdownRegex *regexp.Regexp

//Constant that will define the master key in a JSON response
const gameDescription = "game_"

//Counts the game position in the slice [can be substituted for len(slice)+1 in the nexts implementations]
var gameCount = 0

//totalKills summarize all players kills
var totalKills = 0

//simple struct that represent the content of the JSON
var gamedetail structs.GameDetail

//game is a representation of a mapping (dictionary) using as value "structs.GameDetail" and "game_" as key
var game map[string]structs.GameDetail

//Slice of game (described above)
var games []map[string]structs.GameDetail

//PlayerAux -> Auxiliar to control how many players enters in a server
var playerAux map[string]interface{}

//function that do the regex setup
func initRegexCompile() (err error) {
	//initRegex -> [0] == 0:00"
	initRegex, err = regexp.Compile(`[0-9:]{4}`)
	if err != nil {
		return
	}
	//userRegex -> [1], "---", [2], "---", [3], "---", [4] == 20:34 --- ClientUserinfoChanged --- 2 --- Isgalamido
	userRegex, err = regexp.Compile(`.*?((?:(?:[0-1][0-9])|(?:[2][0-3])|(?:[0-9])):(?:[0-5][0-9])(?::[0-5][0-9])?(?:\s\s)?).*?(ClientUserinfoChanged):\s(.)\sn\\(.*)\\t\\`)
	if err != nil {
		return
	}
	//killRegex -> [2], "---", [4], "---", [6] == <world> --- Isgalamido --- MOD_TRIGGER_HURT
	killRegex, err = regexp.Compile(`Kill:\s(.*?):\s(.*)\s\b(killed)\s(.*)\s\b(by)\s(.*)`)
	if err != nil {
		return
	}
	//shutdownRegex -> [0] == 0:00"
	shutdownRegex, err = regexp.Compile(`[0-9:]{4}`)
	if err != nil {
		return
	}
	return
}

//ParsetoJSON -> Function that reads a file (game.log) and convert it to json
func ParsetoJSON() ([]map[string]structs.GameDetail, error) {
	//Initialize the regex setups
	err := initRegexCompile()
	//in case some regex has failed in the above setup
	if err != nil {
		return games, err
	}
	//Opens games.log file
	file, err := os.Open("./games.log")
	//Close the "./games.log" when the fuction is ended
	defer file.Close()
	//In case some error occurs during the oppening stage
	if err != nil {
		return games, err
	}
	//declare a scanner that will read every single line in a file
	scanner := bufio.NewScanner(file)
	//Reads line a line to parse log file
	for scanner.Scan() {
		//in fact parses the lines to its respectives representations
		parseLine(scanner.Text())
	}
	//If some error is caught in the Scan() process
	if err := scanner.Err(); err != nil {
		return games, err
	}
	return games, err
}

//Function that in fact parse and aggregate the game log file to a single object
func parseLine(text string) {
	//Those variables will be resposible of receiving the submatch of each regex
	var userDetail, killDetail []string
	//Control of key search finding
	var ok bool

	//Does nothing yet
	// if strings.Contains(text, "InitGame") {
	// 	//
	// }

	//In case is reading a line that has "ClientUserinfoChanged" on it
	if strings.Contains(text, "ClientUserinfoChanged") {
		//Splits the line in submatches based on the regex setup
		userDetail = userRegex.FindStringSubmatch(text)
		//Search for a player already registered with same name
		_, ok = playerAux[userDetail[4]]
		//IF there is no key related to a playerStats name
		if !ok {
			//if the map is empty
			if playerAux == nil {
				//create new map with the same attributes to refer to playerAux later
				m := make(map[string]interface{})
				//Setup nil as a value (we will not use this value)
				m[userDetail[4]] = nil
				//refer playerAux to m
				playerAux = m
				//Append that new player to the main structure
				gamedetail.Players = append(gamedetail.Players, userDetail[4])
			} else { //In case the map is not empty
				//Sets the player on the list
				playerAux[userDetail[4]] = nil
				//Append that new player to the main structure
				gamedetail.Players = append(gamedetail.Players, userDetail[4])
			}
		}
	}
	//In case is reading a line that has "Kill" on it
	if strings.Contains(text, "Kill:") {
		//Splits the line in submatches based on the regex setup
		killDetail = killRegex.FindStringSubmatch(text)
		//Check if the player name is <world>
		if killDetail[2] != "<world>" {
			//Initializes "typeOfKill" which will define if was a death or a suicide kill
			var typeOfKill int
			//If the player who killed was different that the player who died
			if killDetail[2] != killDetail[4] {
				//Positive point
				typeOfKill = 1
				//Only adds a kill to total count if 1 player kill another player - in case of suicide this is not added
				totalKills++
			} else {
				//in case of suicide negative point
				typeOfKill = -1
			}
			//checks if there is the killer player on the kill feed list
			_, ok = gamedetail.PlayerKill[killDetail[2]]
			//IF there is no key related to a playerStats name
			if !ok { //if there is not the player in the list
				//If the list is null
				if gamedetail.PlayerKill == nil {
					//initialize the map with the same arguments to refer it to the main map
					m := make(map[string]int)
					//Sets the tipe of the kill and at this point is the first kill/suicide for the player
					m[killDetail[2]] = typeOfKill
					//referes "m"  to "gamedetail.PlayerKill"
					gamedetail.PlayerKill = m
				} else { //In case the structure is not null and is a new player in the list
					gamedetail.PlayerKill[killDetail[2]] = typeOfKill
				}
			} else { //in case the player is already in the list summarize the kill
				gamedetail.PlayerKill[killDetail[2]] = gamedetail.PlayerKill[killDetail[2]] + typeOfKill
			}
		}
	}
	//In case is reading a line that has "ShutdownGame" on it
	if strings.Contains(text, "ShutdownGame") {
		//Add +1 to game counting eg.: "game_1","game_2"
		gameCount++
		//Initialize the variable m which is a support map to be reffered to game later
		m := make(map[string]structs.GameDetail)
		//append the totalkill to the structure
		gamedetail.TotalKills = totalKills
		//make the key "game_" combined with the game position
		m["game_"+strconv.Itoa(gameCount)] = gamedetail
		//reffers "m" to "game"
		game = m
		//append this single element to the slice games
		games = append(games, game)

		//Reseting all variables
		game = nil
		playerAux = nil
		gamedetail = structs.GameDetail{}
		totalKills = 0
	}
}
