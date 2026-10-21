package frontend

import (
	"strings"
	"unicode"

	"github.com/egoist/mygo/ui"

	"godgo/backend/modules"
)

type paletteItem struct {
	title  string
	detail string
	keys   string
	run    func()
	score  int
	kind   string
}

func (a *App) paletteItems() []paletteItem {
	var items []paletteItem
	for _, m := range modules.All() {
		id := m.ID
		items = append(items, paletteItem{
			title: a.t(m.Title), detail: a.t(m.Summary), kind: "module",
			run: func() { a.selectModule(id) },
		})
	}
	for _, p := range settingsPages {
		id := p.id
		items = append(items, paletteItem{
			title: a.t("palette.settings") + " · " + a.t(p.title), kind: "settings",
			run: func() { a.openSettings(id) },
		})
	}
	items = append(items,
		paletteItem{title: a.t("sidebar.hide"), keys: a.set.Shortcut("sidebar.toggle"), kind: "command", run: a.toggleSidebar},
		paletteItem{title: a.t("menu.checkUpdates"), keys: a.set.Shortcut("updates.check"), kind: "command", run: func() {
			if a.onCheckUpdates != nil {
				a.onCheckUpdates()
			}
		}},
	)
	return items
}

func fuzzy(query, s string) int {
	if query == "" {
		return 0
	}
	q := []rune(strings.ToLower(query))
	r := []rune(strings.ToLower(s))
	score, qi, run := 0, 0, 0
	for i := 0; i < len(r) && qi < len(q); i++ {
		if q[qi] == ' ' {
			qi++
			continue
		}
		if r[i] == q[qi] {
			bonus := 1
			if i == 0 || !unicode.IsLetter(r[i-1]) && !unicode.IsDigit(r[i-1]) {
				bonus += 6
			}
			run++
			score += bonus + run*2
			qi++
		} else {
			run = 0
		}
	}
	if qi < len(q) {
		return -1
	}
	return score*10 - len(r)
}

func (a *App) paletteView(c *ui.Context, pal *palette) {
	if !a.pal.open {
		return
	}
	var items []paletteItem
	for _, it := range a.paletteItems() {
		sc := fuzzy(a.pal.query, it.title+" "+it.detail)
		if t := fuzzy(a.pal.query, it.title); t >= 0 {
			sc = max(sc, t+50)
		}
		if sc < 0 {
			continue
		}
		it.score = sc
		items = append(items, it)
	}
	if a.pal.query != "" {
		// Keep the catalog order when scores tie: modules, then settings.
	}
	a.pal.index = min(max(a.pal.index, 0), max(len(items)-1, 0))
	var chosen func()
	ui.DialogBase(c.Key("palette"), &a.pal.open, func(backdrop, panel ui.Element) {
		backdrop.Background(pal.overlay).Justify(ui.Start).Padding(84, 0, 0, 0)
		floating(panel.Width(560).MaxHeight(440).Radius(14).Background(pal.popover), pal).Clip().Label(a.t("palette.placeholder"))
		ui.Row(c).Height(50).Shrink(0).Padding(0, 16).Gap(10).AlignItems(ui.Center).Children(func() {
			ui.Icon(c, searchIcon()).Size(16, 16).TextColor(pal.textMuted)
			in := ui.TextInputBase(c.Key("q"), &a.pal.query).Grow(1).FontSize(15).Placeholder(a.t("palette.placeholder")).AutoFocus()
			if in.Changed() {
				a.pal.index = 0
			}
			if c.Shortcut(0, ui.KeyDown) {
				a.pal.index = min(a.pal.index+1, max(len(items)-1, 0))
			}
			if c.Shortcut(0, ui.KeyUp) {
				a.pal.index = max(a.pal.index-1, 0)
			}
			if (in.Submitted() || c.Shortcut(0, ui.KeyEnter)) && len(items) > 0 {
				chosen = items[a.pal.index].run
			}
		})
		ui.Scroll(c.Key("results")).Grow(1).MinHeight(0).Padding(6).TrackScroll(&a.pal.scroll).Children(func() {
			if len(items) == 0 {
				ui.Text(c, a.t("palette.empty")).FontSize(13).TextColor(pal.textFaint).Padding(12, 12)
			}
			last := ""
			for i, it := range items {
				if a.pal.query == "" && it.kind != last {
					label := a.t("palette.modules")
					if it.kind == "settings" {
						label = a.t("palette.settings")
					} else if it.kind == "command" {
						label = a.t("palette.commands")
					}
					ui.Text(c, label).FontSize(11.5).FontWeight(600).TextColor(pal.textFaint).Padding(10, 10, 4, 10)
				}
				last = it.kind
				row := ui.Row(c.Key(i)).Height(36).Padding(0, 10).Gap(10).AlignItems(ui.Center).Radius(8)
				if i == a.pal.index {
					row.Background(pal.accent.Alpha(0.12))
					row.ScrollIntoView()
				} else if row.Hovered() {
					row.Background(pal.rowHover)
				}
				if row.Clicked() {
					chosen = it.run
				}
				row.Children(func() {
					ui.Text(c, it.title).FontSize(13.5).TextColor(pal.text).SingleLine().Ellipsis("…").Grow(1).MinWidth(0)
					if it.keys != "" {
						ui.Text(c, it.keys).FontSize(11).TextColor(pal.textMuted).Padding(1, 5).Radius(4).Background(pal.kbd).SingleLine()
					}
				})
			}
		})
	})
	if chosen != nil {
		a.pal.open = false
		a.pal.query = ""
		chosen()
	}
}
