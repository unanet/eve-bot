package slack

import "time"

type EveEvent struct {
	Group     string    `json:"group"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
