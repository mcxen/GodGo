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
	for _, h := range []string{"# GodGo 界面规范", "## 窗口边框", "## 窗口控制栏", "## 侧栏", "## 设置", "## 字体", "## 工作流", "## 文档"} {
		if !strings.Contains(string(spec), h) {
			t.Fatalf("spec missing %s", h)
		}
	}
	if !strings.Contains(string(spec), "TrafficLightPosition") || !strings.Contains(string(spec), "(18, 23)") {
		t.Fatal("spec missing the traffic-light placement")
	}
	if !strings.Contains(string(spec), "Frameless") {
		t.Fatal("spec missing the window-frame constraint")
	}
	agents, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(agents), "resources/docs/ui-spec.md") || !strings.Contains(string(agents), "resources/docs/reuse.md") {
		t.Fatal("AGENTS.md does not point at both specs")
	}
	if !strings.Contains(string(agents), ".grok/skills/godgo/SKILL.md") {
		t.Fatal("AGENTS.md does not point at the skill")
	}
	skill, err := os.ReadFile(filepath.Join(root, ".grok", "skills", "godgo", "SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{"name: godgo", "AGENTS.md", "resources/docs/reuse.md", "resources/docs/ui-spec.md", "mygo-maintenance"} {
		if !strings.Contains(string(skill), s) {
			t.Fatalf("skill missing %s", s)
		}
	}
	if strings.Contains(string(skill), "(18, 23)") {
		t.Fatal("skill copied a spec measurement")
	}
	reuse, err := os.ReadFile(filepath.Join(root, "resources", "docs", "reuse.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{"# 复用规范", "禁止过度设计", "禁止重复造轮子", "Onevra", "plugins/updater", "plugins/terminal"} {
		if !strings.Contains(string(reuse), s) {
			t.Fatalf("reuse spec missing %s", s)
		}
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
