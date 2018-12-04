package structs

//Game -> Struct related to a game statistics
type Game struct {
	Gamenumber map[string]GameDetail
}

//GameDetail -> Struct related to a game statistics [Detail]
type GameDetail struct {
	TotalKills int            `json:"total_kills,omitempty"`
	Players    []string       `json:"players,omitempty"`
	PlayerKill map[string]int `json:"kills,omitempty"`
}

//PlayerStats -> Struct resposible for recording a player kill performance or player id and name -> map[PlayerName]kills || map[PlayerName]ConnectionID
type PlayerStats struct {
	PlayerID   map[string]int //NAME ID
	PlayerStat int            //kills
}
