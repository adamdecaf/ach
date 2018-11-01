// Copyright 2018 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package ach

import "fmt"

// BatchRCK holds the BatchHeader and BatchControl and all EntryDetail for RCK Entries.
//
// Represented Check Entries (RCK). A physical check that was presented but returned because of
// insufficient funds may be represented as an ACH entry.
type BatchRCK struct {
	batch
}

// NewBatchRCK returns a *BatchRCK
func NewBatchRCK(bh *BatchHeader) *BatchRCK {
	batch := new(BatchRCK)
	batch.SetControl(NewBatchControl())
	batch.SetHeader(bh)
	return batch
}

// Validate checks valid NACHA batch rules. Assumes properly parsed records.
func (batch *BatchRCK) Validate() error {
	// basic verification of the batch before we validate specific rules.
	if err := batch.verify(); err != nil {
		return err
	}

	// Add configuration and type specific validation for this type.
	if batch.Header.StandardEntryClassCode != "RCK" {
		msg := fmt.Sprintf(msgBatchSECType, batch.Header.StandardEntryClassCode, "RCK")
		return &BatchError{BatchNumber: batch.Header.BatchNumber, FieldName: "StandardEntryClassCode", Msg: msg}
	}

	// RCK detail entries can only be a debit, ServiceClassCode must allow debits
	switch batch.Header.ServiceClassCode {
	case 200, 220, 280:
		msg := fmt.Sprintf(msgBatchServiceClassCode, batch.Header.ServiceClassCode, "RCK")
		return &BatchError{BatchNumber: batch.Header.BatchNumber, FieldName: "ServiceClassCode", Msg: msg}
	}

	// CompanyEntryDescription is required to be REDEPCHECK
	if batch.Header.CompanyEntryDescription != "REDEPCHECK" {
		msg := fmt.Sprintf(msgBatchCompanyEntryDescription, batch.Header.CompanyEntryDescription, "RCK")
		return &BatchError{BatchNumber: batch.Header.BatchNumber, FieldName: "CompanyEntryDescription", Msg: msg}
	}

	for _, entry := range batch.Entries {
		// RCK detail entries must be a debit
		if entry.CreditOrDebit() != "D" {
			msg := fmt.Sprintf(msgBatchTransactionCodeCredit, entry.TransactionCode)
			return &BatchError{BatchNumber: batch.Header.BatchNumber, FieldName: "TransactionCode", Msg: msg}
		}

		// // Amount must be 2,500 or less
		if entry.Amount > 250000 {
			msg := fmt.Sprintf(msgBatchAmount, "2,500", "RCK")
			return &BatchError{BatchNumber: batch.Header.BatchNumber, FieldName: "Amount", Msg: msg}
		}

		// CheckSerialNumber underlying IdentificationNumber, must be defined
		if entry.IdentificationNumber == "" {
			msg := fmt.Sprintf(msgBatchCheckSerialNumber, "RCK")
			return &BatchError{BatchNumber: batch.Header.BatchNumber, FieldName: "CheckSerialNumber", Msg: msg}
		}
		// RCK cannot have Addenda02 or Addenda05.  There can be a NOC (98) or Return (99)
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
func (batch *BatchRCK) Create() error {
	// generates sequence numbers and batch control
	if err := batch.build(); err != nil {
		return err
	}
	// Additional steps specific to batch type
	// ...

	return batch.Validate()
}
