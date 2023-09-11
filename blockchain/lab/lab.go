package lab

import (
	"fmt"
	"kobla/blockchain/core/chain"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const backgroundColor = tcell.ColorBlack

type Lab struct {
	bc    *chain.Blockchain
	app   *tview.Application
	mflex *tview.Flex

	header   *tview.TextView
	commands *tview.List
	main     *tview.TextView
}

func newLab(bc *chain.Blockchain) *Lab {
	return &Lab{
		bc:    bc,
		app:   tview.NewApplication(),
		mflex: tview.NewFlex(),

		header:   tview.NewTextView(),
		commands: tview.NewList(),
		main:     tview.NewTextView(),
	}
}

func Run(bc *chain.Blockchain) error {
	lab := newLab(bc)

	lab.configureApp()
	lab.configureHeader()
	lab.configureCommands()
	lab.configureMain()

	if err := lab.run(); err != nil {
		return fmt.Errorf("run app: %w", err)
	}

	return nil
}

func (lab *Lab) configureApp() {
	flex := tview.NewFlex().
		AddItem(lab.commands, 0, 1, true).
		AddItem(lab.mflex.SetDirection(tview.FlexRow).
			AddItem(lab.header, 0, 1, false).
			AddItem(lab.main, 0, 4, false), 0, 2, false)

	lab.app.SetRoot(flex, true)
}

func (lab *Lab) configureHeader() {

	lab.header.SetChangedFunc(func() {
		_, _, _, height := lab.header.GetRect()
		lab.header.SetBorderPadding(height/4, 0, 0, 0)
		lab.app.Draw()
	})

	lab.header.
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetBorder(true).
		SetTitle("KOBLA - KOLTSA OPEN BLOCKCHAIN LEARNING AID")

	lab.header.SetBackgroundColor(backgroundColor)
}

func (lab *Lab) configureCommands() {

	lab.commands.SetSelectedFunc(func(idx int, command string, _ string, _ rune) {
		switch Command(idx) {
		case blockByNumber:
			// lab.processMineBlockCommand()
		case quit:
			lab.app.Stop()
		default:
			log.Panicf("invalid command %s", command)
		}
	})

	lab.addAllCommands().
		SetShortcutColor(tcell.ColorGreenYellow).
		SetBorder(true).
		SetTitle("Команды")

	lab.commands.SetBackgroundColor(backgroundColor)
}

func (lab *Lab) configureMain() {
	lab.main.SetChangedFunc(func() {
		lab.app.Draw()
	})

	lab.main.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEsc {
			lab.app.SetFocus(lab.commands)
		}
	})

	lab.main.SetDynamicColors(true).SetBorder(true)
	lab.main.SetBackgroundColor(backgroundColor)
}

func (lab *Lab) run() error {
	return lab.app.Run()
}
