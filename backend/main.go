package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/go-traq"

	traqwsbot "github.com/traPtitech/traq-ws-bot"
	payload "github.com/traPtitech/traq-ws-bot/payload"
)

type TimesUbugoe struct {
	Content string `json:"message"`
	UserID  string `json:"userId"`
}

type TrueUbugoe struct {
	Content   string    `json:"message"`
	Channel   string    `json:"channel"`
	CreatedAt time.Time `json:"createdAt"`
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

	bot, err := traqwsbot.NewBot(&traqwsbot.Options{
		AccessToken: token,
	})
	if err != nil {
		panic(err)
	}

	bot.OnMessageCreated(func(p *payload.MessageCreated) {
		log.Println("Received MESSAGE_CREATED event: " + p.Message.Text)
		_, _, err := bot.API().
			MessageAPI.
			PostMessage(context.Background(), p.Message.ChannelID).
			PostMessageRequest(traq.PostMessageRequest{
				Content: "oisu-",
			}).
			Execute()
		if err != nil {
			log.Println(err)
		}
	})

	if err := bot.Start(); err != nil {
		panic(err)
	}
	e := echo.New()

	e.GET("/api/messages/:userId", handler.GETTimesUbugoe)
	e.GET("/api/messages/true/:username", handler.GETTrueUbugoe)

	// e.Start(":8080")
}

func (h *Handler) GETTimesUbugoe(c echo.Context) error {
	auth := h.makeAuth(c.Request().Context())
	userID := c.Param("userId")
	fmt.Println("User ID:", userID)

	channelList, _, err := h.client.ChannelAPI.GetChannels(auth).Path("gps/times/" + userID).Execute()
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(500, "something went wrong")
	}
	if len(channelList.Public) == 0 {
		return c.JSON(404, "No channels found")
	}
	channel := channelList.Public[0]
	channelID := channel.Id

	messages, _, _ := h.client.MessageAPI.GetMessages(auth, channelID).
		Since(time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local)).
		Limit(10).
		Order("asc").
		Execute()

	res := make([]TimesUbugoe, 0, len(messages))
	for _, message := range messages {

		userNameList, _, err := h.client.UserAPI.GetUser(auth, message.UserId).Execute()
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(500, "something went wrong")
		}
		userName := userNameList.Name

		res = append(res, TimesUbugoe{
			Content: message.Content,
			UserID:  userName,
		})
	}

	return c.JSON(200, res)
}

func (h *Handler) GETTrueUbugoe(c echo.Context) error {
	auth := h.makeAuth(c.Request().Context())

	name := c.Param("username")
	fmt.Println("User Name:", name)

	userID, _, err := h.client.UserAPI.GetUsers(auth).
		Name(name).
		Execute()

	if err != nil {
		c.Logger().Error(err)
		return c.JSON(500, "something went wrong")
	}
	if len(userID) == 0 {
		return c.JSON(404, "No user found")
	}

	messages, _, err := h.client.MessageAPI.SearchMessages(auth).
		From([]string{userID[0].Id}).
		Limit(10).
		Sort("-createdAt").
		Execute()

	if err != nil {
		c.Logger().Error(err)
		return c.JSON(500, "something went wrong")
	}

	res := make([]TrueUbugoe, 0, len(messages.Hits))
	for _, message := range messages.Hits {
		channelID := message.ChannelId
		content := message.Content

		channel, _, err := h.client.ChannelAPI.GetChannel(auth, channelID).Execute()
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(500, "something went wrong")
		}
		channelName := channel.Name
		res = append(res, TrueUbugoe{
			Content:   content,
			Channel:   channelName,
			CreatedAt: message.CreatedAt,
		})
	}

	return c.JSON(200, res)
}
