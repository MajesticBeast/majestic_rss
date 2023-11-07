package main

import (
	"context"
	"github.com/gtuk/discordwebhook"
	"github.com/majesticbeast/majestic_rss/internal/database"
	"github.com/mmcdole/gofeed"
	"log"
	"os"
	"sync"
	"time"
)

type webhook struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	FeedName    string `json:"feed_name"`
}

func fetchAndParseFeed(url string) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	rssFeed, err := fp.ParseURL(url)
	if err != nil {
		return nil, err
	}

	return rssFeed, nil
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	feedData, err := fetchAndParseFeed(feed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", feed.Name, err)
		return
	}
	for _, item := range feedData.Items {
		log.Println("Found post", item.Title)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Items))

	if len(feedData.Items) == 0 {
		log.Println("Feed is empty!")
		return
	}
	if feedData.Items[0].Title == feed.LastPostTitle {
		log.Println("Feed is up to date!")
		return
	}

	err = db.UpdateFeedLastPostTitle(context.Background(), database.UpdateFeedLastPostTitleParams{
		LastPostTitle: feedData.Items[0].Title,
		UpdatedAt:     time.Now().UTC(),
		ID:            feed.ID,
	})
	if err != nil {
		log.Println("Couldn't update feed:", err)
		return
	}

	hook := webhook{
		Title:       feedData.Items[0].Title,
		Description: feedData.Items[0].Description,
		URL:         feedData.Items[0].Link,
		FeedName:    feed.Name,
	}

	postToDiscord(hook)
}

func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Collecting feeds every %s on %v goroutines...", timeBetweenRequest, concurrency)
	ticker := time.NewTicker(timeBetweenRequest)

	for ; ; <-ticker.C {
		feeds, err := db.GetFeeds(context.Background())
		if err != nil {
			log.Println("Couldn't get feed URLs to update:", err)
			continue
		}
		log.Printf("Found %v feeds to fetch!", len(feeds))

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func postToDiscord(hook webhook) {
	discordBotUser := hook.FeedName
	content := hook.URL
	webhookUrl := os.Getenv("WEBHOOK_URL")

	message := discordwebhook.Message{
		Username: &discordBotUser,
		Content:  &content,
	}

	if err := sendWebhook(webhookUrl, message); err != nil {
		log.Println("Couldn't send webhook:", err)
	}
}

func sendWebhook(url string, message discordwebhook.Message) error {
	err := discordwebhook.SendMessage(url, message)
	if err != nil {
		return err
	}

	return nil
}
