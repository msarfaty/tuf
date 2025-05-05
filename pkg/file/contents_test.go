package file

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestFileEndsWith(t *testing.T) {
	type args struct {
		name     string
		endsWith string
		content  string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "endswith works when should be true",
			args: args{
				name:     "example.txt",
				endsWith: "\n",
				content:  "hello world!\n",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "endswith works when should be false",
			args: args{
				name:     "example.txt",
				endsWith: "!",
				content:  "hello world!\n",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "endswith works when should be true (long chars)",
			args: args{
				name:     "example.txt",
				endsWith: "hello world!\n",
				content:  "hello world!\n",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "endswith works when should be false (long chars)",
			args: args{
				name:     "example.txt",
				endsWith: "hello world!",
				content:  "hello world!\n",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "endswith works when file is empty",
			args: args{
				name:     "example.txt",
				endsWith: "\n",
				content:  "",
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			td := t.TempDir()
			tf := filepath.Join(td, tt.args.name)
			err := os.WriteFile(tf, []byte(tt.args.content), 0644)
			if err != nil {
				t.Fatal(err)
			}
			got, err := FileEndsWith(tf, tt.args.endsWith)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileEndsWith() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FileEndsWith() = %v, want %v", got, tt.want)
			}
		})
	}
}

var (
	tfSnippet1 string = `resource "aws_security_group" "foo" {
	count = 1
	name  = "test"
}`
	tfSnippet2 = `resource "aws_security_group" "bar" {
	count = 2
	name = baz
}`
	tfInput1 = fmt.Sprintf("%s\n\n\n\n%s\n", tfSnippet1, tfSnippet2)
	tfWant1  = fmt.Sprintf("%s\n\n%s\n", tfSnippet1, tfSnippet2)
	tfInput2 = fmt.Sprintf("%s\n\n\n\n", tfSnippet1)
	tfWant2  = fmt.Sprintf("%s\n\n", tfSnippet1)
)

func TestDeleteOverNOccurrences(t *testing.T) {
	type args struct {
		content     []byte
		chars       []byte
		start       int
		occurrences int
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "works simply",
			args: args{
				content:     []byte("aaaa"),
				chars:       []byte("a"),
				start:       0,
				occurrences: 2,
			},
			want: []byte("aa"),
		},
		{
			name: "works more complexly",
			args: args{
				content:     []byte("bcaxyzxyzxyzdd"),
				chars:       []byte("xyz"),
				start:       5,
				occurrences: 1,
			},
			want: []byte("bcaxyzdd"),
		},
		{
			name: "works with impossible input (chars > content)",
			args: args{
				content:     []byte("xyz"),
				chars:       []byte("xyza"),
				start:       0,
				occurrences: 1,
			},
			want: []byte("xyz"),
		},
		{
			name: "works with impossible input (start position too high)",
			args: args{
				content:     []byte("abc"),
				chars:       []byte("a"),
				start:       3,
				occurrences: 1,
			},
			want: []byte("abc"),
		},
		{
			name: "works with impossible input (chars never occurs)",
			args: args{
				content:     []byte("bbb"),
				chars:       []byte("a"),
				start:       0,
				occurrences: 1,
			},
			want: []byte("bbb"),
		},
		{
			name: "works with after 0 occurrences (chars never occurs)",
			args: args{
				content:     []byte("ababab"),
				chars:       []byte("ab"),
				start:       0,
				occurrences: 0,
			},
			want: []byte(""),
		},
		{
			name: "works with late start position",
			args: args{
				content:     []byte("aaaa"),
				chars:       []byte("aaa"),
				start:       3,
				occurrences: 0,
			},
			want: []byte("a"),
		},
		{
			name: "removes far before start pos if possible",
			args: args{
				content:     []byte("ababababab"),
				chars:       []byte("ab"),
				start:       9,
				occurrences: 1,
			},
			want: []byte("ab"),
		},
		{
			name: "works when the start position is at the end of the substring",
			args: args{
				content:     []byte("abcxyzxyz"),
				chars:       []byte("xyz"),
				start:       5,
				occurrences: 0,
			},
			want: []byte("abc"),
		},
		{
			name: "works with an example terraform context",
			args: args{
				content:     []byte(tfInput1),
				chars:       []byte("\n"),
				start:       len([]byte(tfSnippet1)) + 1,
				occurrences: 2,
			},
			want: []byte(tfWant1),
		},
		{
			name: "works with an example terraform context 2",
			args: args{
				content:     []byte(tfInput1),
				chars:       []byte("\n\n"),
				start:       len([]byte(tfSnippet1)),
				occurrences: 1,
			},
			want: []byte(tfWant1),
		},
		{
			name: "works with an example terraform context 3",
			args: args{
				content:     []byte(tfInput2),
				chars:       []byte("\n\n"),
				start:       len([]byte(tfSnippet1)),
				occurrences: 1,
			},
			want: []byte(tfWant2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DeleteOverNOccurrences(tt.args.content, tt.args.chars, tt.args.start, tt.args.occurrences)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteOverNOccurrences() = %s, want %s", string(got), string(tt.want))
			}
		})
	}
}
