package request

import "time"

type UserToken struct {
	Token string `json:"token"`
}

type UserName struct {
	Username string `json:"name"`
}

type UserNameUserToken struct {
	Username string `json:"name"`
	Token    string `json:"token"`
}

type AudioStatus struct {
	Name   string    `json:"name"`
	Status string    `json:"status"`
	Date   time.Time `json:"published"`
	Token  string    `json:"token"`
}

type UserTokenAudioToken struct {
	Token string `json:"token"`
	Audio string `json:"audio"`
}

type UserNameAudioToken struct {
	Username string `json:"name"`
	Token    string `json:"audio"`
}

// type CodeRequest struct {
// 	Token string `json:"token"`
// 	URL   string `json:"url"`
// 	Key   string `json:"key"`
// }
