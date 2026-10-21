package frontend

import (
	"os"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/plugins/updater"
	"github.com/egoist/mygo/plugins/updater/native"
	"github.com/egoist/mygo/ui"
)

// Run starts the shell. dir is where settings.json is kept; an empty dir
// uses MyGo's per-user data directory.
func Run() error {
	mygo.Use(native.Plugin)
	if !mygo.App.RequestSingleInstanceLock() {
		return nil
	}
	dir, err := mygo.App.Path(mygo.PathUserData)
	if err != nil {
		return err
	}
	loadFonts()
	a := newApp(dir)
	a.live = true
	a.onCheckUpdates = updater.CheckForUpdates
	a.onLogin = mygo.App.SetOpenAtLogin
	a.onPushUpdates = func() {
		updater.SetAutomaticChecks(a.set.AutoCheck)
		updater.SetAutomaticDownloads(a.set.AutoInstall)
	}
	mygo.Theme.SetSource(mygo.ThemeSource(a.set.Appearance))
	mygo.App.SetLocale(a.langTag())
	mygo.App.SetMenu(a.menu())
	mygo.App.OnLocaleChanged(func(tag string) {
		a.update(func() {
			if a.set.Language == "system" {
				a.locale = tag
			}
			mygo.App.SetMenu(a.menu())
			if a.setWin != nil {
				a.setWin.SetTitle(a.t("settings.title"))
			}
		})
	})
	mygo.App.OnSecondInstance(func([]string, string) {
		if a.win != nil {
			a.win.Restore()
			a.win.Focus()
		}
	})
	mygo.App.OnBeforeQuit(func(e *mygo.QuitEvent) {
		a.save()
		if a.set.ConfirmQuit && !a.quitReady {
			e.PreventDefault()
			a.update(func() { a.confirmQuit = true })
		}
	})
	off := updater.OnChange(func() {
		a.update(func() {
			a.ss.updates.lastChecked = updater.LastCheck()
			a.set.AutoCheck = updater.AutomaticChecks()
			a.set.AutoInstall = updater.AutomaticDownloads()
			a.applied.AutoCheck = a.set.AutoCheck
			a.applied.AutoInstall = a.set.AutoInstall
			a.save()
		})
	})
	defer off()
	mygo.App.WhenReady(func() {
		a.notifyOK = mygo.NotificationsSupported()
		a.locale = mygo.App.Locale()
		a.onPushUpdates()
		a.ss.updates.enabled = mygo.Updater.Enabled()
		a.ss.updates.lastChecked = updater.LastCheck()
		a.win = mygo.NewWindow(framed(mygo.WindowOptions{
			Title:     "GodGo",
			Width:     1100,
			Height:    720,
			MinWidth:  760,
			MinHeight: 480,
			StateKey:  "main",
			Vibrancy:  vibrancy(a.set),
			Content:   ui.View(a.view),
		}))
	})
	return mygo.App.Run()
}

func lastCheck() time.Time { return updater.LastCheck() }

func (a *App) menu() *mygo.Menu {
	return mygo.NewMenu([]*mygo.MenuItem{
		{Role: mygo.RoleAppMenu, Submenu: []*mygo.MenuItem{
			{Label: a.t("menu.settings"), Accelerator: a.set.Shortcut("settings.open"), Click: func(*mygo.MenuItem, *mygo.Window) { a.openSettings("") }},
			{Label: a.t("menu.palette"), Accelerator: a.set.Shortcut("palette.open"), Click: func(*mygo.MenuItem, *mygo.Window) {
				a.update(func() { a.pal.open = true })
			}},
			{Label: a.t("menu.checkUpdates"), Accelerator: a.set.Shortcut("updates.check"), Click: func(*mygo.MenuItem, *mygo.Window) {
				if a.onCheckUpdates != nil {
					a.onCheckUpdates()
				}
			}},
			mygo.Separator(),
			{Role: mygo.RoleQuit},
		}},
		{Role: mygo.RoleEditMenu},
		{Role: mygo.RoleWindowMenu},
	})
}

// UserDataDir is where a development run keeps settings when the app
// path is not ready yet. Run uses MyGo's own directory.
func UserDataDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return dir + "/GodGo", nil
}
