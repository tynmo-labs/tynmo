package subscribe

import (
	"bytes"
	"fmt"

	"tynmo/command/helper"
	txpoolProto "tynmo/txpool/proto"
)

type TxPoolEventResult struct {
	EventType txpoolProto.EventType `json:"event_type"`
	TxHash    string                `json:"tx_hash"`
}

func (r *TxPoolEventResult) GetOutput() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n[TXPOOL EVENT]\n")
	buffer.WriteString(helper.FormatKV([]string{
		fmt.Sprintf("TYPE|%s", r.EventType),
		fmt.Sprintf("HASH|%s", r.TxHash),
	}))
	buffer.WriteString("\n")

	return buffer.String()
}
