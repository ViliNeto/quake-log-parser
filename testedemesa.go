package main

import "regexp"

func main() {
	str1 := `20:34 ClientUserinfoChanged: 2 n\Isgalamidoasas\t\0\model\xian/default\hmodel\xian/default\g_redteam\g_blueteam\c1\4\c2\5\hc\100\w\0\l\0\tt\0\tl\0`
	// userRegEx, err := regexp.Compile(`^.{0,10}([a-z A-Z][^:]*)`)
	userRegex, err := regexp.Compile(`(ClientUserinfoChanged):\s(.)\sn\\(.*)\\t\\`)
	if err != nil {
		println(err)
	}
	test := userRegex.FindStringSubmatch(str1)
	println(test[1], "---", test[2], "---", test[3])

	str2 := ` 20:54 Kill: 1022 2 22: <world> killed Isgalamido by MOD_TRIGGER_HURT`
	killRegex, err := regexp.Compile(`Kill:\s(.*?):\s(.*)\s\b(killed)\s(.*)\s\b(by)\s(.*)`)
	if err != nil {
		println(err)
	}
	test2 := killRegex.FindStringSubmatch(str2)
	println(test2[2], "---", test2[4], "---", test2[6])

}
