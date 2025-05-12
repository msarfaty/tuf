package file

import (
	"path/filepath"
	"reflect"
	"testing"

	"mikesarfaty.com/tuf/internal/testutils"
)

func TestGenerateMd5ForFiles(t *testing.T) {
	type args struct {
		contents map[string]string
		names    []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "works with one file",
			args: args{
				contents: map[string]string{
					"hi.txt": "foo",
				},
				names: []string{"hi.txt"},
			},
			want: map[string]string{
				"hi.txt": "acbd18db4cc2f85cedef654fccc4a4d8",
			},
			wantErr: false,
		},
		{
			name: "works with one file",
			args: args{
				contents: map[string]string{
					"hi.txt": "foo",
				},
				names: []string{"hi.txt"},
			},
			want: map[string]string{
				"hi.txt": "acbd18db4cc2f85cedef654fccc4a4d8",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup test; write file contents and get temp dir
			dir := testutils.MakeDirectory(t, &testutils.TempDirOpts{
				Contents: tt.args.contents,
			})
			paths := []string{}
			for _, fname := range tt.args.names {
				paths = append(paths, filepath.Join(dir, fname))
			}

			got, err := GenerateMd5ForFiles(paths)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateMd5ForFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateMd5ForFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}
