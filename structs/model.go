package structs

//GameDetail -> Struct related to a game statistics [Detail]
type GameDetail struct {
	TotalKills int            `json:"total_kills"`
	Players    []string       `json:"players"`
	PlayerKill map[string]int `json:"kills"`
}
