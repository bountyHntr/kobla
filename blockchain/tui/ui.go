package tui

import (
	"fmt"
	"kobla/blockchain/core/chain"
	"kobla/blockchain/core/types"
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TerminalUI struct {
	bc  *chain.Blockchain
	app *tview.Application

	header   *tview.TextView
	commands *tview.List
	main     *tview.TextView
	mempool  *tview.List
}

func Run(bc *chain.Blockchain) error {

	tui := &TerminalUI{
		bc:       bc,
		app:      tview.NewApplication(),
		header:   tview.NewTextView(),
		commands: tview.NewList(),
		main:     tview.NewTextView(),
		mempool:  tview.NewList(),
	}

	tui.configureApp()
	tui.configureHeader()
	tui.configureCommands()
	tui.configureMain()
	tui.configureMempool()

	tui.updateLastBlock()

	if err := tui.run(); err != nil {
		return fmt.Errorf("run app: %w", err)
	}

	return nil
}

func (tui *TerminalUI) configureApp() {
	flex := tview.NewFlex().
		AddItem(tui.commands, 0, 1, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
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
}

func (tui *TerminalUI) configureCommands() {

	tui.commands.SetSelectedFunc(func(_ int, command, _ string, _ rune) {
		switch command {
		case "block_by_number":
			tui.processBlockByNumberCommand()
		case "quit":
			tui.app.Stop()
		default:
			log.Panicf("invalid command %s", command)
		}
	})

	tui.addAllCommands().
		SetBorder(true).
		SetTitle("Commands")
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

	tui.main.SetBorder(true)
}

func (tui *TerminalUI) configureMempool() {
	tui.mempool.SetBorder(true).SetTitle("Mempool")
}

func (tui *TerminalUI) run() error {
	return tui.app.Run()
}

func (tui *TerminalUI) updateLastBlock() {
	go func() {
		blockSub := make(chan *types.Block, 1)

		subID := tui.bc.SubscribeNewBlocks(blockSub)
		defer tui.bc.UnsubscribeBlocks(subID)

		for block := range blockSub {
			fmt.Fprintf(tui.header.Clear(),
				"[red]ПОСЛЕДНИЙ БЛОК:[white]\n[red]НОМЕР:[white] %d [red]ВРЕМЯ:[white] %s [red]NONCE:[white] %d\n[red]HASH:[white] 0x%s",
				block.Number, time.Unix(block.Timestamp, 0).String(), block.Nonce, block.Hash.Hex(),
			)
		}
	}()
}
