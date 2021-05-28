package server_test

import (
	"testing"

	"github.com/mflinn92/loggerj/internal/server"
)

func TestPostRecord(t *testing.T) {

}

type logSpy struct {
	appendCalled bool
	readCalled   bool
}

func (s *logSpy) Read(offset uint64) (server.Record, error) {
	s.readCalled = true
	return server.Record{}, nil
}

func (s *logSpy) Append(record server.Record) (uint64, error) {
	s.appendCalled = true
}
