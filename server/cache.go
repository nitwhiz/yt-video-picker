package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func getCacheFilePath(playlistId string) string {
	s := sha1.New()
	s.Write([]byte(playlistId))
	cacheId := hex.EncodeToString(s.Sum(nil))

	return fmt.Sprintf("/tmp/%s.ytvp-cache", cacheId)
}

func isCached(cacheFile string) (bool, int64) {
	stat, err := os.Stat(cacheFile)

	if err == nil {
		modTime := stat.ModTime()

		if time.Now().Sub(modTime) < 6*time.Hour {
			return true, modTime.Unix()
		}
	}

	return false, 0
}

func ReadCache(playlistId string) ([]VideoInfo, error) {
	cacheFile := getCacheFilePath(playlistId)

	var videos []VideoInfo

	if cached, _ := isCached(cacheFile); !cached {
		return videos, fmt.Errorf("not cached")
	}

	data, err := ioutil.ReadFile(cacheFile)

	if err != nil {
		return videos, err
	}

	err = json.Unmarshal(data, &videos)

	if err != nil {
		return videos, err
	}

	return videos, nil
}

func PassThruCache(videos []VideoInfo, playlistId string) (int64, error) {
	cacheFile := getCacheFilePath(playlistId)

	cached, cacheTime := isCached(cacheFile)

	if cached {
		return cacheTime, nil
	}

	f, err := os.Create(cacheFile)

	if err != nil {
		return 0, err
	}

	defer f.Close()

	cacheContent, err := json.Marshal(videos)

	if err != nil {
		return 0, err
	}

	_, err = f.Write(cacheContent)

	if err != nil {
		return 0, err
	}

	return 0, nil
}
