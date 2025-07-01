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

type Handler struct {
	client *traq.APIClient
	token  string
}

func NewHandler(token string, client *traq.APIClient) *Handler {
	return &Handler{
		token:  token,
		client: client,
	}
}

func (h *Handler) makeAuth(ctx context.Context) context.Context {
	return context.WithValue(ctx, traq.ContextAccessToken, h.token)
}

func main() {
	token := os.Getenv("TRAQ_TOKEN")
	client := traq.NewAPIClient(traq.NewConfiguration())

	handler := NewHandler(token, client)

	e := echo.New()

	e.GET("/api/ubugoe/:userId", handler.GetMessages)

	e.Start(":8080")
}

func (h *Handler) GetMessages(c echo.Context) error {
	auth := h.makeAuth(c.Request().Context())
	userID := c.Param("userId")
	fmt.Println("User ID:", userID)

	channelList, _, err := h.client.ChannelApi.GetChannels(auth).Path("gps/times/" + userID).Execute()
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

	messages, _, _ := h.client.MessageApi.GetMessages(auth, channelID).
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
}
