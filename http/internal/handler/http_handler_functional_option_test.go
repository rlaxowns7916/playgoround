package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func strptr(s string) *string { return &s }

func TestFunctionalOptionApply(t *testing.T) {
	testCases := []struct {
		name     string
		given    *string
		expected string
	}{
		{
			name:     "Functional Option패턴을 통해 설정 적용이 가능하다",
			given:    strptr("foo"),
			expected: "foo",
		},
		{
			name:     "Functional Option 패턴을 적용하지 않으면 기본 값이 적용된다",
			given:    nil,
			expected: "ok",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var options []Option
			if tc.given != nil {
				options = append(options, WithHealthCheckResponse(*tc.given))
			}

			handler := NewMonitorHandler(options...)

			mux := http.NewServeMux()
			handler.Register(mux)

			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			body := rr.Body.String()
			if body != tc.expected {
				t.Fatalf("body = %q, want %q", body, tc.expected)
			}
		})
	}
}
