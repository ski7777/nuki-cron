package config

type Team struct {
	Name        string       `json:"name"`
	Permissions []string     `json:"permissions"`
	Schedules   []Schedule   `json:"schedules"`
	ExtraEvents []ExtraEvent `json:"extra_events"`
}
