package parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
)

func TestNew(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    BlockDescription
		wantErr bool
	}{
		{
			name: "factory creates module block description",
			args: args{address: "module.foo"},
			want: &ModuleBlockDescription{
				name: "foo",
			},
			wantErr: false,
		},
		{
			name: "factory creates resource block description",
			args: args{address: "aws_security_group.foo"},
			want: &ResourceBlockDescription{
				rType: "aws_security_group",
				name:  "foo",
			},
			wantErr: false,
		},
		{
			name:    "factory does not create anything with data block",
			args:    args{address: "data.aws_security_group.foo"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "factory fails creating module block desc with too many parts",
			args:    args{address: "module.toomany.parts"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModuleBlockDescription_Matches(t *testing.T) {
	type fields struct {
		name string
	}
	type args struct {
		filename string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "block matches (given module block)",
			fields: fields{
				name: "foobar",
			},
			args: args{
				filename: "example1.tf",
			},
			want: true,
		},
		{
			name: "block doesn't match (given module block)",
			fields: fields{
				name: "foobaz",
			},
			args: args{
				filename: "example1.tf",
			},
			want: false,
		},
		{
			name: "block doesn't match (given resource block)",
			fields: fields{
				name: "foobar",
			},
			args: args{
				filename: "example2.tf",
			},
			want: false,
		},
		{
			name: "block doesn't match (given module block)",
			fields: fields{
				name: "foobaz",
			},
			args: args{
				filename: "example2.tf",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &ModuleBlockDescription{
				name: tt.fields.name,
			}
			hclp := hclparse.NewParser()
			f, diags := hclp.ParseHCLFile(fmt.Sprintf("testdata/%s", tt.args.filename))
			if diags.HasErrors() {
				panic(diags.Error())
			}
			blocks := f.BlocksAtPos(hcl.InitialPos)
			if got := m.Matches(*blocks[0]); got != tt.want {
				t.Errorf("ModuleBlockDescription.Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}
