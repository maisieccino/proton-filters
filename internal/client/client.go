package client

import (
	"context"
	"net/http/cookiejar"

	"github.com/ProtonMail/go-proton-api"
	"github.com/maisieccino/proton-filters/internal/cookies"
	"golang.org/x/net/publicsuffix"
)

func NewManager(ctx context.Context) (*proton.Manager, *cookies.Jar, error) {
	j, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		return nil, nil, err
	}
	jar, err := cookies.New(j, cookies.CookieJarFilename)
	if err != nil {
		return nil, nil, err
	}

	return proton.New(
		proton.WithCookieJar(jar),
		proton.WithAppVersion("other"),
	), jar, nil
}
