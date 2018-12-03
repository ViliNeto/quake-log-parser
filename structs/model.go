package structs

//Game -> Struct related to a game statistics
type Game struct {
	TotalKills int         `json:"total_kills"`
	Players    []string    `json:"players"`
	PlayerKill PlayerKills `json:"kills"`
}

//PlayerStats -> Struct resposible for recording a player kill performance or player id and name -> map[PlayerName]kills || map[PlayerName]ConnectionID
type PlayerStats struct {
	PlayerID   map[string]int //NAME ID
	PlayerStat int            //kills
}

//PlayerKills -> record player kills in a match
type PlayerKills struct {
	PlayerKill map[string]int //kills
}
