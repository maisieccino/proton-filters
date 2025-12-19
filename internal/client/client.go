package client

import (
	"context"
	"fmt"
	"net/http/cookiejar"
	"os"

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

func NewClient(ctx context.Context) (*proton.Client, error) {
	manager, jar, err := NewManager(ctx)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return nil, err
	}

	var (
		c *proton.Client
		a proton.Auth
	)

	if jar.HasToken() {
		c, a, err = manager.NewClientWithRefresh(ctx, jar.UID, jar.Ref)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return nil, err
		}
		println("Reloaded auth from disk")
	} else {
		c, a, err = manager.NewClientWithLogin(ctx,
			os.Getenv("PROTON_USER"),
			[]byte(os.Getenv("PROTON_PASS")),
		)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return nil, err
		}
	}

	jar.UID = a.UID
	jar.Ref = a.RefreshToken

	err = jar.Persist()
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return nil, err
	}

	return c, nil
}
