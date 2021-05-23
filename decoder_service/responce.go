package main

import "time"

type UserToken struct {
	Token string `json:"token"`
}

type UserName struct {
	Username string `json:"name"`
}

type Audio struct {
	Username string `json:"name"`
	Token    string `json:"token"`
}

type Status struct {
	Name   string    `json:"name"`
	Status string    `json:"status"`
	Date   time.Time `json:"published"`
	Token  string    `json:"token"`
}

type CodeRequest struct {
	Token string `json:"token"`
	URL   string `json:"url"`
	Key   string `json:"key"`
}
