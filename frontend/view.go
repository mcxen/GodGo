// Sidebar layout adapted from github.com/egoist/kero (MIT).
package frontend

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"godgo/backend/modules"
)

func (a *App) view(c *ui.Context) {
	pal := paletteFor(c)
	c.SetTheme(a.theme(c, pal))
	if a.toast != "" {
		c.Toast(a.toast)
		a.toast = ""
	}
	a.shortcuts(c)

	bar := c.TitleBar()
	row := ui.Row(c).Fill().AlignItems(ui.Stretch)
	backdrop(c, pal, row, a.set)
	progress := row.AnimateTo("sidebar", b2f(!a.set.SidebarHidden), sidebarMotion)
	a.chrome = chromeFor(progress, a.set.SidebarWidth, bar)
	shown := a.chrome.shown
	row.Children(func() {
		w := a.chrome.sideW
		side := ui.Row(c).Width(w).Shrink(0).AlignItems(ui.Stretch).ClipX()
		side.Children(func() {
			ui.Column(c).Width(a.set.SidebarWidth).Shrink(0).Margin(0, 0, 0, w-a.set.SidebarWidth).
				Opacity(0.4 + 0.6*shown).Children(func() { a.sidebar(c, pal, bar) })
		})
		main := ui.Column(c).Grow(1).MinWidth(0).Padding(cardGap, cardGap, cardGap, a.chrome.mainPad).DragWindow()
		main.Children(func() {
			if shown > 0.01 {
				a.resizeHandle(c)
			}
			a.card(c, pal, bar, shown)
		})
	})
	a.paletteView(c, pal)
	a.extendDialog(c, pal)
	a.quitDialog(c, pal)
	if a.sideDrag == 0 && a.changed() {
		a.apply(a.set)
	}
}

func (a *App) resizeHandle(c *ui.Context) {
	h := ui.Box(c.Key("resize")).Absolute().Left(-4).Top(0).Bottom(0).Width(9).Cursor(ui.CursorResizeColumn)
	if dx, _, ok := h.Dragged(); ok {
		if a.sideDrag == 0 {
			a.sideDrag = a.set.SidebarWidth
		}
		a.sideDrag += dx
		a.set.SidebarWidth = min(max(a.sideDrag, minSidebar), maxSidebar)
	} else if a.sideDrag != 0 && !h.Pressed() {
		a.sideDrag = 0
		a.save()
		a.applied.SidebarWidth = a.set.SidebarWidth
	}
	if h.DoubleClicked() {
		a.set.SidebarWidth = defaultSidebar
		a.sideDrag = 0
		a.save()
		a.applied.SidebarWidth = a.set.SidebarWidth
	}
}

func (a *App) sidebar(c *ui.Context, pal *palette, bar ui.TitleBar) {
	ui.Row(c).Height(topBar).Shrink(0).AlignItems(ui.Center).Padding(0, 12, 0, lead(bar)).Gap(2).DragWindow().Children(func() {
		if !ownControls {
			ui.Spacer(c)
		}
		iconButton(c, "hide-sidebar", sidebarIcon(), a.t("sidebar.hide"), pal, a.toggleSidebar)
		iconButton(c, "extend", plusIcon(), a.t("sidebar.extend"), pal, func() { a.extendOpen = true })
	})
	ui.Scroll(c.Key("modules")).Grow(1).Padding(0, 10, 16, 10).Children(func() {
		for gi, g := range modules.Groups() {
			a.groupView(c, pal, g, gi)
		}
	})
	a.sidebarFooter(c, pal)
}

func (a *App) groupView(c *ui.Context, pal *palette, g modules.Group, gi int) {
	items := modules.In(g.ID)
	collapsed := a.set.Collapsed[g.ID]
	top := float32(18)
	if gi == 0 {
		top = 2
	}
	title := a.t(g.Title)
	hdr := ui.Row(c.Key("group-"+g.ID)).Height(28).Margin(top, 0, 3, 0).Padding(0, 6, 0, 10).Gap(6).AlignItems(ui.Center).Radius(6)
	hovered := hdr.Hovered()
	if hdr.Clicked() {
		a.toggleGroup(g.ID)
	}
	hdr.Children(func() {
		ui.Text(c, title).FontSize(12.5).FontWeight(600).TextColor(pal.textMuted).SingleLine().Ellipsis("…").MinWidth(0).Shrink(1)
		ui.Text(c, fmt.Sprintf("%d", len(items))).FontSize(12.5).TextColor(pal.textFaint)
		chev := ui.Icon(c, chevronIcon()).Size(11, 11).TextColor(pal.textFaint)
		turn := chev.AnimateTo("turn", 90*(1-b2f(collapsed)), ui.Snappy)
		chev.Rotate(turn).Opacity(chev.Animate("shown", b2f(hovered || collapsed), 120*time.Millisecond))
		ui.Spacer(c)
	})
	if collapsed {
		for _, m := range items {
			if m.ID == a.selected {
				a.moduleRow(c, pal, m)
			}
		}
		return
	}
	for _, m := range items {
		a.moduleRow(c, pal, m)
	}
}

func (a *App) moduleRow(c *ui.Context, pal *palette, m modules.Module) {
	selected := m.ID == a.selected
	title := a.t(m.Title)
	row := ui.Row(c.Key(m.ID)).Height(rowH).Padding(0, 10, 0, 12).Gap(10).AlignItems(ui.Center).Radius(8).Margin(0, 0, 1, 0).Label(title)
	if pal.dark {
		row.Border(1, ui.Transparent)
	}
	bg := ui.Transparent
	switch {
	case selected:
		bg = pal.rowSelected
		chosen(row, pal)
	case row.Hovered():
		bg = pal.rowHover
	}
	row.Background(row.AnimateColor("bg", bg, ui.Animation{Duration: 110 * time.Millisecond}))
	if (row.Pressed() || row.Clicked()) && !selected {
		a.selectModule(m.ID)
	}
	row.Tooltip(a.t(m.Summary))
	row.Children(func() {
		ui.Text(c, title).FontSize(13.5).SingleLine().Ellipsis("…").Grow(1).MinWidth(0).TextColor(pal.text)
		ui.Box(c).Size(6, 6).Radius(3).Background(pal.textFaint)
	})
}

func (a *App) sidebarFooter(c *ui.Context, pal *palette) {
	ui.Row(c).Height(36).Shrink(0).Padding(0, 12, 4, 12).Gap(2).AlignItems(ui.Center).DragWindow().Children(func() {
		iconButton(c, "reveal", folderIcon(), a.t("sidebar.reveal"), pal, func() {
			if a.root != "" && a.live {
				_ = openPath(a.root)
			}
		})
		iconButton(c, "palette", searchIcon(), a.t("sidebar.palette"), pal, func() { a.pal.open = true })
		iconButton(c, "settings", gearIcon(), a.t("sidebar.settings"), pal, func() { a.openSettings("") })
		ui.Spacer(c)
		ui.Row(c).Gap(7).Padding(0, 6).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Size(6, 6).Radius(3).Background(pal.done)
			ui.Text(c, a.t("sidebar.footer")).FontSize(12).TextColor(pal.textFaint).SingleLine()
		})
	})
}
