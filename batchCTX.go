// Copyright 2018 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package ach

import (
	"fmt"
	"strconv"
)

// BatchCTX holds the BatchHeader and BatchControl and all EntryDetail for CTX Entries.
//
// The Corporate Trade Exchange (CTX) application provides the ability to collect and disburse
// funds and information between companies. Generally it is used by businesses paying one another
// for goods or services. These payments replace checks with an electronic process of debiting and
// crediting invoices between the financial institutions of participating companies.
type BatchCTX struct {
	batch
}

var (
	msgBatchCTXAddendaCount = "%v entry detail addenda records not equal to addendum %v"
)

// NewBatchCTX returns a *BatchCTX
func NewBatchCTX(bh *BatchHeader) *BatchCTX {
	batch := new(BatchCTX)
	batch.SetControl(NewBatchControl())
	batch.SetHeader(bh)
	return batch
}

// Validate checks valid NACHA batch rules. Assumes properly parsed records.
func (batch *BatchCTX) Validate() error {
	// basic verification of the batch before we validate specific rules.
	if err := batch.verify(); err != nil {
		return err
	}

	// Add configuration and type specific validation for this type.
	if batch.Header.StandardEntryClassCode != "CTX" {
		msg := fmt.Sprintf(msgBatchSECType, batch.Header.StandardEntryClassCode, "CTX")
		return &BatchError{BatchNumber: batch.Header.BatchNumber, FieldName: "StandardEntryClassCode", Msg: msg}
	}

	for _, entry := range batch.Entries {

		// Trapping this error, as entry.CTXAddendaRecordsField() can not be greater than 9999
		if len(entry.Addendum) > 9999 {
			msg := fmt.Sprintf(msgBatchAddendaCount, len(entry.Addendum), 9999, batch.Header.StandardEntryClassCode)
			return &BatchError{BatchNumber: batch.Header.BatchNumber, FieldName: "AddendaCount", Msg: msg}
		}

		// validate CTXAddendaRecord Field is equal to the actual number of Addenda records
		// use 0 value if there is no Addenda records
		addendaRecords, _ := strconv.Atoi(entry.CATXAddendaRecordsField())
		if len(entry.Addendum) != addendaRecords {
			msg := fmt.Sprintf(msgBatchCTXAddendaCount, addendaRecords, len(entry.Addendum))
			return &BatchError{BatchNumber: batch.Header.BatchNumber, FieldName: "Addendum", Msg: msg}
		}

		if len(entry.Addendum) > 0 {

			switch entry.TransactionCode {
			// Prenote credit  23, 33, 43, 53
			// Prenote debit 28, 38, 48
			case 23, 28, 33, 38, 43, 48, 53:
				msg := fmt.Sprintf(msgBatchTransactionCodeAddenda, entry.TransactionCode, "CTX")
				return &BatchError{BatchNumber: batch.Header.BatchNumber, FieldName: "Addendum", Msg: msg}
			default:
			}

			// CTX can have up to 9999 Addenda Record TypeCode = 05, or there can be a NOC (98) or Return (99)
			for _, addenda := range entry.Addendum {
				switch entry.Category {
				case CategoryForward:
					if err := batch.categoryForwardAddenda05(entry, addenda); err != nil {
						return err
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
	}
	return nil
}

// Create takes Batch Header and Entries and builds a valid batch
func (batch *BatchCTX) Create() error {
	// generates sequence numbers and batch control
	if err := batch.build(); err != nil {
		return err
	}
	// Additional steps specific to batch type
	// ...
	return batch.Validate()
}
