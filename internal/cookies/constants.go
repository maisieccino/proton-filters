package cookies

import (
	"os"
	"path"
)

var (
	CacheDir = ""
	// File for storing session cookies
	CookieJarFilename = ""
)

func init() {
	dir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}
	CacheDir = dir
	appDir := path.Join(CacheDir, "proton-filters")
	if _, err := os.Stat(appDir); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(appDir, 0o755); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	CookieJarFilename = path.Join(appDir, "cookies.json")
}
