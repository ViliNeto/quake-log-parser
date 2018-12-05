package parsers

import (
	"bufio"
	"encoding/json"
	"os"
	"regexp"
	"strconv"
	"strings"

	"../structs"
)

var userRegex, killRegex, initRegex, shutdownRegex *regexp.Regexp

const gameDescription = "game_"

var gameCount = 0

var totalKills = 0
var game map[string]structs.GameDetail
var gamedetail structs.GameDetail
var games []map[string]structs.GameDetail

//PlayerAux -> Auxiliar to control how many players enters in a server
var playerAux map[string]interface{}

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
func ParsetoJSON() (err error) {
	err = initRegexCompile()
	if err != nil {
		return err
	}
	//Opens games.log file
	file, err := os.Open("./games.log")
	defer file.Close()
	if err != nil {
		return
	}
	//Reads line a line to parse log file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parseLine(scanner.Text())
	}
	//Convert the object to JSON
	c, _ := json.Marshal(games)
	//Print json on console
	println(string(c))
	if err := scanner.Err(); err != nil {
		return err
	}
	return
}

func parseLine(text string) {
	var userDetail, killDetail []string
	var ok bool

	if strings.Contains(text, "InitGame") {
		//
	}

	if strings.Contains(text, "ClientUserinfoChanged") {
		userDetail = userRegex.FindStringSubmatch(text)
		//Search for a player already registered with same
		_, ok = playerAux[userDetail[4]]
		//IF there is no key related to a playerStats name
		if !ok {
			if playerAux == nil {
				m := make(map[string]interface{})
				m[userDetail[4]] = nil
				playerAux = m
				gamedetail.Players = append(gamedetail.Players, userDetail[4])
			} else {
				playerAux[userDetail[4]] = nil
				gamedetail.Players = append(gamedetail.Players, userDetail[4])
			}
		}
	}
	if strings.Contains(text, "Kill:") {
		killDetail = killRegex.FindStringSubmatch(text)
		//Search for a player already registered with same
		if killDetail[2] != "<world>" {
			var typeOfKill int
			if killDetail[2] != killDetail[4] {
				typeOfKill = 1
				//Only adds a kill to total count if 1 player kill another player - in case of suicide this is not added
				totalKills++
			} else {
				typeOfKill = -1
			}
			_, ok = gamedetail.PlayerKill[killDetail[2]]
			//IF there is no key related to a playerStats name
			if !ok {
				if gamedetail.PlayerKill == nil {
					m := make(map[string]int)
					m[killDetail[2]] = typeOfKill
					gamedetail.PlayerKill = m
				} else {
					gamedetail.PlayerKill[killDetail[2]] = typeOfKill
				}
			} else {
				gamedetail.PlayerKill[killDetail[2]] = gamedetail.PlayerKill[killDetail[2]] + typeOfKill
			}
		}
	}
	if strings.Contains(text, "ShutdownGame") {
		gameCount++
		m := make(map[string]structs.GameDetail)
		gamedetail.TotalKills = totalKills
		m["game_"+strconv.Itoa(gameCount)] = gamedetail
		game = m
		games = append(games, game)

		//Reseting all variables
		game = nil
		playerAux = nil
		gamedetail = structs.GameDetail{}
		totalKills = 0
	}
}
