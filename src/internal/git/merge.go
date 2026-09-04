package git

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func MergeCurrentBranch(directory string) error {
	repo, worktree, err := openRepositoryAndWorktree(directory)
	if err != nil {
		return err
	}

	status, err := worktree.Status()
	if err != nil {
		return err
	}
	if !status.IsClean() {
		return fmt.Errorf("merge requires a clean worktree")
	}

	headRef, err := repo.Head()
	if err != nil {
		return err
	}
	if !headRef.Name().IsBranch() {
		return fmt.Errorf("merge requires HEAD to point to a branch")
	}

	headCommit, err := repo.CommitObject(headRef.Hash())
	if err != nil {
		return err
	}

	branchCfg, upstreamRef, upstreamCommit, err := resolveUpstream(repo, headRef)
	if err != nil {
		return err
	}

	if headCommit.Hash == upstreamCommit.Hash {
		return nil
	}

	isUpstreamAncestor, err := upstreamCommit.IsAncestor(headCommit)
	if err != nil {
		return err
	}
	if isUpstreamAncestor {
		return nil
	}

	isHeadAncestor, err := headCommit.IsAncestor(upstreamCommit)
	if err != nil {
		return err
	}
	if isHeadAncestor {
		if err := repo.Merge(*upstreamRef, git.MergeOptions{}); err != nil {
			return err
		}

		return worktree.Reset(&git.ResetOptions{Mode: git.HardReset})
	}

	mergeBases, err := headCommit.MergeBase(upstreamCommit)
	if err != nil {
		return err
	}
	if len(mergeBases) == 0 {
		return fmt.Errorf("merge requires a common ancestor")
	}
	if len(mergeBases) > 1 {
		return fmt.Errorf("multiple merge bases are unsupported")
	}

	baseCommit := mergeBases[0]
	baseTree, err := baseCommit.Tree()
	if err != nil {
		return err
	}
	headTree, err := headCommit.Tree()
	if err != nil {
		return err
	}
	upstreamTree, err := upstreamCommit.Tree()
	if err != nil {
		return err
	}

	paths, err := mergeChangedPaths(baseTree, headTree, upstreamTree)
	if err != nil {
		return err
	}

	resolver := mergeResolver{
		repo:           repo,
		baseCommit:     baseCommit,
		headCommit:     headCommit,
		upstreamCommit: upstreamCommit,
	}

	for _, path := range paths {
		resolved, err := resolver.resolvePath(baseTree, headTree, upstreamTree, path)
		if err != nil {
			return err
		}

		current, err := snapshotFromTree(headTree, path)
		if err != nil {
			return err
		}
		if snapshotsEqual(current, resolved) {
			continue
		}

		if err := applyResolvedPath(directory, path, current, resolved); err != nil {
			return err
		}

		if resolved.exists {
			if _, err := worktree.Add(path); err != nil {
				return err
			}
		} else if _, err := worktree.Remove(path); err != nil {
			return err
		}
	}

	status, err = worktree.Status()
	if err != nil {
		return err
	}
	if status.IsClean() {
		return nil
	}

	author, committer, err := commitSignatures(repo)
	if err != nil {
		return err
	}

	_, err = worktree.Commit(
		fmt.Sprintf(
			"Merge remote-tracking branch '%s/%s' into %s",
			branchCfg.Remote,
			branchCfg.Merge.Short(),
			headRef.Name().Short(),
		),
		&git.CommitOptions{
			Author:    author,
			Committer: committer,
			Parents:   []plumbing.Hash{headCommit.Hash, upstreamCommit.Hash},
		},
	)
	if errors.Is(err, git.ErrEmptyCommit) {
		return commitEmptyTreeMerge(repo, worktree, headCommit, upstreamCommit, author, committer, fmt.Sprintf(
			"Merge remote-tracking branch '%s/%s' into %s",
			branchCfg.Remote,
			branchCfg.Merge.Short(),
			headRef.Name().Short(),
		))
	}
	return err
}

