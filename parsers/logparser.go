package parsers

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

var userRegex, killRegex, initRegex, shutdownRegex *regexp.Regexp

const gameDescription = "game_"

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

//ParsetoJson -> Function that reads a file (game.log) and convert it to json
func ParsetoJson() (err error) {
	err = initRegexCompile()
	if err != nil {
		return err
	}

	file, err := os.Open("./games.log")
	defer file.Close()
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parseLine(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	//err = writeFile(games)

	return
}

func parseLine(text string) {
	if strings.Contains(text, "InitGame") {
		println("InitGame")
	}
	if strings.Contains(text, "ClientUserinfoChanged") {
		println("ClientUserinfoChanged")
	}
	if strings.Contains(text, "Kill:") {
		println("Kill")
	}
	if strings.Contains(text, "ShutdownGame") {
		println("ShutdownGame")
	}
}
