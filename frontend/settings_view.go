// Settings window adapted from github.com/egoist/kero (MIT).
package frontend

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/plugins/glass"
	"github.com/egoist/mygo/ui"

	"godgo/backend/i18n"
	"godgo/backend/settings"
)

type settingsPage struct {
	id    string
	title string
}

var settingsPages = []settingsPage{
	{"general", "page.general"},
	{"appearance", "page.appearance"},
	{"updates", "page.updates"},
	{"workflow", "page.workflow"},
	{"notifications", "page.notifications"},
	{"keyboard", "page.keyboard"},
	{"advanced", "page.advanced"},
}

func (a *App) openSettings(page string) {
	if page != "" {
		a.ss.page = page
	}
	if a.ss.page == "" {
		a.ss.page = "general"
	}
	if a.live {
		a.ss.loginOpen = mygo.App.OpenAtLogin()
		a.ss.updates.enabled = mygo.Updater.Enabled()
		a.ss.updates.lastChecked = lastCheck()
	}
	if a.setWin != nil {
		a.setWin.Show()
		a.setWin.Focus()
		return
	}
	if !a.live {
		return
	}
	a.setWin = mygo.NewWindow(framed(mygo.WindowOptions{
		Title: a.t("settings.title"), Width: 860, Height: 640, MinWidth: 720, MinHeight: 480,
		StateKey: "settings",
		Vibrancy: vibrancy(a.set),
		Content:  ui.View(a.settingsView),
	}))
	a.setWin.OnClosed(func() {
		a.setWin = nil
		a.ss.query = ""
	})
}

func vibrancy(s settings.Settings) mygo.Vibrancy {
	if s.Vibrancy {
		return mygo.VibrancySidebar
	}
	return mygo.VibrancyNone
}

func (a *App) settingsView(c *ui.Context) {
	pal := paletteFor(c)
	c.SetTheme(a.theme(c, pal))
	if a.ss.problems != nil && a.toast == "" && len(a.ss.problems) > 0 {
		// Shown once, from the main window's toast path the next time it draws.
	}
	found := a.searchSettings()
	c.OnShortcut(ui.Cmd, ui.KeyF, func() { a.ss.searchRef.Focus(c) })
	row := ui.Row(c).Fill().AlignItems(ui.Stretch)
	backdrop(c, pal, row, a.set)
	row.Children(func() {
		top := float32(topBar)
		if ownControls {
			top = barMid - (controlH+2)/2
		}
		ui.Column(c).Width(212).Shrink(0).Padding(top, 10, 12, 10).Gap(2).DragWindow().Children(func() {
			a.settingsSearch(c, pal)
			for _, p := range settingsPages {
				if found == nil || found[p.id] != nil {
					a.settingsNav(c, pal, p)
				}
			}
			if found != nil && len(found) == 0 {
				ui.Text(c, a.t("settings.noResults")).FontSize(13).TextColor(pal.textFaint).Padding(8, 10)
			}
			ui.Spacer(c)
			link := ui.ButtonBase(c.Key("open-config")).Padding(6, 10).Gap(6).Radius(7).AlignItems(ui.Center).Label(a.t("settings.openConfig"))
			fg := pal.textMuted
			if link.Hovered() {
				fg = pal.text
				link.Background(pal.rowHover)
			}
			link.TextColor(fg).Children(func() {
				ui.Icon(c, glyph("f7-chevron-left-slash-chevron-right")).Size(13, 13)
				ui.Text(c, a.t("settings.openConfig")).FontSize(12.5).SingleLine()
			})
			if link.Clicked() && a.live {
				_ = openText(filepath.Join(a.dir, "settings.json"))
			}
		})
		ui.Column(c).Grow(1).MinWidth(0).Padding(cardGap, cardGap, cardGap, 0).Children(func() {
			card := raised(ui.Column(c).Grow(1).MinHeight(0).Radius(14).Background(pal.card), pal).Clip()
			card.Children(func() {
				page := settingsPages[0]
				for _, p := range settingsPages {
					if p.id == a.ss.page {
						page = p
					}
				}
				if a.ss.scrolls == nil {
					a.ss.scrolls = map[string]*ui.ScrollState{}
				}
				sc := a.ss.scrolls[page.id]
				if sc == nil {
					sc = &ui.ScrollState{}
					a.ss.scrolls[page.id] = sc
				}
				box := ui.Box(c.Key("page")).Grow(1).MinHeight(0)
				_, py, over := box.PointerPosition()
				edgeOn := sc.Y > 0.5 || over && py < headerH
				box.Children(func() {
					ui.Scroll(c.Key("page-"+page.id)).Fill().Padding(headerH+16, 24, 28, 24).ScrollbarInsets(headerH, 0, 0).TrackScroll(sc).Children(func() {
						ui.Column(c).Gap(22).MaxWidth(620).Children(func() {
							switch page.id {
							case "general":
								a.generalPage(c, pal)
							case "appearance":
								a.appearancePage(c, pal)
							case "updates":
								a.updatesPage(c, pal)
							case "workflow":
								a.workflowPage(c, pal)
							case "notifications":
								a.notificationsPage(c, pal)
							case "keyboard":
								a.keyboardPage(c, pal)
							case "advanced":
								a.advancedPage(c, pal)
							}
						})
					})
					edge := ui.Box(c.Key("edge")).Absolute().Top(0).Left(0).Right(0).Height(headerH).PassThrough().
						Material(glass.ScrollEdge{Hard: true, Background: pal.card})
					edge.Opacity(edge.Animate("on", b2f(edgeOn), 160*time.Millisecond))
					ui.Row(c).Absolute().Top(0).Left(0).Right(0).Height(headerH).AlignItems(ui.Center).Padding(0, 10, 0, 24).DragWindow().Children(func() {
						ui.Text(c, a.t(page.title)).FontSize(15).FontWeight(650).TextColor(pal.text)
						ui.Spacer(c)
						windowControls(c, pal, a.setWin, a.t)
					})
				})
			})
		})
	})
	if a.changed() {
		a.apply(a.set)
	}
}

