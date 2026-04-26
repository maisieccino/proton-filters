package types

import "github.com/ProtonMail/go-proton-api"

type (
	FiltersMsg []proton.Filter
	FilterMsg  proton.Filter
)

type ToggleMsg struct {
	FilterID string
}
