// Copyright 2018 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package ach

import "fmt"

// BatchARC holds the BatchHeader and BatchControl and all EntryDetail for ARC Entries.
//
// Accounts Receivable Entry (ARC). A consumer check converted to a one-time ACH debit.
// The Accounts Receivable (ARC) Entry provides billers the opportunity to initiate single-entry ACH
// debits to customer accounts by converting checks at the point of receipt through the U.S. mail, at
// a drop box location or in-person for payment of a bill at a manned location. The biller is required
// to provide the customer with notice prior to the acceptance of the check that states the receipt of
// the customer’s check will be deemed as the authorization for an ARC debit entry to the customer’s
// account. The provision of the notice and the receipt of the check together constitute authorization
// for the ARC entry. The customer’s check is solely be used as a source document to obtain the routing
// number, account number and check serial number.
//
// The difference between ARC and POP is that ARC can result from a check mailed in whereas POP is in-person.
type BatchARC struct {
	batch
}

// NewBatchARC returns a *BatchARC
func NewBatchARC(bh *BatchHeader) *BatchARC {
	batch := new(BatchARC)
	batch.SetControl(NewBatchControl())
	batch.SetHeader(bh)
	return batch
}

// Validate checks valid NACHA batch rules. Assumes properly parsed records.
func (batch *BatchARC) Validate() error {
	// basic verification of the batch before we validate specific rules.
	if err := batch.verify(); err != nil {
		return err
	}
	// Add configuration and type specific validation for this type.

	if batch.Header.StandardEntryClassCode != "ARC" {
		msg := fmt.Sprintf(msgBatchSECType, batch.Header.StandardEntryClassCode, "ARC")
		return &BatchError{BatchNumber: batch.Header.BatchNumber, FieldName: "StandardEntryClassCode", Msg: msg}
	}

	// ARC detail entries can only be a debit, ServiceClassCode must allow debits
	switch batch.Header.ServiceClassCode {
	case 200, 220, 280:
		msg := fmt.Sprintf(msgBatchServiceClassCode, batch.Header.ServiceClassCode, "ARC")
		return &BatchError{BatchNumber: batch.Header.BatchNumber, FieldName: "ServiceClassCode", Msg: msg}
	}

	for _, entry := range batch.Entries {
		// ARC detail entries must be a debit
		if entry.CreditOrDebit() != "D" {
			msg := fmt.Sprintf(msgBatchTransactionCodeCredit, entry.TransactionCode)
			return &BatchError{BatchNumber: batch.Header.BatchNumber, FieldName: "TransactionCode", Msg: msg}
		}

		// Amount must be 25,000 or less
		if entry.Amount > 2500000 {
			msg := fmt.Sprintf(msgBatchAmount, "25,000", "ARC")
			return &BatchError{BatchNumber: batch.Header.BatchNumber, FieldName: "Amount", Msg: msg}
		}

		// CheckSerialNumber underlying IdentificationNumber, must be defined
		if entry.IdentificationNumber == "" {
			msg := fmt.Sprintf(msgBatchCheckSerialNumber, "ARC")
			return &BatchError{BatchNumber: batch.Header.BatchNumber, FieldName: "CheckSerialNumber", Msg: msg}
		}
		// ARC cannot have Addenda02 or Addenda05.  There can be a NOC (98) or Return (99)
		for _, addenda := range entry.Addendum {
			switch entry.Category {
			case CategoryForward:
				if len(entry.Addendum) > 0 {
					msg := fmt.Sprintf(msgBatchAddendaCount, len(entry.Addendum), 0, batch.Header.StandardEntryClassCode)
					return &BatchError{BatchNumber: batch.Header.BatchNumber, FieldName: "AddendaCount", Msg: msg}
				}
			case CategoryNOC:
				if err := batch.categoryNOCAddenda98(entry, addenda); err != nil {
					return err
				}
			case CategoryReturn:
				if err := batch.categoryReturnAddenda99(entry, addenda); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Create takes Batch Header and Entries and builds a valid batch
func (batch *BatchARC) Create() error {
	// generates sequence numbers and batch control
	if err := batch.build(); err != nil {
		return err
	}
	// Additional steps specific to batch type
	// ...

	return batch.Validate()
}
