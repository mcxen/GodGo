package frontend

import (
	"strings"

	"github.com/egoist/mygo/ui"
)

func (a *App) shortcuts(c *ui.Context) {
	a.bind(c, a.set.Shortcut("sidebar.toggle"), a.toggleSidebar)
	a.bind(c, a.set.Shortcut("settings.open"), func() { a.openSettings("") })
	a.bind(c, a.set.Shortcut("palette.open"), func() { a.pal.open = true; a.pal.query = "" })
	a.bind(c, a.set.Shortcut("updates.check"), func() {
		if a.onCheckUpdates != nil {
			a.onCheckUpdates()
		}
	})
}

func (a *App) bind(c *ui.Context, accel string, fn func()) {
	mods, key, ok := parseAccel(accel)
	if !ok {
		return
	}
	c.OnShortcut(mods, key, fn)
}

func parseAccel(s string) (ui.Modifiers, ui.Key, bool) {
	if s == "" {
		return 0, 0, false
	}
	parts := strings.Split(s, "+")
	if len(parts) < 2 {
		return 0, 0, false
	}
	var mods ui.Modifiers
	for _, p := range parts[:len(parts)-1] {
		switch p {
		case "CmdOrCtrl", "Cmd":
			mods |= ui.Cmd
		case "Ctrl":
			mods |= ui.Ctrl
		case "Alt":
			mods |= ui.Alt
		case "Shift":
			mods |= ui.Shift
		case "Super":
			mods |= ui.Super
		default:
			return 0, 0, false
		}
	}
	key, ok := keyByName(parts[len(parts)-1])
	return mods, key, ok
}

func keyByName(name string) (ui.Key, bool) {
	switch strings.ToLower(name) {
	case ",", "comma":
		return ui.KeyComma, true
	case "b":
		return ui.KeyB, true
	case "k":
		return ui.KeyK, true
	case "u":
		return ui.KeyU, true
	case "f":
		return ui.KeyF, true
	default:
		if len(name) == 1 && name[0] >= 'a' && name[0] <= 'z' {
			return ui.KeyA + ui.Key(name[0]-'a'), true
		}
		return 0, false
	}
}
