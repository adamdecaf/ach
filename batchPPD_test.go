// Copyright 2018 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package ach

import (
	"fmt"
	"testing"
	"time"
)

// mockBatchPPDHeader creates a PPD batch header
func mockBatchPPDHeader() *BatchHeader {
	bh := NewBatchHeader()
	bh.ServiceClassCode = 220
	bh.StandardEntryClassCode = "PPD"
	bh.CompanyName = "ACME Corporation"
	bh.CompanyIdentification = "121042882"
	bh.CompanyEntryDescription = "PAYROLL"
	bh.EffectiveEntryDate = time.Now()
	bh.ODFIIdentification = "12104288"
	return bh
}

// mockPPDEntryDetail creates a PPD Entry Detail
func mockPPDEntryDetail() *EntryDetail {
	entry := NewEntryDetail()
	entry.TransactionCode = 22
	entry.SetRDFI("231380104")
	entry.DFIAccountNumber = "123456789"
	entry.Amount = 100000000
	entry.IndividualName = "Wade Arnold"
	entry.SetTraceNumber(mockBatchPPDHeader().ODFIIdentification, 1)
	entry.Category = CategoryForward
	entry.AddendaRecordIndicator = 0

	// addenda := NewAddenda05()
	// addenda.PaymentRelatedInformation = "EXAMPLE PAYMENT"
	// addenda.SequenceNumber = 1
	// addenda.EntryDetailSequenceNumber = 1
	// entry.AddAddenda05(addenda)
	return entry
}

// mockBatchPPDHeader2 creates a 2nd PPD batch header
func mockBatchPPDHeader2() *BatchHeader {
	bh := NewBatchHeader()
	bh.ServiceClassCode = 200
	bh.CompanyName = "MY BEST COMP."
	bh.CompanyDiscretionaryData = "INCLUDES OVERTIME"
	bh.CompanyIdentification = "121042882"
	bh.StandardEntryClassCode = "PPD"
	bh.CompanyEntryDescription = "PAYROLL"
	bh.EffectiveEntryDate = time.Now()
	bh.ODFIIdentification = "12104288"
	return bh
}

// mockPPDEntryDetail2 creates a 2nd PPD entry detail
func mockPPDEntryDetail2() *EntryDetail {
	entry := NewEntryDetail()
	entry.TransactionCode = 22 // ACH Credit
	entry.SetRDFI("231380104")
	entry.DFIAccountNumber = "62292250"         // account number
	entry.Amount = 100000                       // 1k dollars
	entry.IdentificationNumber = "658-888-2468" // Unique ID for payment
	entry.IndividualName = "Wade Arnold"
	entry.SetTraceNumber(mockBatchPPDHeader2().ODFIIdentification, 1)
	entry.Category = CategoryForward
	entry.AddendaRecordIndicator = 0

	// addenda := NewAddenda05()
	// addenda.PaymentRelatedInformation = "EXAMPLE PAYMENT"
	// addenda.SequenceNumber = 1
	// addenda.EntryDetailSequenceNumber = 1
	// entry.AddAddenda05(addenda)
	return entry
}

// mockBatchPPD creates a PPD batch
func mockBatchPPD() *BatchPPD {
	mockBatch := NewBatchPPD(mockBatchPPDHeader())
	mockBatch.AddEntry(mockPPDEntryDetail())
	if err := mockBatch.Create(); err != nil {
		panic(fmt.Sprintf("batch.Control=%#v, %v", mockBatch.Control, err))
	}
	return mockBatch
}

// testBatchError validates batch error handling
func testBatchError(t testing.TB) {
	err := &BatchError{BatchNumber: 1, FieldName: "mock", Msg: "test message"}
	if err.Error() != "BatchNumber 1 mock test message" {
		t.Error("BatchError Error has changed formatting")
	}
}

// TestBatchError tests validating batch error handling
func TestBatchError(t *testing.T) {
	testBatchError(t)
}

// BenchmarkBatchError benchmarks validating batch error handling
func BenchmarkBatchError(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testBatchError(b)
	}
}

// testBatchServiceClassCodeEquality validates service class code equality
func testBatchServiceClassCodeEquality(t testing.TB) {
	mockBatch := mockBatchPPD()
	mockBatch.GetControl().ServiceClassCode = 225
	if err := mockBatch.Validate(); err != nil {
		if e, ok := err.(*BatchError); ok {
			if e.FieldName != "ServiceClassCode" {
				t.Errorf("%T: %s", err, err)
			}
		} else {
			t.Errorf("%T: %s", err, err)
		}
	}
}

// TestBatchServiceClassCodeEquality tests validating service class code equality
func TestBatchServiceClassCodeEquality(t *testing.T) {
	testBatchServiceClassCodeEquality(t)
}