type searchEntry struct {
	page, label, words string
}

func (a *App) searchEntries() []searchEntry {
	t := a.t
	return []searchEntry{
		{"general", t("set.language"), "language locale 语言"},
		{"general", t("set.openAtLogin"), "login startup"},
		{"general", t("set.confirmQuit"), "quit exit"},
		{"general", t("set.restore"), "restore session"},
		{"appearance", t("set.appearance"), "theme dark light"},
		{"appearance", t("set.vibrancy"), "translucent material"},
		{"appearance", t("set.background"), "background wash"},
		{"appearance", t("set.font"), "font inter noto"},
		{"appearance", t("set.scale"), "size scale text"},
		{"updates", t("set.version"), "version"},
		{"updates", t("set.autoCheck"), "update check"},
		{"updates", t("set.autoInstall"), "install download"},
		{"workflow", t("set.runOnSave"), "test workflow"},
		{"notifications", t("set.notifyUpdate"), "notify"},
		{"notifications", t("set.notifyDone"), "done"},
		{"notifications", t("set.notifyFail"), "fail"},
		{"notifications", t("set.notifySound"), "sound"},
		{"keyboard", t("shortcut.sidebar"), "shortcut keyboard"},
		{"keyboard", t("shortcut.settings"), "shortcut"},
		{"keyboard", t("shortcut.palette"), "palette"},
		{"keyboard", t("shortcut.updates"), "update"},
		{"advanced", t("set.reset"), "reset default"},
		{"advanced", t("set.configPath"), "config path"},
	}
}

func (a *App) searchSettings() map[string][]string {
	terms := strings.Fields(strings.ToLower(a.ss.query))
	a.searchHits, a.searchReveal = nil, ""
	if len(terms) == 0 {
		a.ss.revealed = ""
		return nil
	}
	found := map[string][]string{}
	for _, p := range settingsPages {
		if matches(a.t(p.title), terms) {
			found[p.id] = []string{}
		}
	}
	for _, e := range a.searchEntries() {
		if matches(e.label+" "+e.words, terms) {
			found[e.page] = append(found[e.page], e.label)
		}
	}
	if found[a.ss.page] == nil {
		for _, p := range settingsPages {
			if found[p.id] != nil {
				a.ss.page = p.id
				break
			}
		}
	}
	a.searchHits = map[string]bool{}
	for _, l := range found[a.ss.page] {
		a.searchHits[l] = true
	}
	if at := a.ss.query + "\x00" + a.ss.page; len(found[a.ss.page]) > 0 && a.ss.revealed != at {
		a.ss.revealed, a.searchReveal = at, found[a.ss.page][0]
	}
	return found
}

func matches(s string, terms []string) bool {
	s = strings.ToLower(s)
	for _, term := range terms {
		if !strings.Contains(s, term) {
			return false
		}
	}
	return true
}

