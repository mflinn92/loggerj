package server_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mflinn92/loggerj/internal/server"
)

func TestPostRecord(t *testing.T) {
	t.Run("it posts a valid record", func(t *testing.T) {
		log := &logSpy{}
		// value is hello b64 encoded
		res := newServerRequest(log, http.MethodPost, `{"record":{"value":"aGVsbG8=","offset":0}}`)

		assertAppendCalled(t, log)
		assertStatusCode(t, res, http.StatusOK)
	})

	t.Run("post returns 400 with invalid json", func(t *testing.T) {
		log := &logSpy{}
		res := newServerRequest(log, http.MethodPost, `{"invalid_json": "blah"`)

		assertAppendNotCalled(t, log)
		assertStatusCode(t, res, http.StatusBadRequest)

	})
}

func TestGetRecord(t *testing.T) {
	t.Run("it retrieves a valid record", func(t *testing.T) {
		log := newLogWithRecord(server.Record{Value: []byte("hello"), Offset: 0})
		res := newServerRequest(log, http.MethodGet, `{"offset":0}`)

		assertStatusCode(t, res, http.StatusOK)
		// record value is b64 encoded
		wantRes := `{"record":{"value":"aGVsbG8=","offset":0}}`
		assertResponseBody(t, res, wantRes)
	})
	t.Run("get returns 400 with invalid json", func(t *testing.T) {
		log := &logSpy{}
		res := newServerRequest(log, http.MethodGet, `{"invalid_json": "bah`)

		assertStatusCode(t, res, http.StatusBadRequest)
		assertReadNotCalled(t, log)
	})

	t.Run("it returns 404 not found for invalid offset", func(t *testing.T) {
		log := &logSpy{}
		res := newServerRequest(log, http.MethodGet, `{"offset": 3}`)

		assertStatusCode(t, res, http.StatusNotFound)
		assertReadCalled(t, log)
	})
}

func TestPostGetRecordIntegration(t *testing.T) {
	log := &logSpy{}
	var postResponse *httptest.ResponseRecorder

	t.Run("post a record", func(t *testing.T) {
		// value should be a valid b64 encoding
		postResponse = newServerRequest(log, http.MethodPost, `{"record":{"value":"AS232fwf"}}`)

		assertStatusCode(t, postResponse, http.StatusOK)
		assertResponseBody(t, postResponse, `{"offset":0}`)
	})

	t.Run("get previously posted record", func(t *testing.T) {
		res := newServerRequest(log, http.MethodGet, postResponse.Body.String())
		wantRes := `{"record":{"value":"AS232fwf","offset":0}}`

		assertStatusCode(t, res, http.StatusOK)
		assertResponseBody(t, res, wantRes)
	})

	t.Run("get returns 400 with invalid json", func(t *testing.T) {
		log := &logSpy{}
		res := newServerRequest(log, http.MethodGet, `{"invalid_json": "bah`)

		assertStatusCode(t, res, http.StatusBadRequest)
		assertReadNotCalled(t, log)
	})
}

func TestPostGetRecordIntegration(t *testing.T) {
	log := &logSpy{}
	var postResponse *httptest.ResponseRecorder

	t.Run("post a record", func(t *testing.T) {
		// value should be a valid b64 encoding
		postResponse = newServerRequest(log, http.MethodPost, `{"record":{"value":"AS232fwf"}}`)

		assertStatusCode(t, postResponse, http.StatusOK)
		assertResponseBody(t, postResponse, `{"offset":0}`)
	})

	t.Run("get previously posted record", func(t *testing.T) {
		res := newServerRequest(log, http.MethodGet, postResponse.Body.String())
		wantRes := `{"record":{"value":"AS232fwf","offset":0}}`

		assertStatusCode(t, res, http.StatusOK)
		assertResponseBody(t, res, wantRes)
	})
}

type logSpy struct {
	appendCalled bool
	readCalled   bool
	records      []server.Record
}

func (s *logSpy) Read(offset uint64) (server.Record, error) {
	s.readCalled = true
	if offset >= uint64(len(s.records)) {
		return server.Record{}, server.ErrOffsetNotFound
	}
	return s.records[offset], nil
}

func (s *logSpy) Append(record server.Record) (uint64, error) {
	s.appendCalled = true
	s.records = append(s.records, record)
	return 0, nil
}

func assertAppendCalled(t testing.TB, log *logSpy) {
	t.Helper()
	if !log.appendCalled {
		t.Errorf("expected an append to log ")
	}
}

func assertAppendNotCalled(t testing.TB, log *logSpy) {
	t.Helper()
	if log.appendCalled {
		t.Errorf("append should not be called with invalid request body")
	}
}

func assertReadCalled(t testing.TB, log *logSpy) {
	t.Helper()
	if !log.readCalled {
		t.Errorf("expected a read from the log")
	}
}

func assertReadNotCalled(t testing.TB, log *logSpy) {
	t.Helper()
	if log.readCalled {
		t.Errorf("read should not have been called with invalid request body")
	}
}

func assertStatusCode(t testing.TB, res *httptest.ResponseRecorder, want int) {
	t.Helper()
	if res.Code != want {
		t.Errorf("got status code %d wanted %d", res.Code, want)
	}
}

func assertResponseBody(t testing.TB, res *httptest.ResponseRecorder, want string) {
	t.Helper()
	body := strings.TrimSuffix(res.Body.String(), "\n")
	if body != want {
		t.Errorf("got response body %q, want %q", body, want)
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

func newLogWithRecord(r server.Record) *logSpy {
	return &logSpy{
		records: []server.Record{r},
	}
}
