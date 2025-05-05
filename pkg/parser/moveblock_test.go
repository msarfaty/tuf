package parser

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"testing"
)

func TestMoveHclBlock(t *testing.T) {
	type args struct {
		testDir string
		mo      *MoveOptions
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "can move a single resource",
			args: args{
				testDir: path.Join("testdata", "moves", "test1"),
				mo: &MoveOptions{
					Address:  "aws_iam_role.eks_auto",
					FromFile: "original.tf",
					ToFile:   "moved.tf",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origDir := t.TempDir()
			newDir := t.TempDir()
			input, err := os.ReadFile(path.Join(tt.args.testDir, "before.txt"))
			if err != nil {
				t.Fatal(err)
			}
			os.WriteFile(path.Join(origDir, tt.args.mo.FromFile), input, 0644)
			tt.args.mo.ToFile = path.Join(newDir, tt.args.mo.ToFile)
			tt.args.mo.FromFile = path.Join(origDir, tt.args.mo.FromFile)
			wantMoved, err := os.ReadFile(path.Join(tt.args.testDir, "after_moved.txt"))
			if err != nil {
				t.Fatal(err)
			}
			wantRemoved, err := os.ReadFile(path.Join(tt.args.testDir, "after_original.txt"))
			if err != nil {
				t.Fatal(err)
			}

			if err := MoveHclBlock(tt.args.mo); (err != nil) != tt.wantErr {
				t.Errorf("MoveHclBlock() error = %v, wantErr %v", err, tt.wantErr)
			}
			gotMoved, err := os.ReadFile(tt.args.mo.ToFile)
			if err != nil {
				t.Fatal(fmt.Errorf("test failure: %w", err))
			}
			gotRemoved, err := os.ReadFile(tt.args.mo.FromFile)
			if err != nil {
				t.Fatal(fmt.Errorf("test failure: %w", err))
			}

			if !reflect.DeepEqual(gotMoved, wantMoved) {
				t.Errorf("MoveHclBlock() moved file =\nSTART%sEOF\nwant\nSTART%sEOF", string(gotMoved), string(wantMoved))
			}
			if !reflect.DeepEqual(gotRemoved, wantRemoved) {
				t.Errorf("MoveHclBlock() removed file =\nSTART%sEOF, want\nSTART%sEOF", string(gotRemoved), string(wantRemoved))
			}
		})
	}
}