// BenchmarkBatchServiceClassCodeEquality benchmarks validating service class code equality
func BenchmarkBatchServiceClassCodeEquality(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testBatchServiceClassCodeEquality(b)
	}
}

// BatchPPDCreate validates batch create for an invalid service code
func testBatchPPDCreate(t testing.TB) {
	mockBatch := mockBatchPPD()
	// can not have default values in Batch Header to build batch
	mockBatch.GetHeader().ServiceClassCode = 0
	mockBatch.Create()
	if err := mockBatch.Validate(); err != nil {
		if e, ok := err.(*BatchError); ok {
			if e.FieldName != "ServiceClassCode" {
				t.Errorf("%T: %s", err, err)
			}
		} else {
			t.Errorf("%T: %s", err, err)
		}
	}
}

// TestBatchPPDCreate tests validating batch create for an invalid service code
func TestBatchPPDCreate(t *testing.T) {
	testBatchPPDCreate(t)
}

// BenchmarkBatchPPDCreate benchmarks validating batch create for an invalid service code
func BenchmarkBatchPPDCreate(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testBatchPPDCreate(b)
	}
}

// testBatchPPDTypeCode validates batch PPD type code
func testBatchPPDTypeCode(t testing.TB) {
	mockBatch := mockBatchPPD()
	// change an addendum to an invalid type code
	a := mockAddenda05()
	a.TypeCode = "63"
	mockBatch.GetEntries()[0].AddAddenda05(a)
	mockBatch.Create()
	if err := mockBatch.Validate(); err != nil {
		if e, ok := err.(*BatchError); ok {
			if e.FieldName != "TypeCode" {
				t.Errorf("%T: %s", err, err)
			}
		} else {
			t.Errorf("%T: %s", err, err)
		}
	}
}

// TestBatchPPDTypeCode tests validating batch PPD type code
func TestBatchPPDTypeCode(t *testing.T) {
	testBatchPPDTypeCode(t)
}

// BenchmarkBatchPPDTypeCode benchmarks validating batch PPD type code
func BenchmarkBatchPPDTypeCode(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testBatchPPDTypeCode(b)
	}
}

// testBatchCompanyIdentification validates batch PPD company identification
func testBatchCompanyIdentification(t testing.TB) {
	mockBatch := mockBatchPPD()
	mockBatch.GetControl().CompanyIdentification = "XYZ Inc"
	if err := mockBatch.Validate(); err != nil {
		if e, ok := err.(*BatchError); ok {
			if e.FieldName != "CompanyIdentification" {
				t.Errorf("%T: %s", err, err)
			}
		} else {
			t.Errorf("%T: %s", err, err)
		}
	}
}

// TestBatchCompanyIdentification tests validating batch PPD company identification
func TestBatchCompanyIdentification(t *testing.T) {
	testBatchCompanyIdentification(t)
}

// BenchmarkBatchCompanyIdentification benchmarks validating batch PPD company identification
func BenchmarkBatchCompanyIdentification(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testBatchCompanyIdentification(b)
	}
}

// testBatchODFIIDMismatch validates ODFIIdentification mismatch
func testBatchODFIIDMismatch(t testing.TB) {
	mockBatch := mockBatchPPD()
	mockBatch.GetControl().ODFIIdentification = "987654321"
	if err := mockBatch.Validate(); err != nil {
		if e, ok := err.(*BatchError); ok {
			if e.FieldName != "ODFIIdentification" {
				t.Errorf("%T: %s", err, err)
			}
		} else {
			t.Errorf("%T: %s", err, err)
		}
	}
}

// TestBatchODFIIDMismatch tests validating ODFIIdentification mismatch
func TestBatchODFIIDMismatch(t *testing.T) {
	testBatchODFIIDMismatch(t)
}

// BenchmarkBatchODFIIDMismatch benchmarks validating ODFIIdentification mismatch
func BenchmarkBatchODFIIDMismatch(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testBatchODFIIDMismatch(b)
	}
}

// testBatchBuild builds a PPD batch
func testBatchBuild(t testing.TB) {
	mockBatch := NewBatchPPD(mockBatchPPDHeader2())
	entry := mockPPDEntryDetail2()
	addenda05 := NewAddenda05()
	entry.AddendaRecordIndicator = 1
	entry.AddAddenda05(addenda05)
	mockBatch.AddEntry(entry)
	if err := mockBatch.Create(); err != nil {
		t.Errorf("%T: %s", err, err)
	}
}

// TestBatchBuild tests building a PPD batch
func TestBatchBuild(t *testing.T) {
	testBatchBuild(t)
}

