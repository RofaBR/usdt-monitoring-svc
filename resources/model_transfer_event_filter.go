/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type TransferEventFilter struct {
	// Amount of the transaction
	Amount *string `json:"amount,omitempty"`
	// Address involved in the transaction
	Counterparty *string `json:"counterparty,omitempty"`
	// Sender's address
	From *string `json:"from,omitempty"`
	// Recipient's address
	To *string `json:"to,omitempty"`
}
