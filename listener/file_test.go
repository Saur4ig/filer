package main

import (
	"os"
	"testing"
)

func Test_chunk_createNewFileName(t *testing.T) {
	tests := []struct {
		name    string
		chunk   *chunk
		want    string
		wantErr bool
	}{
		{
			name: "Example 1",
			chunk: &chunk{
				ID:            1,
				totalChunks:   2,
				chunkTemplate: "chunk-*.txt",
				baseFilename:  "test1.txt",
			},
			wantErr: false,
			want:    "test1.txt_dir/new_test1.txt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := getDir(tt.chunk)
			if err != nil {
				t.Error(err)
				return
			}

			got, err := tt.chunk.createNewFileName(dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("createNewFileName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("createNewFileName() got = %v, want %v", got, tt.want)
			}
			err = os.RemoveAll(dir)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
