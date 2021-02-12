package main

import (
	"context"
	"fmt"
	"google.golang.org/api/youtube/v3"
	"math"
)

type VideoInfo struct {
	ID           string
	Title        string
	ThumbnailURL string
	Duration     int64
}

func FetchVideos(youtubeService *youtube.Service, playlistId string) ([]VideoInfo, error) {
	var videos []VideoInfo

	if playlistId == "" {
		return videos, fmt.Errorf("playlist id is empty")
	}

	playListItemsCall := youtubeService.PlaylistItems.List([]string{"snippet"}).PlaylistId(playlistId)

	err := playListItemsCall.Pages(context.Background(), func(playlistItemListResponse *youtube.PlaylistItemListResponse) error {
		var ids []string

		for _, item := range playlistItemListResponse.Items {
			ids = append(ids, item.Snippet.ResourceId.VideoId)
		}

		videoListService := youtubeService.Videos.List([]string{"snippet", "contentDetails"})

		for len(ids) > 0 {
			i := int(math.Min(float64(len(ids)), 50))

			fetchIdsSlice := ids[:i]
			ids = ids[i:]

			videoListResponse, err := videoListService.Id(fetchIdsSlice...).Do()

			if err != nil {
				return err
			}

			for _, video := range videoListResponse.Items {
				vi := VideoInfo{
					ID:           video.Id,
					Title:        video.Snippet.Title,
					ThumbnailURL: video.Snippet.Thumbnails.High.Url,
					Duration:     DurationToSeconds(video.ContentDetails.Duration),
				}

				videos = append(videos, vi)
			}
		}

		return nil
	})

	if err != nil {
		return videos, err
	}

	return videos, nil
}
