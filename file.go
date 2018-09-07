// Copyright 2018 The ACH Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package ach

import (
	"fmt"
	"strconv"

	"github.com/moov-io/ach/errors"
)

// First position of all Record Types. These codes are uniquely assigned to
// the first byte of each row in a file.
const (
	fileHeaderPos   = "1"
	batchHeaderPos  = "5"
	entryDetailPos  = "6"
	entryAddendaPos = "7"
	batchControlPos = "8"
	fileControlPos  = "9"

	// RecordLength character count of each line representing a letter in a file
	RecordLength = 94
)

// Errors strings specific to parsing a Batch container
var (
	msgFileCalculatedControlEquality = "calculated %v is out-of-balance with control %v"
	// specific messages
	msgRecordLength      = "must be 94 characters and found %d"
	msgFileBatchOutside  = "outside of current batch"
	msgFileBatchInside   = "inside of current batch"
	msgFileControl       = "none or more than one file control exists"
	msgFileHeader        = "none or more than one file headers exists"
	msgUnknownRecordType = "%s is an unknown record type"
	msgFileNoneSEC       = "%v Standard Entry Class Code is not implemented"
	msgFileIATSEC        = "%v Standard Entry Class Code should use iatBatch"
)

// File contains the structures of a parsed ACH File.
type File struct {
	ID         string      `json:"id"`
	Header     FileHeader  `json:"fileHeader"`
	Batches    []Batcher   `json:"batches"`
	IATBatches []IATBatch  `json:"IATBatches"`
	Control    FileControl `json:"fileControl"`

	// NotificationOfChange (Notification of change) is a slice of references to BatchCOR in file.Batches
	NotificationOfChange []*BatchCOR
	// ReturnEntries is a slice of references to file.Batches that contain return entries
	ReturnEntries []Batcher

	converters
}

// NewFile constructs a file template.
func NewFile() *File {
	return &File{
		Header:  NewFileHeader(),
		Control: NewFileControl(),
	}
}

// Create creates a valid file and requires that the FileHeader and at least one Batch
func (f *File) Create() error {
	// Requires a valid FileHeader to build FileControl
	if err := f.Header.Validate(); err != nil {
		return err
	}
	// Requires at least one Batch in the new file.
	if len(f.Batches) <= 0 && len(f.IATBatches) <= 0 {
		return errors.File("Batches", strconv.Itoa(len(f.Batches)), "must have []*Batches to be built")
	}
	// add 2 for FileHeader/control and reset if build was called twice do to error
	totalRecordsInFile := 2
	batchSeq := 1
	fileEntryAddendaCount := 0
	fileEntryHashSum := 0
	totalDebitAmount := 0
	totalCreditAmount := 0
	for i, batch := range f.Batches {
		// create ascending batch numbers
		f.Batches[i].GetHeader().BatchNumber = batchSeq
		f.Batches[i].GetControl().BatchNumber = batchSeq
		batchSeq++
		// sum file entry and addenda records. Assume batch.Create() batch properly calculated control
		fileEntryAddendaCount = fileEntryAddendaCount + batch.GetControl().EntryAddendaCount
		// add 2 for Batch header/control + entry added count
		totalRecordsInFile = totalRecordsInFile + 2 + batch.GetControl().EntryAddendaCount
		// sum hash from batch control. Assume Batch.Build properly calculated field.
		fileEntryHashSum = fileEntryHashSum + batch.GetControl().EntryHash
		totalDebitAmount = totalDebitAmount + batch.GetControl().TotalDebitEntryDollarAmount
		totalCreditAmount = totalCreditAmount + batch.GetControl().TotalCreditEntryDollarAmount

	}
	for i, iatBatch := range f.IATBatches {
		// create ascending batch numbers
		f.IATBatches[i].GetHeader().BatchNumber = batchSeq
		f.IATBatches[i].GetControl().BatchNumber = batchSeq
		batchSeq++
		// sum file entry and addenda records. Assume batch.Create() batch properly calculated control
		fileEntryAddendaCount = fileEntryAddendaCount + iatBatch.GetControl().EntryAddendaCount
		// add 2 for Batch header/control + entry added count
		totalRecordsInFile = totalRecordsInFile + 2 + iatBatch.GetControl().EntryAddendaCount
		// sum hash from batch control. Assume Batch.Build properly calculated field.
		fileEntryHashSum = fileEntryHashSum + iatBatch.GetControl().EntryHash
		totalDebitAmount = totalDebitAmount + iatBatch.GetControl().TotalDebitEntryDollarAmount
		totalCreditAmount = totalCreditAmount + iatBatch.GetControl().TotalCreditEntryDollarAmount

	}
	// create FileControl from calculated values
	fc := NewFileControl()
	fc.BatchCount = batchSeq - 1
	// blocking factor of 10 is static default value in f.Header.blockingFactor.
	if (totalRecordsInFile % 10) != 0 {
		fc.BlockCount = totalRecordsInFile/10 + 1
	} else {
		fc.BlockCount = totalRecordsInFile / 10
	}
	fc.EntryAddendaCount = fileEntryAddendaCount
	fc.EntryHash = fileEntryHashSum
	fc.TotalDebitEntryDollarAmountInFile = totalDebitAmount
	fc.TotalCreditEntryDollarAmountInFile = totalCreditAmount
	f.Control = fc

	return nil
}

