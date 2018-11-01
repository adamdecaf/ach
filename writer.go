// Copyright 2018 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package ach

import (
	"bufio"
	"io"
	"strings"
)

// A Writer writes an ach.file to a NACHA encoded file.
//
// As returned by NewWriter, a Writer writes ach.file structs into
// NACHA formatted files.
//
type Writer struct {
	w       *bufio.Writer
	lineNum int //current line being written
}

// NewWriter returns a new Writer that writes to w.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w: bufio.NewWriter(w),
	}
}

// Writer writes a single ach.file record to w
func (w *Writer) Write(file *File) error {
	if err := file.Validate(); err != nil {
		return err
	}

	w.lineNum = 0
	// Iterate over all records in the file
	if _, err := w.w.WriteString(file.Header.String() + "\n"); err != nil {
		return err
	}
	w.lineNum++

	if err := w.writeBatch(file); err != nil {
		return err
	}

	if err := w.writeIATBatch(file); err != nil {
		return err
	}

	if _, err := w.w.WriteString(file.Control.String() + "\n"); err != nil {
		return err
	}
	w.lineNum++

	// pad the final block
	for i := 0; i < (10-(w.lineNum%10)) && w.lineNum%10 != 0; i++ {
		if _, err := w.w.WriteString(strings.Repeat("9", 94) + "\n"); err != nil {
			return err
		}
	}

	return w.w.Flush()
}

// Flush writes any buffered data to the underlying io.Writer.
func (w *Writer) Flush() error {
	return w.w.Flush()
}

func (w *Writer) writeBatch(file *File) error {
	for _, batch := range file.Batches {
		if _, err := w.w.WriteString(batch.GetHeader().String() + "\n"); err != nil {
			return err
		}
		w.lineNum++
		for _, entry := range batch.GetEntries() {
			if _, err := w.w.WriteString(entry.String() + "\n"); err != nil {
				return err
			}
			w.lineNum++
			for _, addenda := range entry.Addendum {
				if _, err := w.w.WriteString(addenda.String() + "\n"); err != nil {
					return err
				}
				w.lineNum++
			}
		}
		if _, err := w.w.WriteString(batch.GetControl().String() + "\n"); err != nil {
			return err
		}
		w.lineNum++
	}
	return nil
}

func (w *Writer) writeIATBatch(file *File) error {
	for _, iatBatch := range file.IATBatches {
		if _, err := w.w.WriteString(iatBatch.GetHeader().String() + "\n"); err != nil {
			return err
		}
		w.lineNum++
		for _, entry := range iatBatch.GetEntries() {
			if _, err := w.w.WriteString(entry.String() + "\n"); err != nil {
				return err
			}
			w.lineNum++
			if _, err := w.w.WriteString(entry.Addenda10.String() + "\n"); err != nil {
				return err
			}
			w.lineNum++
			if _, err := w.w.WriteString(entry.Addenda11.String() + "\n"); err != nil {
				return err
			}
			w.lineNum++
			if _, err := w.w.WriteString(entry.Addenda12.String() + "\n"); err != nil {
				return err
			}
			w.lineNum++
			if _, err := w.w.WriteString(entry.Addenda13.String() + "\n"); err != nil {
				return err
			}
			w.lineNum++
			if _, err := w.w.WriteString(entry.Addenda14.String() + "\n"); err != nil {
				return err
			}
			w.lineNum++
			if _, err := w.w.WriteString(entry.Addenda15.String() + "\n"); err != nil {
				return err
			}
			w.lineNum++
			if _, err := w.w.WriteString(entry.Addenda16.String() + "\n"); err != nil {
				return err
			}
			w.lineNum++
			// IAT Addenda17 and IAT Addenda18 records
			for _, IATaddenda := range entry.Addendum {
				if _, err := w.w.WriteString(IATaddenda.String() + "\n"); err != nil {
					return err
				}
				w.lineNum++
			}
		}
		if _, err := w.w.WriteString(iatBatch.GetControl().String() + "\n"); err != nil {
			return err
		}
		w.lineNum++
	}
	return nil
}
