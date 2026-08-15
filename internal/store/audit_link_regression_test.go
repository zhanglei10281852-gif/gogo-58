package store

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestVerifyAuditRejectsSplicedChain replaces the middle entry of an audit log
// with an entry taken from an unrelated log. Every entry is internally
// consistent and the sequence numbers still run 1, 2, 3, so only the chaining
// between neighbouring entries can expose the splice.
func TestVerifyAuditRejectsSplicedChain(t *testing.T) {
	root := t.TempDir()
	original := auditLines(t, filepath.Join(root, "original"), "records=1", "records=2", "records=3")
	foreign := auditLines(t, filepath.Join(root, "foreign"), "other=1", "other=2", "other=3")

	handle, err := Open(filepath.Join(root, "spliced"))
	if err != nil {
		t.Fatalf("Open returned %v", err)
	}
	spliced := []string{original[0], foreign[1], original[2]}
	if err := os.WriteFile(handle.Path(AuditFile), []byte(strings.Join(spliced, "\n")+"\n"), 0o600); err != nil {
		t.Fatalf("writing audit log: %v", err)
	}

	entries, err := handle.LoadAudit()
	if err != nil {
		t.Fatalf("LoadAudit returned %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("the spliced log holds %d entries", len(entries))
	}
	for position, entry := range entries {
		if entry.Seq != position+1 {
			t.Fatalf("entry %d declares sequence %d, the fixture is not what the test needs", position+1, entry.Seq)
		}
	}

	verification, err := handle.VerifyAudit()
	if err != nil {
		t.Fatalf("VerifyAudit returned %v", err)
	}
	if verification.Valid {
		t.Fatalf("a spliced audit log verified as a valid chain: %+v", verification)
	}
	if verification.BrokenAt != 2 {
		t.Fatalf("verification is %+v, want the second entry reported as broken", verification)
	}
	if verification.Head != "" {
		t.Fatalf("a broken chain reported head %s", verification.Head)
	}

	intact, err := Open(filepath.Join(root, "intact"))
	if err != nil {
		t.Fatalf("Open returned %v", err)
	}
	if err := os.WriteFile(intact.Path(AuditFile), []byte(strings.Join(original, "\n")+"\n"), 0o600); err != nil {
		t.Fatalf("writing audit log: %v", err)
	}
	good, err := intact.VerifyAudit()
	if err != nil {
		t.Fatalf("VerifyAudit returned %v", err)
	}
	if !good.Valid || good.EntryCount != 3 || good.Head == "" {
		t.Fatalf("an untouched chain verified as %+v", good)
	}
}

// auditLines records one audit entry per detail in a fresh store and returns the
// raw log lines.
func auditLines(t *testing.T, dir string, details ...string) []string {
	t.Helper()
	handle, err := Open(dir)
	if err != nil {
		t.Fatalf("Open returned %v", err)
	}
	for _, detail := range details {
		if _, err := handle.Record("import", LedgerFile, detail, DigestString(detail)); err != nil {
			t.Fatalf("Record returned %v", err)
		}
	}
	raw, err := os.ReadFile(handle.Path(AuditFile))
	if err != nil {
		t.Fatalf("reading audit log: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(raw)), "\n")
	if len(lines) != len(details) {
		t.Fatalf("the audit log holds %d lines, want %d", len(lines), len(details))
	}
	return lines
}
