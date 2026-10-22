package frontend

import (
	"path/filepath"
	"strconv"

	"github.com/egoist/mygo/ui"

	"godgo/backend/settings"
)

func (a *App) generalPage(c *ui.Context, pal *palette) {
	s := &a.set
	a.settingsGroup(c, pal, "", "", func() {
		a.settingRow(c, pal, a.t("set.language"), "", func() {
			ids := append([]string{"system"}, settings.Languages...)
			popup(c, "language", &s.Language, ids, func(tag string) string {
				if tag == "system" {
					return a.t("set.language.system")
				}
				if tag == "zh-Hans" {
					return "简体中文"
				}
				return "English"
			}, 180, pal)
		})
		a.settingRow(c, pal, a.t("set.openAtLogin"), a.t("set.openAtLogin.detail"), func() {
			toggle(c.Key("login"), &a.ss.loginOpen, pal).OnChange(func() {
				if a.onLogin == nil {
					return
				}
				if err := a.onLogin(a.ss.loginOpen); err != nil {
					a.ss.loginOpen = !a.ss.loginOpen
					a.ss.notice = err.Error()
				}
			})
		})
		a.settingRow(c, pal, a.t("set.confirmQuit"), a.t("set.confirmQuit.detail"), func() {
			toggle(c.Key("quit"), &s.ConfirmQuit, pal)
		})
		a.settingRow(c, pal, a.t("set.restore"), a.t("set.restore.detail"), func() {
			toggle(c.Key("restore"), &s.RestorePage, pal)
		})
	})
}

func (a *App) appearancePage(c *ui.Context, pal *palette) {
	s := &a.set
	a.settingsGroup(c, pal, "", "", func() {
		a.settingRow(c, pal, a.t("set.appearance"), "", func() {
			segmented(c, "appearance", &s.Appearance, []string{"system", "light", "dark"}, []string{a.t("set.theme.system"), a.t("set.theme.light"), a.t("set.theme.dark")})
		})
		a.settingRow(c, pal, a.t("set.vibrancy"), a.t("set.vibrancy.detail"), func() {
			toggle(c.Key("vibrancy"), &s.Vibrancy, pal)
		})
		a.settingRow(c, pal, a.t("set.font"), a.t("set.font.detail"), func() {
			segmented(c, "font", &s.FontFamily, []string{"", "system"}, []string{a.t("set.font.bundled"), a.t("set.font.system")})
		})
		a.settingRow(c, pal, a.t("set.scale"), a.t("set.scale.detail"), func() {
			popup(c, "scale", &s.FontScale, settings.FontScales, func(v float64) string {
				return strconv.FormatFloat(v, 'f', -1, 64) + "×"
			}, 90, pal)
		})
	})
	ui.Column(c).Gap(8).Children(func() {
		ui.Text(c, a.t("set.background")).FontSize(12).FontWeight(600).TextColor(pal.textMuted).Padding(0, 4)
		ui.Row(c).Gap(12).Wrap().Children(func() {
			for _, id := range settings.Backgrounds {
				a.backgroundTile(c, pal, id)
			}
		})
	})
}

func (a *App) backgroundTile(c *ui.Context, pal *palette, id string) {
	name := a.t("bg." + id)
	sel := a.set.Background == id
	tile := ui.ButtonBase(c.Key("bg-" + id)).Column().Width(104).Gap(8).AlignItems(ui.Stretch).Label(name)
	if tile.Clicked() {
		a.set.Background = id
	}
	tile.Children(func() {
		sw := raised(ui.Box(c).Size(104, 68).Radius(12), pal).Clip()
		base, glow := washes(id, pal.dark)
		sw.LinearGradient(base)
		sw.Children(func() {
			ui.Box(c).Absolute().Top(0).Left(0).Right(0).Bottom(0).LinearGradient(glow)
			ui.Box(c).Absolute().Top(12).Left(26).Right(8).Bottom(-6).Radius(7).Background(pal.card)
			if sel {
				ui.Box(c).Absolute().Top(6).Right(6).Size(18, 18).Radius(9).Center().Background(pal.accent).Children(func() {
					ui.Icon(c, checkIcon()).Size(11, 11).TextColor(ui.Hex("#ffffff"))
				})
			}
		})
		lbl := ui.Text(c, name).FontSize(12.5).TextColor(pal.textMuted).TextAlign(ui.Center)
		if sel {
			lbl.TextColor(pal.text).FontWeight(600)
		}
	})
}

func (a *App) updatesPage(c *ui.Context, pal *palette) {
	s := &a.set
	a.settingsGroup(c, pal, "", a.t("set.updatesNote"), func() {
		a.settingRow(c, pal, a.t("set.version"), a.versionLine(), func() {
			if button(c, "check-updates", a.t("set.check"), "", false, pal) && a.onCheckUpdates != nil {
				a.onCheckUpdates()
			}
		})
		a.settingRow(c, pal, a.t("set.autoCheck"), a.t("set.autoCheck.detail"), func() {
			toggle(c.Key("auto-check"), &s.AutoCheck, pal).OnChange(func() {
				if a.onPushUpdates != nil {
					a.onPushUpdates()
				}
			})
		})
		a.settingRow(c, pal, a.t("set.autoInstall"), a.t("set.autoInstall.detail"), func() {
			toggle(c.Key("auto-install"), &s.AutoInstall, pal).OnChange(func() {
				if a.onPushUpdates != nil {
					a.onPushUpdates()
				}
			})
		})
	})
}

