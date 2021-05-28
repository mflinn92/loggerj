package server_test

import (
	"reflect"
	"testing"

	"github.com/mflinn92/loggerj/internal/server"
)

func TestLog(t *testing.T) {
	cases := []server.Record{
		{Value: []byte("first"), Offset: 0},
		{Value: []byte("second"), Offset: 1},
		{Value: []byte("third"), Offset: 2},
	}
	t.Run("it appends new records and reads them", func(t *testing.T) {
		log := server.NewLog()

		assertRecordsAppended(t, cases, log)

		// Attempt to read back each test case and assure they were added to the log
		for offset, record := range cases {
			assertRecordAtOffset(t, log, uint64(offset), record)
		}
	})

	t.Run("it errors with out of bounds offset", func(t *testing.T) {
		log := server.NewLog()

		assertRecordsAppended(t, cases, log)
		assertOffsetNotFoundErr(t, log, 4)
	})
}

// takes a slice of records and a log and assures Append is called successfully
// for each record
func assertRecordsAppended(t testing.TB, records []server.Record, log *server.Log) {
	t.Helper()
	for _, record := range records {
		_, err := log.Append(record)
		if err != nil {
			t.Errorf("unexpected error appending record %v %v", record, err)
		}
	}
}

// takes a log offset and desired record and confirms the log contains the desired record
// at specified offset
func assertRecordAtOffset(t testing.TB, log *server.Log, offset uint64, want server.Record) {
	t.Helper()
	got, err := log.Read(offset)
	if err != nil {
		t.Fatalf("unexpected error reading record at offset %d %v", offset, err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got record %v at offset %d, wanted %v", got, offset, want)
	}
}

// asserts that an offset not found err is returned by log.Read and that
// an empty record is returned
func assertOffsetNotFoundErr(t testing.TB, log *server.Log, offset uint64) {
	record, err := log.Read(offset)
	if err != server.ErrOffsetNotFound {
		t.Errorf("expected offset not found error, got %v", err)
	}

	if !reflect.DeepEqual(record, server.Record{}) {
		t.Errorf("expected an empty record, got %v", record)
	}
}