func (a *App) settingsSearch(c *ui.Context, pal *palette) {
	var in ui.Element
	frame := ui.Row(c.Key("search")).Cursor(ui.CursorText).Margin(0, 0, 10, 0)
	label := a.t("settings.search")
	frame.Children(func() {
		ui.Icon(c, searchIcon()).Size(13, 13).TextColor(pal.textFaint).Shrink(0)
		in = ui.TextInputBase(c.Key("in"), &a.ss.query).Bind(&a.ss.searchRef).Grow(1).MinWidth(0).FontSize(13).
			Placeholder(label).Label(label)
		in.OnShortcut(0, ui.KeyEscape, func() { a.ss.query = "" })
		if a.ss.query != "" {
			x := ui.ButtonBase(c.Key("clear")).Size(16, 16).Radius(8).Center().Shrink(0).Label(a.t("settings.clear"))
			x.Children(func() { ui.Icon(c, closeIcon()).Size(8, 8) })
			if x.Clicked() {
				a.ss.query = ""
				a.ss.searchRef.Focus(c)
			}
		}
	})
	field(frame, in.Focused(), pal)
}

func (a *App) settingsNav(c *ui.Context, pal *palette, p settingsPage) {
	sel := a.ss.page == p.id
	title := a.t(p.title)
	item := ui.Row(c.Key(p.id)).Height(32).Padding(0, 10).Gap(10).AlignItems(ui.Center).Radius(8).Label(title)
	if pal.dark {
		item.Border(1, ui.Transparent)
	}
	bg := ui.Transparent
	switch {
	case sel:
		bg = pal.rowSelected
		chosen(item, pal)
	case item.Hovered():
		bg = pal.rowHover
	}
	item.Background(item.AnimateColor("bg", bg, ui.Animation{Duration: 110 * time.Millisecond}))
	if item.Pressed() && !sel {
		a.ss.page = p.id
	}
	item.Children(func() {
		ic := ui.Icon(c, pageIcon(p.id)).Size(15, 15).TextColor(pal.textMuted)
		if sel {
			ic.TextColor(pal.accent)
		}
		ui.Text(c, title).FontSize(13.5).TextColor(pal.text).SingleLine()
	})
}

func (a *App) settingsGroup(c *ui.Context, pal *palette, title, note string, fn func()) {
	ui.Column(c).Gap(8).Children(func() {
		if title != "" {
			ui.Text(c, title).FontSize(12).FontWeight(600).TextColor(pal.textMuted).Padding(0, 4)
		}
		if fn != nil {
			a.groupRows = 0
			ui.Column(c).Radius(12).Background(pal.chip).Padding(4, 16).Children(fn)
		}
		if note != "" {
			ui.Text(c, note).FontSize(12).TextColor(pal.textFaint).LineHeight(1.4).Padding(0, 4)
		}
	})
}

func (a *App) settingRow(c *ui.Context, pal *palette, label, detail string, control func()) {
	if a.groupRows++; a.groupRows > 1 {
		ui.Box(c).Height(1).Shrink(0).Background(pal.text.Alpha(0.07))
	}
	row := ui.Row(c.Key("row-"+label)).Padding(11, 0).Gap(16).AlignItems(ui.Center)
	if a.searchHits[label] {
		row.Background(pal.accent.Alpha(0.13)).Radius(8)
		if label == a.searchReveal {
			row.ScrollIntoView()
		}
	}
	row.Children(func() {
		ui.Column(c).Grow(1).MinWidth(0).Gap(3).Children(func() {
			ui.Text(c, label).FontSize(13.5).TextColor(pal.text)
			if detail != "" {
				ui.Text(c, detail).FontSize(12).TextColor(pal.textMuted).LineHeight(1.35)
			}
		})
		ui.Row(c).Shrink(0).Gap(8).AlignItems(ui.Center).Children(control)
	})
}

func (a *App) versionLine() string {
	v := "dev"
	if a.live {
		if got := mygo.App.Version(); got != "" {
			v = got
		}
	}
	line := "GodGo " + v
	if !a.ss.updates.enabled {
		line += " · " + a.t("set.version.dev")
		return line
	}
	ago := i18n.Ago(a.displayLang(), a.ss.updates.lastChecked, time.Now())
	return line + " · " + ago
}

func (a *App) displayLang() string {
	if a.set.Language != "" && a.set.Language != "system" {
		return a.set.Language
	}
	if a.locale != "" {
		return a.locale
	}
	return i18n.SystemTag()
}
