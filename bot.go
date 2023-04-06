package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"github.com/gin-gonic/gin"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type message struct {
	Message string `json:"message"`
}

var messages = []message{}
var DANIELS_ID string = "1536307849"
var MATTHEWS_ID string = "1504084599"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	router := gin.Default()
	b, err := bot.New("6092711168:AAFQ1PfJnYGW3_IjOOR-1bdXJo1QF6X4J04", bot.WithDefaultHandler(handler))
	if err != nil {
		panic(err)
	}

	router.POST("/data", func(c *gin.Context) {
		var m message

		if err := c.ShouldBindJSON(&m); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		messages = append(messages, m)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: DANIELS_ID,
			Text:   m.Message,
		})
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: MATTHEWS_ID,
			Text:   m.Message,
		})
		c.IndentedJSON(http.StatusCreated, m)
	})

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		b.Start(ctx)
		wg.Done()
	}()
	go func() {
		router.Run("localhost:8080")
		wg.Done()
	}()
	wg.Wait()
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	var lastMessage message
	if len(messages) > 0 {
		lastMessage = messages[len(messages)-1]
	} else {
		lastMessage = message{
			Message: "No Last Message Found !",
		}

	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   lastMessage.Message,
	})
	// fmt.Println(update.Message.Chat.ID)
}
