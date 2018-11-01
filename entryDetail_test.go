// Copyright 2018 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package ach

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// mockEntryDetail creates an entry detail
func mockEntryDetail() *EntryDetail {
	entry := NewEntryDetail()
	entry.TransactionCode = 22
	entry.SetRDFI("121042882")
	entry.DFIAccountNumber = "123456789"
	entry.Amount = 100000000
	entry.IndividualName = "Wade Arnold"
	entry.SetTraceNumber(mockBatchHeader().ODFIIdentification, 1)
	entry.IdentificationNumber = "ABC##jvkdjfuiwn"
	entry.Category = CategoryForward
	return entry
}

// testMockEntryDetail validates an entry detail record
func testMockEntryDetail(t testing.TB) {
	entry := mockEntryDetail()
	if err := entry.Validate(); err != nil {
		t.Error("mockEntryDetail does not validate and will break other tests")
	}
	if entry.TransactionCode != 22 {
		t.Error("TransactionCode dependent default value has changed")
	}
	if entry.DFIAccountNumber != "123456789" {
		t.Error("DFIAccountNumber dependent default value has changed")
	}
	if entry.Amount != 100000000 {
		t.Error("Amount dependent default value has changed")
	}
	if entry.IndividualName != "Wade Arnold" {
		t.Error("IndividualName dependent default value has changed")
	}
	if entry.TraceNumber != 121042880000001 {
		t.Errorf("TraceNumber dependent default value has changed %v", entry.TraceNumber)
	}
}

// TestMockEntryDetail tests validating an entry detail record
func TestMockEntryDetail(t *testing.T) {
	testMockEntryDetail(t)
}

// BenchmarkMockEntryDetail benchmarks validating an entry detail record
func BenchmarkMockEntryDetail(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testMockEntryDetail(b)
	}
}

// testParseEntryDetail parses a known entry detail record string.
func testParseEntryDetail(t testing.TB) {
	var line = "62705320001912345            0000010500c-1            Arnold Wade           DD0076401255655291"
	r := NewReader(strings.NewReader(line))
	r.addCurrentBatch(NewBatchPPD(mockBatchPPDHeader()))
	r.currentBatch.SetHeader(mockBatchHeader())
	r.line = line
	if err := r.parseEntryDetail(); err != nil {
		t.Errorf("%T: %s", err, err)
	}
	record := r.currentBatch.GetEntries()[0]

	if record.recordType != "6" {
		t.Errorf("RecordType Expected '6' got: %v", record.recordType)
	}
	if record.TransactionCode != 27 {
		t.Errorf("TransactionCode Expected '27' got: %v", record.TransactionCode)
	}
	if record.RDFIIdentificationField() != "05320001" {
		t.Errorf("RDFIIdentification Expected '05320001' got: '%v'", record.RDFIIdentificationField())
	}
	if record.CheckDigit != "9" {
		t.Errorf("CheckDigit Expected '9' got: %v", record.CheckDigit)
	}
	if record.DFIAccountNumberField() != "12345            " {
		t.Errorf("DfiAccountNumber Expected '12345            ' got: %v", record.DFIAccountNumberField())
	}
	if record.AmountField() != "0000010500" {
		t.Errorf("Amount Expected '0000010500' got: %v", record.AmountField())
	}

	if record.IdentificationNumber != "c-1            " {
		t.Errorf("IdentificationNumber Expected 'c-1            ' got: %v", record.IdentificationNumber)
	}
	if record.IndividualName != "Arnold Wade           " {
		t.Errorf("IndividualName Expected 'Arnold Wade           ' got: %v", record.IndividualName)
	}
	if record.DiscretionaryData != "DD" {
		t.Errorf("DiscretionaryData Expected 'DD' got: %v", record.DiscretionaryData)
	}
	if record.AddendaRecordIndicator != 0 {
		t.Errorf("AddendaRecordIndicator Expected '0' got: %v", record.AddendaRecordIndicator)
	}
	if record.TraceNumberField() != "076401255655291" {
		t.Errorf("TraceNumber Expected '076401255655291' got: %v", record.TraceNumberField())
	}
}

// TestParseEntryDetail tests parsing a known entry detail record string.
func TestParseEntryDetail(t *testing.T) {
	testParseEntryDetail(t)
}

// BenchmarkParseEntryDetail benchmarks parsing a known entry detail record string.
func BenchmarkParseEntryDetail(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testParseEntryDetail(b)
	}
}