func commitEmptyTreeMerge(
	repo *git.Repository,
	worktree *git.Worktree,
	headCommit *object.Commit,
	upstreamCommit *object.Commit,
	author *object.Signature,
	committer *object.Signature,
	message string,
) error {
	idx, err := repo.Storer.Index()
	if err != nil {
		return err
	}
	if len(idx.Entries) != 0 {
		return git.ErrEmptyCommit
	}

	emptyTreeHash, err := ensureEmptyTree(repo)
	if err != nil {
		return err
	}
	if headCommit.TreeHash == emptyTreeHash {
		return nil
	}

	commit := &object.Commit{
		Author:       *author,
		Committer:    *committer,
		Message:      message,
		TreeHash:     emptyTreeHash,
		ParentHashes: []plumbing.Hash{headCommit.Hash, upstreamCommit.Hash},
	}

	obj := repo.Storer.NewEncodedObject()
	if err := commit.Encode(obj); err != nil {
		return err
	}

	commitHash, err := repo.Storer.SetEncodedObject(obj)
	if err != nil {
		return err
	}

	headRef, err := repo.Storer.Reference(plumbing.HEAD)
	if err != nil {
		return err
	}

	refName := plumbing.HEAD
	if headRef.Type() != plumbing.HashReference {
		refName = headRef.Target()
	}

	if err := repo.Storer.SetReference(plumbing.NewHashReference(refName, commitHash)); err != nil {
		return err
	}

	return worktree.Reset(&git.ResetOptions{Commit: commitHash, Mode: git.HardReset})
}

func ensureEmptyTree(repo *git.Repository) (plumbing.Hash, error) {
	obj := repo.Storer.NewEncodedObject()
	tree := &object.Tree{}
	if err := tree.Encode(obj); err != nil {
		return plumbing.ZeroHash, err
	}

	return repo.Storer.SetEncodedObject(obj)
}

type mergeResolver struct {
	repo           *git.Repository
	baseCommit     *object.Commit
	headCommit     *object.Commit
	upstreamCommit *object.Commit

	headLatestCache     map[string]time.Time
	upstreamLatestCache map[string]time.Time
}

type pathSnapshot struct {
	exists   bool
	mode     filemode.FileMode
	contents []byte
}

type lineEdit struct {
	start       int
	end         int
	replacement []string
}

func resolveUpstream(repo *git.Repository, headRef *plumbing.Reference) (*config.Branch, *plumbing.Reference, *object.Commit, error) {
	branchCfg, err := repo.Branch(headRef.Name().Short())
	if err != nil {
		return nil, nil, nil, err
	}
	if branchCfg.Remote == "" {
		return nil, nil, nil, fmt.Errorf("current branch has no configured remote")
	}
	if branchCfg.Merge == "" {
		return nil, nil, nil, fmt.Errorf("current branch has no configured upstream branch")
	}

	upstreamName := plumbing.NewRemoteReferenceName(branchCfg.Remote, branchCfg.Merge.Short())
	upstreamRef, err := repo.Reference(upstreamName, true)
	if err != nil {
		return nil, nil, nil, err
	}

	upstreamCommit, err := repo.CommitObject(upstreamRef.Hash())
	if err != nil {
		return nil, nil, nil, err
	}

	return branchCfg, upstreamRef, upstreamCommit, nil
}

func mergeChangedPaths(baseTree *object.Tree, headTree *object.Tree, upstreamTree *object.Tree) ([]string, error) {
	pathSet := map[string]struct{}{}
	for _, trees := range [][2]*object.Tree{{baseTree, headTree}, {baseTree, upstreamTree}} {
		changes, err := object.DiffTree(trees[0], trees[1])
		if err != nil {
			return nil, err
		}
		for _, change := range changes {
			pathSet[changeName(change)] = struct{}{}
		}
	}

	paths := make([]string, 0, len(pathSet))
	for path := range pathSet {
		paths = append(paths, path)
	}
	slices.Sort(paths)

	return paths, nil
}

func changeName(change *object.Change) string {
	if change.From.Name != "" {
		return change.From.Name
	}

	return change.To.Name
}

