package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_getChunkAndValidate(t *testing.T) {
	//httptest
	tests := []struct {
		name    string
		req     *http.Request
		want    *chunk
		wantErr bool
	}{
		{
			name:    "Example 1",
			req:     httptest.NewRequest("PUT", "http://test.test/test", nil),
			want:    nil,
			wantErr: true,
		},
		{
			name: "Example 2",
			req:  getAllHeadersTest(),
			want: &chunk{
				baseFilename:  "test",
				totalChunks:   2,
				ID:            1,
				chunkTemplate: "chunk-*.txt",
			},
			wantErr: false,
		},
		{
			name:    "Example 3",
			req:     getMissedHeadersTest(),
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getChunkAndValidate(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("getChunkAndValidate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getChunkAndValidate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func getAllHeadersTest() *http.Request {
	req := httptest.NewRequest("PUT", "http://test.test/test", nil)
	req.Header.Add("Base-Filename", "test")
	req.Header.Add("Chunk-Id", "1")
	req.Header.Add("Total-Chunks", "2")
	return req
}

func getMissedHeadersTest() *http.Request {
	req := httptest.NewRequest("PUT", "http://test.test/test", nil)
	req.Header.Add("Base-Filename", "test")
	return req
}
