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
					files: []*WorkspaceFile{
						{
							name: "hi.tf",
							md5:  "acbd18db4cc2f85cedef654fccc4a4d8",
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
					files: []*WorkspaceFile{
						{
							name: "hi.tf",
							md5:  "acbd18db4cc2f85cedef654fccc4a4d8",
						},
						{
							name: "config.tf",
							md5:  "6df23dc03f9b54cc38a0fc1483df6e21",
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
				workspaces: tt.wantWorkspaces,
			}

			gotWsmgr := NewWorkspaceMgr()
			if err := gotWsmgr.AddWorkspace(dir); (err != nil) != tt.wantErr {
				t.Errorf("WorkspaceMgr.AddWorkspace() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(wantWsmgr.workspaces) != len(gotWsmgr.workspaces) {
				t.Fatalf("WorkspaceMgr.AddWorkspace() workspaces count mismatch; want = %v, got = %v", wantWsmgr, gotWsmgr)
			}

			// check workspace files have same md5 hashes
			for i, gotWs := range gotWsmgr.workspaces {
				wantWs := tt.wantWorkspaces[i]
				if len(gotWs.uuid) != len(uuid.NewString()) || len(gotWs.uuid) == 0 {
					t.Errorf("WorkspaceMgr.AddWorkspace() got malformed uuid for workspace: %v", gotWs)
				}
				if gotWs.abspath != dir {
					t.Errorf("did not set workspace abspath correctly; got %s, want %s", gotWs.abspath, dir)
				}
				gotWs.uuid = ""
				gotWs.abspath = ""
				sort.Slice(wantWs.files, func(i, j int) bool {
					return wantWs.files[i].name < wantWs.files[j].name
				})
				sort.Slice(gotWs.files, func(i, j int) bool {
					return gotWs.files[i].name < gotWs.files[j].name
				})

				if !reflect.DeepEqual(wantWs, gotWs) {
					t.Errorf("WorkspaceMgr.AddWorkspace() want = %v; got = %v", wantWs, gotWs)
				}
			}
		})
	}
}
