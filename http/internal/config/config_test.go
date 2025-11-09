package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironmentValueReplaceMiddleware(t *testing.T) {
	testCases := []struct {
		name  string
		given string
	}{
		{"PORT에 환경변수 설정 시, Config에서 정상적으로 읽을 수 있다", "2212"},
	}
	_ = os.Setenv("APP_PHASE", "test")
	defer func() {
		_ = os.Unsetenv("APP_PHASE")
		_ = os.Unsetenv("HTTP_PORT")
	}()

	for _, tc := range testCases {
		_ = os.Setenv("HTTP_PORT", tc.given)
		t.Run(tc.name, func(t *testing.T) {
			config := Load()
			assert.Equal(t, 2212, config.HTTP.Port)
		})
	}
}
