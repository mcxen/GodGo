---
name: godgo
description: >
  维护 GodGo 技术外壳：MyGo 界面、设置、窗口外框和红绿灯、工作流、签名更新。
  改 GodGo、加侧栏模块、加设置页、加工作流任务，或用户说外壳、规范、窗口控制栏时使用。
  Use when editing the GodGo repo or running /godgo.
metadata:
  short-description: 按 GodGo 规范改外壳
---

# GodGo

改仓库前按这个顺序读。尺寸、颜色、文案键和扩展点以这些文件为准。不要抄进本技能，也不要写进 Go 字符串。

1. `AGENTS.md`。怎么改、三处扩展、验证命令。
2. `resources/docs/reuse.md`。加组件或能力之前先读。先查 MyGo、Onevra、Kero。
3. `resources/docs/ui-spec.md`。已经实现的界面。窗口、侧栏、设置、字体或快捷键变了，同一次改动里更新它。
4. `README.md`。目录，以及发布前要换的公钥和仓库。

MyGo 的控件生命周期、`c.Key`、绑定值和 `ui.NewTester` 用 mygo-maintenance 里的规则。这里不重写。

## 动手

1. 对上 `AGENTS.md` 里的一处扩展：侧栏模块、设置页，或工作流任务。不要为同一件事再开一个登记处。
2. 用户能看见的字符串放进 `resources/locales/zh-Hans.json` 和 `en.json`，键保持一致。界面测试把语言设成 `zh-Hans`。
3. `backend/` 不导入 `github.com/egoist/mygo`。测试不设 `App.live`，不调用 updater。
4. 窗口外框和红绿灯沿用 `frontend/chrome.go` 的 `framed`。不设 `Frameless`，不自绘红绿灯。位置和边距以界面规范为准。
5. 侧栏和设置的层次继续改编 Kero。文件头保留 MIT 出处。不搬终端、daemon 或 PTY。
6. 用 `ui.NewTester` 走一条用户会点的路径。没有打开原生窗口，就不要说已经看过 `mygo dev`。

## 验证

跑 `AGENTS.md`「验证」里的三条命令。界面规范或复用规范的标题变了，同时改 `backend/assets/fonts_test.go` 的 `TestDocsAgree`。
