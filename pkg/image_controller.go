package pkg

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

type UnsplashResponse struct {
	URLs `json:"urls"`
}

type URLs struct {
	Regular string `json:"regular"`
}

type ImageController struct {
	UnsplashKey string
	mu          sync.Mutex
}

func (c *ImageController) GetPhoto(ChatID int64) tgbotapi.PhotoConfig {
	response, err := http.Get(fmt.Sprintf("https://api.unsplash.com/photos/random?client_id=%s", c.UnsplashKey))
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	var unsplashResponse UnsplashResponse
	if err := json.NewDecoder(response.Body).Decode(&unsplashResponse); err != nil {
		log.Fatal(err)
	}

	imageURL := unsplashResponse.URLs.Regular

	response, err = http.Get(imageURL)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	fileName := "image.jpg"

	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	imageData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	photoConfig := tgbotapi.NewPhoto(ChatID, tgbotapi.FileBytes{Name: fileName, Bytes: imageData})

	return photoConfig
}
