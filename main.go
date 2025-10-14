package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ProtonMail/go-proton-api"
)

func main() {
	manager := proton.New(proton.WithAppVersion("Other"))
	ctx := context.Background()

	c, _, err := manager.NewClientWithLogin(ctx,
		os.Getenv("PROTON_USER"),
		[]byte(os.Getenv("PROTON_PASS")),
	)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return
	}

	addresses, err := c.GetAddresses(ctx)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return
	}
	fmt.Println(addresses)

	// filters, err := c.GetAllFilters(ctx)
	// if err != nil {
	// 	fmt.Printf("error: %s\n", err.Error())
	// 	return
	// }

	// fmt.Println(filters)
}
