package exec

import (
	"reflect"
	"testing"
)

func TestDiffFiles(t *testing.T) {
	before := map[string]struct{}{
		"M  internal/model/iam/adminusermodel.go": {},
	}
	after := map[string]struct{}{
		"M  internal/model/iam/adminusermodel.go": {},
		"?? db/migrations/create_table_order.sql": {},
		"?? db/migrations/init_order.sql":         {},
	}

	got := diffFiles(before, after)
	want := []string{"db/migrations/create_table_order.sql", "db/migrations/init_order.sql"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("diffFiles() = %v, want %v", got, want)
	}
}

func TestDiffFilesNoChange(t *testing.T) {
	snapshot := map[string]struct{}{
		"M  internal/model/iam/adminusermodel.go": {},
	}
	got := diffFiles(snapshot, snapshot)
	if len(got) != 0 {
		t.Fatalf("diffFiles() with identical snapshots = %v, want empty", got)
	}
}

func TestDiffFilesRename(t *testing.T) {
	before := map[string]struct{}{}
	after := map[string]struct{}{
		"R  db/migrations/old_name.sql -> db/migrations/new_name.sql": {},
	}
	got := diffFiles(before, after)
	want := []string{"db/migrations/new_name.sql"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("diffFiles() rename = %v, want %v", got, want)
	}
}
