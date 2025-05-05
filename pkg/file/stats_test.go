package file

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileIsEmpty(t *testing.T) {
	type args struct {
		name    string
		content string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "file is empty",
			args: args{
				name:    "example.txt",
				content: "",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "file is not empty",
			args: args{
				name:    "example.txt",
				content: "hello world",
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			td := t.TempDir()
			tf := filepath.Join(td, tt.name)
			os.WriteFile(tf, []byte(tt.args.content), 0644)

			got, err := FileIsEmpty(tf)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileIsEmpty() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FileIsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}
