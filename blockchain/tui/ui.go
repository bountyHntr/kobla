package tui

import (
	"fmt"
	"kobla/blockchain/core/chain"
	"kobla/blockchain/core/types"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TerminalUI struct {
	bc    *chain.Blockchain
	app   *tview.Application
	mflex *tview.Flex

	header   *tview.TextView
	commands *tview.List
	main     *tview.TextView
	mempool  *tview.List
}

func newTUI(bc *chain.Blockchain) *TerminalUI {
	return &TerminalUI{
		bc:    bc,
		app:   tview.NewApplication(),
		mflex: tview.NewFlex(),

		header:   tview.NewTextView(),
		commands: tview.NewList(),
		main:     tview.NewTextView(),
		mempool:  tview.NewList(),
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
		SetShortcutColor(tcell.ColorRed).
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
				"[red]ПОСЛЕДНИЙ БЛОК:[white]\n[red]НОМЕР:[white] %d [red]ВРЕМЯ:[white] %s [red]NONCE:[white] %d\n[red]HASH:[white] %s",
				block.Number, time.Unix(block.Timestamp, 0).String(), block.Nonce, block.Hash.String(),
			)
		}
	}()
}
