package frontend

import (
	"strings"

	"github.com/egoist/mygo/ui"

	"godgo/backend/modules"
	"godgo/backend/workflow"
)

func (a *App) card(c *ui.Context, pal *palette, bar ui.TitleBar, shown float32) {
	card := raised(ui.Column(c).Grow(1).MinHeight(0).Radius(14).Background(pal.card), pal).Clip()
	card.Children(func() {
		m, _ := modules.ByID(a.selected)
		title := a.t(m.Title)
		if title == m.Title {
			title = a.t("app.name")
		}
		left := a.chrome.headerPad
		ui.Row(c).Height(headerH).Shrink(0).AlignItems(ui.Center).Padding(0, 10, 0, left).DragWindow().Children(func() {
			if shown < 0.5 {
				ui.Row(c).Absolute().Left(a.chrome.buttonsLeft(bar)).Top(0).Bottom(0).AlignItems(ui.Center).Gap(2).Opacity(1 - 2*shown).Children(func() {
					iconButton(c, "show-sidebar", sidebarIcon(), a.t("sidebar.show"), pal, a.toggleSidebar)
					iconButton(c, "extend-2", plusIcon(), a.t("sidebar.extend"), pal, func() { a.extendOpen = true })
				})
			}
			ui.Text(c, title).FontSize(14).FontWeight(600).TextColor(pal.text)
			ui.Spacer(c)
			windowControls(c, pal, a.win, a.t)
		})
		ui.Scroll(c.Key("body")).Grow(1).MinHeight(0).Padding(8, 28, 28, 28).Children(func() {
			ui.Column(c).Gap(16).MaxWidth(640).Children(func() {
				ui.Text(c, a.t(m.Summary)).FontSize(13).TextColor(pal.textMuted)
				ui.Text(c, a.t(m.Body)).FontSize(14.5).TextColor(pal.text).LineHeight(1.5)
				switch m.ID {
				case "spec":
					a.specBody(c, pal)
				case "workflow":
					a.workflowBody(c, pal)
				case "fonts":
					a.fontsBody(c, pal)
				case "testdata":
					a.testdataBody(c, pal)
				case "overview":
					ui.Text(c, a.t("home.ready")).FontSize(13).TextColor(pal.textMuted).LineHeight(1.45)
				}
			})
		})
	})
}

func (a *App) workflowBody(c *ui.Context, pal *palette) {
	if a.root == "" {
		ui.Text(c, a.t("task.noRoot")).FontSize(13).TextColor(pal.attention)
	}
	for _, g := range a.tasks.Groups {
		ui.Text(c, a.t(g.TitleKey)).FontSize(12).FontWeight(600).TextColor(pal.textMuted)
		for _, task := range a.tasks.TasksIn(g.ID) {
			a.taskRow(c, pal, task)
		}
	}
}

func (a *App) taskRow(c *ui.Context, pal *palette, task workflow.Task) {
	ui.Column(c.Key("task-"+task.ID)).Padding(12, 14).Gap(6).Radius(12).Background(pal.chip).Children(func() {
		ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
			ui.Column(c).Grow(1).MinWidth(0).Gap(2).Children(func() {
				ui.Text(c, a.t(task.TitleKey)).FontSize(13.5).FontWeight(600).TextColor(pal.text)
				ui.Text(c, a.t(task.DetailKey)).FontSize(12).TextColor(pal.textMuted)
				ui.Text(c, strings.Join(task.Argv, " ")).Font("monospace").FontSize(12).TextColor(pal.textFaint).SingleLine().Ellipsis("…")
			})
			label := a.t("task.run")
			if a.running[task.ID] {
				label = a.t("task.running")
			}
			if button(c, "run-"+task.ID, label, "", true, pal) && a.live && !a.running[task.ID] {
				a.startTask(task)
			}
		})
		if res, ok := a.runs[task.ID]; ok {
			state := a.t("task.passed")
			col := pal.done
			if !res.OK {
				state, col = a.t("task.failed"), pal.danger
			}
			ui.Text(c, state).FontSize(12).FontWeight(600).TextColor(col)
			if res.Err != "" {
				ui.Text(c, res.Err).FontSize(12).TextColor(pal.danger)
			}
			if res.Output != "" {
				ui.Text(c, res.Output).Font("monospace").FontSize(11.5).TextColor(pal.textMuted).LineHeight(1.35)
			}
		}
	})
}

func (a *App) fontsBody(c *ui.Context, pal *palette) {
	ui.Column(c).Gap(8).Padding(12, 14).Radius(12).Background(pal.chip).Children(func() {
		for _, line := range []string{fontInter + "  ·  Regular, Medium, SemiBold, Bold", fontNoto + "  ·  Regular, Medium, Bold"} {
			ui.Text(c, line).FontSize(13).TextColor(pal.text)
		}
		sample := "GodGo 外壳 Ag"
		ui.Text(c, sample).Font(a.fontFamily()).FontSize(22).TextColor(pal.text).Margin(8, 0, 0, 0)
	})
}

func (a *App) testdataBody(c *ui.Context, pal *palette) {
	for _, line := range []string{
		"testdata/settings",
		"testdata/workflow",
		"testdata/i18n",
		"testdata/fonts",
		"testdata/ui",
	} {
		ui.Text(c, line).Font("monospace").FontSize(13).TextColor(pal.text)
	}
}
