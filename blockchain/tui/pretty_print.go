package tui

import (
	"fmt"
	"kobla/blockchain/core/types"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/btcsuite/btcutil/base58"
)

func sprintBlock(block *types.Block) string {
	str := fmt.Sprintf("[greenyellow]НОМЕР:[white] %d\n[greenyellow]ВРЕМЯ:[white] %s\n[greenyellow]NONCE:[white] %d\n[greenyellow]ХЕШ:[white] %s\n[greenyellow]ТРАНЗАКЦИИ:[white]\n",
		block.Number, time.Unix(block.Timestamp, 0), block.Nonce, block.Hash)

	txs := make([]string, 0, len(block.Transactions))
	for i, tx := range block.Transactions {
		txs = append(txs, fmt.Sprintf("[greenyellow](%d)[white] %s", i+1, tx.Hash))
	}

	return str + strings.Join(txs, "\n")
}

func sprintTx(tx *types.Transaction) string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("[greenyellow]ХЕШ:[white] %s\n", tx.Hash))

	if tx.Sender != types.ZeroAddress {
		s.WriteString(fmt.Sprintf("[greenyellow]ОТПРАВИТЕЛЬ:[white] %s\n", tx.Sender))
	}

	if tx.Amount != 0 {
		s.WriteString(fmt.Sprintf("[greenyellow]ПОЛУЧАТЕЛЬ:[white] %s\n", tx.Receiver))
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
