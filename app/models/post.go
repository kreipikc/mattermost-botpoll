package models

type Post struct {
	Id        string `json:"id"`
	ChannelId string `json:"channel_id"`
	UserId    string `json:"user_id"`
	Message   string `json:"message"`
	CreateAt  int64  `json:"create_at"`
}
