# GodGo

GodGo 是一层技术外壳，做法接近 Spring Boot 的初始化工程：窗口、设置、工作流和签名更新先接好，后面的业务只往约定的位置加。

前端是 [MyGo](https://github.com/egoist/mygo) 0.6。后端是不依赖窗口的 Go 包。侧栏和设置窗口的层次改编自 [Kero](https://github.com/egoist/kero)（MIT）。

## 目录

```
main.go                  入口，只调用 frontend.Run
mygo.json                应用名、标识、版本、更新公钥
frontend/                MyGo 界面
backend/settings/        配置的读写和校验
backend/workflow/        任务目录和执行
backend/modules/         侧栏模块注册
backend/i18n/            文案
resources/fonts/         随包字体
resources/icons/         随包图标
resources/locales/       zh-Hans 与 en，键必须一致
resources/workflow/      任务目录 tasks.json
testdata/settings/       合法配置、未知键、非法枚举
testdata/workflow/       分组正确、未知分组
testdata/i18n/           缺键回退
testdata/fonts/          字体校验和
testdata/ui/             界面测试用的侧栏标题
```

`resources/` 会随 `mygo dev` 和 `mygo build` 打进应用。不要在 `mygo.json` 里再列一遍。

规范有两份，侧栏的「规范」渲染它们。`resources/docs/ui-spec.md` 是已经实现的界面。`resources/docs/reuse.md` 规定加组件前先查 MyGo、Onevra 和 Kero，并且禁止过度设计、过度防御和重复造轮子。改仓库的 agent 先读根目录的 `AGENTS.md`。Grok 用配套技能 `.grok/skills/godgo/SKILL.md`，技能里不重写这两份规范。

## 三处扩展

1. 侧栏模块：`backend/modules/modules.go`，并在 `resources/locales` 的两份文案里加上同样的键。
2. 设置页：`frontend/settings_pages.go` 的 `settingsPages`，再写对应的页面函数。
3. 工作流任务：`resources/workflow/tasks.json`。分组必须先在文件里声明。命令按参数数组执行，不经过 shell。

## 字体

界面字体是 `Inter, Noto Sans SC`。Inter 只覆盖拉丁字母，简体中文落到 Noto Sans SC。两套都是 SIL OFL 1.1 的静态子集，和 AvaMedia 桌面端用的是同一批文件。设置里可以改回系统字体。

## 开发

```
go tool mygo dev
CGO_ENABLED=0 go test ./...
go tool mygo vet ./...
```

开发运行的副本不能给自己装更新。设置里的「这个副本不能自己更新」说的就是这件事。

## 发布前要换的两处

`mygo.json` 里的公钥是占位，只能证明格式合法，校验不了真正的更新包。发布前在本机执行：

```
go tool mygo keygen
```

把生成的公钥写进 `mygo.json` 的 `updates.publicKey`。私钥不要提交。CI 里把它放在密钥 `MYGO_UPDATER_PRIVATE_KEY`。

`updates.github` 现在写的是 `openspring/GodGo`。改成实际的 `owner/repo`，否则检查更新会指向错误的仓库。

macOS 签名不写死在配置里。发布工作流只在仓库配了 `MACOS_CERTIFICATE` 和 `MACOS_SIGNING_IDENTITY` 时才签名，配了 Apple 公证账号时才公证。
