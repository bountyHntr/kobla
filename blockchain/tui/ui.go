package tui

import (
	"fmt"

	"github.com/rivo/tview"
)

func Run() error {
	// box := tview.NewBox().SetBorder(true)

	flex := tview.NewFlex().
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Commands"), 0, 1, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("KOBLA - KOLTSA OPEN BLOCKCHAIN LEARNING AID"), 0, 1, false).
			AddItem(tview.NewBox().SetBorder(true), 0, 4, false), 0, 2, false).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Mempool"), 0, 1, false)

	if err := tview.NewApplication().SetRoot(flex, true).Run(); err != nil {
		return fmt.Errorf("run app: %w", err)
	}

	return nil
}
