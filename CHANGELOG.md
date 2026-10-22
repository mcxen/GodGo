## 0.1.0

- 项目壳：侧栏、设置窗口、工作流任务，以及签名更新的接入点。
- 随包字体为 Inter 与 Noto Sans SC，界面先用 Inter，缺字再用 Noto Sans SC。
- 侧栏和设置窗口的层次改编自 [Kero](https://github.com/egoist/kero)（MIT）。
- 复用规范：组件先查 MyGo、Onevra、Kero，禁止过度设计和重复实现。
- 界面规范写明窗口控制栏：macOS 用系统红绿灯，其他系统在卡片标题栏自绘按钮。
- 界面规范写明窗口边框：保留系统缩放边、阴影和圆角，不用 Frameless，内容离内缘 8。
- 设置页：路径单行省略，资源目录指向实际字体目录，恢复默认按钮不拉满，提示留在设置窗口。
- 配套技能 `.grok/skills/godgo/SKILL.md`：改外壳时先读 AGENTS.md 和两份规范，不在技能里重写尺寸。
