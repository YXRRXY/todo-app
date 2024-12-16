package model

import "time"

type Todo struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}
