package commands

import (
	"fmt"
	"log"

	"github.com/mattermost/mattermost-server/v6/model"
)

func Hello(client *model.Client4, post *model.Post) {
	user, _, err := client.GetUser(post.UserId, "")
	username := "user"
	if err == nil {
		username = user.Username
	}

	reply := &model.Post{
		ChannelId: post.ChannelId,
		Message:   fmt.Sprintf("Привет, @%s!", username),
	}
	_, _, err = client.CreatePost(reply)
	if err != nil {
		log.Printf("Ошибка отправки ответа: %v", err)
	}
}
