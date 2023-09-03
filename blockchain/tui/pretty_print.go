package tui

import (
	"fmt"
	"kobla/blockchain/core/types"
	"strings"
	"unicode/utf8"

	"github.com/btcsuite/btcutil/base58"
)

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
