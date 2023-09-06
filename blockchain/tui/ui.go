package tui

import (
	"fmt"
	"kobla/blockchain/core/chain"
	"kobla/blockchain/core/types"

	log "github.com/sirupsen/logrus"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const backgroundColor = tcell.ColorBlack

type TerminalUI struct {
	bc    *chain.Blockchain
	app   *tview.Application
	mflex *tview.Flex

	header   *tview.TextView
	commands *tview.List
	main     *tview.TextView
	mempool  *tview.TextView
}

func newTUI(bc *chain.Blockchain) *TerminalUI {
	return &TerminalUI{
		bc:    bc,
		app:   tview.NewApplication(),
		mflex: tview.NewFlex(),

		header:   tview.NewTextView(),
		commands: tview.NewList(),
		main:     tview.NewTextView(),
		mempool:  tview.NewTextView(),
	}
}

func Run(bc *chain.Blockchain) error {

	tui := newTUI(bc)

	tui.configureApp()
	tui.configureHeader()
	tui.configureCommands()
	tui.configureMain()
	tui.configureMempool()

	tui.updateLastBlock()
	tui.updateMempool()

	if err := tui.run(); err != nil {
		return fmt.Errorf("run app: %w", err)
	}

	return nil
}

func (tui *TerminalUI) configureApp() {
	flex := tview.NewFlex().
		AddItem(tui.commands, 0, 1, true).
		AddItem(tui.mflex.SetDirection(tview.FlexRow).
			AddItem(tui.header, 0, 1, false).
			AddItem(tui.main, 0, 4, false), 0, 2, false).
		AddItem(tui.mempool, 0, 1, false)

	tui.app.SetRoot(flex, true)
}

func (tui *TerminalUI) configureHeader() {

	tui.header.SetChangedFunc(func() {
		_, _, _, height := tui.header.GetRect()
		tui.header.SetBorderPadding(height/4, 0, 0, 0)
		tui.app.Draw()
	})

	tui.header.
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetBorder(true).
		SetTitle("KOBLA - KOLTSA OPEN BLOCKCHAIN LEARNING AID")

	tui.header.SetBackgroundColor(backgroundColor)
}

func (tui *TerminalUI) configureCommands() {

	tui.commands.SetSelectedFunc(func(idx int, command string, _ string, _ rune) {
		switch Command(idx) {
		case blockByNumber:
			tui.processBlockByNumberCommand()
		case blockByHash:
			tui.processBlockByHashCommand()
		case txByHash:
			tui.processTxByHash()
		case balance:
			tui.processBalance()
		case newAccount:
			tui.processNewAccount()
		case sendTransaction:
			tui.processSendTxCommand()
		case mineBlock:
			tui.processMineBlockCommand()
		case quit:
			tui.app.Stop()
		default:
			log.Panicf("invalid command %s", command)
		}
	})

	tui.commands.SetChangedFunc(func(idx int, _, _ string, _ rune) {
		switch Command(idx) {
		case blockByNumber:
			fmt.Fprint(tui.main.Clear(), blockByNumberTheory)
		case blockByHash:
			fmt.Fprint(tui.main.Clear(), blockByHashTheory)
		case txByHash:
		case balance:
		case newAccount:
		case sendTransaction:
		case mineBlock:
		case quit:
		}
	})

	tui.addAllCommands().
		SetShortcutColor(tcell.ColorGreenYellow).
		SetBorder(true).
		SetTitle("Команды")

	tui.commands.SetBackgroundColor(backgroundColor)
}

func (tui *TerminalUI) configureMain() {
	tui.main.SetChangedFunc(func() {
		tui.app.Draw()
	})

	tui.main.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEsc {
			tui.app.SetFocus(tui.commands)
		}
	})

	tui.main.SetDynamicColors(true).SetBorder(true)
	tui.main.SetBackgroundColor(backgroundColor)
}

func (tui *TerminalUI) configureMempool() {
	tui.mempool.SetChangedFunc(func() {
		tui.app.Draw()
	})

	tui.mempool.
		SetDynamicColors(true).
		SetBorder(true).
		SetTitle("Мемпул")

	tui.mempool.SetBackgroundColor(backgroundColor)
}

func (tui *TerminalUI) run() error {
	return tui.app.Run()
}

func (tui *TerminalUI) updateMempool() {
	const topN = 20

	go func() {
		updates := make(chan *types.Void, 1)

		subID := tui.bc.SubscribeMempoolUpdates(updates)
		defer tui.bc.UnsubscribeMempoolUpdates(subID)

		for range updates {
			txs := tui.bc.TopMempoolTxs(topN)
			tui.mempool.Clear()

			for i, tx := range txs {
				fmt.Fprintf(tui.mempool, "[greenyellow](%d)[white] %s\n", i+1, tx.Hash.String())
			}
		}
	}()
}
