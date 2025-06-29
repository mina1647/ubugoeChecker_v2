package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/go-traq"
)

type Message struct {
	Content string `json:"message"`
	UserID  string `json:"userId"`
}

func main() {
	token := os.Getenv("TRAQ_TOKEN")

	client := traq.NewAPIClient(traq.NewConfiguration())
	auth := context.WithValue(context.Background(), traq.ContextAccessToken, token)

	e := echo.New()
	e.GET("/api/ubugoe/:userId", func(c echo.Context) error {
		userID := c.Param("userId")
		fmt.Println("User ID:", userID)

		channelList, _, err := client.ChannelApi.GetChannels(auth).Path("gps/times/" + userID).Execute()
		//_, _, err := client.UserApi.GetUser(context.Background(), userID).Execute()
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(500, "something went wrong")
		}
		if len(channelList.Public) == 0 {
			return c.JSON(404, "No channels found")
		}
		channel := channelList.Public[0]
		channelID := channel.Id

		messages, _, _ := client.MessageApi.GetMessages(auth, channelID).
			Since(time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local)).
			Limit(10).
			Order("asc").
			Execute()

		res := make([]Message, 0, len(messages))
		for _, message := range messages {

			res = append(res, Message{
				Content: message.Content,
				UserID:  message.UserId,
			})
		}

		return c.JSON(200, res)
	})

	e.Start(":8080")
}
