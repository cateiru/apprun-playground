package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// TestHandleHealth はヘルスチェックハンドラーのユニットテスト
func TestHandleHealth(t *testing.T) {
	// Echo インスタンスを作成
	e := echo.New()

	// テスト用のリクエストを作成
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// ハンドラーを実行
	if assert.NoError(t, handleHealth(c)) {
		// ステータスコードの確認
		assert.Equal(t, http.StatusOK, rec.Code)

		// Content-Type の確認
		assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")

		// レスポンスボディの確認
		var response HealthResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ok", response.Status)
	}
}

// TestHealthEndpoint はヘルスチェックエンドポイントの統合テスト
func TestHealthEndpoint(t *testing.T) {
	// Echo インスタンスを作成し、ルートを設定
	e := echo.New()
	e.GET("/", handleHealth)

	// テスト用のリクエストを作成
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	// リクエストを実行
	e.ServeHTTP(rec, req)

	// ステータスコードの確認
	assert.Equal(t, http.StatusOK, rec.Code)

	// Content-Type の確認
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")

	// レスポンスボディの確認
	var response HealthResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response.Status)

	// レスポンスボディが期待通りの JSON であることを確認
	expectedJSON := `{"status":"ok"}`
	assert.JSONEq(t, expectedJSON, rec.Body.String())
}

// TestHealthEndpointWithMiddleware はミドルウェアを含めた統合テスト
func TestHealthEndpointWithMiddleware(t *testing.T) {
	// Echo インスタンスを作成し、ミドルウェアとルートを設定
	e := echo.New()
	e.Use(echo.MiddlewareFunc(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// ミドルウェアが実行されることを確認するためのヘッダーを追加
			c.Response().Header().Set("X-Test-Middleware", "executed")
			return next(c)
		}
	}))
	e.GET("/", handleHealth)

	// テスト用のリクエストを作成
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	// リクエストを実行
	e.ServeHTTP(rec, req)

	// ステータスコードの確認
	assert.Equal(t, http.StatusOK, rec.Code)

	// ミドルウェアが実行されたことを確認
	assert.Equal(t, "executed", rec.Header().Get("X-Test-Middleware"))

	// レスポンスの確認
	var response HealthResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response.Status)
}
