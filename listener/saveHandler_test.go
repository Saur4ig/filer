package main

import (
	"os"
	"testing"
)

func Test_getDir(t *testing.T) {
	tests := []struct {
		name    string
		chunk   *chunk
		want    string
		wantErr bool
	}{
		{
			name: "Example 1",
			chunk: &chunk{
				baseFilename:  "test",
				totalChunks:   1,
				ID:            1,
				chunkTemplate: "chunk-*.txt",
			},
			want:    "files",
			wantErr: false,
		},
		{
			name: "Example 2",
			chunk: &chunk{
				baseFilename:  "test.txt",
				totalChunks:   2,
				ID:            1,
				chunkTemplate: "chunk-*.txt",
			},
			want:    "test.txt_dir",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getDir(tt.chunk)
			if tt.wantErr && err == nil {
				t.Errorf("getDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getDir() got = %v, want %v", got, tt.want)
			}

			if got != "files" {
				err := os.RemoveAll(got)
				if err != nil {
					t.Error(err)
				}
			}
		})
	}
}