func (r *mergeResolver) resolvePath(baseTree *object.Tree, headTree *object.Tree, upstreamTree *object.Tree, path string) (pathSnapshot, error) {
	base, err := snapshotFromTree(baseTree, path)
	if err != nil {
		return pathSnapshot{}, err
	}
	head, err := snapshotFromTree(headTree, path)
	if err != nil {
		return pathSnapshot{}, err
	}
	upstream, err := snapshotFromTree(upstreamTree, path)
	if err != nil {
		return pathSnapshot{}, err
	}

	switch {
	case snapshotsEqual(head, upstream):
		return head, nil
	case snapshotsEqual(head, base):
		return upstream, nil
	case snapshotsEqual(upstream, base):
		return head, nil
	}

	if !head.exists || !upstream.exists {
		return r.pickLatest(path, head, upstream)
	}

	if head.mode != upstream.mode {
		return r.pickLatest(path, head, upstream)
	}

	if head.mode == filemode.Symlink {
		return r.pickLatest(path, head, upstream)
	}

	if !head.mode.IsFile() || !upstream.mode.IsFile() {
		return r.pickLatest(path, head, upstream)
	}

	merged, ok := mergeTextContents(base.contents, head.contents, upstream.contents)
	if !ok {
		return r.pickLatest(path, head, upstream)
	}

	return pathSnapshot{
		exists:   true,
		mode:     head.mode,
		contents: merged,
	}, nil
}

func (r *mergeResolver) pickLatest(path string, head pathSnapshot, upstream pathSnapshot) (pathSnapshot, error) {
	headTime, err := r.latestChangeTime(path, true)
	if err != nil {
		return pathSnapshot{}, err
	}
	upstreamTime, err := r.latestChangeTime(path, false)
	if err != nil {
		return pathSnapshot{}, err
	}

	if upstreamTime.After(headTime) || upstreamTime.Equal(headTime) {
		return upstream, nil
	}

	return head, nil
}

func (r *mergeResolver) latestChangeTime(path string, headSide bool) (time.Time, error) {
	var cache map[string]time.Time
	var from *object.Commit
	if headSide {
		if r.headLatestCache == nil {
			r.headLatestCache = map[string]time.Time{}
		}
		cache = r.headLatestCache
		from = r.headCommit
	} else {
		if r.upstreamLatestCache == nil {
			r.upstreamLatestCache = map[string]time.Time{}
		}
		cache = r.upstreamLatestCache
		from = r.upstreamCommit
	}

	if when, ok := cache[path]; ok {
		return when, nil
	}

	iter, err := r.repo.Log(&git.LogOptions{
		From:     from.Hash,
		Order:    git.LogOrderCommitterTime,
		FileName: &path,
	})
	if err != nil {
		return time.Time{}, err
	}
	defer iter.Close()

	for {
		commit, err := iter.Next()
		if err == io.EOF {
			cache[path] = time.Time{}
			return time.Time{}, nil
		}
		if err != nil {
			return time.Time{}, err
		}

		isBaseAncestor, err := commit.IsAncestor(r.baseCommit)
		if err != nil {
			return time.Time{}, err
		}
		if isBaseAncestor {
			cache[path] = time.Time{}
			return time.Time{}, nil
		}

		cache[path] = commit.Committer.When
		return commit.Committer.When, nil
	}
}

func snapshotFromTree(tree *object.Tree, path string) (pathSnapshot, error) {
	file, err := tree.File(path)
	if err == object.ErrFileNotFound {
		return pathSnapshot{}, nil
	}
	if err != nil {
		return pathSnapshot{}, err
	}

	contents, err := readObjectFile(file)
	if err != nil {
		return pathSnapshot{}, err
	}

	return pathSnapshot{
		exists:   true,
		mode:     file.Mode,
		contents: contents,
	}, nil
}

func readObjectFile(file *object.File) ([]byte, error) {
	reader, err := file.Reader()
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return io.ReadAll(reader)
}

func snapshotsEqual(a pathSnapshot, b pathSnapshot) bool {
	if a.exists != b.exists {
		return false
	}
	if !a.exists {
		return true
	}
	if a.mode != b.mode {
		return false
	}

	return bytes.Equal(a.contents, b.contents)
}

