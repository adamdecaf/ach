// Copyright 2018 The ACH Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package ach

import (
	"testing"

	"github.com/moov-io/ach/errors"
)

// mockFilePPD creates an ACH file with PPD batch and entry
func mockFilePPD() *File {
	mockFile := NewFile()
	mockFile.SetHeader(mockFileHeader())
	mockBatch := mockBatchPPD()
	mockFile.AddBatch(mockBatch)
	if err := mockFile.Create(); err != nil {
		panic(err)
	}
	return mockFile
}

// testFileBatchCount validates if calculated count is different from control
func testFileBatchCount(t testing.TB) {
	file := mockFilePPD()

	// More batches than the file control count.
	file.AddBatch(mockBatchPPD())
	if err := file.Validate(); err != nil {
		if !errors.Contains(err, "BatchCount") {
			t.Errorf("expected BatchCount, but got %v", err)
		}
	}
}

// TestFileBatchCount tests validating if calculated count is different from control
func TestFileBatchCount(t *testing.T) {
	testFileBatchCount(t)
}

// BenchmarkFileBatchCount benchmarks validating if calculated count is different from control
func BenchmarkFileBatchCount(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testFileBatchCount(b)
	}
}

// testFileEntryAddenda validates an addenda entry
func testFileEntryAddenda(t testing.TB) {
	file := mockFilePPD()

	// more entries than the file control
	file.Control.EntryAddendaCount = 5
	if err := file.Validate(); err != nil {
		if !errors.Contains(err, "EntryAddendaCount") {
			t.Error(err.Error())
		}
	}
}

// TestFileEntryAddenda tests validating an addenda entry
func TestFileEntryAddenda(t *testing.T) {
	testFileEntryAddenda(t)
}

// BenchmarkFileEntryAddenda benchmarks validating an addenda entry
func BenchmarkFileEntryAddenda(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testFileEntryAddenda(b)
	}
}

// testFileDebitAmount validates file total debit amount
func testFileDebitAmount(t testing.TB) {
	file := mockFilePPD()

	// inequality in total debit amount
	file.Control.TotalDebitEntryDollarAmountInFile = 63
	if err := file.Validate(); err != nil {
		if !errors.Contains(err, "TotalDebitEntryDollarAmountInFile") {
			t.Errorf("expected TotalDebitEntryDollarAmountInFile, but got: %v", err)
		}
	}
}

// TestFileDebitAmount tests validating file total debit amount
func TestFileDebitAmount(t *testing.T) {
	testFileDebitAmount(t)
}

// BenchmarkFileDebitAmount benchmarks validating file total debit amount
func BenchmarkFileDebitAmount(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testFileDebitAmount(b)
	}
}

// testFileCreditAmount validates file total credit amount
func testFileCreditAmount(t testing.TB) {
	file := mockFilePPD()

	// inequality in total credit amount
	file.Control.TotalCreditEntryDollarAmountInFile = 63
	if err := file.Validate(); err != nil {
		if !errors.Contains(err, "TotalCreditEntryDollarAmountInFile") {
			t.Errorf("expected TotalCreditEntryDollarAmountInFile, but got: %v", err)
		}
	}
}

// TestFileCreditAmount tests validating file total credit amount
func TestFileCreditAmount(t *testing.T) {
	testFileCreditAmount(t)
}

// BenchmarkFileCreditAmount benchmarks validating file total credit amount
func BenchmarkFileCreditAmount(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testFileCreditAmount(b)
	}
}

// testFileEntryHash validates entry hash
func testFileEntryHash(t testing.TB) {
	file := mockFilePPD()
	file.AddBatch(mockBatchPPD())
	file.Create()
	file.Control.EntryHash = 63
	if err := file.Validate(); err != nil {
		if !errors.Contains(err, "EntryHash") {
			t.Errorf("expected EntryHash, but got: %v", err)
		}
	}
}

// TestFileEntryHash tests validating entry hash
func TestFileEntryHash(t *testing.T) {
	testFileEntryHash(t)
}

// BenchmarkFileEntryHash benchmarks validating entry hash
func BenchmarkFileEntryHash(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testFileEntryHash(b)
	}
}

