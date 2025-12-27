// cookies is used to build a persistent cookie jar for Proton API
// authentication.
package cookies

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
)

type cookieMap map[string][]*http.Cookie

// Jar implements http.CookieJar by wrapping an existing cookie jar with
// functions for saving to disk and loading from disk on startup.
// UID and Ref are used for the Proton API to be able to re-authenticate without
// a password.
type Jar struct {
	jar           http.CookieJar
	cookiesByHost cookieMap
	UID           string
	Ref           string
	filename      string
	sync.RWMutex
}

func New(jar http.CookieJar, filename string) (*Jar, error) {
	j := &Jar{}

	err := j.loadCookies(filename)
	if err != nil {
		return nil, err
	}
	for urlString, cookies := range j.cookiesByHost {
		u, err := url.Parse(urlString)
		if err != nil {
			fmt.Printf("Error parsing url %s: %v\n", urlString, err)
			return nil, err
		}
		jar.SetCookies(u, cookies)
	}
	j.jar = jar
	j.filename = filename
	if j.cookiesByHost == nil {
		j.cookiesByHost = make(cookieMap)
	}
	return j, nil
}

func cookieKey(url *url.URL) string {
	return fmt.Sprintf("%s://%s", url.Scheme, url.Host)
}

func (j *Jar) SetCookies(url *url.URL, cookies []*http.Cookie) {
	j.Lock()
	defer j.Unlock()
	j.jar.SetCookies(url, cookies)

	j.cookiesByHost[cookieKey(url)] = cookies
}

func (j *Jar) Cookies(u *url.URL) []*http.Cookie {
	j.Lock()
	defer j.Unlock()
	return j.jar.Cookies(u)
}

type ExportData struct {
	Cookies cookieMap
	UID     string
	Ref     string
}

// Persist saves cookie jar to disk, as a JSON file.
func (j *Jar) Persist() error {
	data := ExportData{
		Cookies: j.cookiesByHost,
		UID:     j.UID,
		Ref:     j.Ref,
	}
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if err := os.WriteFile(j.filename, raw, 0o600); err != nil {
		return fmt.Errorf("writing to %s: %v", j.filename, err)
	}
	return nil
}

func (j *Jar) loadCookies(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return err
	}
	data := ExportData{}
	if err := json.Unmarshal(file, &data); err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}
	j.cookiesByHost = data.Cookies
	j.UID = data.UID
	j.Ref = data.Ref
	return nil
}

func (j *Jar) HasToken() bool {
	return j.Ref != ""
}
