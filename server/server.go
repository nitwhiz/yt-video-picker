package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"os"
	"strings"
	"time"
)

func DurationToSeconds(duration string) int64 {
	parts := strings.Split(duration, "T")

	d, err := time.ParseDuration(strings.ToLower(parts[1]))

	if err != nil {
		return 0
	}

	return int64(d.Seconds())
}

func main() {
	apiKey, apiKeyIsSet := os.LookupEnv("YTVP_GOOGLE_API_KEY")

	if !apiKeyIsSet {
		panic(fmt.Errorf("YTVP_GOOGLE_API_KEY not set"))
	}

	youtubeService, err := youtube.NewService(context.Background(), option.WithAPIKey(apiKey))

	if err != nil {
		panic(err)
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.Use(cors.Default())

	r.GET("/lists/:listId/videos", func(c *gin.Context) {
		listId := c.Param("listId")

		cacheTime := int64(0)

		var videos []VideoInfo

		videos, err = ReadCache(listId)

		if err != nil {
			videos, err = FetchVideos(youtubeService, listId)

			if err != nil {
				c.Status(500)
				return
			}
		}

		cacheTime, _ = PassThruCache(videos, listId)

		c.JSON(200, gin.H{
			"Videos":    videos,
			"Count":     len(videos),
			"CacheTime": cacheTime,
		})
	})

	_ = r.Run("0.0.0.0:8085")
}
