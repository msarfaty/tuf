package file

import (
	"path"
	"reflect"
	"testing"

	"github.com/msarfaty/tuf/internal/testutils"
)

func defaultArgMutator(dir string, name string) string {
	return path.Join(dir, name)
}

func TestGenerateMd5ForFiles(t *testing.T) {
	type args struct {
		contents   map[string]string
		names      []string
		argMutator func(dir string, name string) string
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
				names:      []string{"hi.txt"},
				argMutator: defaultArgMutator,
			},
			want: map[string]string{
				"hi.txt": "acbd18db4cc2f85cedef654fccc4a4d8",
			},
			wantErr: false,
		},
		{
			name: "returns nothing if not checking a listed file",
			args: args{
				contents: map[string]string{
					"hi.txt": "foo",
				},
				names:      []string{},
				argMutator: defaultArgMutator,
			},
			want:    map[string]string{},
			wantErr: false,
		},
		{
			name: "works with many files",
			args: args{
				contents: map[string]string{
					"hi.txt":      "foo",
					"goodbye.txt": "foobar",
					"hi.tf":       "different",
				},
				names:      []string{"hi.txt", "hi.tf"},
				argMutator: defaultArgMutator,
			},
			want: map[string]string{
				"hi.txt": "acbd18db4cc2f85cedef654fccc4a4d8",
				"hi.tf":  "29e4b66fa8076de4d7a26c727b8dbdfa",
			},
			wantErr: false,
		},
		{
			name: "errors if trying to check a non-existent file",
			args: args{
				contents: map[string]string{
					"hi.txt": "foo",
				},
				names:      []string{"goodbye.txt"},
				argMutator: defaultArgMutator,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "errors if one of the attempted md5'd files is a dir",
			args: args{
				contents: map[string]string{
					"hi.txt":  "foo",
					"bye.txt": "foo",
				},
				names: []string{"goodbye.txt"},
				argMutator: func(d string, n string) string {
					if n == "hi.txt" {
						return d
					} else {
						return path.Join(d, n)
					}
				},
			},
			want:    nil,
			wantErr: true,
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
				paths = append(paths, tt.args.argMutator(dir, fname))
			}

			if tt.want != nil {
				wantWithPaths := map[string]string{}
				for k, v := range tt.want {
					wantWithPaths[path.Join(dir, k)] = v
				}
				tt.want = wantWithPaths
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
