package tui

import (
	"fmt"
	"kobla/blockchain/core/types"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Command int

const (
	blockByNumber Command = iota
	blockByHash
	txByHash
	balance
	sendTransaction
	mineBlock
	quit
)

var allCommands = []struct {
	command Command
	name    string
}{
	{blockByNumber, "БЛОК ПО НОМЕРУ"},
	{blockByHash, "БЛОК ПО ХЕШУ"},
	{txByHash, "ТРАНЗАКЦИЯ ПО ХЕШУ"},
	{balance, "БАЛАНС АДРЕСА"},
	{sendTransaction, "ОТПРАВИТЬ ТРАНЗАКЦИЮ"},
	{mineBlock, "СФОРМИРОВАТЬ БЛОК"},
	{quit, "ВЫЙТИ"},
}

func (tui *TerminalUI) addAllCommands() *tview.List {
	for _, cmd := range allCommands {
		tui.commands.AddItem(cmd.name, "", '1'+rune(cmd.command), nil).ShowSecondaryText(false)
	}
	return tui.commands
}

func (tui *TerminalUI) processBlockByNumberCommand() {
	inputField := tview.NewInputField().
		SetFieldWidth(20).
		SetLabel("ВВЕДИТЕ НОМЕР БЛОКА:").
		SetAcceptanceFunc(tview.InputFieldInteger)

	tui.addInputField(inputField, func(input string) string {
		defer tui.app.SetFocus(tui.commands)

		number, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			return fmt.Sprintf("Error: invalid block number: %s: %s", input, err)
		}

		block, err := tui.bc.BlockByNumber(number)
		if err != nil {
			return fmt.Sprintf("Error: can't get block %d: %s", number, err)
		}

		return sprintBlock(block)
	})
}

func (tui *TerminalUI) processBlockByHashCommand() {
	inputField := tview.NewInputField().
		SetFieldWidth(128).
		SetLabel("ВВЕДИТЕ ХЕШ БЛОКА:")

	tui.addInputField(inputField, func(input string) (s string) {
		defer func() {
			if r := recover(); r != nil {
				s = fmt.Sprintf("Error: %s: %s", input, r)
			}
		}()

		defer tui.app.SetFocus(tui.commands)

		hash := types.HashFromString(input)
		block, err := tui.bc.BlockByHash(hash)
		if err != nil {
			return fmt.Sprintf("Error: can't get block %s: %s", hash, err)
		}

		return sprintBlock(block)
	})
}

func (tui *TerminalUI) processTxByHash() {
	inputField := tview.NewInputField().
		SetFieldWidth(128).
		SetLabel("ВВЕДИТЕ ХЕШ ТРАНЗАКЦИИ:")

	tui.addInputField(inputField, func(input string) (s string) {
		defer func() {
			if r := recover(); r != nil {
				s = fmt.Sprintf("Error: %s: %s", input, r)
			}
		}()

		defer tui.app.SetFocus(tui.commands)

		hash := types.HashFromString(input)
		tx, err := tui.bc.TxByHash(hash)
		if err != nil {
			return fmt.Sprintf("Error: can't get tx %s: %s", hash, err)
		}

		return sprintTx(tx)
	})
}

func (tui *TerminalUI) processBalance() {
	inputField := tview.NewInputField().
		SetFieldWidth(256).
		SetLabel("ВВЕДИТЕ АДРЕС КОШЕЛЬКА:")

	tui.addInputField(inputField, func(input string) (s string) {
		defer tui.app.SetFocus(tui.commands)

		address := types.AddressFromString(input)
		balance, err := tui.bc.Balance(address)
		if err != nil {
			return fmt.Sprintf("Error: can't get balance %s: %s", address, err)
		}

		return fmt.Sprintf("[greenyellow]АДРЕС:[white] %s\n[greenyellow]БАЛАНС:[white] %d", address, balance)
	})
}

func (tui *TerminalUI) processSendTxCommand() {
	const backgroundColor = 0x000003

	form := tview.NewForm().
		AddInputField("ОТПРАВИТЕЛЬ:", "", 256, nil, nil).
		AddInputField("ПОЛУЧАТЕЛЬ:", "", 256, nil, nil).
		AddInputField("СУММА ПЕРЕВОДА:", "", 40, tview.InputFieldInteger, nil).
		AddInputField("ДАННЫЕ:", "", 2048, nil, nil).
		AddButton("ОТПРАВИТЬ", func() {
		})

	form.SetFieldBackgroundColor(backgroundColor).
		SetFieldTextColor(tcell.ColorWhite).
		SetButtonBackgroundColor(backgroundColor).
		SetButtonTextColor(tcell.ColorDarkSeaGreen)

	form.SetLabelColor(tcell.ColorGreenYellow).
		SetBackgroundColor(tcell.ColorBlack).
		SetBorder(true)

	tui.mflex.AddItem(form, 0, 2, false)
	tui.app.SetFocus(form)
}

func (tui *TerminalUI) addInputField(inputField *tview.InputField, f func(string) string) {
	inputField.SetLabelColor(tcell.ColorGreenYellow).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetBorder(true)

	inputField.SetDoneFunc(func(key tcell.Key) {
		defer tui.mflex.RemoveItem(inputField)

		switch key {
		case tcell.KeyEsc:
			tui.app.SetFocus(tui.commands)
		case tcell.KeyEnter:
			output := f(strings.TrimSpace(inputField.GetText()))
			fmt.Fprint(tui.main.Clear(), output)
		}
	})

	tui.mflex.AddItem(inputField, 0, 1, false)
	tui.app.SetFocus(inputField)
}
