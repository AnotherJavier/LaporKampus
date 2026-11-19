//go:build integration

// Paket ini berisi uji integrasi yang memakai aplikasi Fiber sungguhan
// (tanpa server jaringan) sehingga kontrak HTTP dapat diuji cepat.
//
// Jalankan dengan:
//
//	go test -tags=integration ./tests/integration/...
package integration

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"

	"github.com/vyerru/lapor-kampus/app/model"
	"github.com/vyerru/lapor-kampus/config"
	"github.com/vyerru/lapor-kampus/route"
)

func newTestApp() *fiber.App {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return config.NewApp(logger, route.Dependencies{
		AppName: "LaporKampus",
		Env:     "test",
		Version: config.Version,
	})
}

func TestHealth(t *testing.T) {
	app := newTestApp()

	resp, err := app.Test(newRequest(http.MethodGet, "/api/v1/health", ""), -1)
	if err != nil {
		t.Fatalf("request gagal: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, ingin %d", resp.StatusCode, http.StatusOK)
	}

	var body model.WebResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("body bukan JSON: %v", err)
	}
	if !body.Success {
		t.Fatal("success harus true")
	}
	if body.Message != "aplikasi sehat" {
		t.Fatalf("message = %q, ingin %q", body.Message, "aplikasi sehat")
	}
}

func TestUnknownRouteReturnsEnvelope(t *testing.T) {
	app := newTestApp()

	resp, err := app.Test(newRequest(http.MethodGet, "/api/v1/tidak-ada", ""), -1)
	if err != nil {
		t.Fatalf("request gagal: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d, ingin %d", resp.StatusCode, http.StatusNotFound)
	}

	var body model.ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("body bukan JSON: %v", err)
	}
	if body.Success {
		t.Fatal("success harus false")
	}
	if body.Code != "NOT_FOUND" {
		t.Fatalf("code = %q, ingin %q", body.Code, "NOT_FOUND")
	}
	if body.RequestID == "" {
		t.Fatal("request_id harus terisi")
	}
}

// newRequest membungkus request tanpa menulis URL lengkap di setiap uji.
// httptest.NewRequest mengisi Host, yang wajib sejak fasthttp v1.70.
func newRequest(method, path, body string) *http.Request {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	return req
}
