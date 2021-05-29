package server_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mflinn92/loggerj/internal/server"
)

func TestPostRecord(t *testing.T) {
	t.Run("post calls append with arbitrary json body and responds 200", func(t *testing.T) {
		log := &logSpy{}
		res := newServerRequest(log, http.MethodPost, "{}")
		assertAppendCalled(t, log)
		assertStatusCode(t, res, http.StatusOK)
	})

	t.Run("get calls read with arbitrary json body and responds 200", func(t *testing.T) {
		log := &logSpy{}
		res := newServerRequest(log, http.MethodGet, "{}")
		assertReadCalled(t, log)
		assertStatusCode(t, res, http.StatusOK)
	})
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

func newServerRequest(log *logSpy, method, body string) *httptest.ResponseRecorder {
	server := server.NewHTTPServer(":8000", log)
	reqBody := strings.NewReader(body)
	req, _ := http.NewRequest(method, "/", reqBody)
	res := httptest.NewRecorder()

	server.Handler.ServeHTTP(res, req)
	return res
}
