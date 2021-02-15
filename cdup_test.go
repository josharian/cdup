package cdup

import (
	"errors"
	"io/fs"
	"os"
	"testing"
	"testing/fstest"
)

func TestFind(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	_, err = Find(pwd, "no-way-this-file-exists-on-your-computer-by-accident")
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected os.ErrNotExist, got %v", err)
	}
	found, err := Find(pwd, "cdup_test.go")
	if err != nil {
		t.Fatalf("expected to find this test file, got %v", err)
	}
	if found != pwd {
		t.Fatalf("looking for this test file, expected %v got %v", pwd, found)
	}
}

func TestFindIn(t *testing.T) {
	tests := []struct {
		fsys fs.FS
		want string
	}{
		// Success cases
		{
			fsys: fstest.MapFS{"a/b/c/x": &fstest.MapFile{}},
			want: "a/b/c",
		},
		{
			fsys: fstest.MapFS{"a/b/x": &fstest.MapFile{}},
			want: "a/b",
		},
		{
			fsys: fstest.MapFS{"a/x": &fstest.MapFile{}},
			want: "a",
		},
		{
			fsys: fstest.MapFS{"x": &fstest.MapFile{}},
			want: ".",
		},
		// Error cases
		{
			fsys: fstest.MapFS{"a/b/c/d/x": &fstest.MapFile{}},
		},
		{
			fsys: fstest.MapFS{"a/b/c": &fstest.MapFile{}},
		},
		{
			fsys: fstest.MapFS{"/x": &fstest.MapFile{}},
		},
		{
			fsys: fstest.MapFS{},
		},
	}

	for _, tt := range tests {
		got, err := FindIn(tt.fsys, "a/b/c", "x")
		if tt.want != "" && err != nil {
			t.Errorf("FindIn(%v): %v", tt.fsys, err)
			continue
		}
		if got != tt.want {
			t.Errorf("FindIn(%v) = %q want %q", tt.fsys, got, tt.want)
		}
	}
}
