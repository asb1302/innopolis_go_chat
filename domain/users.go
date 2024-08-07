package domain

import "time"

type UserDialog struct {
	UIDs     []string  `json:"uid"`
	ChID     string    `json:"ch_id"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Mid    string
	Body   string    `json:"body"`
	TDate  time.Time `json:"t_date"`
	FromID string    `json:"from_id"`
}

type MessageRequest struct {
	Msg  Message `json:"msg"`
	Type string  // delete / upd / add
	ChID string  `json:"ch_id"`
}
