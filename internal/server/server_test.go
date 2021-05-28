package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPostRecord(t *testing.T) {
	t.Run("it calls append with arbitrary json body and responds 200", func(t *testing.T) {
		log := &logSpy{}
		server := NewLogServer(log)
		body := strings.NewReader("{}")
		req, _ := http.NewRequest(http.MethodPost, "/", body)
		res := httptest.NewRecorder()

		server.handleProduce(res, req)
		assertAppendCalled(t, log)
		assertStatusCode(t, res, http.StatusOK)
	})

	t.Run("it calls read with arbitrary json body and responds 200", func(t *testing.T) {
		log := &logSpy{}
		server := NewLogServer(log)
		body := strings.NewReader("{}")
		req, _ := http.NewRequest(http.MethodGet, "/", body)
		res := httptest.NewRecorder()

		server.handleConsume(res, req)
		assertReadCalled(t, log)
		assertStatusCode(t, res, http.StatusOK)
	})
}

type logSpy struct {
	appendCalled bool
	readCalled   bool
}

func (s *logSpy) Read(offset uint64) (Record, error) {
	s.readCalled = true
	return Record{}, nil
}

func (s *logSpy) Append(record Record) (uint64, error) {
	s.appendCalled = true
	return 0, nil
}

func assertAppendCalled(t testing.TB, log *logSpy) {
	t.Helper()
	if !log.appendCalled {
		t.Errorf("expected an append to log ")
	}
}

func assertReadCalled(t testing.TB, log *logSpy) {
	t.Helper()
	if !log.readCalled {
		t.Errorf("expected a read from the log")
	}
}

func assertStatusCode(t testing.TB, res *httptest.ResponseRecorder, want int) {
	t.Helper()
	if res.Code != want {
		t.Errorf("got status code %d wanted %d", res.Code, want)
	}
}
