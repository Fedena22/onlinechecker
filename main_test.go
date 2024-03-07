package main

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_checkConnection(t *testing.T) {
	tests := []struct {
		name     string
		httpcode int
		want     bool
		wantErr  bool
	}{
		{
			name:     "host not reachable",
			httpcode: http.StatusBadGateway,
			want:     false,
			wantErr:  false,
		},
		{
			name:     "host rechable",
			httpcode: http.StatusOK,
			want:     true,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		logger := slog.Default()

		t.Run(tt.name, func(t *testing.T) {

			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.httpcode)
			}))
			defer svr.Close()
			check := httpClient{Client: svr.Client(), baseUrl: svr.URL}
			got, err := check.checkConnection(logger)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkConnection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkConnection() = %v, want %v", got, tt.want)
			}
		})
	}
}
