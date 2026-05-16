package filter

import (
	"os"
	"os/exec"

	tea "charm.land/bubbletea/v2"
	"github.com/ProtonMail/go-proton-api"
	"github.com/maisieccino/proton-filters/internal/types"
)

func EditFilter(filter proton.Filter) tea.Cmd {
	// Get editor
	editor, ok := os.LookupEnv("EDITOR")
	if !ok {
		// return proton.Filter{}, errors.New("EDITOR is blank")
		return nil
	}

	// Create temp file.
	f, err := os.CreateTemp("", "*.siv")
	if err != nil {
		// return proton.Filter{}, err
		return nil
	}

	_, err = f.WriteString(filter.Sieve)
	if err != nil {
		// return proton.Filter{}, err
		return nil
	}
	if err = f.Close(); err != nil {
		// return proton.Filter{}, err
		return nil
	}

	// Create EDITOR command, link up stdin/stdout.
	cmd := exec.Command(editor, f.Name())
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		defer os.Remove(f.Name())
		if err != nil {
			// return proton.Filter{}, err
			return nil
		}

		// Re-read file, grab the contents of the file.
		f, err = os.Open(f.Name())
		if err != nil {
			// return proton.Filter{}, err
			return nil
		}

		buf := make([]byte, 0)
		_, err = f.Read(buf)
		if err != nil {
			// return proton.Filter{}, err
			return nil
		}

		filter.Sieve = string(buf)

		return types.FilterMsg(filter)
	})
}