// testEDString validates that a known parsed entry
// detail can be returned to a string of the same value
func testEDString(t testing.TB) {
	var line = "62705320001912345            0000010500c-1            Arnold Wade           DD0076401255655291"
	r := NewReader(strings.NewReader(line))
	r.addCurrentBatch(NewBatchPPD(mockBatchPPDHeader()))
	r.currentBatch.SetHeader(mockBatchHeader())
	r.line = line
	if err := r.parseEntryDetail(); err != nil {
		t.Errorf("%T: %s", err, err)
	}
	record := r.currentBatch.GetEntries()[0]

	if record.String() != line {
		t.Errorf("Strings do not match")
	}
}

// TestEDString tests validating that a known parsed entry
// detail can be returned to a string of the same value
func TestEDString(t *testing.T) {
	testEDString(t)
}

// BenchmarkEDString benchmarks validating that a known parsed entry
// detail can be returned to a string of the same value
func BenchmarkEDString(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testEDString(b)
	}
}

// testValidateEDRecordType validates error if recordType is not 6
func testValidateEDRecordType(t testing.TB) {
	ed := mockEntryDetail()
	ed.recordType = "2"
	if err := ed.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if e.FieldName != "recordType" {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestValidateEDRecordType tests validating error if recordType is not 6
func TestValidateEDRecordType(t *testing.T) {
	testValidateEDRecordType(t)
}

// BenchmarkValidateEDRecordType benchmarks validating error if recordType is not 6
func BenchmarkValidateEDRecordType(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testValidateEDRecordType(b)
	}
}

// testValidateEDTransactionCode validates error if transaction code is not valid
func testValidateEDTransactionCode(t testing.TB) {
	ed := mockEntryDetail()
	ed.TransactionCode = 63
	if err := ed.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if e.FieldName != "TransactionCode" {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestValidateEDTransactionCode tests validating error if transaction code is not valid
func TestValidateEDTransactionCode(t *testing.T) {
	testValidateEDTransactionCode(t)
}

// BenchmarkValidateEDTransactionCode benchmarks validating error if transaction code is not valid
func BenchmarkValidateEDTransactionCode(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testValidateEDTransactionCode(b)
	}
}

// testEDFieldInclusion validates entry detail field inclusion
func testEDFieldInclusion(t testing.TB) {
	ed := mockEntryDetail()
	ed.Amount = 0
	if err := ed.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if !strings.Contains(e.Msg, msgFieldInclusion) {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestEDFieldInclusion tests validating entry detail field inclusion
func TestEDFieldInclusion(t *testing.T) {
	testEDFieldInclusion(t)
}

// BenchmarkEDFieldInclusion benchmarks validating entry detail field inclusion
func BenchmarkEDFieldInclusion(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testEDFieldInclusion(b)
	}
}

// testEDdfiAccountNumberAlphaNumeric validates DFI account number is alpha numeric
func testEDdfiAccountNumberAlphaNumeric(t testing.TB) {
	ed := mockEntryDetail()
	ed.DFIAccountNumber = "®"
	if err := ed.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if e.FieldName != "DFIAccountNumber" {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestEDdfiAccountNumberAlphaNumeric tests validating DFI account number is alpha numeric
func TestEDdfiAccountNumberAlphaNumeric(t *testing.T) {
	testEDdfiAccountNumberAlphaNumeric(t)
}

// BenchmarkEDdfiAccountNumberAlphaNumeric benchmarks validating DFI account number is alpha numeric
func BenchmarkEDdfiAccountNumberAlphaNumeric(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testEDdfiAccountNumberAlphaNumeric(b)
	}
}

// testEDIdentificationNumberAlphaNumeric validates identification number is alpha numeric
func testEDIdentificationNumberAlphaNumeric(t testing.TB) {
	ed := mockEntryDetail()
	ed.IdentificationNumber = "®"
	if err := ed.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if e.FieldName != "IdentificationNumber" {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestEDIdentificationNumberAlphaNumeric tests validating identification number is alpha numeric
func TestEDIdentificationNumberAlphaNumeric(t *testing.T) {
	testEDIdentificationNumberAlphaNumeric(t)
}

// BenchmarkEDIdentificationNumberAlphaNumeric benchmarks validating identification number is alpha numeric
func BenchmarkEDIdentificationNumberAlphaNumeric(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testEDIdentificationNumberAlphaNumeric(b)
	}
}

// testEDIndividualNameAlphaNumeric validates individual name is alpha numeric
func testEDIndividualNameAlphaNumeric(t testing.TB) {
	ed := mockEntryDetail()
	ed.IndividualName = "W®DE"
	if err := ed.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if e.FieldName != "IndividualName" {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestEDIndividualNameAlphaNumeric tests validating individual name is alpha numeric
func TestEDIndividualNameAlphaNumeric(t *testing.T) {
	testEDIndividualNameAlphaNumeric(t)
}

// BenchmarkEDIndividualNameAlphaNumeric benchmarks validating individual name is alpha numeric
func BenchmarkEDIndividualNameAlphaNumeric(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testEDIndividualNameAlphaNumeric(b)
	}
}

// testEDDiscretionaryDataAlphaNumeric validates discretionary data is alpha numeric
func testEDDiscretionaryDataAlphaNumeric(t testing.TB) {
	ed := mockEntryDetail()
	ed.DiscretionaryData = "®!"
	if err := ed.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if e.FieldName != "DiscretionaryData" {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestEDDiscretionaryDataAlphaNumeric tests validating discretionary data is alpha numeric
func TestEDDiscretionaryDataAlphaNumeric(t *testing.T) {
	testEDDiscretionaryDataAlphaNumeric(t)
}

// BenchmarkEDDiscretionaryDataAlphaNumeric benchmarks validating discretionary data is alpha numeric
func BenchmarkEDDiscretionaryDataAlphaNumeric(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testEDDiscretionaryDataAlphaNumeric(b)
	}
}

// testEDisCheckDigit validates check digit
func testEDisCheckDigit(t testing.TB) {
	ed := mockEntryDetail()
	ed.CheckDigit = "1"
	if err := ed.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if e.FieldName != "RDFIIdentification" {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestEDisCheckDigit tests validating check digit
func TestEDisCheckDigit(t *testing.T) {
	testEDisCheckDigit(t)
}

// BenchmarkEDSetRDFI benchmarks validating check digit
func BenchmarkEDisCheckDigit(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testEDisCheckDigit(b)
	}
}

// testEDSetRDFI validates setting RDFI
func testEDSetRDFI(t testing.TB) {
	ed := NewEntryDetail()
	ed.SetRDFI("810866774")
	if ed.RDFIIdentification != "81086677" {
		t.Error("RDFI identification")
	}
	if ed.CheckDigit != "4" {
		t.Error("Unexpected check digit")
	}
}

// TestEDSetRDFI tests validating setting RDFI
func TestEDSetRDFI(t *testing.T) {
	testEDSetRDFI(t)
}

// BenchmarkEDSetRDFI benchmarks validating setting RDFI
func BenchmarkEDSetRDFI(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testEDSetRDFI(b)
	}
}

// testEDFieldInclusionRecordType validates record type field inclusion
func testEDFieldInclusionRecordType(t testing.TB) {
	entry := mockEntryDetail()
	entry.recordType = ""
	if err := entry.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if !strings.Contains(e.Msg, msgFieldInclusion) {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestEDFieldInclusionRecordType tests validating record type field inclusion
func TestEDFieldInclusionRecordType(t *testing.T) {
	testEDFieldInclusionRecordType(t)
}

// BenchmarkEDFieldInclusionRecordType benchmarks validating record type field inclusion
func BenchmarkEDFieldInclusionRecordType(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testEDFieldInclusionRecordType(b)
	}
}

// testEDFieldInclusionTransactionCode validates transaction code field inclusion
func testEDFieldInclusionTransactionCode(t testing.TB) {
	entry := mockEntryDetail()
	entry.TransactionCode = 0
	if err := entry.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if !strings.Contains(e.Msg, msgFieldInclusion) {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestEDFieldInclusionTransactionCode tests validating transaction code field inclusion
func TestEDFieldInclusionTransactionCode(t *testing.T) {
	testEDFieldInclusionTransactionCode(t)
}

// BenchmarkEDFieldInclusionTransactionCode benchmarks validating transaction code field inclusion
func BenchmarkEDFieldInclusionTransactionCode(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testEDFieldInclusionTransactionCode(b)
	}
}

// testEDFieldInclusionRDFIIdentification validates RDFI identification field inclusion
func testEDFieldInclusionRDFIIdentification(t testing.TB) {
	entry := mockEntryDetail()
	entry.RDFIIdentification = ""
	if err := entry.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if !strings.Contains(e.Msg, msgFieldInclusion) {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestEDFieldInclusionRDFIIdentification tests validating RDFI identification field inclusion
func TestEDFieldInclusionRDFIIdentification(t *testing.T) {
	testEDFieldInclusionRDFIIdentification(t)
}

// BenchmarkEDFieldInclusionRDFIIdentification benchmarks validating RDFI identification field inclusion
func BenchmarkEDFieldInclusionRDFIIdentification(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testEDFieldInclusionRDFIIdentification(b)
	}
}

// testEDFieldInclusionDFIAccountNumber validates DFI account number field inclusion
func testEDFieldInclusionDFIAccountNumber(t testing.TB) {
	entry := mockEntryDetail()
	entry.DFIAccountNumber = ""
	if err := entry.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if !strings.Contains(e.Msg, msgFieldInclusion) {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestEDFieldInclusionDFIAccountNumber tests validating DFI account number field inclusion
func TestEDFieldInclusionDFIAccountNumber(t *testing.T) {
	testEDFieldInclusionDFIAccountNumber(t)
}

// BenchmarkEDFieldInclusionDFIAccountNumber benchmarks validating DFI account number field inclusion
func BenchmarkEDFieldInclusionDFIAccountNumber(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testEDFieldInclusionDFIAccountNumber(b)
	}
}

// testEDFieldInclusionIndividualName validates individual name field inclusion
func testEDFieldInclusionIndividualName(t testing.TB) {
	entry := mockEntryDetail()
	entry.IndividualName = ""
	if err := entry.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if !strings.Contains(e.Msg, msgFieldInclusion) {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestEDFieldInclusionIndividualName tests validating individual name field inclusion
func TestEDFieldInclusionIndividualName(t *testing.T) {
	testEDFieldInclusionIndividualName(t)
}

// BenchmarkEDFieldInclusionIndividualName benchmarks validating individual name field inclusion
func BenchmarkEDFieldInclusionIndividualName(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testEDFieldInclusionIndividualName(b)
	}
}

// testEDFieldInclusionTraceNumber validates trace number field inclusion
func testEDFieldInclusionTraceNumber(t testing.TB) {
	entry := mockEntryDetail()
	entry.TraceNumber = 0
	if err := entry.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if !strings.Contains(e.Msg, msgFieldInclusion) {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestEDFieldInclusionTraceNumber tests validating trace number field inclusion
func TestEDFieldInclusionTraceNumber(t *testing.T) {
	testEDFieldInclusionTraceNumber(t)
}

// BenchmarkEDFieldInclusionTraceNumber benchmarks validating trace number field inclusion
func BenchmarkEDFieldInclusionTraceNumber(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testEDFieldInclusionTraceNumber(b)
	}
}

// testEDAddAddenda99 validates adding Addenda99 to an entry detail
func testEDAddAddenda99(t testing.TB) {
	entry := mockEntryDetail()
	entry.AddAddenda(mockAddenda99())
	if entry.Category != CategoryReturn {
		t.Error("Addenda99 added and isReturn is false")
	}
	if entry.AddendaRecordIndicator != 1 {
		t.Error("Addenda99 added and record indicator is not 1")
	}

}

// TestEDAddAddenda99 tests validating adding Addenda99 to an entry detail
func TestEDAddAddenda99(t *testing.T) {
	testEDAddAddenda99(t)
}

// BenchmarkEDAddAddenda99 benchmarks validating adding Addenda99 to an entry detail
func BenchmarkEDAddAddenda99(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testEDAddAddenda99(b)
	}
}

// testEDAddAddenda99Twice validates only one Addenda99 is added to an entry detail
func testEDAddAddenda99Twice(t testing.TB) {
	entry := mockEntryDetail()
	entry.AddAddenda(mockAddenda99())
	entry.AddAddenda(mockAddenda99())
	if entry.Category != CategoryReturn {
		t.Error("Addenda99 added and Category is not CategoryReturn")
	}

	if len(entry.Addendum) != 1 {
		t.Error("Addenda99 added and isReturn is false")
	}
}

// TestEDAddAddenda99Twice tests validating only one Addenda99 is added to an entry detail
func TestEDAddAddenda99Twice(t *testing.T) {
	testEDAddAddenda99Twice(t)
}

// BenchmarkEDAddAddenda99Twice benchmarks validating only one Addenda99 is added to an entry detail
func BenchmarkEDAddAddenda99Twice(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testEDAddAddenda99Twice(b)
	}
}

// testEDCreditOrDebit validates debit and credit transaction code
func testEDCreditOrDebit(t testing.TB) {
	entry := mockEntryDetail()
	if entry.CreditOrDebit() != "C" { // our mock's default
		t.Errorf("TransactionCode %v expected a Credit(C) got %v", entry.TransactionCode, entry.CreditOrDebit())
	}

	// TransactionCode -> C or D
	var cases = map[int]string{
		// invalid
		-1:  "",
		00:  "", // invalid
		1:   "",
		108: "",
		// valid
		22: "C",
		23: "C",
		27: "D",
		28: "D",
		32: "C",
		33: "C",
		37: "D",
		38: "D",
	}
	for code, expected := range cases {
		entry.TransactionCode = code
		if v := entry.CreditOrDebit(); v != expected {
			t.Errorf("TransactionCode %d expected %s, got %s", code, expected, v)
		}
	}
}

// TestEDCreditOrDebit tests validating debit and credit transaction code
func TestEDCreditOrDebit(t *testing.T) {
	testEDCreditOrDebit(t)
}

// BenchmarkEDCreditOrDebit benchmarks validating debit and credit transaction code
func BenchmarkEDCreditOrDebit(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testEDCreditOrDebit(b)
	}
}

// testValidateEDCheckDigit validates CheckDigit error
func testValidateEDCheckDigit(t testing.TB) {
	ed := mockEntryDetail()
	ed.CheckDigit = "XYZ"
	if err := ed.Validate(); err != nil {
		if e, ok := err.(*FieldError); ok {
			if e.FieldName != "CheckDigit" {
				t.Errorf("%T: %s", err, err)
			}
		}
	}
}

// TestValidateEDCheckDigit tests validating validates CheckDigit error
func TestValidateEDCheckDigit(t *testing.T) {
	testValidateEDCheckDigit(t)
}

// BenchmarkValidateEDCheckDigit benchmarks validating CheckDigit error
func BenchmarkValidateEDCheckDigit(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testValidateEDCheckDigit(b)
	}
}

func TestEntryDetail__json(t *testing.T) {
	f, err := os.Open(filepath.Join("test", "testdata", "entrydetail-valid.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var ed EntryDetail
	if err := json.NewDecoder(f).Decode(&ed); err != nil {
		t.Fatal(err)
	}

	if ed.ID != "test" {
		t.Error(ed.ID)
	}
	if ed.TransactionCode != 27 {
		t.Error(ed)
	}
	if ed.RDFIIdentification != "RDFI" {
		t.Error(ed.RDFIIdentification)
	}
	if ed.CheckDigit != "0" {
		t.Error(ed.CheckDigit)
	}
	if ed.DFIAccountNumber != "132" {
		t.Error(ed.DFIAccountNumber)
	}
	if ed.Amount != 10000 {
		t.Error(ed.Amount)
	}
	if ed.IndividualName != "jane doe" {
		t.Error(ed.IndividualName)
	}
	if len(ed.Addendum) != 2 {
		t.Errorf("got %d addenda records", len(ed.Addendum))
	}

	// addenda record order matters
	rec05, ok := ed.Addendum[0].(*Addenda05)
	if !ok || rec05 == nil {
		t.Fatalf("%#v", rec05)
	}
	if rec05.ID != "test-05" {
		t.Error(rec05.ID)
	}
	if rec05.TypeCode != "05" {
		t.Error(rec05.TypeCode)
	}
	if rec05.PaymentRelatedInformation != "lottery winnings" {
		t.Error(rec05.PaymentRelatedInformation)
	}
	if rec05.SequenceNumber != 1 {
		t.Errorf("got %d", rec05.SequenceNumber)
	}

	// second addenda record
	rec98, ok := ed.Addendum[1].(*Addenda98)
	if !ok || rec98 == nil {
		t.Fatalf("%#v", rec98)
	}
	if rec98.TypeCode != "98" {
		t.Error(rec98.TypeCode)
	}
	if rec98.ChangeCode != "C01" {
		t.Error(rec98.ChangeCode)
	}
	if rec98.OriginalTrace != 18571 {
		t.Error(rec98.OriginalTrace)
	}
	if rec98.OriginalDFI != "123456789" {
		t.Error(rec98.OriginalDFI)
	}
	if rec98.CorrectedData != "fix me" {
		t.Error(rec98.CorrectedData)
	}
	if rec98.TraceNumber != 87198212121 {
		t.Errorf("traceNumber: %d", rec98.TraceNumber)
	}
}
