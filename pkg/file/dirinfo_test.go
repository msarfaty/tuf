package file

import (
	"path"
	"reflect"
	"slices"
	"testing"

	"mikesarfaty.com/tuf/internal/testutils"
)

func TestGetAllTerraformFilesInDirectory(t *testing.T) {
	type args struct {
		contents map[string]string
		mutator  func(dir string) string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "works with one file",
			args: args{
				contents: map[string]string{"main.tf": "hi"},
				mutator:  func(dir string) string { return dir },
			},
			want:    []string{"main.tf"},
			wantErr: false,
		},
		{
			name: "works with a mix of files",
			args: args{
				contents: map[string]string{"main.tf": "abc", "data.tf": "def", "backend.tf": "some backend config"},
				mutator:  func(dir string) string { return dir },
			},
			want:    []string{"main.tf", "data.tf", "backend.tf"},
			wantErr: false,
		},
		{
			name: "errors when providing a non-existent directory",
			args: args{
				contents: map[string]string{"main.tf": "abc", "data.tf": "def", "backend.tf": "some backend config"},
				mutator:  func(dir string) string { return "/var/some/nonexistent/directory" },
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "errors when providing a file",
			args: args{
				contents: map[string]string{"main.tf": "abc", "data.tf": "def", "backend.tf": "some backend config"},
				mutator:  func(dir string) string { return path.Join(dir, "main.tf") },
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup test
			dir := testutils.MakeDirectory(t, &testutils.TempDirOpts{
				Contents: tt.args.contents,
			})
			dir = tt.args.mutator(dir)

			for i, want := range tt.want {
				tt.want[i] = path.Join(dir, want)
			}
			slices.Sort(tt.want)

			// got should return sorted; no need here
			got, err := GetAllTerraformFilesInDirectory(dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllTerraformFilesInDirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllTerraformFilesInDirectory() = %v, want %v", got, tt.want)
			}
		})
	}
}