func (a *App) workflowPage(c *ui.Context, pal *palette) {
	a.settingsGroup(c, pal, "", "", func() {
		a.settingRow(c, pal, a.t("set.runOnSave"), a.t("set.runOnSave.detail"), func() {
			toggle(c.Key("run-on-save"), &a.set.RunOnSave, pal)
		})
	})
}

func (a *App) notificationsPage(c *ui.Context, pal *palette) {
	note := ""
	if a.live && !a.notifyOK {
		note = a.t("set.notifyUnsupported")
	}
	s := &a.set
	a.settingsGroup(c, pal, "", note, func() {
		a.settingRow(c, pal, a.t("set.notifyUpdate"), a.t("set.notifyUpdate.detail"), func() {
			toggle(c.Key("nu"), &s.NotifyUpdate, pal)
		})
		a.settingRow(c, pal, a.t("set.notifyDone"), a.t("set.notifyDone.detail"), func() {
			toggle(c.Key("nd"), &s.NotifyDone, pal)
		})
		a.settingRow(c, pal, a.t("set.notifyFail"), a.t("set.notifyFail.detail"), func() {
			toggle(c.Key("nf"), &s.NotifyFail, pal)
		})
		a.settingRow(c, pal, a.t("set.notifySound"), a.t("set.notifySound.detail"), func() {
			toggle(c.Key("ns"), &s.NotifySound, pal)
		})
	})
}

func (a *App) keyboardPage(c *ui.Context, pal *palette) {
	if a.ss.draft == nil {
		a.ss.draft = map[string]*string{}
	}
	rows := []struct{ id, key string }{
		{settings.ShortcutSidebar, "shortcut.sidebar"},
		{settings.ShortcutSettings, "shortcut.settings"},
		{settings.ShortcutPalette, "shortcut.palette"},
		{settings.ShortcutUpdates, "shortcut.updates"},
	}
	a.settingsGroup(c, pal, "", a.t("set.shortcut.detail"), func() {
		for _, row := range rows {
			id := row.id
			if a.ss.draft[id] == nil {
				v := a.set.Shortcut(id)
				a.ss.draft[id] = &v
			}
			a.settingRow(c, pal, a.t(row.key), "", func() {
				in := textField(c, "sc-"+id, a.ss.draft[id], "CmdOrCtrl+K", 180, pal)
				in.OnSubmit(func() {
					v := *a.ss.draft[id]
					if !settings.ValidAccelerator(v) {
						a.ss.notice = a.t("set.shortcut.invalid")
						return
					}
					if a.set.Shortcuts == nil {
						a.set.Shortcuts = map[string]string{}
					}
					a.set.Shortcuts[id] = v
				})
			})
		}
	})
	ui.Row(c).Children(func() {
		if button(c, "reset-keys", a.t("set.shortcut.reset"), "", false, pal) {
			a.set.Shortcuts = nil
			a.ss.draft = map[string]*string{}
		}
	})
}

func (a *App) advancedPage(c *ui.Context, pal *palette) {
	a.settingsGroup(c, pal, "", "", func() {
		a.settingRow(c, pal, a.t("set.configPath"), a.dir, func() {})
		a.settingRow(c, pal, a.t("set.resourcePath"), resourceLabel(), func() {})
		root := a.root
		if root == "" {
			root = a.t("set.moduleRoot.none")
		}
		a.settingRow(c, pal, a.t("set.moduleRoot"), root, func() {})
		a.settingRow(c, pal, a.t("set.reset"), a.t("set.reset.detail"), func() {
			if button(c, "reset", a.t("set.reset.button"), "", false, pal) {
				keep := a.set
				next := settings.Default()
				next.SidebarHidden = keep.SidebarHidden
				next.SidebarWidth = keep.SidebarWidth
				next.Collapsed = keep.Collapsed
				next.LastModule = keep.LastModule
				a.set = next
				a.ss.draft = nil
			}
		})
	})
	if len(a.ss.problems) > 0 {
		ui.Text(c, joinProblems(a.ss.problems)).FontSize(12).TextColor(pal.attention).LineHeight(1.4)
	}
}

func resourceLabel() string {
	path := resourceFile("fonts/inter/OFL.txt")
	if path == "" {
		dirs := resourceDirs()
		if len(dirs) == 0 {
			return ""
		}
		return dirs[0]
	}
	// fonts/inter/OFL.txt sits three directories under the resources root.
	return filepath.Dir(filepath.Dir(filepath.Dir(path)))
}

func joinProblems(ps []string) string {
	out := ""
	for i, p := range ps {
		if i > 0 {
			out += "\n"
		}
		out += p
	}
	return out
}
