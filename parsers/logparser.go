package parsers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"../structs"
)

var userRegex, killRegex, initRegex, shutdownRegex *regexp.Regexp

const gameDescription = "game_"

var games []structs.Game
var game structs.Game
var gamedetail structs.GameDetail
var playerStats structs.PlayerStats

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
	control := false
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
		control = parseLine(scanner.Text())
		if control {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	//err = writeFile(games)

	return
}

func parseLine(text string) bool {
	var userDetail, killDetail []string
	var ok bool

	if strings.Contains(text, "InitGame") {
		// println("InitGame")
	}

	if strings.Contains(text, "ClientUserinfoChanged") {
		userDetail = userRegex.FindStringSubmatch(text)
		//Search for a player already registered with same
		_, ok = playerStats.PlayerID[userDetail[4]]
		//IF there is no key related to a playerStats name
		if !ok {
			if playerStats.PlayerID == nil {
				m := make(map[string]int)
				m[userDetail[4]], _ = strconv.Atoi(userDetail[3])
				playerStats.PlayerID = m
				gamedetail.Players = append(gamedetail.Players, userDetail[4])
			} else {
				playerStats.PlayerID[userDetail[4]], _ = strconv.Atoi(userDetail[3])
				gamedetail.Players = append(gamedetail.Players, userDetail[4])
			}
		}
		// fmt.Printf("%v", game.Players)
	}
	if strings.Contains(text, "Kill:") {
		killDetail = killRegex.FindStringSubmatch(text)
		//Search for a player already registered with same
		if killDetail[2] != "<world>" {
			_, ok = gamedetail.PlayerKill[killDetail[2]]
			//IF there is no key related to a playerStats name
			if !ok {
				if gamedetail.PlayerKill == nil {
					m := make(map[string]int)
					m[killDetail[2]] = 1
					gamedetail.PlayerKill = m
				} else {
					gamedetail.PlayerKill[killDetail[2]] = 1
				}
			} else {
				gamedetail.PlayerKill[killDetail[2]] = gamedetail.PlayerKill[killDetail[2]] + 1
			}
		}
	}
	if strings.Contains(text, "ShutdownGame") {
		// println("ShutdownGame").
		m := make(map[string]structs.GameDetail)
		m["game_1"] = gamedetail
		game.Gamenumber = m
		fmt.Printf("%v", game)
		b, _ := json.Marshal(game.Gamenumber)

		println(string(b))

		game = structs.Game{}
		return true
	}

	return false
}
