package components

import (
	"github.com/rivo/tview"

	"github.com/jorgerojas26/lazysql/internal/i18n"
)

func NewErrorModal(message string) *tview.Modal {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{i18n.T("OK")})
	return modal
}
