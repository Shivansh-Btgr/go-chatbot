//So in this project. We will use 4 things :

// 1) Golang as the programming language

// 2) Slack to interact with the bot

// 3) Wit to understand the request

// 4) Wolfram to etch the answers from the web

package main

import (
	// This package is used to pass deadlines n stuff. It can cancel and/or timeout requests.. i think
	"context"

	// Encoding and decoding json files.. self explainatory
	"encoding/json"

	// The printing n stuff
	"fmt"

	// Logging stuff
	"log"

	// Dealing with os related stuff like file handling or smth
	"os"

	// This helps in dealing with the env files
	"github.com/joho/godotenv"

	// This helps in dealing with wolfram api
	"github.com/krognol/go-wolfram"

	// This is slacker
	"github.com/shomali11/slacker"

	// This thing makes json more bearable
	"github.com/tidwall/gjson"

	// witai
	witai "github.com/wit-ai/wit-go/v2"
)

// Pointer to help us interact with wolfram
var wolframClient *wolfram.Client

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Command Events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
		fmt.Println()
	}
}

func main() {
	godotenv.Load(".env")

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	client := witai.NewClient(os.Getenv("WIT_AI_TOKEN"))
	wolframClient := &wolfram.Client{AppID: os.Getenv("WOLFRAM_APP_ID")}
	go printCommandEvents(bot.CommandEvents())

	bot.Command("query for bot - <message>", &slacker.CommandDefinition{
		Description: "send any question to wolfram",
		Example:     "which is the worst language to ever exist?",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			query := request.Param("message")

			msg, _ := client.Parse(&witai.MessageRequest{
				Query: query,
			})
			data, _ := json.MarshalIndent(msg, "", "    ")
			rough := string(data[:])
			value := gjson.Get(rough, "entities.wit$wolfram_search_query:wolfram_search_query.0.value")
			answer := value.String()
			res, err := wolframClient.GetSpokentAnswerQuery(answer, wolfram.Metric, 1000)
			if err != nil {
				fmt.Println("there is an error")
			}
			fmt.Println(value)
			response.Reply(res)
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)

	if err != nil {
		log.Fatal(err)
	}
}