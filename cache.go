package mangadex

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

var cacheEnabled bool = true

// EnableCache enables request caching for reads
func EnableCache() {
	cacheEnabled = true
}

// DisableCache disables cache reads
func DisableCache() {
	cacheEnabled = false
}

// getBaseCachePath returns the base path for the cache storage
func getBaseCachePath() string {
	userCacheDir, errCache := os.UserCacheDir()
	if errCache != nil {
		logrus.Fatalf("Unable to retrieve cache directory: %s", errCache)
	}
	return filepath.Join(userCacheDir, "go-mangadex")
}

// getCachePath returns the absolute path to the files cache for the provided URL
func getCachePath(mangadexURL *url.URL) string {
	fileName := getCacheFilename(mangadexURL)
	return filepath.Join(getBaseCachePath(), fileName)
}

// getCacheFilename generates a cache filename based on the URL of the request
// TODO: Use query arguments as well
func getCacheFilename(mangadexURL *url.URL) string {
	return strings.ReplaceAll(mangadexURL.Path, "/", "_")
}

// cacheExists checks that the cache for a certain URL exists or not
func cacheExists(mangadexURL *url.URL) bool {
	stat, err := os.Stat(getCachePath(mangadexURL))
	if os.IsNotExist(err) {
		return false
	}
	return !stat.IsDir()
}

// initCache makes sure that the cache directory exists so reads and writes on cache folders won't fail
func initCache() {
	cachePath := getBaseCachePath()
	_, err := os.Stat(cachePath)
	if os.IsNotExist(err) {
		logrus.Infof("Cache directory does not exist, creating. [%s]", cachePath)
		errCachePath := os.MkdirAll(cachePath, 0755)
		if errCachePath != nil {
			logrus.Errorf("Cache directory couldn't be generated, caching will likely fail: %s", errCachePath)
		}
	}
}
