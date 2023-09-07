package tui

import (
	"fmt"
	"kobla/blockchain/core/types"
	"math/big"
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
	newAccount
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
	{balance, "БАЛАНС КОШЕЛЬКА"},
	{newAccount, "НОВЫЙ КОШЕЛЕК"},
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

func (tui *TerminalUI) processNewAccount() {
	defer tui.app.SetFocus(tui.commands)
	tui.main.Clear()

	account, err := types.NewAccount()
	if err != nil {
		fmt.Fprintf(tui.main, "Error: can't generate new account: %s", err)
		return
	}

	fmt.Fprintf(tui.main, "[greenyellow]АДРЕС:[white] %s\n[greenyellow]ПРИВАТНЫЙ КЛЮЧ:[white] %s",
		account.Address(), account.PrivateKey())
}

func (tui *TerminalUI) processSendTxCommand() {
	form := newForm().
		AddInputField("ПОЛУЧАТЕЛЬ:", "", 256, nil, nil).
		AddInputField("СУММА ПЕРЕВОДА:", "", 40, tview.InputFieldInteger, nil).
		AddInputField("ДАННЫЕ:", "", 2048, nil, nil).
		AddInputField("ПРИВАТНЫЙ КЛЮЧ:", "", 128, nil, nil)

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			tui.mflex.RemoveItem(form)
			tui.app.SetFocus(tui.commands)
			return nil
		}

		return event
	})

	form.AddButton("ОТПРАВИТЬ", func() {
		defer tui.app.SetFocus(tui.commands)
		defer tui.mflex.RemoveItem(form)
		tui.main.Clear()

		receiver := types.AddressFromString(getFormInput(form, 0))

		var amount big.Int
		if _, ok := amount.SetString(getFormInput(form, 1), 10); !ok {
			fmt.Fprintf(tui.main, "Error: invalid amount value")
			return
		}

		data := []byte(getFormInput(form, 2))

		signer, err := types.AccountFromPrivKey(getFormInput(form, 3))
		if err != nil {
			fmt.Fprintf(tui.main, "Error: invalid private key: %s", err)
			return
		}

		tx := types.NewTransaction(signer.Address(), receiver, amount.Uint64(), data)
		if err := tx.WithSignature(signer); err != nil {
			fmt.Fprintf(tui.main, "Error: can't sign transaction: %s", err)
			return
		}

		tui.bc.SendTx(tx)
	})

	tui.mflex.AddItem(form, 0, 2, false)
	tui.app.SetFocus(form)
}

func (tui *TerminalUI) processMineBlockCommand() {
	form := newForm()

	addTxInputField := func() { form.AddInputField("ВВЕДИТЕ ХЕШ ТРАНЗАКЦИИ:", "", 128, nil, nil) }
	addTxInputField()

	done := false
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			if lastItemIdx := form.GetFormItemCount() - 1; !done && lastItemIdx >= 0 {
				if input := getFormInput(form, lastItemIdx); input == "" {
					done = true
					form.RemoveFormItem(lastItemIdx)
					form.AddInputField("ВВЕДИТЕ СВОЙ АДРЕС:", "", 128, nil, nil)
					tui.app.SetFocus(form)
					form.SetFocus(lastItemIdx)
					return nil
				}
			}

			tui.mflex.RemoveItem(form)
			tui.app.SetFocus(tui.commands)
			return nil
		case tcell.KeyEnter:
			if !done {
				addTxInputField()
			}
		}

		return event
	})

	form.AddButton("СФОРМИРОВАТЬ БЛОК", func() {
		defer tui.app.SetFocus(tui.commands)
		defer tui.mflex.RemoveItem(form)
		tui.main.Clear()

		lastItemIdx := form.GetFormItemCount() - 1
		coinbase := types.AddressFromString(getFormInput(form, lastItemIdx))

		var txs []*types.Transaction
		for i := 0; i < lastItemIdx; i++ {
			hash := types.HashFromString(getFormInput(form, i))
			tx, err := tui.bc.TxByHashFromMempool(hash)
			if err != nil {
				fmt.Fprintf(tui.main, "Error: can't get %s tx from mempool: %s", hash, err)
				return
			}

			txs = append(txs, tx)
		}

		if err := tui.bc.MineBlock(txs, coinbase); err != nil {
			fmt.Fprintf(tui.main, "Error: can't mine block: %s", err)
			return
		}

		fmt.Fprint(tui.main, "Новый блок успешно сформирован!")
	})

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
			output := f(trimInput(inputField))
			fmt.Fprint(tui.main.Clear(), output)
		}
	})

	tui.mflex.AddItem(inputField, 0, 1, false)
	tui.app.SetFocus(inputField)
}

func trimInput(f *tview.InputField) string {
	return strings.TrimSpace(f.GetText())
}

func getFormInput(form *tview.Form, idx int) string {
	return trimInput(form.GetFormItem(idx).(*tview.InputField))
}

func newForm() *tview.Form {
	form := tview.NewForm().
		SetFieldBackgroundColor(backgroundColor).
		SetFieldTextColor(tcell.ColorWhite).
		SetButtonBackgroundColor(backgroundColor).
		SetButtonTextColor(tcell.ColorDarkSeaGreen)

	form.SetLabelColor(tcell.ColorGreenYellow).
		SetBackgroundColor(tcell.ColorBlack).
		SetBorder(true)

	return form
}
