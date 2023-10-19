package tui

import (
	"errors"
	"fmt"
	"kobla/blockchain/core/types"
	"math/big"
	"runtime/debug"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/btcsuite/btcutil/base58"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	log "github.com/sirupsen/logrus"
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
			return tui.printError(err, "Ошибка: неверный номер блока: %s", input)
		}

		block, err := tui.bc.BlockByNumber(number)
		if err != nil {
			return tui.printError(err, "Ошибка: не удалость получить блок %d", number)
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
				s = tui.printError(fmt.Errorf("%s", r), "Ошибка: %s", input)
			}
		}()

		defer tui.app.SetFocus(tui.commands)

		hash := types.HashFromString(input)
		block, err := tui.bc.BlockByHash(hash)
		if err != nil {
			return tui.printError(err, "Ошибка: не удалось получить блок: %s", hash)
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
			return tui.printError(err, "Ошибка: не удалось получить транзакцию %s", hash)
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
			return tui.printError(err, "Ошибка: не удалось получить баланс %s", address)
		}

		return fmt.Sprintf("[greenyellow]АДРЕС:[white] %s\n[greenyellow]БАЛАНС:[white] %d", address, balance)
	})
}

func (tui *TerminalUI) processNewAccount() {
	defer tui.app.SetFocus(tui.commands)
	tui.main.Clear()

	account, err := types.NewAccount()
	if err != nil {
		tui.printError(err, "Ошибка: не удалось сгенирировать новый аккаунт")
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
			tui.printError(errors.New("invalid amount value"), "Ошибка: неверная величина перевода %s", getFormInput(form, 1))
			return
		}

		data := []byte(getFormInput(form, 2))

		signer, err := types.AccountFromPrivKey(getFormInput(form, 3))
		if err != nil {
			tui.printError(err, "Ошибка: некорректный приватный ключ")
			return
		}

		tx := types.NewTransaction(signer.Address(), receiver, amount.Uint64(), data)
		if err := tx.WithSignature(signer); err != nil {
			tui.printError(err, "Ошибка: не удалось подписать транзакцию")
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
				tui.printError(err, "Ошибка: не удалось получить транзакцию %s из мемпула", hash)
				return
			}

			txs = append(txs, tx)
		}

		if err := tui.bc.MineBlock(txs, coinbase); err != nil {
			tui.printError(err, "Ошибка: не удалось сформировать блок")
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
			if output != "" {
				fmt.Fprint(tui.main.Clear(), output)
			}
		}
	})

	tui.mflex.AddItem(inputField, 0, 1, false)
	tui.app.SetFocus(inputField)
}

func (tui *TerminalUI) printError(err error, msg string, opt ...any) string {
	fmt.Fprintf(tui.main.Clear(), msg, opt...)
	log.Error(err)
	debug.PrintStack()
	return ""
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

func sprintTx(tx *types.Transaction) string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("[greenyellow]ХЕШ:[white] %s\n", tx.Hash))

	if tx.Sender != types.ZeroAddress {
		s.WriteString(fmt.Sprintf("[greenyellow]ОТПРАВИТЕЛЬ:[white] %s\n", tx.Sender))
	}

	if tx.Amount != 0 {
		if tx.Receiver != types.ZeroAddress {
			s.WriteString(fmt.Sprintf("[greenyellow]ПОЛУЧАТЕЛЬ:[white] %s\n", tx.Receiver))
		}
		s.WriteString(fmt.Sprintf("[greenyellow]СУММА ПЕРЕВОДА:[white]: %d\n", tx.Amount))
	}

	s.WriteString("[greenyellow]СТАТУС:[white] ")
	if tx.Status == types.TxSuccess {
		s.WriteString("ВАЛИДНА\n")
	} else {
		s.WriteString("НЕВАЛИДНА\n")
	}

	s.WriteString("[greenyellow]ДАННЫЕ:[white] ")
	if utf8.Valid(tx.Data) {
		s.WriteString(string(tx.Data))
	} else {
		s.WriteString(base58.Encode(tx.Data))
	}
	s.WriteRune('\n')

	if tx.Sender != types.ZeroAddress {
		s.WriteString(fmt.Sprintf("[greenyellow]ПОДПИСЬ:[white] %s", base58.Encode(tx.Signature)))
	}

	return s.String()
}
