package frontend

import (
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/markdown"
)

// specBody renders the two spec files. The text stays in resources/docs;
// this view does not keep a second copy.
func (a *App) specBody(c *ui.Context, pal *palette) {
	if a.specWhich == "" {
		a.specWhich = "ui"
	}
	segmented(c, "spec-which", &a.specWhich, []string{"ui", "reuse"}, []string{a.t("spec.tab.ui"), a.t("spec.tab.reuse")})
	src := a.spec
	if a.specWhich == "reuse" {
		src = a.specReuse
	}
	if src == "" {
		ui.Text(c, a.t("spec.missing")).FontSize(13).TextColor(pal.attention).LineHeight(1.4)
		return
	}
	markdown.View(c.Key("spec-doc-"+a.specWhich), src, nil)
}
