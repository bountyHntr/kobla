package tui

import (
	"fmt"
	"kobla/blockchain/core/types"
	"strings"
	"time"
)

func sprintBlock(block *types.Block) string {
	str := fmt.Sprintf("[red]НОМЕР:[white] %d\n[red]ВРЕМЯ:[white] %s\n[red]NONCE:[white] %d\n[red]ХЕШ:[white] %s\n[red]ТРАНЗАКЦИИ:[white]\n",
		block.Number, time.Unix(block.Timestamp, 0).Format(time.UnixDate), block.Nonce, block.Hash.String())

	txs := make([]string, 0, len(block.Transactions))
	for i, tx := range block.Transactions {
		txs = append(txs, fmt.Sprintf("[red](%d)[white] %s", i+1, tx.Hash.String()))
	}

	return str + strings.Join(txs, "\n")
}
