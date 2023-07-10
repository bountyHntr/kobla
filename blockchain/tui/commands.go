package tui

import (
	"fmt"

	"github.com/rivo/tview"
	log "github.com/sirupsen/logrus"
)

func (tui *TerminalUI) addAllCommands() *tview.List {
	tui.commands.AddItem("block_by_number", "", '1', nil).ShowSecondaryText(false)
	tui.commands.AddItem("block_by_hash", "", '2', nil).ShowSecondaryText(false)
	tui.commands.AddItem("quit", "", '3', nil).ShowSecondaryText(false)
	return tui.commands
}

func (tui *TerminalUI) processBlockByNumberCommand() {
	tui.app.SetFocus(tui.main)

	// inputField := tview.NewInputField().
	// 	SetLabel("Введите номер блока: ").
	// 	SetFieldWidth(15).
	// 	SetAcceptanceFunc(tview.InputFieldInteger).
	// 	SetDoneFunc(func(key tcell.Key) {
	// 		app.Stop()
	// 	})

	block, err := tui.bc.BlockByNumber(-1)
	if err != nil {
		log.WithField("command", "block_by_number").
			WithField("block_number", -1).
			WithError(err).
			Error("get block by number")
		return
	}

	fmt.Fprint(tui.main.Clear(), block.PrettyPrintString())
}
