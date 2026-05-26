package git

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func TestAddStagesNewFile(t *testing.T) {
	dir, _, _ := newTestRepo(t)
	writeFile(t, dir, "note.md", "hello")

	if err := Add(dir, "note.md"); err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	status := worktreeStatus(t, dir)
	if got := status.File("note.md").Staging; got != gogit.Added {
		t.Fatalf("staging status = %v, want %v", got, gogit.Added)
	}
}

func TestRemoveDeletesFromWorktreeAndIndex(t *testing.T) {
	dir, _, w := newTestRepo(t)
	writeFile(t, dir, "note.md", "hello")
	commitFile(t, w, "note.md")

	if err := Remove(dir, "note.md"); err != nil {
		t.Fatalf("Remove() error = %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, "note.md")); !os.IsNotExist(err) {
		t.Fatalf("removed file stat error = %v, want not exist", err)
	}

	status := worktreeStatus(t, dir)
	if got := status.File("note.md").Staging; got != gogit.Deleted {
		t.Fatalf("staging status = %v, want %v", got, gogit.Deleted)
	}
}

func TestMoveRenamesInWorktreeAndIndex(t *testing.T) {
	dir, _, w := newTestRepo(t)
	writeFile(t, dir, "note.md", "hello")
	commitFile(t, w, "note.md")

	if err := Move(dir, "note.md", "renamed.md"); err != nil {
		t.Fatalf("Move() error = %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, "note.md")); !os.IsNotExist(err) {
		t.Fatalf("old file stat error = %v, want not exist", err)
	}
	if got := readFile(t, dir, "renamed.md"); got != "hello" {
		t.Fatalf("new file content = %q, want %q", got, "hello")
	}

	status := worktreeStatus(t, dir)
	if got := status.File("note.md").Staging; got != gogit.Deleted {
		t.Fatalf("old path staging status = %v, want %v", got, gogit.Deleted)
	}
	if got := status.File("renamed.md").Staging; got != gogit.Added {
		t.Fatalf("new path staging status = %v, want %v", got, gogit.Added)
	}
}

func TestMoveFailsForMissingSource(t *testing.T) {
	dir, _, _ := newTestRepo(t)

	if err := Move(dir, "missing.md", "renamed.md"); err == nil {
		t.Fatal("Move() error = nil, want error")
	}
}

func TestMoveFailsForExistingDestination(t *testing.T) {
	dir, _, w := newTestRepo(t)
	writeFile(t, dir, "note.md", "hello")
	commitFile(t, w, "note.md")
	writeFile(t, dir, "renamed.md", "existing")

	if err := Move(dir, "note.md", "renamed.md"); err == nil {
		t.Fatal("Move() error = nil, want error")
	}
}

func TestResetHardRestoresCurrentHead(t *testing.T) {
	dir, _, w := newTestRepo(t)
	writeFile(t, dir, "note.md", "original")
	commitFile(t, w, "note.md")
	writeFile(t, dir, "note.md", "changed")

	if err := ResetHard(dir); err != nil {
		t.Fatalf("ResetHard() error = %v", err)
	}

	if got := readFile(t, dir, "note.md"); got != "original" {
		t.Fatalf("file content = %q, want %q", got, "original")
	}
}

func TestResetHardToMovesHeadAndWorktree(t *testing.T) {
	dir, repo, w := newTestRepo(t)
	writeFile(t, dir, "note.md", "one")
	first := commitFile(t, w, "note.md")
	writeFile(t, dir, "note.md", "two")
	commitFile(t, w, "note.md")

	if err := ResetHardTo(dir, first.String()); err != nil {
		t.Fatalf("ResetHardTo() error = %v", err)
	}

	head, err := repo.Head()
	if err != nil {
		t.Fatalf("Head() error = %v", err)
	}
	if got := head.Hash(); got != first {
		t.Fatalf("HEAD = %s, want %s", got, first)
	}
	if got := readFile(t, dir, "note.md"); got != "one" {
		t.Fatalf("file content = %q, want %q", got, "one")
	}
}

func TestResetHardToFailsForInvalidHash(t *testing.T) {
	dir, _, w := newTestRepo(t)
	writeFile(t, dir, "note.md", "one")
	commitFile(t, w, "note.md")

	if err := ResetHardTo(dir, "not-a-hash"); err == nil {
		t.Fatal("ResetHardTo() error = nil, want error")
	}
}

func TestSwitchSwitchesExistingLocalBranches(t *testing.T) {
	dir, _, w := newTestRepo(t)
	writeFile(t, dir, "note.md", "master")
	commitFile(t, w, "note.md")

	if err := w.Checkout(&gogit.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName("feature"),
		Create: true,
	}); err != nil {
		t.Fatalf("create feature branch error = %v", err)
	}
	writeFile(t, dir, "note.md", "feature")
	commitFile(t, w, "note.md")

	if err := Switch(dir, "master"); err != nil {
		t.Fatalf("Switch(master) error = %v", err)
	}
	if got := readFile(t, dir, "note.md"); got != "master" {
		t.Fatalf("master content = %q, want %q", got, "master")
	}

	if err := Switch(dir, "feature"); err != nil {
		t.Fatalf("Switch(feature) error = %v", err)
	}
	if got := readFile(t, dir, "note.md"); got != "feature" {
		t.Fatalf("feature content = %q, want %q", got, "feature")
	}
}

func TestSwitchFailsForMissingBranch(t *testing.T) {
	dir, _, w := newTestRepo(t)
	writeFile(t, dir, "note.md", "master")
	commitFile(t, w, "note.md")

	if err := Switch(dir, "missing"); err == nil {
		t.Fatal("Switch() error = nil, want error")
	}
}

func TestSwitchFailsForBlankBranch(t *testing.T) {
	dir, _, w := newTestRepo(t)
	writeFile(t, dir, "note.md", "master")
	commitFile(t, w, "note.md")

	if err := Switch(dir, " \t\n"); err == nil {
		t.Fatal("Switch() error = nil, want error")
	}
}

func TestPullFetchesThenMerges(t *testing.T) {
	var calls []string
	fetchFn := func(remote string, directory string, privateKey []byte, password string) error {
		calls = append(calls, "fetch:"+remote+":"+directory+":"+string(privateKey)+":"+password)
		return nil
	}
	mergeFn := func(directory string) error {
		calls = append(calls, "merge:"+directory)
		return nil
	}

	if err := pull(fetchFn, mergeFn, "origin", "/repo", []byte("pem"), "secret"); err != nil {
		t.Fatalf("pull() error = %v", err)
	}

	want := []string{"fetch:origin:/repo:pem:secret", "merge:/repo"}
	if len(calls) != len(want) {
		t.Fatalf("calls = %v, want %v", calls, want)
	}
	for i := range want {
		if calls[i] != want[i] {
			t.Fatalf("calls[%d] = %q, want %q", i, calls[i], want[i])
		}
	}
}

func TestPullReturnsFetchError(t *testing.T) {
	wantErr := errors.New("fetch failed")
	fetchFn := func(string, string, []byte, string) error {
		return wantErr
	}
	mergeCalled := false
	mergeFn := func(string) error {
		mergeCalled = true
		return nil
	}

	err := pull(fetchFn, mergeFn, "origin", "/repo", nil, "")
	if !errors.Is(err, wantErr) {
		t.Fatalf("pull() error = %v, want %v", err, wantErr)
	}
	if mergeCalled {
		t.Fatal("merge was called after fetch failure")
	}
}

func TestPullReturnsMergeError(t *testing.T) {
	wantErr := errors.New("merge failed")
	fetchFn := func(string, string, []byte, string) error {
		return nil
	}
	mergeFn := func(string) error {
		return wantErr
	}

	err := pull(fetchFn, mergeFn, "origin", "/repo", nil, "")
	if !errors.Is(err, wantErr) {
		t.Fatalf("pull() error = %v, want %v", err, wantErr)
	}
}

func newTestRepo(t *testing.T) (string, *gogit.Repository, *gogit.Worktree) {
	t.Helper()

	dir := t.TempDir()
	repo, err := gogit.PlainInit(dir, false)
	if err != nil {
		t.Fatalf("PlainInit() error = %v", err)
	}

	w, err := repo.Worktree()
	if err != nil {
		t.Fatalf("Worktree() error = %v", err)
	}

	return dir, repo, w
}

func commitFile(t *testing.T, w *gogit.Worktree, path string) plumbing.Hash {
	t.Helper()

	if _, err := w.Add(path); err != nil {
		t.Fatalf("Add(%q) error = %v", path, err)
	}

	hash, err := w.Commit("commit "+path, &gogit.CommitOptions{
		Author: &object.Signature{
			Name:  "GitJournal",
			Email: "gitjournal@example.com",
			When:  time.Unix(1, 0),
		},
	})
	if err != nil {
		t.Fatalf("Commit() error = %v", err)
	}

	return hash
}

func worktreeStatus(t *testing.T, dir string) gogit.Status {
	t.Helper()

	repo, err := gogit.PlainOpen(dir)
	if err != nil {
		t.Fatalf("PlainOpen() error = %v", err)
	}
	w, err := repo.Worktree()
	if err != nil {
		t.Fatalf("Worktree() error = %v", err)
	}
	status, err := w.Status()
	if err != nil {
		t.Fatalf("Status() error = %v", err)
	}

	return status
}

func writeFile(t *testing.T, dir string, path string, content string) {
	t.Helper()

	fullPath := filepath.Join(dir, path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
}

func readFile(t *testing.T, dir string, path string) string {
	t.Helper()

	bytes, err := os.ReadFile(filepath.Join(dir, path))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	return string(bytes)
}
