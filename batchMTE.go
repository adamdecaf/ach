// Licensed to The Moov Authors under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. The Moov Authors licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package ach

import (
	"strings"

	"github.com/moov-io/ach/internal/usabbrev"
)

// BatchMTE holds the BatchHeader, BatchControl, and EntryDetail for Machine Transfer Entry (MTE) entries.
//
// A MTE transaction is created when a consumer uses their debit card at an Automated Teller Machine (ATM) to withdraw cash.
// MTE transactions cannot be aggregated together under a single Entry.
type BatchMTE struct {
	Batch
}

// NewBatchMTE returns a *BatchMTE
func NewBatchMTE(bh *BatchHeader) *BatchMTE {
	batch := new(BatchMTE)
	batch.SetControl(NewBatchControl())
	batch.SetHeader(bh)
	batch.SetID(bh.ID)
	return batch
}

// Validate checks properties of the ACH batch to ensure they match NACHA guidelines.
// This includes computing checksums, totals, and sequence orderings.
//
// Validate will never modify the batch.
func (batch *BatchMTE) Validate() error {
	if batch.validateOpts != nil && batch.validateOpts.SkipAll {
		return nil
	}

	// basic verification of the batch before we validate specific rules.
	if err := batch.verify(); err != nil {
		return err
	}

	if batch.Header.StandardEntryClassCode != MTE {
		return batch.Error("StandardEntryClassCode", ErrBatchSECType, MTE)
	}

	invalidEntries := batch.InvalidEntries()
	if len(invalidEntries) > 0 {
		return invalidEntries[0].Error // return the first invalid entry's error
	}

	return nil
}

// InvalidEntries returns entries with validation errors in the batch
func (batch *BatchMTE) InvalidEntries() []InvalidEntry {
	var out []InvalidEntry

	for _, entry := range batch.Entries {
		if entry.Amount <= 0 {
			out = append(out, InvalidEntry{
				Entry: entry,
				Error: batch.Error("Amount", ErrBatchAmountZero, entry.Amount),
			})
		}
		// Verify the Amount is valid for SEC code and TransactionCode
		if err := batch.ValidAmountForCodes(entry); err != nil {
			out = append(out, InvalidEntry{
				Entry: entry,
				Error: err,
			})
		}
		// Verify the TransactionCode is valid for a ServiceClassCode
		if err := batch.ValidTranCodeForServiceClassCode(entry); err != nil {
			out = append(out, InvalidEntry{
				Entry: entry,
				Error: err,
			})
		}
		// Verify Addenda* FieldInclusion based on entry.Category and batchHeader.StandardEntryClassCode
		if err := batch.addendaFieldInclusion(entry); err != nil {
			out = append(out, InvalidEntry{
				Entry: entry,
				Error: err,
			})
		}
		if entry.Category == CategoryForward {
			if !usabbrev.Valid(entry.Addenda02.TerminalState) {
				out = append(out, InvalidEntry{
					Entry: entry,
					Error: batch.Error("TerminalState", ErrValidState, entry.Addenda02.TerminalState),
				})
			}
		}

		// MTE entries cannot have an identification number that is all spaces or all zeros
		if strings.Trim(entry.IdentificationNumber, " 0") == "" {
			out = append(out, InvalidEntry{
				Entry: entry,
				Error: batch.Error("IdentificationNumber", ErrIdentificationNumber, entry.IdentificationNumber),
			})
		}
	}

	return out
}

// Create will tabulate and assemble an ACH batch into a valid state. This includes
// setting any posting dates, sequence numbers, counts, and sums.
//
// Create implementations are free to modify computable fields in a file and should
// call the Batch's Validate function at the end of their execution.
func (batch *BatchMTE) Create() error {
	// generates sequence numbers and batch control
	if err := batch.build(); err != nil {
		return err
	}
	return batch.Validate()
}
