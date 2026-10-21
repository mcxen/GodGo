package frontend

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestSidebarModules(t *testing.T) {
	a := testApp(t)
	view := ui.NewTester(a.view, 1100, 700)
	root := moduleRoot()
	b, err := os.ReadFile(filepath.Join(root, "testdata", "ui", "sidebar.txt"))
	if err != nil {
		t.Fatal(err)
	}
	for line := range strings.SplitSeq(strings.TrimSpace(string(b)), "\n") {
		if !view.HasText(line) {
			t.Fatalf("missing %q\n%s", line, strings.Join(view.Texts(), "\n"))
		}
	}
	if err := view.Click("工作流"); err != nil {
		t.Fatal(err)
	}
	view.Frame()
	if !view.HasText("go test ./...") {
		t.Fatalf("workflow argv missing\n%s", strings.Join(view.Texts(), "\n"))
	}
}

func TestSettingsPages(t *testing.T) {
	a := testApp(t)
	view := ui.NewTester(a.settingsView, 860, 640)
	if err := view.Click("外观"); err != nil {
		t.Fatal(err)
	}
	view.Frame()
	if !view.HasText("跟随系统") || !view.HasText("随包字体") {
		t.Fatalf("appearance page\n%s", strings.Join(view.Texts(), "\n"))
	}
	if _, ok := view.Find("搜索设置"); !ok {
		t.Fatal("search field has no label")
	}
	if err := view.Click("搜索设置"); err != nil {
		t.Fatal(err)
	}
	view.Type("键盘")
	view.Frame()
	if !view.HasText("显示或隐藏边栏") {
		t.Fatalf("keyboard page\n%s", strings.Join(view.Texts(), "\n"))
	}
}

func TestSpecPage(t *testing.T) {
	a := testApp(t)
	if !strings.Contains(a.spec, "# GodGo 界面规范") {
		t.Fatal("spec was not loaded")
	}
	view := ui.NewTester(a.view, 1100, 700)
	if err := view.Click("规范"); err != nil {
		t.Fatal(err)
	}
	view.Frame()
	if !view.HasText("GodGo 界面规范") || !view.HasText("这份文件是界面契约") {
		t.Fatalf("spec page\n%s", strings.Join(view.Texts(), "\n"))
	}
}

func testApp(t *testing.T) *App {
	t.Helper()
	a := newApp(t.TempDir())
	a.set.Language = "zh-Hans"
	if a.live {
		t.Fatal("tests must not talk to the updater")
	}
	return a
}
