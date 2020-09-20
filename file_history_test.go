package file_history_test

import (
	"crypto/sha1"
	"encoding/hex"
	"reflect"
	"testing"

	git "github.com/go-git/go-git/v5"

	filehistory "github.com/tasdomas/go-git-file-history"
)

func TestFileHistory(t *testing.T) {
	repo, err := git.PlainOpenWithOptions("./", &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	sums, err := filehistory.FileHistory(repo, "testdata/test.txt", sha1.New())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(stringSums(sums), []string{"e6c4fbd4fe7607f3e6ebf68b2ea4ef694da7b4fe"}) {
		t.Errorf("unexpected file hash sums")
	}

	// Calling with a non-existent file returns error.
	sums, err = filehistory.FileHistory(repo, "testdata/no_such_file.txt", sha1.New())
	if err != filehistory.ErrNotFound {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
	if len(sums) != 0 {
		t.Errorf("expected no result, got %v", sums)
	}

	// Calling with a non-existent file that was in git history returns no error.
	sums, err = filehistory.FileHistory(repo, "testdata/removed.txt", sha1.New())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(stringSums(sums), []string{"ed104b931644491e2a0055e4e3b2d6de364c82a5"}) {
		t.Errorf("unexpected file hash sums")
	}
}

func stringSums(sums [][]byte) []string {
	result := make([]string, len(sums))
	for i, s := range sums {
		result[i] = hex.EncodeToString(s)
	}
	return result
}
