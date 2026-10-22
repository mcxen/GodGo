package frontend

import (
	"path/filepath"
	"reflect"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"godgo/backend/i18n"
	"godgo/backend/modules"
	"godgo/backend/settings"
	"godgo/backend/workflow"
)

// App is the shell: one main window, one settings window, and the
// configuration both of them edit.
type App struct {
	dir     string
	root    string
	set     settings.Settings
	applied settings.Settings
	live    bool

	win      *mygo.Window
	setWin   *mygo.Window
	chrome   chrome
	sideDrag float32

	selected    string
	ss          settingsState
	pal         paletteState
	extendOpen  bool
	confirmQuit bool
	quitReady   bool

	toast        string
	runs         map[string]workflow.Result
	running      map[string]bool
	searchHits   map[string]bool
	searchReveal string
	groupRows    int

	catalogs  i18n.Catalogs
	tasks     workflow.Catalog
	spec      string
	specReuse string
	specWhich string
	locale    string
	notifyOK  bool

	onCheckUpdates func()
	onLogin        func(bool) error
	onPushUpdates  func()
}

type settingsState struct {
	page      string
	query     string
	searchRef ui.Handle
	revealed  string
	draft     map[string]*string
	loginOpen bool
	notice    string
	updates   updatePrefs
	scrolls   map[string]*ui.ScrollState
	problems  []string
}

type updatePrefs struct {
	enabled     bool
	lastChecked time.Time
}

type paletteState struct {
	open   bool
	query  string
	index  int
	scroll ui.ScrollState
}

func newApp(dir string) *App {
	catalogs, _ := i18n.LoadDir(filepath.Dir(resourceFile("locales/zh-Hans.json")))
	set, probs, _ := settings.Load(dir)
	tasks, _ := workflow.LoadFile(resourceFile("workflow/tasks.json"))
	specText, _ := readResource("docs/ui-spec.md")
	reuseText, _ := readResource("docs/reuse.md")
	if _, ok := modules.ByID(set.LastModule); !ok || !set.RestorePage {
		set.LastModule = "overview"
	}
	a := &App{
		dir:       dir,
		root:      moduleRoot(),
		set:       set,
		catalogs:  catalogs,
		tasks:     tasks,
		spec:      string(specText),
		specReuse: string(reuseText),
		selected:  set.LastModule,
		runs:      map[string]workflow.Result{},
		running:   map[string]bool{},
		ss:        settingsState{page: "general", problems: probs},
	}
	a.applied = a.set.Clone()
	return a
}

func (a *App) t(key string) string {
	lang := a.set.Language
	if lang == "" || lang == "system" {
		if a.locale != "" {
			lang = a.locale
		}
	}
	return a.catalogs.Text(lang, key)
}

func (a *App) langTag() string {
	if a.set.Language == "system" {
		return ""
	}
	return a.set.Language
}

func (a *App) save() {
	_ = settings.Save(a.dir, a.set)
}

func (a *App) apply(s settings.Settings) {
	if a.live {
		mygo.Theme.SetSource(mygo.ThemeSource(s.Appearance))
		v := mygo.VibrancyNone
		if s.Vibrancy {
			v = mygo.VibrancySidebar
		}
		if a.win != nil {
			a.win.SetVibrancy(v)
		}
		if a.setWin != nil {
			a.setWin.SetVibrancy(v)
			a.setWin.SetTitle(a.t("settings.title"))
		}
		mygo.App.SetLocale(a.langTag())
		mygo.App.SetMenu(a.menu())
	}
	prev := a.applied
	a.applied = s.Clone()
	a.set = s
	a.save()
	// Sidebar width and the open module are remembered on their own.
	// Only a change to the settings the user edits starts the verify task.
	if a.live && s.RunOnSave && !a.running["test"] && !reflect.DeepEqual(material(prev), material(s)) {
		if task, ok := a.tasks.ByID("test"); ok {
			a.startTask(task)
		}
	}
}

// material drops window geometry so dragging the sidebar does not look
// like an edited setting.
func material(s settings.Settings) settings.Settings {
	s.SidebarHidden = false
	s.SidebarWidth = 0
	s.Collapsed = nil
	s.LastModule = ""
	return s
}

func (a *App) changed() bool {
	return !reflect.DeepEqual(a.set, a.applied)
}

func (a *App) update(fn func()) {
	if a.win != nil {
		a.win.Update(fn)
		return
	}
	fn()
}

func (a *App) selectModule(id string) {
	if _, ok := modules.ByID(id); !ok {
		return
	}
	a.selected = id
	if a.set.RestorePage {
		a.set.LastModule = id
	}
}

func (a *App) toggleSidebar() {
	a.set.SidebarHidden = !a.set.SidebarHidden
}

func (a *App) toggleGroup(id string) {
	if a.set.Collapsed == nil {
		a.set.Collapsed = map[string]bool{}
	}
	a.set.Collapsed[id] = !a.set.Collapsed[id]
	if !a.set.Collapsed[id] {
		delete(a.set.Collapsed, id)
	}
}
