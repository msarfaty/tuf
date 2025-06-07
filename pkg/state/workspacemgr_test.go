package state

import (
	"reflect"
	"sort"
	"testing"

	"github.com/google/uuid"
	"github.com/msarfaty/tuf/internal/testutils"
)

func TestWorkspaceMgr_AddWorkspace(t *testing.T) {
	type args struct {
		contents map[string]string
	}
	tests := []struct {
		name           string
		args           args
		wantWorkspaces []*Workspace
		wantErr        bool
	}{
		{
			name: "works with a single-file workspace",
			args: args{contents: map[string]string{"hi.tf": "foo"}},
			wantWorkspaces: []*Workspace{
				{
					Files: []*WorkspaceFile{
						{
							Name: "hi.tf",
							Md5:  "acbd18db4cc2f85cedef654fccc4a4d8",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "works with a many tf files",
			args: args{contents: map[string]string{"hi.tf": "foo", "config.tf": "foobarbaz"}},
			wantWorkspaces: []*Workspace{
				{
					Files: []*WorkspaceFile{
						{
							Name: "hi.tf",
							Md5:  "acbd18db4cc2f85cedef654fccc4a4d8",
						},
						{
							Name: "config.tf",
							Md5:  "6df23dc03f9b54cc38a0fc1483df6e21",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := testutils.MakeDirectory(t, &testutils.TempDirOpts{
				Contents: tt.args.contents,
			})
			wantWsmgr := &WorkspaceMgr{
				Workspaces: tt.wantWorkspaces,
			}

			gotWsmgr := NewWorkspaceMgr()
			if err := gotWsmgr.AddWorkspace(dir); (err != nil) != tt.wantErr {
				t.Errorf("WorkspaceMgr.AddWorkspace() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(wantWsmgr.Workspaces) != len(gotWsmgr.Workspaces) {
				t.Fatalf("WorkspaceMgr.AddWorkspace() workspaces count mismatch; want = %v, got = %v", wantWsmgr, gotWsmgr)
			}

			// check workspace files have same md5 hashes
			for i, gotWs := range gotWsmgr.Workspaces {
				wantWs := tt.wantWorkspaces[i]
				if len(gotWs.Uuid) != len(uuid.NewString()) || len(gotWs.Uuid) == 0 {
					t.Errorf("WorkspaceMgr.AddWorkspace() got malformed uuid for workspace: %v", gotWs)
				}
				if gotWs.Abspath != dir {
					t.Errorf("did not set workspace abspath correctly; got %s, want %s", gotWs.Abspath, dir)
				}
				gotWs.Uuid = ""
				gotWs.Abspath = ""
				sort.Slice(wantWs.Files, func(i, j int) bool {
					return wantWs.Files[i].Name < wantWs.Files[j].Name
				})
				sort.Slice(gotWs.Files, func(i, j int) bool {
					return gotWs.Files[i].Name < gotWs.Files[j].Name
				})

				if !reflect.DeepEqual(wantWs, gotWs) {
					t.Errorf("WorkspaceMgr.AddWorkspace() want = %v; got = %v", wantWs, gotWs)
				}
			}
		})
	}
}
