/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources
import "gitlab.com/tokend/go/xdr"
import (
	"time"
)

type TransferEvent struct {
	// Number of tokens
	Amount string `json:"amount"`
	// Block number
	BlockNumber *int32 `json:"blockNumber,omitempty"`
	// Sender's address
	From string `json:"from"`
	// Transaction creation time
	Timestamp *date-time `json:"timestamp,omitempty"`
	// Recipient's address
	To string `json:"to"`
	// Transaction hash
	TransactionHash string `json:"transactionHash"`
}
