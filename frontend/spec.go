package frontend

import (
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/markdown"
)

// specBody renders resources/docs/ui-spec.md. That file is the spec;
// this view does not keep a second copy of it.
func (a *App) specBody(c *ui.Context, pal *palette) {
	if a.spec == "" {
		ui.Text(c, a.t("spec.missing")).FontSize(13).TextColor(pal.attention).LineHeight(1.4)
		return
	}
	markdown.View(c.Key("spec-doc"), a.spec, nil)
}
