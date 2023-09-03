//go:build pow

package tui

import (
	"fmt"
	"kobla/blockchain/core/types"
	"strings"
	"time"
)

func (tui *TerminalUI) updateLastBlock() {
	go func() {
		blockSub := make(chan *types.Block, 1)

		subID := tui.bc.SubscribeNewBlocks(blockSub)
		defer tui.bc.UnsubscribeBlocks(subID)

		for block := range blockSub {
			fmt.Fprintf(tui.header.Clear(),
				"[greenyellow]ПОСЛЕДНИЙ БЛОК:[white]\n[greenyellow]НОМЕР:[white] %d [greenyellow]ВРЕМЯ:[white] %s [greenyellow]NONCE:[white] %d\n[greenyellow]ХЕШ:[white] %s",
				block.Number, time.Unix(block.Timestamp, 0), block.Nonce, block.Hash.String(),
			)
		}
	}()
}

func sprintBlock(block *types.Block) string {
	str := fmt.Sprintf("[greenyellow]НОМЕР:[white] %d\n[greenyellow]ВРЕМЯ:[white] %s\n[greenyellow]NONCE:[white] %d\n[greenyellow]ХЕШ:[white] %s\n[greenyellow]ТРАНЗАКЦИИ:[white]\n",
		block.Number, time.Unix(block.Timestamp, 0), block.Nonce, block.Hash)

	txs := make([]string, 0, len(block.Transactions))
	for i, tx := range block.Transactions {
		txs = append(txs, fmt.Sprintf("[greenyellow](%d)[white] %s", i+1, tx.Hash))
	}

	return str + strings.Join(txs, "\n")
}
