package mangadex

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

var cacheEnabled bool = true

func EnableCache() {
	cacheEnabled = true
}

func DisableCache() {
	cacheEnabled = false
}

func getCachePath() string {
	userCacheDir, errCache := os.UserCacheDir()
	if errCache != nil {
		logrus.Fatalf("Unable to retrieve cache directory: %s", errCache)
	}
	return filepath.Join(userCacheDir, "go-mangadex")
}

func getCachePathFor(mangadexURL *url.URL) string {
	fileName := getCacheFilename(mangadexURL)
	return filepath.Join(getCachePath(), fileName)
}

func getCacheFilename(mangadexURL *url.URL) string {
	return strings.ReplaceAll(mangadexURL.Path, "/", "_")
}

func cacheExists(mangadexURL *url.URL) bool {
	stat, err := os.Stat(getCachePathFor(mangadexURL))
	if os.IsNotExist(err) {
		return false
	}
	return !stat.IsDir()
}

func initCache() {
	cachePath := getCachePath()
	_, err := os.Stat(cachePath)
	if os.IsNotExist(err) {
		logrus.Infof("Cache directory does not exist, creating. [%s]", cachePath)
		errCachePath := os.MkdirAll(cachePath, 0755)
		if errCachePath != nil {
			logrus.Errorf("Cache directory couldn't be generated, caching will likely fail: %s", errCachePath)
		}
	}
}
