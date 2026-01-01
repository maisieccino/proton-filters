package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/ProtonMail/go-proton-api"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/maisieccino/proton-filters/internal/client"
	"github.com/spf13/cobra"
)

func Check(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	c, err := client.NewClient(ctx)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	f, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Println("error reading file")
		return nil
	}

	res, err := c.CheckFilter(ctx, string(f))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if len(res) == 0 {
		fmt.Println("✅ Sieve filter has no issues.")
		return nil
	}
	fmt.Println("❌ Issues found while checking sieve filter:")
	for _, item := range res {
		fmt.Printf("%s Line %d: %s\n", sev(item), item.From.Line, item.Message)
	}
	quick.Highlight(os.Stdout, string(f), "sieve", "terminal256", "monokai")
	return nil
}

func sev(i proton.FilterIssue) string {
	switch i.Severity {
	case "error":
		return "⛔️"
	case "warning":
		return "⚠️"
	case "hint":
		return "💡"
	case "info":
	default:
		return "ℹ️"
	}
	return ""
}
