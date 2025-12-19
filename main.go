package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ProtonMail/go-proton-api"
	"github.com/maisieccino/proton-filters/internal/client"
)

func main() {
	ctx := context.Background()
	manager, jar, err := client.NewManager(ctx)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return
	}

	var (
		c *proton.Client
		a proton.Auth
	)

	if jar.HasToken() {
		c, a, err = manager.NewClientWithRefresh(ctx, jar.UID, jar.Ref)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
		println("Reloaded auth from disk")
	} else {
		c, a, err = manager.NewClientWithLogin(ctx,
			os.Getenv("PROTON_USER"),
			[]byte(os.Getenv("PROTON_PASS")),
		)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
	}

	jar.UID = a.UID
	jar.Ref = a.RefreshToken

	err = jar.Persist()
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return
	}
	defer c.Close()

	filters, err := c.GetAllFilters(ctx)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return
	}

	fmt.Println(filters)
}