// BenchmarkBatchBuild benchmarks building a PPD batch
func BenchmarkBatchBuild(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testBatchBuild(b)
	}
}

// testBatchPPDAddendaCount validates BatchPPD Addendum count of 2
func testBatchPPDAddendaCount(t testing.TB) {
	mockBatch := mockBatchPPD()
	mockBatch.Entries[0].AddendaRecordIndicator = 1
	mockBatch.GetEntries()[0].AddAddenda05(mockAddenda05())
	mockBatch.GetEntries()[0].AddAddenda05(mockAddenda05())
	mockBatch.Create()
	if err := mockBatch.Validate(); err != nil {
		if e, ok := err.(*BatchError); ok {
			if e.FieldName != "AddendaCount" {
				t.Errorf("%T: %s", err, err)
			}
		} else {
			t.Errorf("%T: %s", err, err)
		}
	}
}

// TestBatchPPDAddendaCount tests validating BatchPPD Addendum count of 2
func TestBatchPPDAddendaCount(t *testing.T) {
	testBatchPPDAddendaCount(t)
}

// BenchmarkBatchPPDAddendaCount benchmarks validating BatchPPD Addendum count of 2
func BenchmarkBatchPPDAddendaCount(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testBatchPPDAddendaCount(b)
	}
}

// TestBatchPPDAddendum98 validates Addenda98 returns an error
func TestBatchPPDAddendum98(t *testing.T) {
	mockBatch := NewBatchPPD(mockBatchPPDHeader())
	mockBatch.AddEntry(mockPPDEntryDetail())
	mockAddenda98 := mockAddenda98()
	mockAddenda98.TypeCode = "05"
	mockBatch.GetEntries()[0].Category = CategoryNOC
	mockBatch.GetEntries()[0].Addenda98 = mockAddenda98
	if err := mockBatch.Create(); err != nil {
		if e, ok := err.(*BatchError); ok {
			if e.FieldName != "TypeCode" {
				t.Errorf("%T: %s", err, err)
			}
		} else {
			t.Errorf("%T: %s", err, err)
		}
	}
}

// TestBatchPPDAddendum99 validates Addenda99 returns an error
func TestBatchPPDAddendum99(t *testing.T) {
	mockBatch := NewBatchPPD(mockBatchPPDHeader())
	mockBatch.AddEntry(mockPPDEntryDetail())
	mockAddenda99 := mockAddenda99()
	mockAddenda99.TypeCode = "05"
	mockBatch.GetEntries()[0].Category = CategoryReturn
	mockBatch.GetEntries()[0].Addenda99 = mockAddenda99
	if err := mockBatch.Create(); err != nil {
		if e, ok := err.(*BatchError); ok {
			if e.FieldName != "TypeCode" {
				t.Errorf("%T: %s", err, err)
			}
		} else {
			t.Errorf("%T: %s", err, err)
		}
	}
}

// TestBatchPPDSEC validates that the standard entry class code is PPD for batch Web
func TestBatchPPDSEC(t *testing.T) {
	mockBatch := mockBatchPPD()
	mockBatch.Header.StandardEntryClassCode = "RCK"
	if err := mockBatch.Validate(); err != nil {
		if e, ok := err.(*BatchError); ok {
			if e.FieldName != "StandardEntryClassCode" {
				t.Errorf("%T: %s", err, err)
			}
		} else {
			t.Errorf("%T: %s", err, err)
		}
	}
}

// TestBatchPPDValidTranCodeForServiceClassCode validates a transactionCode based on ServiceClassCode
func TestBatchPPDValidTranCodeForServiceClassCode(t *testing.T) {
	mockBatch := mockBatchPPD()
	mockBatch.GetHeader().ServiceClassCode = 225
	if err := mockBatch.Create(); err != nil {
		if e, ok := err.(*BatchError); ok {
			if e.FieldName != "TransactionCode" {
				t.Errorf("%T: %s", err, err)
			}
		} else {
			t.Errorf("%T: %s", err, err)
		}
	}
}

// TestBatchPPDAddenda02 validates BatchPPD cannot have Addenda02
func TestBatchPPDAddenda02(t *testing.T) {
	mockBatch := mockBatchPPD()
	mockBatch.Entries[0].AddendaRecordIndicator = 1
	mockBatch.GetEntries()[0].Addenda02 = mockAddenda02()
	mockBatch.Create()
	if err := mockBatch.Validate(); err != nil {
		if e, ok := err.(*BatchError); ok {
			if e.FieldName != "Addenda02" {
				t.Errorf("%T: %s", err, err)
			}
		} else {
			t.Errorf("%T: %s", err, err)
		}
	}
}
