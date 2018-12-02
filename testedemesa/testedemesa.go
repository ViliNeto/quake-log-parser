package testedemesa

import "regexp"

func main() {
	str1 := `20:34 ClientUserinfoChanged: 2 n\Isgalamido\t\0\model\xian/default\hmodel\xian/default\g_redteam\g_blueteam\c1\4\c2\5\hc\100\w\0\l\0\tt\0\tl\0`
	// userRegEx, err := regexp.Compile(`^.{0,10}([a-z A-Z][^:]*)`)
	userRegex, err := regexp.Compile(`.*?((?:(?:[0-1][0-9])|(?:[2][0-3])|(?:[0-9])):(?:[0-5][0-9])(?::[0-5][0-9])?(?:\s\s)?).*?(ClientUserinfoChanged):\s(.)\sn\\(.*)\\t\\`)
	if err != nil {
		println(err)
	}
	test := userRegex.FindStringSubmatch(str1)
	println(test[1], "---", test[2], "---", test[3], "---", test[4])

	str2 := ` 20:54 Kill: 1022 2 22: <world> killed Isgalamido by MOD_TRIGGER_HURT`
	killRegex, err := regexp.Compile(`Kill:\s(.*?):\s(.*)\s\b(killed)\s(.*)\s\b(by)\s(.*)`)
	if err != nil {
		println(err)
	}
	test2 := killRegex.FindStringSubmatch(str2)
	println(test2[2], "---", test2[4], "---", test2[6])

	str3 := `  0:00 InitGame: \sv_floodProtect\1\sv_maxPing\0\sv_minPing\0\sv_maxRate\10000\sv_minRate\0\sv_hostname\Code Miner`
	initRegex, err := regexp.Compile(`[0-9:]{4}`)
	if err != nil {
		println(err)
	}
	test3 := initRegex.FindStringSubmatch(str3)
	println(test3[0], "--- Init Time")

	str4 := `                      1:47 ShutdownGame:`
	shutdownRegex, err := regexp.Compile(`[0-9:]{4}`)
	if err != nil {
		println(err)
	}
	test4 := shutdownRegex.FindStringSubmatch(str4)
	println(test4[0], "--- Shutdown Time")

	//SEPARAR IDENTIFICADOR PARA FAZER O SELECT SWITCH
}
