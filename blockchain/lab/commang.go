package lab

import "github.com/rivo/tview"

type Command int

const (
	nextPage Command = iota
	prevPage
	blockByNumber
	blockByHash
	txByHash
	balance
	calcHash
	calcSign
	quit
)

var allCommands = []struct {
	command Command
	name    string
}{
	{nextPage, ">>>"},
	{prevPage, "<<<"},
	{blockByNumber, "БЛОК ПО НОМЕРУ"},
	{blockByHash, "БЛОК ПО ХЕШУ"},
	{txByHash, "ТРАНЗАКЦИЯ ПО ХЕШУ"},
	{balance, "БАЛАНС КОШЕЛЬКА"},
	{calcHash, "ВЫЧИСЛИТЬ ХЕШ"},
	{calcSign, "ВЫЧИСЛИТЬ ПОДПИСЬ"},
	{quit, "ВЫЙТИ"},
}

func (lab *Lab) addAllCommands() *tview.List {
	for _, cmd := range allCommands {
		lab.commands.AddItem(cmd.name, "", '1'+rune(cmd.command), nil).ShowSecondaryText(false)
	}
	return lab.commands
}
