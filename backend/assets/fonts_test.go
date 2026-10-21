package assets

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFontChecksums(t *testing.T) {
	root := repoRoot(t)
	b, err := os.ReadFile(filepath.Join(root, "testdata", "fonts", "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var entries []struct {
		File   string `json:"file"`
		SHA256 string `json:"sha256"`
	}
	if err := json.Unmarshal(b, &entries); err != nil {
		t.Fatal(err)
	}
	if len(entries) != 7 {
		t.Fatalf("manifest entries: %d", len(entries))
	}
	for _, e := range entries {
		data, err := os.ReadFile(filepath.Join(root, "resources", "fonts", e.File))
		if err != nil {
			t.Fatal(err)
		}
		sum := sha256.Sum256(data)
		if hex.EncodeToString(sum[:]) != e.SHA256 {
			t.Fatalf("%s checksum changed", e.File)
		}
	}
}

func TestDocsAgree(t *testing.T) {
	root := repoRoot(t)
	spec, err := os.ReadFile(filepath.Join(root, "resources", "docs", "ui-spec.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, h := range []string{"# GodGo 界面规范", "## 侧栏", "## 设置", "## 字体", "## 工作流", "## 文档"} {
		if !strings.Contains(string(spec), h) {
			t.Fatalf("spec missing %s", h)
		}
	}
	agents, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(agents), "resources/docs/ui-spec.md") {
		t.Fatal("AGENTS.md does not point at the spec")
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for range 8 {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.Fatal("go.mod not found")
	return ""
}
