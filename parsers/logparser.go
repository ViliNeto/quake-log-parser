package parsers

import (
	"bufio"
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
var playerStats structs.PlayerStats
var playerKills structs.PlayerKills

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
			m := make(map[string]int)
			id, _ := strconv.Atoi(userDetail[3])

			m[userDetail[4]] = id
			playerStats.PlayerID = m

			game.Players = append(game.Players, userDetail[4])
		}
		// fmt.Printf("%v", game.Players)
	}
	if strings.Contains(text, "Kill:") {
		killDetail = killRegex.FindStringSubmatch(text)
		//Search for a player already registered with same
		if killDetail[2] != "<world>" {
			_, ok = playerKills.PlayerKill[killDetail[2]]
			//IF there is no key related to a playerStats name
			if !ok {
				if playerKills.PlayerKill == nil {
					m := make(map[string]int)
					m[killDetail[2]] = 1
					playerKills.PlayerKill = m
					game.PlayerKill = playerKills
				} else {
					playerKills.PlayerKill[killDetail[2]] = 1
					game.PlayerKill = playerKills
				}
			} else {
				playerKills.PlayerKill[killDetail[2]] = playerKills.PlayerKill[killDetail[2]] + 1
				game.PlayerKill = playerKills
			}
		}
	}
	if strings.Contains(text, "ShutdownGame") {
		// println("ShutdownGame")
		fmt.Printf("%v", game)
		return true
	}

	return false
}
