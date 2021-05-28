package request

import "time"

// UserToken .
type UserToken struct {
	Token string `json:"token"`
}

// AudioToken .
type AudioToken struct {
	Token string `json:"token"`
}

// UserName .
type UserName struct {
	Username string `json:"name"`
}

// UserNameUserToken .
type UserNameUserToken struct {
	Username string `json:"name"`
	Token    string `json:"token"`
}

// AudioStatus .
type AudioStatus struct {
	Name   string    `json:"name"`
	Status string    `json:"status"`
	Date   time.Time `json:"published"`
	Token  string    `json:"token"`
}

// UserTokenAudioToken .
type UserTokenAudioToken struct {
	Token string `json:"token"`
	Audio string `json:"audio"`
}

// UserNameAudioToken .
type UserNameAudioToken struct {
	Username string `json:"name"`
	Token    string `json:"audio"`
}

// AudioInit .
type AudioInit struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Key  string `json:"key"`
	KID  string `json:"kid,omitempty"`
}

// UserAudio .
type UserAudio struct {
	Username  string `json:"name"`
	AudioInit `json:"audio"`
}

// UserTokenAudio .
type UserTokenAudio struct {
	Token     string `json:"token"`
	AudioInit `json:"audio"`
}
