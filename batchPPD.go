// Copyright 2018 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package ach

import "fmt"

// BatchPPD holds the Batch Header and Batch Control and all Entry Records for PPD Entries
type BatchPPD struct {
	batch
}

// NewBatchPPD returns a *BatchPPD
func NewBatchPPD(bh *BatchHeader) *BatchPPD {
	batch := new(BatchPPD)
	batch.SetControl(NewBatchControl())
	batch.SetHeader(bh)
	return batch
}

// Validate checks valid NACHA batch rules. Assumes properly parsed records.
func (batch *BatchPPD) Validate() error {
	// basic verification of the batch before we validate specific rules.
	if err := batch.verify(); err != nil {
		return err
	}
	// Add configuration and type specific validation for this type.

	if batch.Header.StandardEntryClassCode != "PPD" {
		msg := fmt.Sprintf(msgBatchSECType, batch.Header.StandardEntryClassCode, "PPD")
		return &BatchError{BatchNumber: batch.Header.BatchNumber, FieldName: "StandardEntryClassCode", Msg: msg}
	}

	for _, entry := range batch.Entries {
		// PPD can have up to one Record TypeCode = 05, or there can be a NOC (98) or Return (99)
		for _, addenda := range entry.Addendum {
			switch entry.Category {
			case CategoryForward:
				if err := batch.categoryForwardAddenda05(entry, addenda); err != nil {
					return err
				}
				if len(entry.Addendum) > 1 {
					msg := fmt.Sprintf(msgBatchAddendaCount, len(entry.Addendum), 1, batch.Header.StandardEntryClassCode)
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
func (batch *BatchPPD) Create() error {
	// generates sequence numbers and batch control
	if err := batch.build(); err != nil {
		return err
	}
	// Additional steps specific to batch type
	// ...

	return batch.Validate()
}
