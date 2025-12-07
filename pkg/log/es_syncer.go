package log

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap/zapcore"
)

// elasticsearchSyncer implements zapcore.WriteSyncer and ships logs to Elasticsearch using the Bulk API.
type elasticsearchSyncer struct {
	bulkURL   string
	indexName string
	client    *http.Client

	mu     sync.Mutex
	buf    bytes.Buffer // accumulates NDJSON bulk payload
	count  int          // number of actions in current buffer
	closed bool

	flushEvery time.Duration
	maxActions int
}

func newElasticsearchSyncer(baseURL, index string) zapcore.WriteSyncer {
	es := &elasticsearchSyncer{
		bulkURL:    strings.TrimRight(baseURL, "/") + "/_bulk",
		indexName:  index,
		client:     &http.Client{Timeout: 5 * time.Second},
		flushEvery: 1 * time.Second,
		maxActions: 200, // pairs of lines (index + doc) counts as 1 action
	}
	// start periodic flusher
	go es.flushLoop()
	return es
}

func (e *elasticsearchSyncer) Write(p []byte) (int, error) {
	// Each zap entry is a single JSON object. Bulk requires action line + source line.
	// We add: {"index":{"_index":"<index>"}}\n<json>\n
	if len(p) == 0 {
		return 0, nil
	}

	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed {
		return 0, io.ErrClosedPipe
	}

	// ensure newline-trimmed JSON then re-add newline for NDJSON
	jsonLine := strings.TrimSpace(string(p))
	meta := fmt.Sprintf("{\"index\":{\"_index\":\"%s\"}}\n", e.indexName)
	e.buf.WriteString(meta)
	e.buf.WriteString(jsonLine)
	e.buf.WriteByte('\n')
	e.count++

	if e.count >= e.maxActions {
		// flush without holding the lock for network call
		buf := e.swapBufferLocked()
		go e.post(buf)
	}
	return len(p), nil
}

func (e *elasticsearchSyncer) Sync() error {
	e.mu.Lock()
	buf := e.swapBufferLocked()
	e.mu.Unlock()
	if buf == nil {
		return nil
	}
	return e.post(buf)
}

func (e *elasticsearchSyncer) flushLoop() {
	ticker := time.NewTicker(e.flushEvery)
	defer ticker.Stop()
	for range ticker.C {
		e.mu.Lock()
		if e.closed {
			e.mu.Unlock()
			return
		}
		buf := e.swapBufferLocked()
		e.mu.Unlock()
		if buf != nil {
			_ = e.post(buf) // best-effort; errors are swallowed to avoid crashing the app
		}
	}
}

func (e *elasticsearchSyncer) swapBufferLocked() *bytes.Buffer {
	if e.count == 0 || e.buf.Len() == 0 {
		return nil
	}
	old := e.buf
	e.buf = bytes.Buffer{}
	e.count = 0
	return &old
}

func (e *elasticsearchSyncer) post(body *bytes.Buffer) error {
	if body == nil || body.Len() == 0 {
		return nil
	}
	req, err := http.NewRequest(http.MethodPost, e.bulkURL, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-ndjson")
	resp, err := e.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Consider non-2xx a failure but do not propagate to crash; return error
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		io.Copy(io.Discard, resp.Body)
		return fmt.Errorf("elasticsearch bulk post failed: status %d", resp.StatusCode)
	}
	io.Copy(io.Discard, resp.Body)
	return nil
}
