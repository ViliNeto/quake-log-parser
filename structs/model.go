package structs

//Game -> Struct related to a game statistics
type Game struct {
	TotalKills int      `json:"total_kills"`
	Players    []string `json:"players"`
	Kills      struct {
		Players PlayerStats
	} `json:"kills"`
}

//PlayerStats -> Struct resposible for recording a player kill performance -> map[PlayerName]kills
type PlayerStats struct {
	PlayerStat map[string]int
}
