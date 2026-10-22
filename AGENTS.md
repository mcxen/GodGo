# GodGo — Agent 维护说明

改这个仓库之前先读：

Grok 的配套技能是 [`.grok/skills/godgo/SKILL.md`](.grok/skills/godgo/SKILL.md)。它只把人带到下面这些文件，不另写一套尺寸和扩展点。

1. 本文件。它只讲怎么改，不重复规范正文。
2. [`resources/docs/reuse.md`](resources/docs/reuse.md)。加组件或加能力之前先读。先查 MyGo、Onevra、Kero，禁止过度设计、过度防御和重复造轮子。
3. [`resources/docs/ui-spec.md`](resources/docs/ui-spec.md)。这是已经实现的界面。窗口、侧栏、设置、字体、快捷键或文案变了，同一次改动里更新它。
4. [`README.md`](README.md)。目录、三处扩展，以及发布前要换的公钥和仓库。

两份规范都由侧栏「规范」渲染。正文只留在 `resources/docs/`，不要写进 Go 字符串，也不要再抄一份。

## 文档各管一段

三份文档不是同一篇文章的三份拷贝：

| 文件 | 谁读 | 打进应用 |
|---|---|---|
| `AGENTS.md` | 改仓库的 agent | 否 |
| `.grok/skills/godgo/SKILL.md` | 进入仓库的 Grok | 否 |
| `resources/docs/reuse.md` | 加组件之前的人和「规范」页的复用 | 是 |
| `resources/docs/ui-spec.md` | 改界面的人和「规范」页的界面 | 是 |
| `README.md` | 第一次打开仓库的人 | 否 |

## 这个仓库是什么

GodGo 是技术外壳，不是业务产品。后续项目沿用窗口、设置、工作流和签名更新，业务加在约定的位置。

- `backend/` 不导入 `github.com/egoist/mygo`。设置、任务、模块注册和文案在这里，测试不打开窗口。
- `frontend/` 用 MyGo 0.6 画界面。侧栏和设置窗口的层次改编自 [Kero](https://github.com/egoist/kero)（MIT），文件头要注明。不要搬 Kero 的终端、daemon 或 PTY。
- `resources/` 是随包的字体、图标、文案、工作流和界面规范。
- `testdata/` 按设置、工作流、文案、字体、界面分开。

`main.go` 只调用 `frontend.Run`。`mygo.json` 不要加 `i18n` 块，文案用我们自己的 JSON。

## 改代码时守住的行为

- `App.live` 只在 `frontend.Run` 里设为 true，且在 `mygo.Use` 之后。测试用 `newApp`，不要设 `live`，不要调用 updater。updater 的 `used()` 在插件未挂上时会 panic。
- 侧栏宽度、折叠和当前模块会记住，但不算「改了设置」。`material()` 去掉这几项之后才决定要不要跑 verify 里的测试任务。拖动侧栏时不要 `apply`。
- `resources/locales/zh-Hans.json` 和 `en.json` 键必须一致。缺键先回退 zh-Hans，再 en，最后显示键本身。界面测试把语言设成 `zh-Hans`。
- 随包字体是 `Inter, Noto Sans SC`。字号是 13 乘 `fontScale`，不要再乘 `c.Theme().FontSize`。
- 工作流命令按参数数组执行，不经过 shell。分组必须先在 `resources/workflow/tasks.json` 里声明。输出保留最后 16KB，单次任务 10 分钟。
- 不要提交私钥。`mygo.json` 里的公钥是占位，发布前用 `go tool mygo keygen` 换成自己的，私钥放进 CI 密钥 `MYGO_UPDATER_PRIVATE_KEY`。

## 扩展点

1. 侧栏模块：`backend/modules/modules.go`，两份文案，`testdata/ui/sidebar.txt`，以及规范里的模块表。
2. 设置页：`frontend/settings_pages.go` 的 `settingsPages`，页面函数，`searchEntries`，以及规范里的设置节。
3. 工作流任务：`resources/workflow/tasks.json`。

## 改界面

- 用户能看见的字符串放进两份文案，不要在 Go 里写死中文或英文。
- 控件先按 `resources/docs/reuse.md` 选用已有的。稳定身份用 `c.Key`。不要把 `ui.Element` 存进 `App`。
- 绑定值在构建之后才落地。`OnChange` 里看到的是这一帧重建后的值。
- 慢的工作放进 goroutine，结果用 `Window.Update` 送回界面。`live` 为 false 时 `update` 直接执行函数，测试不要靠它跨 goroutine 改状态。
- 用 `ui.NewTester` 点一条用户会走的路径。侧栏标题以 `testdata/ui/sidebar.txt` 为准。

## 验证

```
CGO_ENABLED=0 go test ./...
go tool mygo vet .
go tool mygo vet ./frontend
```

没有打开原生窗口，就不要写成已经看过 `mygo dev`。