// testFileBlockCount10 validates file block count
func testFileBlockCount10(t testing.TB) {
	file := NewFile().SetHeader(mockFileHeader())
	batch := NewBatchPPD(mockBatchPPDHeader())
	batch.AddEntry(mockEntryDetail())
	batch.AddEntry(mockEntryDetail())
	batch.AddEntry(mockEntryDetail())
	batch.AddEntry(mockEntryDetail())
	batch.AddEntry(mockEntryDetail())
	batch.AddEntry(mockEntryDetail())
	batch.Create()
	file.AddBatch(batch)
	if err := file.Create(); err != nil {
		t.Errorf("%T: %s", err, err)
	}

	// ensure with 10 records in file we don't get 2 for a block count
	if file.Control.BlockCount != 1 {
		t.Error("BlockCount on 10 records is not equal to 1")
	}
	// make 11th record which should produce BlockCount of 2
	file.Batches[0].AddEntry(mockEntryDetail())
	file.Batches[0].Create() // File.Build does not re-build Batches
	if err := file.Create(); err != nil {
		t.Errorf("%T: %s", err, err)
	}
	if file.Control.BlockCount != 2 {
		t.Error("BlockCount on 11 records is not equal to 2")
	}
}

// TestFileBlockCount10 tests validating file block count
func TestFileBlockCount10(t *testing.T) {
	testFileBlockCount10(t)
}

// BenchmarkFileBlockCount10 benchmarks validating file block count
func BenchmarkFileBlockCount10(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testFileBlockCount10(b)
	}
}

// testFileBuildBadFileHeader validates a bad file header
func testFileBuildBadFileHeader(t testing.TB) {
	file := NewFile().SetHeader(FileHeader{})
	if err := file.Create(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if e.Msg != msgFieldInclusion {
				t.Errorf("%T: %s", err, err)
			}
		} else {
			t.Errorf("%T: %s", err, err)
		}
	}
}

// TestFileBuildBadFileHeader tests validating a bad file header
func TestFileBuildBadFileHeader(t *testing.T) {
	testFileBuildBadFileHeader(t)
}

// BenchmarkFileBuildBadFileHeader benchmarks validating a bad file header
func BenchmarkFileBuildBadFileHeader(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testFileBuildBadFileHeader(b)
	}
}

// testFileBuildNoBatch validates a file with no batches
func testFileBuildNoBatch(t testing.TB) {
	file := NewFile().SetHeader(mockFileHeader())
	if err := file.Create(); err != nil {
		if !errors.Contains(err, "Batches") {
			t.Errorf("expected Batches, but got: %v", err)
		}
	}
}

// TestFileBuildNoBatch tests validating a file with no batches
func TestFileBuildNoBatch(t *testing.T) {
	testFileBuildNoBatch(t)
}

// BenchmarkFileBuildNoBatch benchmarks validating a file with no batches
func BenchmarkFileBuildNoBatch(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testFileBuildNoBatch(b)
	}
}

// testFileNotificationOfChange validates if a file contains
// BatchNOC notification of change
func testFileNotificationOfChange(t testing.TB) {
	file := NewFile().SetHeader(mockFileHeader())
	file.AddBatch(mockBatchPPD())
	bCOR := mockBatchCOR()
	file.AddBatch(bCOR)
	file.Create()

	if file.NotificationOfChange[0] != bCOR {
		t.Error("BatchCOR added to File.AddBatch should exist in NotificationOfChange")
	}
}

// TestFileNotificationOfChange tests validating if a file contains
// BatchNOC notification of change
func TestFileNotificationOfChange(t *testing.T) {
	testFileNotificationOfChange(t)
}

// BenchmarkFileNotificationOfChange benchmarks validating if a file contains
// BatchNOC notification of change
func BenchmarkFileNotificationOfChange(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testFileNotificationOfChange(b)
	}
}

// testFileReturnEntries validates file return entries
func testFileReturnEntries(t testing.TB) {
	// create or copy the entry to be returned record
	entry := mockEntryDetail()
	// Add the addenda return with appropriate ReturnCode and addenda information
	entry.AddAddenda(mockAddenda99())
	// create or copy the previous batch header of the item being returned
	batchHeader := mockBatchHeader()
	// create or copy the batch to be returned

	//batch, err := NewBatch(BatchParam{StandardEntryClass: batchHeader.StandardEntryClassCode})
	batch, err := NewBatch(batchHeader)

	if err != nil {
		t.Error(err.Error())
	}
	// Add the entry to be returned to the batch
	batch.AddEntry(entry)
	// Create the batch
	batch.Create()
	// Add the batch to your file.
	file := NewFile().SetHeader(mockFileHeader())
	file.AddBatch(batch)
	// Create the return file
	if err := file.Create(); err != nil {
		t.Error(err.Error())
	}

	if len(file.ReturnEntries) != 1 {
		t.Errorf("1 file.ReturnEntries added and %v exist", len(file.ReturnEntries))
	}
}

// TestFileReturnEntries tests validating file return entries
func TestFileReturnEntries(t *testing.T) {
	testFileReturnEntries(t)
}

// BenchmarkFileReturnEntries benchmarks validating file return entries
func BenchmarkFileReturnEntries(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testFileReturnEntries(b)
	}
}
