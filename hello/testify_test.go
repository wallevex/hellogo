package hello

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSplitHostPort(t *testing.T) {
	tests := []struct {
		give     string
		wantHost string
		wantPort string
	}{
		{
			give:     "127.0.0.1:8000",
			wantHost: "127.0.0.1",
			wantPort: "8000",
		},
		{
			give:     "192.168.33.186:6500",
			wantHost: "192.168.33.186",
			wantPort: "6500",
		},
		{
			give:     "8.8.8.8:80",
			wantHost: "8.8.8.8",
			wantPort: "80",
		},
	}
	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			host, port, err := net.SplitHostPort(tt.give)
			require.NoError(t, err)
			assert.Equal(t, tt.wantHost, host)
			assert.Equal(t, tt.wantPort, port)
		})
	}
}
