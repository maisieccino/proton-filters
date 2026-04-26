package filter

import (
	"context"
	"errors"
	"os"
	"os/exec"

	"github.com/ProtonMail/go-proton-api"
)

func EditFilter(ctx context.Context, filter proton.Filter) (proton.Filter, error) {
	// Get editor
	editor, ok := os.LookupEnv("EDITOR")
	if !ok {
		return proton.Filter{}, errors.New("EDITOR is blank")
	}

	// Create temp file.
	f, err := os.CreateTemp("", "*.siv")
	if err != nil {
		return proton.Filter{}, err
	}
	defer os.Remove(f.Name())

	_, err = f.WriteString(filter.Sieve)
	if err != nil {
		return proton.Filter{}, err
	}
	if err = f.Close(); err != nil {
		return proton.Filter{}, err
	}

	// Create EDITOR command, link up stdin/stdout.
	cmd := exec.CommandContext(ctx, editor, f.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return proton.Filter{}, err
	}

	// Re-read file, grab the contents of the file.
	f, err = os.Open(f.Name())
	if err != nil {
		return proton.Filter{}, err
	}

	buf := make([]byte, 0)
	_, err = f.Read(buf)
	if err != nil {
		return proton.Filter{}, err
	}

	filter.Sieve = string(buf)

	return filter, nil
}