// AddBatch appends a Batch to the ach.File
func (f *File) AddBatch(batch Batcher) []Batcher {
	switch batch.(type) {
	case *BatchCOR:
		f.NotificationOfChange = append(f.NotificationOfChange, batch.(*BatchCOR))
	}
	if batch.Category() == CategoryReturn {
		f.ReturnEntries = append(f.ReturnEntries, batch)
	}
	f.Batches = append(f.Batches, batch)
	return f.Batches
}

// AddIATBatch appends a IATBatch to the ach.File
func (f *File) AddIATBatch(iatBatch IATBatch) []IATBatch {
	f.IATBatches = append(f.IATBatches, iatBatch)
	return f.IATBatches
}

// SetHeader allows for header to be built.
func (f *File) SetHeader(h FileHeader) *File {
	f.Header = h
	return f
}

// Validate NACHA rules on the entire batch before being added to a File
func (f *File) Validate() error {
	// The value of the Batch Count Field is equal to the number of Company/Batch/Header Records in the file.
	if f.Control.BatchCount != (len(f.Batches) + len(f.IATBatches)) {
		msg := fmt.Sprintf(msgFileCalculatedControlEquality, len(f.Batches), f.Control.BatchCount)
		return errors.File("BatchCount", strconv.Itoa(len(f.Batches)), msg)
	}

	if err := f.isEntryAddendaCount(); err != nil {
		return errors.WithError(err)
	}
	if err := f.isFileAmount(); err != nil {
		return errors.WithError(err)
	}
	return errors.WithError(f.isEntryHash())
}

// isEntryAddendaCount is prepared by hashing the RDFI’s 8-digit Routing Number in each entry.
//The Entry Hash provides a check against inadvertent alteration of data
func (f *File) isEntryAddendaCount() error {
	count := 0
	// we assume that each batch block has already validated the addenda count is accurate in batch control.
	for _, batch := range f.Batches {
		count += batch.GetControl().EntryAddendaCount
	}
	// IAT
	for _, iatBatch := range f.IATBatches {
		count += iatBatch.GetControl().EntryAddendaCount
	}
	if f.Control.EntryAddendaCount != count {
		msg := fmt.Sprintf(msgFileCalculatedControlEquality, count, f.Control.EntryAddendaCount)
		return errors.File("EntryAddendaCount", f.Control.EntryAddendaCountField(), msg)
	}
	return nil
}

// isFileAmount tThe Total Debit and Credit Entry Dollar Amounts Fields contain accumulated
// Entry Detail debit and credit totals within the file
func (f *File) isFileAmount() error {
	debit := 0
	credit := 0
	for _, batch := range f.Batches {
		debit += batch.GetControl().TotalDebitEntryDollarAmount
		credit += batch.GetControl().TotalCreditEntryDollarAmount
	}
	// IAT
	for _, iatBatch := range f.IATBatches {
		debit += iatBatch.GetControl().TotalDebitEntryDollarAmount
		credit += iatBatch.GetControl().TotalCreditEntryDollarAmount
	}
	if f.Control.TotalDebitEntryDollarAmountInFile != debit {
		msg := fmt.Sprintf(msgFileCalculatedControlEquality, debit, f.Control.TotalDebitEntryDollarAmountInFile)
		return errors.File("TotalDebitEntryDollarAmountInFile", f.Control.TotalDebitEntryDollarAmountInFileField(), msg)
	}
	if f.Control.TotalCreditEntryDollarAmountInFile != credit {
		msg := fmt.Sprintf(msgFileCalculatedControlEquality, credit, f.Control.TotalCreditEntryDollarAmountInFile)
		return errors.File("TotalCreditEntryDollarAmountInFile", f.Control.TotalCreditEntryDollarAmountInFileField(), msg)
	}
	return nil
}

// isEntryHash validates the hash by recalculating the result
func (f *File) isEntryHash() error {
	hashField := f.calculateEntryHash()
	if hashField != f.Control.EntryHashField() {
		msg := fmt.Sprintf(msgFileCalculatedControlEquality, hashField, f.Control.EntryHashField())
		return errors.File("EntryHash", f.Control.EntryHashField(), msg)
	}
	return nil
}

// calculateEntryHash This field is prepared by hashing the 8-digit Routing Number in each batch.
// The Entry Hash provides a check against inadvertent alteration of data
func (f *File) calculateEntryHash() string {
	hash := 0
	for _, batch := range f.Batches {
		hash = hash + batch.GetControl().EntryHash
	}
	// IAT
	for _, iatBatch := range f.IATBatches {
		hash = hash + iatBatch.GetControl().EntryHash
	}
	return f.numericField(hash, 10)
}
