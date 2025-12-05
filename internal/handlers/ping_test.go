package handlers

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPingHandler(t *testing.T) {
	// Setup
	app := fiber.New()
	app.Get("/ping", PingHandler)

	t.Run("returns 200 status code", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/ping", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(resp.Body)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("returns JSON content type", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/ping", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(resp.Body)

		contentType := resp.Header.Get("Content-Type")
		assert.Contains(t, contentType, "application/json")
	})

	t.Run("returns valid PingResponse structure", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/ping", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(resp.Body)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var pingResp PingResponse
		err = json.Unmarshal(body, &pingResp)
		require.NoError(t, err)

		assert.Equal(t, "ok", pingResp.Status)
		assert.Equal(t, "pong", pingResp.Message)
		assert.NotEmpty(t, pingResp.Version)
	})

	t.Run("timestamp is recent", func(t *testing.T) {
		before := time.Now().UTC()

		req := httptest.NewRequest("GET", "/ping", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(resp.Body)

		after := time.Now().UTC()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var pingResp PingResponse
		err = json.Unmarshal(body, &pingResp)
		require.NoError(t, err)

		// Timestamp should be between before and after
		assert.True(t, pingResp.Timestamp.After(before.Add(-time.Second)))
		assert.True(t, pingResp.Timestamp.Before(after.Add(time.Second)))
	})

	t.Run("handles multiple concurrent requests", func(t *testing.T) {
		const numRequests = 10
		results := make(chan error, numRequests)

		for i := 0; i < numRequests; i++ {
			go func() {
				req := httptest.NewRequest("GET", "/ping", nil)
				resp, err := app.Test(req)
				if err != nil {
					results <- err
					return
				}

				err = resp.Body.Close()
				if err != nil {
					return
				}

				if resp.StatusCode != fiber.StatusOK {
					results <- assert.AnError
					return
				}
				results <- nil
			}()
		}

		for i := 0; i < numRequests; i++ {
			err := <-results
			assert.NoError(t, err)
		}
	})
}

func TestPingHandlerSimple(t *testing.T) {
	// Setup
	app := fiber.New()
	app.Get("/ping/simple", PingHandlerSimple)

	t.Run("returns 200 status code", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/ping/simple", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(resp.Body)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("returns plain text", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/ping/simple", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(resp.Body)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, "pong", string(body))
	})

	t.Run("content type is text/plain", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/ping/simple", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(resp.Body)

		contentType := resp.Header.Get("Content-Type")
		assert.Contains(t, contentType, "text/plain")
	})
}

func TestPingResponse_Structure(t *testing.T) {
	t.Run("can be marshaled to JSON", func(t *testing.T) {
		response := PingResponse{
			Status:    "ok",
			Message:   "pong",
			Timestamp: time.Now().UTC(),
			Version:   "1.0.0",
		}

		data, err := json.Marshal(response)
		require.NoError(t, err)
		assert.NotEmpty(t, data)

		// Verify it can be unmarshalled
		var decoded PingResponse
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)
		assert.Equal(t, response.Status, decoded.Status)
		assert.Equal(t, response.Message, decoded.Message)
		assert.Equal(t, response.Version, decoded.Version)
	})

	t.Run("timestamp is in ISO8601 format", func(t *testing.T) {
		response := PingResponse{
			Status:    "ok",
			Message:   "pong",
			Timestamp: time.Date(2025, 12, 5, 10, 30, 0, 0, time.UTC),
			Version:   "1.0.0",
		}

		data, err := json.Marshal(response)
		require.NoError(t, err)

		// Check that timestamp is in ISO8601 format
		assert.Contains(t, string(data), "2025-12-05T10:30:00Z")
	})
}

// Benchmark tests
func BenchmarkPingHandler(b *testing.B) {
	app := fiber.New()
	app.Get("/ping", PingHandler)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/ping", nil)
		resp, _ := app.Test(req, -1)
		err := resp.Body.Close()
		if err != nil {
			return
		}
	}
}

func BenchmarkPingHandlerSimple(b *testing.B) {
	app := fiber.New()
	app.Get("/ping/simple", PingHandlerSimple)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/ping/simple", nil)
		resp, _ := app.Test(req, -1)
		err := resp.Body.Close()
		if err != nil {
			return
		}
	}
}

func BenchmarkPingHandler_Parallel(b *testing.B) {
	app := fiber.New()
	app.Get("/ping", PingHandler)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest("GET", "/ping", nil)
			resp, _ := app.Test(req, -1)
			err := resp.Body.Close()
			if err != nil {
				return
			}
		}
	})
}

func BenchmarkPingHandlerSimple_Parallel(b *testing.B) {
	app := fiber.New()
	app.Get("/ping/simple", PingHandlerSimple)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest("GET", "/ping/simple", nil)
			resp, _ := app.Test(req, -1)
			err := resp.Body.Close()
			if err != nil {
				return
			}
		}
	})
}

func BenchmarkPingResponse_Marshal(b *testing.B) {
	response := PingResponse{
		Status:    "ok",
		Message:   "pong",
		Timestamp: time.Now().UTC(),
		Version:   "1.0.0",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(response)
	}
}
