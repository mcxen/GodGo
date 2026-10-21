package frontend

import (
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

func (a *App) extendDialog(c *ui.Context, pal *palette) {
	if !a.extendOpen {
		return
	}
	ui.DialogBase(c.Key("extend"), &a.extendOpen, func(backdrop, panel ui.Element) {
		backdrop.Background(pal.overlay)
		floating(panel.Width(460).Radius(14).Background(pal.popover).Padding(22).Gap(12), pal).Children(func() {
			ui.Text(c, a.t("extend.title")).FontSize(17).FontWeight(650).TextColor(pal.text)
			ui.Text(c, a.t("extend.detail")).FontSize(13).TextColor(pal.textMuted).LineHeight(1.45)
			for _, key := range []string{"extend.modules", "extend.settings", "extend.workflow"} {
				ui.Text(c, a.t(key)).FontSize(13).TextColor(pal.text).LineHeight(1.4)
			}
			ui.Row(c).Justify(ui.End).Margin(8, 0, 0, 0).Children(func() {
				if button(c, "extend-ok", a.t("extend.close"), "", true, pal) {
					a.extendOpen = false
				}
			})
		})
	})
}

func (a *App) quitDialog(c *ui.Context, pal *palette) {
	if !a.confirmQuit {
		return
	}
	ui.DialogBase(c.Key("quit"), &a.confirmQuit, func(backdrop, panel ui.Element) {
		backdrop.Background(pal.overlay)
		floating(panel.Width(420).Radius(14).Background(pal.popover).Padding(22).Gap(12), pal).Children(func() {
			ui.Text(c, a.t("quit.title")).FontSize(17).FontWeight(650).TextColor(pal.text)
			ui.Text(c, a.t("quit.detail")).FontSize(13).TextColor(pal.textMuted).LineHeight(1.45)
			ui.Row(c).Justify(ui.End).Gap(8).Margin(8, 0, 0, 0).Children(func() {
				if button(c, "quit-cancel", a.t("quit.cancel"), "", false, pal) {
					a.confirmQuit = false
				}
				if button(c, "quit-ok", a.t("quit.confirm"), "", true, pal) {
					a.confirmQuit = false
					a.quitReady = true
					a.save()
					if a.live {
						mygo.App.Quit()
					}
				}
			})
		})
	})
}
