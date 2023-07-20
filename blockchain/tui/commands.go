package tui

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (tui *TerminalUI) addAllCommands() *tview.List {
	tui.commands.AddItem("block_by_number", "", '1', nil).ShowSecondaryText(false)
	tui.commands.AddItem("block_by_hash", "", '2', nil).ShowSecondaryText(false)
	tui.commands.AddItem("quit", "", '3', nil).ShowSecondaryText(false)
	return tui.commands
}

func (tui *TerminalUI) processBlockByNumberCommand() {

	inputField := tview.NewInputField().
		SetLabel("ВВЕДИТЕ НОМЕР БЛОКА: ").
		SetLabelColor(tcell.ColorRed).
		SetFieldWidth(20).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetAcceptanceFunc(tview.InputFieldInteger)

	inputField.SetBorder(true)

	inputField.SetDoneFunc(func(key tcell.Key) {

		defer tui.mflex.RemoveItem(inputField)

		switch key {
		case tcell.KeyEsc:
			tui.app.SetFocus(tui.commands)
		case tcell.KeyEnter:
			blockNumStr := inputField.GetText()

			blockNum, err := strconv.ParseInt(blockNumStr, 10, 64)
			if err != nil {
				fmt.Fprintf(tui.main.Clear(), "Error: invalid block number: %s: %s",
					blockNumStr, err)
				tui.app.SetFocus(tui.commands)
				return
			}

			block, err := tui.bc.BlockByNumber(blockNum)
			if err != nil {
				fmt.Fprintf(tui.main.Clear(), "Error: can't get block %d: %s", blockNum, err)
				tui.app.SetFocus(tui.commands)
				return
			}

			fmt.Fprint(tui.main.Clear(), block.PrettyPrintString())
			tui.app.SetFocus(tui.main)
		}
	})

	tui.mflex.AddItem(inputField, 0, 1, false)
	tui.app.SetFocus(inputField)
}
