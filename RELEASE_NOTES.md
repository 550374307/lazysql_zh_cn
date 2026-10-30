# LazySQL I18N 汉化版 v0.7.1-zh

> 本 fork 基于上游 [LazySQL](https://github.com/jorgerojas26/lazysql)（jorgerojas26），仅实现 I18N 国际化翻译（默认简体中文、中英切换、自定义语言文件），未修改上游功能逻辑。

## 本次修复

- **补全框中文乱码**：SQL 编辑器补全弹窗的描述列此前按字节索引渲染，中文字符的 UTF-8 多字节被逐字节拆开显示（如"字段"显示为 `å æ®µ`）。已改为按 `[]rune` 索引绘制，中文正常显示。
- **回车 / Tab 无法选中补全候选**：在表限定场景（如 `SELECT tenant_llm.`）下回车此前只会关闭弹窗、不插入任何内容。已修复为可正常接受补全。
- **表限定补全丢失表名前缀**：`SELECT tenant_llm.` 选中列后此前会变成 `SELECT id`，已修复为保留表名限定符，结果为 `SELECT tenant_llm.id`；`tenant_llm.cre` 选中后为 `tenant_llm.create_time`。
- **Tab 键被误判为下移**：`KeyTab` 同时命中"向下移动"与"选择"两个分支，导致 Tab 永远只移动候选而不确认。已移除冲突，Tab 与 Enter 行为一致。
- 新增 4 个补全回归测试（表限定空前缀、带前缀、真实 InputHandler 回车、无点前缀不回归）。

## 特性

- 默认简体中文；支持中 / 英切换（`--config` 或 `.lazysql.toml` 中设置 `language`）。
- 支持自定义语言文件（`language_file` 配置，JSON 键值覆盖内置翻译，支持 `${env:VAR}` 展开）。
- README 中文化（默认显示中文），并标注 fork 上游来源。