func applyResolvedPath(directory string, path string, current pathSnapshot, resolved pathSnapshot) error {
	fullPath := filepath.Join(directory, filepath.FromSlash(path))
	if !resolved.exists {
		if current.exists {
			if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
				return err
			}
		}

		return nil
	}

	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return err
	}

	if current.exists && current.mode == filemode.Symlink {
		if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	if current.exists && current.mode != filemode.Symlink && resolved.mode == filemode.Symlink {
		if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	switch resolved.mode {
	case filemode.Regular, filemode.Deprecated, filemode.Executable:
		perm := os.FileMode(0o644)
		if resolved.mode == filemode.Executable {
			perm = 0o755
		}
		if err := os.WriteFile(fullPath, resolved.contents, perm); err != nil {
			return err
		}
		return os.Chmod(fullPath, perm)
	case filemode.Symlink:
		return os.Symlink(string(resolved.contents), fullPath)
	default:
		return fmt.Errorf("unsupported file mode for merge: %s", resolved.mode)
	}
}

func commitSignatures(repo *git.Repository) (*object.Signature, *object.Signature, error) {
	cfg, err := repo.Config()
	if err != nil {
		return nil, nil, err
	}

	now := time.Now()
	name := cfg.User.Name
	email := cfg.User.Email
	if name == "" || email == "" {
		name = "GitJournal"
		email = "gitjournal@example.com"
	}

	author := &object.Signature{Name: name, Email: email, When: now}
	committer := &object.Signature{Name: name, Email: email, When: now}

	return author, committer, nil
}

func mergeTextContents(base []byte, head []byte, upstream []byte) ([]byte, bool) {
	if !utf8.Valid(base) || !utf8.Valid(head) || !utf8.Valid(upstream) {
		return nil, false
	}

	mergedLines, ok := mergeTextLines(
		splitLines(string(base)),
		splitLines(string(head)),
		splitLines(string(upstream)),
	)
	if !ok {
		return nil, false
	}

	return []byte(strings.Join(mergedLines, "")), true
}

func mergeTextLines(base []string, head []string, upstream []string) ([]string, bool) {
	switch {
	case slices.Equal(head, upstream):
		return append([]string(nil), head...), true
	case slices.Equal(head, base):
		return append([]string(nil), upstream...), true
	case slices.Equal(upstream, base):
		return append([]string(nil), head...), true
	}

	headEdits := diffLineEdits(base, head)
	upstreamEdits := diffLineEdits(base, upstream)

	var merged []string
	pos, headIndex, upstreamIndex := 0, 0, 0
	for pos < len(base) || headIndex < len(headEdits) || upstreamIndex < len(upstreamEdits) {
		nextPos := len(base)
		if headIndex < len(headEdits) && headEdits[headIndex].start < nextPos {
			nextPos = headEdits[headIndex].start
		}
		if upstreamIndex < len(upstreamEdits) && upstreamEdits[upstreamIndex].start < nextPos {
			nextPos = upstreamEdits[upstreamIndex].start
		}

		if nextPos > pos {
			merged = append(merged, base[pos:nextPos]...)
			pos = nextPos
			continue
		}

		regionStart := pos
		regionEnd := pos
		headLimit, upstreamLimit := headIndex, upstreamIndex
		for {
			consumed := false
			for headLimit < len(headEdits) && editOverlapsRegion(headEdits[headLimit], regionStart, regionEnd) {
				if headEdits[headLimit].end > regionEnd {
					regionEnd = headEdits[headLimit].end
				}
				headLimit++
				consumed = true
			}
			for upstreamLimit < len(upstreamEdits) && editOverlapsRegion(upstreamEdits[upstreamLimit], regionStart, regionEnd) {
				if upstreamEdits[upstreamLimit].end > regionEnd {
					regionEnd = upstreamEdits[upstreamLimit].end
				}
				upstreamLimit++
				consumed = true
			}
			if !consumed {
				break
			}
		}

		baseRegion := append([]string(nil), base[regionStart:regionEnd]...)
		headRegion := materializeRegion(base, regionStart, regionEnd, headEdits[headIndex:headLimit])
		upstreamRegion := materializeRegion(base, regionStart, regionEnd, upstreamEdits[upstreamIndex:upstreamLimit])

		resolved, ok := resolveTextRegion(baseRegion, headRegion, upstreamRegion)
		if !ok {
			return nil, false
		}

		merged = append(merged, resolved...)
		pos = regionEnd
		headIndex = headLimit
		upstreamIndex = upstreamLimit
	}

	return merged, true
}

func resolveTextRegion(base []string, head []string, upstream []string) ([]string, bool) {
	switch {
	case slices.Equal(head, upstream):
		return head, true
	case slices.Equal(head, base):
		return upstream, true
	case slices.Equal(upstream, base):
		return head, true
	}

	prefixLen := sharedPrefixLen(head, upstream)
	suffixLen := sharedSuffixLen(head[prefixLen:], upstream[prefixLen:])

	if prefixLen > 0 || suffixLen > 0 {
		baseStart := min(prefixLen, len(base))
		baseEnd := len(base) - suffixLen
		if baseEnd < baseStart {
			baseEnd = baseStart
		}
		resolvedMiddle, ok := resolveTextRegion(
			base[baseStart:baseEnd],
			head[prefixLen:len(head)-suffixLen],
			upstream[prefixLen:len(upstream)-suffixLen],
		)
		if !ok {
			return nil, false
		}

		resolved := append([]string(nil), head[:prefixLen]...)
		resolved = append(resolved, resolvedMiddle...)
		resolved = append(resolved, head[len(head)-suffixLen:]...)
		return resolved, true
	}

	return nil, false
}

func materializeRegion(base []string, regionStart int, regionEnd int, edits []lineEdit) []string {
	var result []string
	cursor := regionStart
	for _, edit := range edits {
		if edit.start > cursor {
			result = append(result, base[cursor:edit.start]...)
		}
		result = append(result, edit.replacement...)
		cursor = edit.end
	}
	if cursor < regionEnd {
		result = append(result, base[cursor:regionEnd]...)
	}

	return result
}

func diffLineEdits(base []string, target []string) []lineEdit {
	n, m := len(base), len(target)
	lcs := make([][]int, n+1)
	for i := range lcs {
		lcs[i] = make([]int, m+1)
	}

	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			if base[i] == target[j] {
				lcs[i][j] = lcs[i+1][j+1] + 1
				continue
			}
			if lcs[i+1][j] >= lcs[i][j+1] {
				lcs[i][j] = lcs[i+1][j]
			} else {
				lcs[i][j] = lcs[i][j+1]
			}
		}
	}

	var edits []lineEdit
	i, j := 0, 0
	active := -1
	for i < n || j < m {
		if i < n && j < m && base[i] == target[j] {
			active = -1
			i++
			j++
			continue
		}

		if active == -1 {
			edits = append(edits, lineEdit{start: i, end: i})
			active = len(edits) - 1
		}

		if j < m && (i == n || lcs[i][j+1] > lcs[i+1][j]) {
			edits[active].replacement = append(edits[active].replacement, target[j])
			j++
			continue
		}

		if i < n {
			i++
			edits[active].end = i
		}
	}

	return edits
}

func editOverlapsRegion(edit lineEdit, regionStart int, regionEnd int) bool {
	if regionStart == regionEnd {
		return edit.start == regionStart
	}

	return edit.start < regionEnd
}

func splitLines(content string) []string {
	if content == "" {
		return nil
	}

	lines := strings.SplitAfter(content, "\n")
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	return lines
}

func sharedPrefixLen(a []string, b []string) int {
	limit := min(len(a), len(b))
	for i := 0; i < limit; i++ {
		if a[i] != b[i] {
			return i
		}
	}

	return limit
}

func sharedSuffixLen(a []string, b []string) int {
	limit := min(len(a), len(b))
	for i := 1; i <= limit; i++ {
		if a[len(a)-i] != b[len(b)-i] {
			return i - 1
		}
	}

	return limit
}

func min(a int, b int) int {
	if a < b {
		return a
	}

	return b
}
