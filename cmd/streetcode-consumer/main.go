package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

type Response struct {
	Data []Entity `json:"data"`
}

type Entity struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

var ctx = context.Background()

func main() {
	if godotenv.Load(".env") != nil {
		log.Fatal("error loading .env file")
	}

	log.Println("consumer started")

	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_AUTH"),
	})

	defer redisClient.Close()

	// Confirm Redis connection
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal(" unbale to connect to Redis ", err)
	}

	// Connect to Redis Stream
	err = redisClient.XGroupCreate(
		ctx,
		os.Getenv("REDIS_STREAM"),
		os.Getenv("REDIS_GROUP"),
		"0",
	).Err()
	if err != nil {
		log.Println(err)
	}

	for {
		entries, err := redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    os.Getenv("REDIS_GROUP"),
			Consumer: "*",
			Streams:  []string{os.Getenv("REDIS_STREAM"), ">"},
			Count:    1,
			Block:    0,
			NoAck:    false,
		}).Result()

		if err != nil {
			log.Fatal(err)
		}

		processEntries(entries, redisClient)
	}
}

func processEntries(entries []redis.XStream, redisClient *redis.Client) {
	for i := 0; i < len(entries[0].Messages); i++ {
		values := entries[0].Messages[i].Values
		eventName := fmt.Sprintf("%v", values["eventName"])
		href := fmt.Sprintf("%v", values["href"])

		if eventName == "href received" {
			if err := handleNewHref(href); err != nil {
				log.Fatal(err)
			}

			redisClient.XAck(
				ctx,
				os.Getenv("REDIS_STREAM"),
				os.Getenv("REDIS_GROUP"),
				entries[0].Messages[i].ID,
			)
		}
	}
}

func handleNewHref(href string) error {
	// log.Println("checking url", href)

	// Check if link has already been submitted to Streetcode
	// Assemble Streetcode API url that will search for link
	urlTest := fmt.Sprintf("%s%s", os.Getenv("API_FILTER_URL"), href)

	log.Println("checking url", urlTest)

	// Call Streetcode
	response, err := http.Get(urlTest)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Process response from Streetcode or fail
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(" response.Body ", err)
	}

	// Process the JSON response data
	data := Response{}
	json.Unmarshal([]byte(responseData), &data)

	// Finally, no data means we can publish to Streetcode
	if len(data.Data) == 0 {
		// POST to Streetcode
		jsonData := fmt.Sprintf(`{"data":{"type":"post--photo","attributes":{"field_post":{"value":"%s","format":"basic_html"},"field_visibility": "1"},"relationships":{"field_recipient_group":{"data":{"type":"group--public_group","id":"b55fe232-0fbf-4fa8-b697-ff7bb863ae6a"}}}}}`, href)
		request, _ := http.NewRequest("POST", os.Getenv("API_URL"), bytes.NewBuffer([]byte(jsonData)))
		request.Header.Set("Content-Type", "application/vnd.api+json")
		request.Header.Set("Accept", "application/vnd.api+json")
		request.SetBasicAuth(os.Getenv("USERNAME"), os.Getenv("PASSWORD"))

		client := &http.Client{}
		response, error := client.Do(request)
		if error != nil {
			panic(error)
		}
		defer response.Body.Close()

		fmt.Printf("INFO: [response] %s\n", response.Status)
	} else {
		log.Printf("INFO: [exists] %s", href)
	}
	return nil
}