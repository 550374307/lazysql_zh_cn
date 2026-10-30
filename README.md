<a name="readme-top"></a>

<!-- LANGUAGES -->
<div align="center">
  **简体中文** · <a href="README.en.md">English</a>
</div>

> [!NOTE]
> **Fork 汉化版说明**
>
> - **上游来源**：[LazySQL](https://github.com/jorgerojas26/lazysql)（作者 Jorge Rojas）
> - **本 fork 定位**：仅实现 **I18N 国际化翻译**（默认简体中文，支持中/英切换与自定义语言文件），**不修改上游的任何功能逻辑**
> - 上游更新可随时通过 `git fetch upstream` 合并同步

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]
[![LinkedIn][linkedin-shield]][linkedin-url]

<!-- PROJECT LOGO -->
<div align="center">
  <!-- <a href="https://github.com/jorgerojas26/lazysql"> -->
  <!--   <img src="images/logo.png" alt="Logo" width="80" height="80"> -->
  <!-- </a> -->

  <h3 align="center">LAZYSQL</h3>

  <p align="center">
        一个用 Go 编写的跨平台 TUI 数据库管理工具。
  </p>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>目录</summary>
  <ol>
    <li>
      <a href="#关于项目">关于项目</a>
      <ul>
        <li><a href="#技术栈">技术栈</a></li>
      </ul>
    </li>
    <li><a href="#功能特性">功能特性</a></li>
    <li>
      <a href="#快速开始">快速开始</a>
      <ul>
        <li><a href="#安装">安装</a></li>
      </ul>
    </li>
    <li><a href="#可选遥测">可选遥测</a></li>
    <li><a href="#配置">配置</a></li>
    <li><a href="#使用">使用</a></li>
    <li><a href="#支持">支持</a></li>
    <li><a href="#命令">命令</a></li>
    <li><a href="#环境变量">环境变量</a></li>
    <li><a href="#快捷键">快捷键</a></li>
    <li><a href="#示例连接-url">示例连接 URL</a></li>
    <li><a href="#路线图">路线图</a></li>
    <li><a href="#剪贴板支持">剪贴板支持</a></li>
    <li><a href="#参与贡献">参与贡献</a></li>
    <li><a href="#许可证">许可证</a></li>
    <li><a href="#联系方式">联系方式</a></li>
    <li><a href="#备选方案">备选方案</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->

## 关于项目

![Product Name Screen Shot][product-screenshot1]
![Product Name Screen Shot][product-screenshot2]

这个项目的灵感主要来自 [Lazygit](https://github.com/jesseduffield/lazygit)，我认为它是 Git 最好的 TUI 客户端。

我想要一个类似的工具，但用于 SQL。我没有找到符合需求的，所以自己造了一个。

我整天生活在终端里，如果你也像我一样，这个工具对你也会很有用。

这是我的第一个开源项目，也是我的第一个 Go 项目。我不是一个出色的程序员，只是一个想学一门新语言的普通 JavaScript 开发者，同时我也想要一个 TUI SQL 客户端，所以 blanca y en botella, leche!（白瓶装牛奶，顺理成章！）

这个项目还处于 **ALPHA 阶段**，欢迎吐槽我的意大利面条式代码。

我在全职工作中每天都在使用 Lazysql（目前还是带 bug 的状态 xD）。所以计划是在业余时间作为副项目不断改进和修复我的这个小家伙。

### 技术栈

![Golang][golang-shield]
![Golang][tview-shield]

## 功能特性

- [x] 跨平台（macOS、Windows、Linux）
- [x] Vim 键位
- [x] 支持管理多个连接（Backspace）
- [x] 标签页
- [x] SQL 编辑器（CTRL + e）

<!-- GETTING STARTED -->

## 快速开始

### 安装

#### Homebrew（macOS/Linux）

```bash
brew install lazysql
```

#### 使用 Go 包管理器安装

```bash
go install github.com/jorgerojas26/lazysql@latest
```

#### 二进制发行版

对于 Windows、macOS 或 Linux，你可以从 [这里](https://github.com/jorgerojas26/lazysql/releases) 下载二进制发行版。

#### 第三方（社区维护）

Arch Linux 用户可以通过 AUR 安装：

```bash
paru -S lazysql

```

或者

```bash
yay -S lazysql

```

或者手动安装：

```bash
git clone https://aur.archlinux.org/lazysql.git
cd lazysql
makepkg -si
```

<p align="right">(<a href="#readme-top">返回顶部</a>)</p>

<!-- USAGE EXAMPLES -->

## 可选遥测

官方发行版在发送任何数据前，会先披露极简的匿名聚合使用统计。默认**保持开启**；选择 **Disable** 或使用 `--no-telemetry` 即可退出。源码构建不发送任何数据。详见 [收集的数据、退出方式、限制与免费自托管说明](docs/telemetry.md)。

## 配置

如果设置了 `XDG_CONFIG_HOME` 环境变量，配置文件位于：

- `${XDG_CONFIG_HOME}/lazysql/config.toml`

否则，配置文件位于：

- Windows：`%APPDATA%\lazysql\config.toml`
- macOS：`~/Library/Application Support/lazysql/config.toml`
- Linux：`~/.config/lazysql/config.toml`

配置文件是 TOML 格式，可用于定义多个连接。

### 示例配置

```toml
[[database]]
Name = 'Production database'
Provider = 'postgres'
DBName = 'foo'
URL = 'postgres://${user}:urlencodedpassword@localhost:${port}/foo'
ReadOnly = true
Commands = [
  { Command = 'ssh -tt remote-bastion -L ${port}:localhost:5432', WaitForPort = '${port}' },
  { Command = 'whoami', SaveOutputTo = 'user' },
]
[[database]]
Name = 'Development database'
Provider = 'postgres'
DBName = 'foo'
URL = 'postgres://postgres:urlencodedpassword@localhost:5432/foo'
[application]
DefaultPageSize = 300
DisableSidebar = false
SidebarOverlay = false
JSONViewerWordWrap = false
EnterOpensJSONViewer = false
language = "zh"
language_file = ""
schema_bulk_load_threshold = 200
exact_count_threshold = 50000
exact_count_timeout_ms = 200
max_query_rows = 1000
max_open_connections = 8
max_idle_connections = 8
```

`ReadOnly` 字段（可选，默认为 `false`）可以设为 `true` 以启用连接的只读模式。启用后，所有变更类查询（INSERT、UPDATE、DELETE、DROP 等）都会被阻止。

### MySQL login-path 凭据（无需在配置中写明文密码）

MySQL 的 `mysql_config_editor` 工具可以把连接凭据混淆加密后保存到 `~/.mylogin.cnf`（Windows 为 `%APPDATA%\MySQL\.mylogin.cnf`）。LazySQL 可以直接读取这个加密文件，因此 `.lazysql.toml` 里不再需要保存明文密码。

1. 先在终端里生成一个 login-path（按提示输入密码，密码不会进入 shell 历史）：

   ```bash
   mysql_config_editor set --login-path=prod --host=127.0.0.1 --port=3306 --user=root --password
   # 查看已保存的条目（会显示明文，仅用于确认）
   mysql_config_editor print --all
   ```

2. 然后在配置里只引用名字，密码等凭据留空：

   ```toml
   [[database]]
   Name = '生产库'
   Provider = 'mysql'
   LoginPath = 'prod'
   DBName = 'app'
   ```

   加载连接时，`user` / `password` / `host` / `port` 中**留空**的字段会自动从 `prod` 条目补齐；**显式写出的字段优先**（与 `mysql` 客户端一致：命令行选项覆盖 `--login-path`）。只写 `Provider = 'mysql'` + `LoginPath` + `DBName` 时，LazySQL 会自动拼接出 `mysql://user:pass@host:3306/db` 形式的连接串。

   该加密仅用于防止凭据被直接读取（与 MySQL 官方说明一致，并非强加密）；`~/.mylogin.cnf` 自身也应保持仅当前用户可读。

**自动发现（免配置）**：LazySQL 启动时会自动扫描 `~/.mylogin.cnf`（Windows 为 `%APPDATA%\MySQL\.mylogin.cnf`），把里面**每一个** login-path 条目直接作为一条 MySQL 连接显示在连接选择器里——无需在 `.lazysql.toml` 里写 `[[database]]`。选中即可连接，连接后侧边栏默认列出该用户有权限的全部数据库。规则：

- 自动条目只在会话内存中，**不会**被写回配置文件；
- 若配置里已有同名 `LoginPath` 的显式条目，以显式条目为准（自动条目跳过）；
- 条目未指定 host 时默认 `localhost`、端口默认 `3306`；条目带 `socket` 且未指定 host 时走 unix socket 连接；
- 登录文件缺失或无法解析时静默跳过，不影响启动。

数据库条目可以用 `max_open_connections` 和 `max_idle_connections` 覆盖连接池设置。省略的值继承应用设置；显式的 `0` 使用 LazySQL 默认值 `8`。无效组合（例如空闲连接数大于打开连接数）会被拒绝。SQLite 始终使用一条打开连接和一条空闲连接（`1/1`）以保持内存数据库行为，忽略通用连接池设置。

`DBName` 字段（可选）控制打开连接时侧边栏树的填充方式：

- **设置了 `DBName`**（例如 `DBName = 'foo'`）：树只固定显示该单一数据库。这是通过应用内连接表单创建/编辑连接时的默认行为，因为 `DBName` 会从连接 URL 中内嵌的数据库自动填充。
- **`DBName` 为空或省略**：树会列出连接用户拥有 `CONNECT` 权限的*所有*数据库（通过 `GetDatabases()`），而不只是 URL 里的那个。当你在同一台 PostgreSQL 或 MSSQL 服务器上经常需要跨多个数据库/目录浏览或查询时非常有用——尤其是 PostgreSQL，连接串里必须指定一个数据库，尽管同一个登录账号通常对该实例上的多个数据库都有效。

要为某个连接启用多数据库浏览，可以：

- 直接在 `config.toml` 中删除/注释掉该连接的 `DBName` 行，**或者**
- 在应用内连接表单中添加或编辑连接时勾选 **"Show all databases in this instance"**（`a` 添加、`e` 在连接选择器中编辑）。这会覆盖自动填充的 `DBName`，并以留空的方式保存连接，无需手工编辑 TOML 文件。重新打开已有连接编辑时，如果之前是这样保存的，复选框会自动预勾选。

注意：此设置只影响哪些数据库填充侧边栏树。它不会改变初始连接认证所针对的数据库——那仍然完全由连接 URL 中的数据库决定（这是 PostgreSQL 和 MSSQL 线协议的要求）。

`[application]` 段用于定义一些应用设置。并非所有设置都已可用，这是一个持续完善的过程。

### 应用设置

| 设置 | 默认值 | 说明 |
| ------- | ------- | ----------- |
| DefaultPageSize | 300 | 每页获取的记录数 |
| DisableSidebar | false | 禁用侧边栏 |
| SidebarOverlay | false | 将侧边栏显示为覆盖层而不是侧面板 |
| JSONViewerWordWrap | false | 在 JSON 查看器中启用自动换行 |
| EnterOpensJSONViewer | false | 在单元格上按 Enter 时打开 JSON 查看器 |
| DisableUpdateCheck | false | 禁用自动 GitHub 发行版检查（手动检查仍然可用） |
| exact_count_threshold | 50000 | 当未过滤的估算值小于等于该值时自动执行精确计数（0 = 永远不对估算值自动计数） |
| exact_count_timeout_ms | 200 | 自动精确计数的预算时间，单位毫秒（0 = 禁用自动精确计数；手动 `#` 仍然可用） |
| max_query_rows | 1000 | 交互式 SQL 编辑器结果最多显示的行数（0 = 无限制） |
| schema_bulk_load_threshold | 200 | 可见表数量达到该值时批量预加载列（0 = 始终惰性加载） |
| max_open_connections | 8 | MySQL、PostgreSQL 和 MSSQL 的最大打开连接数（0 = 使用默认值 8） |
| max_idle_connections | 8 | MySQL、PostgreSQL 和 MSSQL 的最大空闲连接数（0 = 使用默认值 8） |
| language | zh | UI 语言。内置代码：`zh`（简体中文，默认）和 `en`（英文）。其他任意代码以英文为基础，配合 `language_file` 使用 |
| language_file | "" | 自定义语言文件的路径（JSON，将英文 UI 字符串映射为译文）。条目会覆盖内置语言包，可在不重新编译的情况下实现任意语言或术语定制 |

`schema_bulk_load_threshold` 限制对小型可见 schema 的列做批量预加载。`0` 让所有列保持惰性；较大的 schema 仍然立即暴露表名，并在需要时按需获取指定表。无论阈值如何，已缓存或加载中的列都会被复用。

### 国际化（i18n）

LazySQL 内置简体中文（`zh`，默认）和英文（`en`）两种 UI 语言包。通过 `language` 设置选择语言：

```toml
[application]
language = "zh"   # 或 "en"
```

要使用自己的语言或术语而无需重新编译，请将 `language_file` 指向一个 JSON 文件，把英文 UI 字符串映射为你的译文。自定义文件中的条目会覆盖内置语言包；自定义文件中缺失的键保留内置翻译：

```json
{
  "Yes": "行",
  "Save": "存储",
  "Query failed: %s": "查询出错了：%s"
}
```

自定义语言文件在内置语言包之后加载，因此你只需包含想要修改的字符串。将 `language` 设置为 `zh`/`en` 以外的值（例如 `"custom"`），即以英文为基础，由自定义文件提供整个 UI。默认语言为 `zh`。

### 内置翻译文件位置

内置语言包在编译时内嵌：

| 文件 | 作用 |
| --- | --- |
| `internal/i18n/locales/zh.json` | 简体中文语言包（默认；当前 433 个键） |
| `internal/i18n/locales/en.json` | 英文语言包（保持为 `{}`；所有键回退到英文原文） |

翻译**键就是英文 UI 字符串本身**，另有三个运行时渲染的前缀键族：

| 键类型 | 示例 | 来源 |
| --- | --- | --- |
| 普通 UI 字符串 | `"Save"`、`"Query failed: %s"` | 各组件中的 `i18n.T(...)` / `i18n.Tf(...)` 调用 |
| 命令名 | `"cmd.Quit"`、`"cmd.CopyRowsAs"` | `commands.Command.String()` |
| 键位分组 | `"keymap.group.table"` | `app/keymap.go` 中的分组名 |
| 按键描述 | `"Focus table"`、`"Export to CSV"` | `app/keymap.go` 中的 `Bind.Description` |

要修改中文 UI 的默认文案，编辑 `zh.json` 后重新编译即可。要新增一门语言，复制 `zh.json` 的结构（或从空文件开始，让英文回退补足缺口），再通过 `language_file` 在运行时加载——无需重新编译。

### 升级提示与应用内升级

LazySQL 会在后台检查官方 GitHub 仓库是否有更新的稳定版本，每次成功检查后每天最多一次。提示会出现在底部状态栏，不会打断你的工作。开发构建不检查更新。只请求发行版元数据；不发送任何数据库或连接信息。

在任何位置按 **F10** 检查更新或打开可用升级。快捷键可在 `[keymap.global]` 中用 `CheckForUpdates` 修改。

- **独立的 macOS/Linux 二进制：** 选择 **Upgrade** 下载匹配的发行版，与官方发行版校验和文件核对 SHA-256 校验和，并原子替换可执行文件。安装目录必须可写。LazySQL 永远不会运行 `sudo` 或重启你的活动会话；准备好后请手动重启。
- **Homebrew 和其他受认可的受管安装：** 对话框给出包管理器指令，而不是覆盖受管文件。
- **Windows：** 对话框提供手动/包管理器升级说明；不支持就地自升级。

自动检查可通过在 `[application]` 下设置 `DisableUpdateCheck = true` 或设置环境变量 `LAZYSQL_NO_UPDATE_CHECK=1` 来禁用。手动检查仍然可用。网络失败不会阻止启动；手动检查会显示错误并可重试。关闭升级对话框后你仍可继续工作，同时下载在后台完成。退出 LazySQL 会取消未完成的下载，并等待安装程序文件操作完成。如果校验或暂存失败，现有二进制保持原样。

### 主题

颜色在 `[theme]` 段中设置。选择一个预设，并可选地在预设之上覆盖单个颜色：

```toml
[theme]
Preset = "light"

[theme.Colors]
Border = "#666A7E"
Title = "black"
PrimaryText = "default"
EditorStatusBarBackground = "lightgray"
SQLKeyword = "#005F87"
```

预设：

| 预设 | 说明 |
| ------ | ----------- |
| `default` | 标准 lazysql 颜色，使用终端的背景色 |
| `light` | 受 GitHub 启发的浅色配色，带自定义背景 |
| `dracula` | [Dracula](https://draculatheme.com) |
| `gruvbox-dark` | [Gruvbox](https://github.com/morhetz/gruvbox) 暗色 |
| `nord` | [Nord](https://www.nordtheme.com) |
| `solarized-light` | [Solarized](https://ethanschoonover.com/solarized/) 浅色 |
| `tokyo-night` | [Tokyo Night](https://github.com/enkia/tokyo-night-vscode-theme) |
| `catppuccin-mocha` | [Catppuccin](https://github.com/catppuccin/catppuccin) Mocha |

在 lazysql 中任意位置按 `Ctrl+T` 打开主题选择器。用方向键、`j`/`k` 或 `Ctrl+N`/`Ctrl+P` 在运行中的应用上预览每个主题。按 `Enter` 应用并保存所选预设，或按 `q`/`Esc` 恢复之前的主题。

`default` 使用终端的背景色。其他所有预设都会绘制自己的背景。

颜色可以是名称（`"red"`、`"dodgerblue"`，任意 [W3C 颜色名](https://github.com/gdamore/tcell/blob/v2.7.4/color.go#L851)）、十六进制值（`"#RRGGBB"`），或 `"default"` 使用终端自身的颜色。lazysql 启动时会报告未知的键、预设和颜色。

`[theme.Colors]` 的键：

| 键 | 用途 |
| --- | -------- |
| `PrimitiveBackground` | 所有视图的背景 |
| `ContrastBackground`、`MoreContrastBackground` | 输入框、按钮和下拉框的背景 |
| `Border`、`Title`、`Graphics` | 默认边框、标题和线条颜色 |
| `PrimaryText` | 主要文本 |
| `SecondaryText` | 高亮：选中行、活动标签页、按键提示、侧边栏值 |
| `TertiaryText` | 标签和状态消息 |
| `InverseText` | 未聚焦视图、表单字段的边框和文本 |
| `ContrastSecondaryText` | 绘制在 `SecondaryText` 上的文本（选中行、活动标签页） |
| `SidebarTitleBorder` | 侧边栏字段名旁边的分隔线 |
| `Error` | 错误消息 |
| `ReadOnly` | 只读标记 |
| `TableChange`、`TableInsert`、`TableDelete`、`TableMarked` | 已修改、已插入、已删除和已标记行的背景 |
| `EditorSelection` | SQL 编辑器中的选中文本 |
| `EditorStatusBarBackground`、`EditorStatusBarText` | SQL 编辑器状态栏 |
| `AutocompleteBackground`、`AutocompleteText`、`AutocompleteSelected`、`AutocompleteDescription`、`AutocompleteSeparator` | SQL 编辑器自动补全弹窗 |
| `SQLKeyword`、`SQLString`、`SQLNumber`、`SQLComment`、`SQLFunction`、`SQLOperator`、`SQLType`、`SQLBoolean`、`SQLParameter` | SQL 语法高亮 |
| `JSONKey`、`JSONString`、`JSONBoolean`、`JSONNull`、`JSONNumber` | JSON 查看器 |

### 本地配置

你可以在项目目录中（`.git` 文件夹旁边）放置一个 `.lazysql.toml` 文件，为该项目覆盖全局配置。这对于定义项目专属的数据库连接或设置很有用。

lazysql 从当前工作目录向上查找 `.lazysql.toml`。在 git 仓库根目录（找到 `.git` 的地方）停止。如果找不到本地配置，则直接使用全局配置。

**合并行为：**

| 段 | 行为 |
| ------- | -------- |
| `[application]` | 深度合并——本地值覆盖全局，未设置字段保留全局/默认值 |
| `[[database]]` | 替换——本地连接完全替换全局连接 |
| `[keymap.*]` | 深度合并——同一命令的本地按键覆盖全局 |
| `[theme]` | 深度合并——本地的 `Preset` 或颜色覆盖全局，其他颜色保留 |

**示例 `.lazysql.toml`：**

```toml
[application]
DefaultPageSize = 500

[[database]]
Name = 'Local development'
Provider = 'postgres'
URL = 'postgres://localhost/myproject_dev'
```

使用此本地配置时，`DefaultPageSize` 覆盖全局值，且只有 `Local development` 连接可用（全局连接被替换）。

环境变量（`${env:VAR_NAME}`）在本地配置文件中的用法与全局配置相同。保存连接时，会保留你未修改字段中的占位符，而不是写入展开后的凭据。

注意：从 UI 所做的更改只更新发生变化的设置，并保留目标文件的其余部分。如果本地 `.lazysql.toml` 定义了 `[[database]]`，你添加、编辑或删除的连接会保存到本地 `.lazysql.toml`，否则保存到全局配置文件。当找到本地 `.lazysql.toml` 时，用 `Ctrl+T` 选择的主题会保存到本地文件，否则保存到全局配置文件。lazysql 保存文件时不会保留注释。

### 手动数据库测试环境

本仓库包含一套即插即用、仅供开发的 fixtures，用于 MySQL、PostgreSQL、MSSQL 和 SQLite。该 fixtures 用于手动测试记录、分页、过滤、排序、元数据、外键跳转、SQL 编辑器结果和 CSV 导出。

前置条件为 Docker、Docker Compose v2 和 OpenSSL。该栈不需要宿主机安装任何数据库客户端。

在仓库根目录：

```bash
./scripts/manual-databases.sh validate
./scripts/manual-databases.sh init
./scripts/manual-databases.sh up
lazysql
```

`init` 会生成随机密码和私有的、被 git 忽略的 `.env.manual-databases` 与 `.lazysql.toml` 文件。它拒绝覆盖已有的本地配置。端口只绑定到 `127.0.0.1`。不要在其他地方复用这些开发凭据。如果是从旧固定密码 fixtures 升级，请在 `init` 后运行 `reset` 以用新密码重建卷（这会删除所有 fixtures 数据）。

当 LazySQL 从仓库根目录启动时，会自动发现生成的 `.lazysql.toml`，其中包含四个即用连接：

| 连接 | 提供方 | 宿主机端口 | fixtures 数据库 |
| ---------- | -------- | --------- | ---------------- |
| Docker MySQL | `mysql` | `3307` | `lazysql_test` |
| Docker PostgreSQL | `postgres` | `5433` | `lazysql_test` |
| Docker MSSQL | `sqlserver` | `14331` | `lazysql_test` |
| Docker SQLite | `sqlite3` | 文件 | `testdata/sqlite/lazysql.sqlite3` |

使用辅助脚本管理生命周期：

```bash
./scripts/manual-databases.sh status  # 健康与种子状态
./scripts/manual-databases.sh down    # 停止，保留数据
./scripts/manual-databases.sh reset   # 销毁并重新播种所有内容
```

`up` 会构建小型 SQLite 辅助镜像、初始化服务器 fixtures，并等待所有数据库包含其种子表。MySQL 和 PostgreSQL 使用官方镜像的初始化钩子和命名卷；这些脚本在卷为空时运行。MSSQL 通过第二个容器播种，并保留种子标记，因此重启栈不会覆盖手动更改。SQLite 是本地文件，由辅助容器播种，可安全地多次初始化。需要干净 fixtures 时使用 `reset`。

fixtures 包含 1,200 位客户、300 个产品、3,000 个订单、9,000 个订单项和 2,400 条客户备注。包含主/外键、唯一约束、索引、可空列、日期、数值和布尔值、JSON/文本值、视图以及各提供方特有的目录对象。模型见 [`testdata/README.md`](testdata/README.md)，端口和凭据见 [`docker-compose.yml`](docker-compose.yml)。

MSSQL 镜像是 x86-64 专属，服务明确以 `linux/amd64` 运行。Docker 在 Apple Silicon 上可能会模拟它；微软不支持该模拟路径，因此如果失败，请使用原生 x86 主机或远程 SQL Server。如果宿主机端口已被占用，请在启动栈前同时修改 `docker-compose.yml` 和 `.lazysql.toml` 中的端口。SQLite URL 相对于进程工作目录，因此请从仓库根目录启动 LazySQL 以使用该连接。

## 使用

> 查看键盘快捷键列表请按 `?`

打开 TUI：

```console
$ lazysql
```

从已保存的连接中选择并启动：

```console
$ lazysql [connection_url]
```

连接 `[connection_url]` 处的数据库并启动：

```console
$ lazysql --read-only [connection_url]
```

以只读模式启动。

### 连接数据库

1. 启动 `lazysql`
2. 新建连接（按 `n`）
3. 为连接提供名称以及要连接的 URL（参见 <a href="#示例连接-url">示例连接 URL</a>）
4. 连接数据库（按 `<Enter>`）

如果你已经有连接配置：

1. 启动 `lazysql`
2. 选择正确的连接（按 `j` 和 `h` 导航）
3. 连接数据库（按 `c` 或 `<Enter>`）

### 创建表

目前还没有办法从 TUI 中直接创建表。但你可以在 <a href="#执行-sql-查询">SQL 编辑器</a> 中把建表语句作为 SQL 查询运行。

你可以按 `R` 更新树，就能看到新创建的表。

### 执行 SQL 查询

1. 按 `<Ctrl+E>` 打开内置 SQL 编辑器
2. 编写 SQL 查询
3. 按 `<Ctrl+R>` 执行 SQL 查询

> 要切回表树请按 `H`
>
> 执行 `SELECT` 查询后，SQL 编辑器下方会显示一个包含查询结果的表。\
> 要切回 SQL 编辑器请按 `/`

### 结果与网络语义

记录页按"一页行数 + 一行前瞻"获取。前瞻行让导航无需精确计数即可工作，因此分页标签有三种有意的形式：

- `843 rows` / `1-843 of 843 rows`：**精确**；数据库计数或页尾推断证明了总数。
- `~4.3M rows`：**估算**；数据库提供了有用的估算值，但不保证准确。
- `300+ rows`：**未知更多**；当前页还有更多行，尚不知道精确总数。

#### 精确记录数

在 Records 界面按 `#`（`ExactCount`）开始精确计数。运行时再按 `#` 取消；失败的计数只停留在分页栏，可用 `#` 重试。自动计数永远不会阻塞第一页 Records：LazySQL 优先使用驱动估算（如有），仅当估算值小于等于 `exact_count_threshold`（默认 `50000`）或没有估算值时，才运行精确计数。阈值为 `0` 会禁用基于估算的自动精确计数；过滤计数和无估算值仍受超时约束。自动计数受 `exact_count_timeout_ms`（默认 `200`）限制。将该超时设为 `0` 可禁用自动精确计数；估算值仍会渲染，手动 `#` 计数仍然可用。过滤后的 Records 计数也受自动超时约束。

#### SQL 结果与取消

交互式 SQL 结果会渐进式流式返回，默认上限为 `max_query_rows = 1000`。正数上限最多渲染这么多行，并执行一次前瞻读取，以便 UI 显示 `result truncated`；`0` 表示无限制。该上限只适用于交互式结果，不适用于完整导出。当 SQL 编辑器结果查询正在运行时，按 `Esc` 取消其上下文并保留已渲染的部分行。当没有结果查询运行时，`Esc` 保持编辑器的正常取消聚焦/编辑行为。

#### CSV 导出

- **导出可见结果** 写入已显示的行，绝不重新执行 SQL 语句。
- **导出全部结果** 独立于 `max_query_rows` 流式导出完整的表/查询结果。对于 SQL 结果，它只重新执行保守、可重放的安全只读语句；变更类或未知语句只提供可见结果选项。

两种范围都通过临时文件写入。取消、查询失败或写入/重命名失败都会删除临时文件，并保持已有目标文件不变。

#### 性能诊断

如需本地 JSONL 计时，请以 `--loglevel debug --logfile /path/to/lazysql.jsonl` 启动 LazySQL。日志包含操作耗时、数据库标识、缓存/回退结果、取消/失败，以及 Records 和 SQL 编辑器结果的 `event=first_useful_result`。SQL 文本、参数、行值、凭据和连接 URL 都会被脱敏；不发送任何遥测。

实现矩阵和无凭据 RTT 基准测试工具见 [`docs/performance.md`](docs/performance.md)。

### 打开或查看表

1. 按 `e` 或 `<Enter>` 展开表树
2. 选择要查看的表
    - 下一个节点 `j`
    - 上一个节点 `k`
    - 最后一个节点 `G`
    - 第一个节点 `g`
3. 按 `<Enter>` 打开表

> 要切回表树请按 `H` \
> 要切回表请按 `L`

### 过滤行

1. [打开表](#打开或查看表)
2. 按 `/` 聚焦过滤输入框
3. 编写 `WHERE` 子句过滤表
4. 按 `<Enter>` 提交过滤

> 要移除过滤，聚焦过滤输入框（按 `/`）并按 `<Esc>`。

### 隐藏列

在表的 Records 视图中，按 `V`（**visibility**）或点击分页栏中的 **Columns [V]**。一个紧凑的选择器让你选择显示哪些列：

- **空格** 或点击某列切换其复选框。
- **↑/↓** 或 **j/k** 移动；**/** 按名称查找列。
- **A** 或 **All** 恢复所有列，包括搜索范围之外的。
- **Enter** 或 **Apply** 保存；**Esc** 或 **Cancel** 放弃更改。
- **Tab / Shift+Tab** 在列表、搜索和按钮之间移动。

至少保持一列可见。底部会显示隐藏了多少列。选择会按连接、数据库和表（包括 schema 限定的表名）在重启后记住。新列默认可见。

这只是显示层面：查询仍会获取所有列，编辑、主键、行侧边栏、行复制和 CSV 导出都保留完整数据。SQL 编辑器结果和元数据视图不受影响。

偏好存储在当前活动配置文件（项目 `.lazysql.toml` 存在时用它，否则用全局配置）的 `[hidden_columns]` 下。快捷键可在 `[keymap.Table]` 中用 `ColumnVisibility` 自定义。

### 跳转到被引用的行（外键跳转）

属于外键的列在表头中带下划线，可跟随的值也带下划线。带待定更改（插入、编辑或删除行）的单元格不会被标记。

1. [打开表](#打开或查看表)
2. 按 `1` 切换到记录标签页
3. 移动到你想要跟随的下划线外键单元格
4. 按 `<Enter>` 跳转到被引用的行

被引用的表会打开在新标签页中（如果已打开则聚焦），并预填被引用值的过滤条件，因此只显示匹配的行：`WHERE <referenced_column> = '<value>'`。

> 按 `[` 或 `]` 切换标签页，按 `H` 切回表树。

注意：

- 只能跟随单列外键。复合外键（跨多列的约束）会被忽略。
- 值为 `NULL`、`EMPTY` 或 `DEFAULT` 的单元格没有可跟随的内容，因此对其按 `<Enter>` 不会做任何事。
- 支持的提供方：PostgreSQL、SQLite 和 SQL Server（MSSQL）。MySQL 表不暴露跳转所需的外键信息。
- 在支持外键跳转的单元格上，即使启用了 `EnterOpensJSONViewer`，按 `<Enter>` 也会跟随关系而不是打开 JSON 查看器。

### 插入行

1. [打开表](#打开或查看表)
2. 按 `1` 切换到记录标签页
3. 按 `o` 插入新行
4. 填写所有列
5. 按 `<Ctrl+S>` 保存更改

### 编辑列

1. [打开表](#打开或查看表)
2. 按 `1` 切换到记录标签页
3. 移动到你想要编辑的单元格
4. 编辑值：
    - 按 `c` 内联编辑，然后按 `<Enter>` 提交
    - 按 `e` 在你的编辑器中编辑，保存并退出编辑器
5. 按 `<Ctrl+S>` 保存更改

### 复制行

1. [打开表](#打开或查看表)
2. 移动到一行并按 `<Space>` 标记。对不相邻的行重复（再按 `<Space>` 取消标记）。对于连续范围，按 `v` 并上下移动：高亮选区随光标扩大。按 `<Space>` 将该范围加入标记，或按 `<Esc>`（或 `v`）取消。
3. 在选择范围时按 `Y` 直接复制，无需提交标记。否则，`Y` 复制标记的行；没有标记时复制当前行。选择 SQL（`s`）、Markdown（`m`）、JSON（`j`）或 CSV（`c`）。CSV 和 Markdown 包含列表头；JSON 是对象数组。SQL 在记录表中生成 `INSERT`，在 SQL 编辑器结果中生成独立的 `VALUES` 语句（没有目标表时）。

这在记录表和 SQL 编辑器结果中都有效。`y` 保持不变：它把标记的行复制为制表符分隔的值，没有标记行时复制选中的单元格。结果重新加载时标记会清除。

### 导出为 CSV

#### 从表视图

1. [打开表](#打开或查看表)
2. 按需应用过滤或排序
3. 按 `E` 打开导出对话框
4. 可选地修改文件路径和批量大小
5. 选择导出范围：
   - 导出当前页：只导出当前显示的行
   - 导出全部记录：从表中获取并导出所有记录

> 批量大小（默认：10000）：导出全部记录时，数据按批获取，以避免大表出现超时或内存问题。增大可获得更快的导出，遇到错误时减小。
>
> 默认文件路径是 `~/Downloads/{database}_{table}_{timestamp}.csv`。

#### 从 SQL 编辑器

1. [执行 SQL 查询](#执行-sql-查询)
2. 按 `E` 打开导出对话框
3. 可选地修改文件路径
4. 选择 **导出可见结果** 保存当前显示的行
5. 对于可重放/只读语句，选择 **导出全部结果** 重新执行查询并流式导出每一行

> 导出全部结果可能会重新执行查询，可能需要较长时间。未知或可能变更的语句只提供可见结果导出。导出过程中按 `Esc` 取消；已取消或失败的导出保持目标文件不变。

<p align="right">(<a href="#readme-top">返回顶部</a>)</p>

## 支持

- [x] MySQL
- [x] PostgreSQL
- [x] SQLite
- [x] MSSQL
- [x] ClickHouse
- [ ] MongoDB

对多种 RDBMS 的支持仍在持续完善中。

<!-- COMMANDS -->

## 命令

某些情况下，主要是连接远程数据库时，可能需要在连接数据库之前运行自定义命令。例如，当你只能通过远程堡垒机访问数据库时，可能需要在单独的终端中先打开 SSH 隧道：

```bash
ssh remote-bastion -L 5432:localhost:5432
```

为了更方便地运行这些命令，lazysql 支持在连接数据库前运行自定义命令。你可以在配置文件中这样定义这些命令：

```toml
[[database]]
Name = 'server'
Provider = 'postgres'
DBName = 'foo'
URL = 'postgres://${user}:password@localhost:${port}/foo'
Commands = [
  { Command = 'ssh -tt remote-bastion -L ${port}:localhost:5432', WaitForPort = '${port}' },
  { Command = 'whoami', SaveOutputTo = 'user' },
]
```

`Command` 字段是必需的，可以包含任何你通常在终端中运行的命令。`WaitForPort` 字段可选，可用于在继续之前等待特定端口打开。`SaveOutputTo` 字段可选，可用于创建用户自定义变量。命令的输出（`stdout`）会被保存到该变量中，并可通过 `${VARIABLE}` 语法在 URL 或后续命令中使用。

当你在 URL 字段中定义 `${port}` 变量时，lazysql 会自动用随机空闲端口号替换它。该端口号随后用于连接 URL，并在 `Commands` 字段中可用，以便配置命令。

你甚至可以链式执行命令，例如先连接远程服务器，再连接远程 k8s 集群中运行的 postgres 容器：

```toml
[[database]]
Name = 'container'
Provider = 'postgres'
DBName = 'foo'
URL = 'postgres://postgres:password@localhost:${port}/foo'
Commands = [
  { Command = 'ssh -tt remote-bastion -L 6443:localhost:6443', WaitForPort = '6443' },
  { Command = 'kubectl port-forward service/postgres ${port}:5432 --kubeconfig /path/to/kube.conf', WaitForPort = '${port}' }
]
```

## 环境变量

你可以在配置文件中使用 `${env:VAR_NAME}` 语法引用环境变量。这对于将密码等敏感信息排除在配置文件之外非常有用。

```toml
[[database]]
Name = 'Production'
Provider = 'postgres'
URL = 'postgres://${env:DB_USER}:${env:DB_PASSWORD}@localhost:5432/mydb'
```

```bash
export DB_USER=admin
export DB_PASSWORD=secret
lazysql
```

注意：未定义的环境变量会被替换为空字符串。

<!-- KEYBINDINGS -->

## 快捷键

### 自定义快捷键

你可以通过在 `config.toml` 文件中添加 `[keymap.<Group>]` 段来自定义快捷键。每个条目把一个命令名映射到一个键。

```toml
[keymap.Home]
SwitchToEditorView = "i"
Quit = "Esc"

[keymap.Tree]
GotoTop = "t"
Search = "Ctrl-F"

```

对于单字符键，直接使用该字符（例如 `"q"`、`"G"`、`"1"`、`"/"`）。对于特殊键，使用 [tcell 键名](https://github.com/gdamore/tcell/blob/v2.7.4/key.go#L83)（例如 `"Enter"`、`"Esc"`、`"Ctrl-S"`）。只支持 tcell 中定义的键名。

分组名不区分大小写（`Home`、`home` 和 `HOME` 都有效）。

可用分组：`Home`、`Connection`、`Tree`、`TreeFilter`、`Table`、`Editor`、`Sidebar`、`QueryPreview`、`QueryHistory`、`JSONViewer`。

### 默认快捷键

#### Home

| 默认键 | 命令 | 说明 |
| --- | --- | --- |
| L | MoveRight | 聚焦表格 |
| H | MoveLeft | 聚焦树 |
| Ctrl-E | SwitchToEditorView | 打开 SQL 编辑器 |
| Ctrl-S | Save | 执行待定更改 |
| q | Quit | 退出 |
| Backspace | SwitchToConnectionsView | 切换到连接列表 |
| ? | HelpPopup | 帮助 |
| Ctrl-P | SearchGlobal | 全局搜索 |
| Ctrl-_ | ToggleQueryHistory | 切换查询历史弹窗 |
| T | ToggleTree | 切换文件树 |

#### Connection

| 默认键 | 命令 | 说明 |
| --- | --- | --- |
| n | NewConnection | 新建数据库连接 |
| c | Connect | 连接数据库 |
| Enter | Connect | 连接数据库 |
| e | EditConnection | 编辑数据库连接 |
| d | DeleteConnection | 删除数据库连接 |
| q | Quit | 退出 |

#### Tree

| 默认键 | 命令 | 说明 |
| --- | --- | --- |
| g | GotoTop | 跳到顶部 |
| G | GotoBottom | 跳到底部 |
| Enter | Execute | 打开 |
| j | MoveDown | 向下 |
| Down | MoveDown | 向下 |
| Ctrl-U | PagePrev | 向上翻页 |
| Ctrl-D | PageNext | 向下翻页 |
| k | MoveUp | 向上 |
| Up | MoveUp | 向上 |
| / | Search | 搜索 |
| n | NextFoundNode | 跳到下一个匹配节点 |
| N | PreviousFoundNode | 跳到上一个匹配节点 |
| p | PreviousFoundNode | 跳到上一个匹配节点 |
| P | NextFoundNode | 跳到下一个匹配节点 |
| c | TreeCollapseAll | 全部折叠 |
| e | ExpandAll | 全部展开 |
| R | Refresh | 刷新树 |

#### Tree Filter

| 默认键 | 命令 | 说明 |
| --- | --- | --- |
| Esc | UnfocusTreeFilter | 取消聚焦树过滤器 |
| Enter | CommitTreeFilter | 提交树过滤器搜索 |

#### Table

| 默认键 | 命令 | 说明 |
| --- | --- | --- |
| / | Search | 搜索 |
| c | Edit | 修改单元格 |
| d | Delete | 删除行 |
| w | GotoNext | 跳到下一个单元格 |
| b | GotoPrev | 跳到上一个单元格 |
| $ | GotoEnd | 跳到最后一个单元格 |
| 0 | GotoStart | 跳到第一个单元格 |
| Enter | ForeignKeyJump | 从外键单元格跳转到被引用的行（参见 [外键跳转](#跳转到被引用的行外键跳转)） |
| y | Copy | 复制单元格值到剪贴板（如有标记行则复制标记行） |
| Y | CopyRowsAs | 将标记行（或当前行）复制为 SQL、Markdown、JSON、CSV |
| v | RowRange | 选择实时行范围；Esc 取消 |
| Space | RowSelect | 切换行或提交可视范围 |
| o | AppendNewRow | 追加新行 |
| O | DuplicateRow | 复制行 |
| J | SortDesc | 降序排序 |
| R | Refresh | 刷新当前表 |
| # | ExactCount | 计算或取消精确记录数 |
| K | SortAsc | 升序排序 |
| C | SetValue | 切换值菜单（NULL、EMPTY、DEFAULT） |
| [ | TabPrev | 切换到上一个标签页 |
| ] | TabNext | 切换到下一个标签页 |
| { | TabFirst | 切换到第一个标签页 |
| } | TabLast | 切换到最后一个标签页 |
| X | TabClose | 关闭标签页 |
| > | PageNext | 切换到下一页 |
| < | PagePrev | 切换到上一页 |
| 1 | RecordsMenu | 切换到记录菜单 |
| 2 | ColumnsMenu | 切换到列菜单 |
| 3 | ConstraintsMenu | 切换到约束菜单 |
| 4 | ForeignKeysMenu | 切换到外键菜单 |
| 5 | IndexesMenu | 切换到索引菜单 |
| S | ToggleSidebar | 切换侧边栏 |
| s | FocusSidebar | 聚焦侧边栏 |
| Z | ShowRowJSONViewer | 切换行的 JSON 查看器 |
| z | ShowCellJSONViewer | 切换单元格的 JSON 查看器 |
| f | ReverseForeignKeyJump | 选择引用当前行的表并打开过滤视图 |
| E | ExportCSV | 导出为 CSV |
| e | OpenCellInExternalEditor | 在外部编辑器中编辑单元格 |

> `Enter`（`ForeignKeyJump`）只在表视图的记录标签页（`1`）中生效。它属于 `Table` 分组，因此可以像其他快捷键一样重新映射：`[keymap.Table] ForeignKeyJump = "Ctrl-G"`。

#### Editor

| 默认键 | 命令 | 说明 |
| --- | --- | --- |
| Ctrl-R | Execute | 执行查询 |
| Esc | CancelQuery | 取消当前查询；否则保持编辑器正常的 Esc 行为 |
| Ctrl-Space | OpenInExternalEditor | 在外部编辑器中打开 |

lazysql 的专用编辑器可以通过 `$SQL_EDITOR` 设置。

### JSON Viewer

| 键 | 操作           |
| --- | ---------------- |
| w   | 切换自动换行 |
| y   | 复制到剪贴板 |
| z/Z | 关闭查看器     |

JSON 查看器可以在表格单元格上按 `z`（单元格）或 `Z`（行）打开。如果启用了 `EnterOpensJSONViewer`，在单元格上按 Enter 也会打开 JSON 查看器，除非该单元格支持 [外键跳转](#跳转到被引用的行外键跳转)（优先执行跳转）。

#### Sidebar

| 默认键 | 命令 | 说明 |
| --- | --- | --- |
| s | UnfocusSidebar | 聚焦表格 |
| S | ToggleSidebar | 切换侧边栏 |
| j | MoveDown | 聚焦下一个字段 |
| k | MoveUp | 聚焦上一个字段 |
| g | GotoStart | 聚焦第一个字段 |
| G | GotoEnd | 聚焦最后一个字段 |
| c | Edit | 编辑字段 |
| Enter | CommitEdit | 将修改加入待定更改 |
| Esc | DiscardEdit | 放弃修改 |
| C | SetValue | 切换值菜单（NULL、EMPTY、DEFAULT） |
| y | Copy | 复制值到剪贴板 |

#### Query Preview

| 默认键 | 命令 | 说明 |
| --- | --- | --- |
| Ctrl-S | Save | 执行查询 |
| q | Quit | 退出 |
| y | Copy | 复制查询到剪贴板 |
| d | Delete | 删除查询 |

#### Query History

| 默认键 | 命令 | 说明 |
| --- | --- | --- |
| s | Save | 保存查询 |
| d | Delete | 删除查询 |
| q | Quit | 退出 |
| y | Copy | 复制查询到剪贴板 |
| / | Search | 搜索 |
| Ctrl-_ | ToggleQueryHistory | 切换查询历史弹窗 |
| [ | TabPrev | 切换到上一个标签页 |
| ] | TabNext | 切换到下一个标签页 |

#### JSON Viewer

| 默认键 | 命令 | 说明 |
| --- | --- | --- |
| Z | ShowRowJSONViewer | 切换 JSON 查看器 |
| z | ShowCellJSONViewer | 切换 JSON 查看器 |
| y | Copy | 复制值到剪贴板 |
| w | ToggleJSONViewerWrap | 切换自动换行 |

### 外部编辑器

外部编辑器功能（SQL 编辑器中按 CTRL + Space，表格中按 `e`）使用以下环境变量决定使用哪个编辑器：

- SQL 编辑器：`$SQL_EDITOR` > `$EDITOR` > `$VISUAL` > `vi`
- 表格单元格：`$EDITOR` > `$VISUAL` > `vi`

支持带参数的编辑器命令，例如 `EDITOR="vim -u NONE"`。GUI 编辑器应使用阻塞到文件关闭的参数（例如 VS Code 用 `--wait`，Sublime Text 用 `-w`），否则 lazysql 会在你完成编辑之前就读取文件。

此功能仅适用于 Linux 和 macOS。

## 示例连接 URL

```
postgres://user:pass@localhost/dbname
pg://user:pass@localhost/dbname?sslmode=disable
mysql://user:pass@localhost/dbname
mysql:/var/run/mysqld/mysqld.sock
sqlserver://user:pass@remote-host.com/dbname
mssql://user:pass@remote-host.com/instance/dbname
ms://user:pass@remote-host.com:port/instance/dbname?keepAlive=10
oracle://user:pass@somehost.com/sid
sap://user:pass@localhost/dbname
file:myfile.sqlite3?loc=auto
/path/to/sqlite/file/test.db
odbc+postgres://user:pass@localhost:port/dbname?option1=
clickhouse://user:pass@localhost:9000/dbname
ch://user:pass@remote-host.com:9440/dbname?secure=true
clickhouse+http://user:pass@localhost:8123/dbname
clickhouse+https://user:pass@remote-host.com:8443/dbname
```

### ClickHouse 说明

- `clickhouse://`（别名 `ch://`）使用原生协议（默认端口 9000）；`clickhouse+http://` 和 `clickhouse+https://` 使用 HTTP 接口。查询参数传递给 [clickhouse-go](https://github.com/ClickHouse/clickhouse-go#dsn)，例如原生端口启用 TLS 用 `?secure=true`。
- Constraints 标签页显示表引擎及其分区、排序键和主键。Indexes 标签页显示数据跳过索引。ClickHouse 没有外键。
- 行编辑和删除以突变（mutation）方式执行（`ALTER TABLE ... UPDATE/DELETE`）并等待突变完成。它们只对支持突变的表有效（例如 MergeTree 系列），键列不能更新，ClickHouse 主键不唯一：与编辑行共享主键值的每一行都会被修改。待定更改不在事务中运行。

<!-- ROADMAP -->

## 路线图

- [ ] 支持 NoSQL 数据库
- [ ] 通过 TUI 创建列和索引
- [x] 表树输入过滤
- [x] 自定义快捷键
- [x] 在弹窗中显示快捷键
- [x] 重写行的 `create`、`update` 和 `delete` 逻辑

有关完整的功能建议（和已知问题）列表，请参见 [open issues](https://github.com/jorgerojas26/lazysql/issues)。

<p align="right">(<a href="#readme-top">返回顶部</a>)</p>

## 剪贴板支持

我们使用 [atotto/clipboard](https://github.com/atotto/clipboard?tab=readme-ov-file#clipboard-for-go) 复制到剪贴板。

平台：

- OSX
- Windows 7（其他 Windows 版本大概率可用）
- Linux、Unix（需要安装 'xclip' 或 'xsel' 命令）

<!-- CONTRIBUTING -->

## 参与贡献

欢迎贡献、提 issue 和提交 pull request！

<p align="right">(<a href="#readme-top">返回顶部</a>)</p>

<!-- LICENSE -->

## 许可证

基于 MIT 许可证分发。更多信息请参见 `LICENSE.txt`。

<p align="right">(<a href="#readme-top">返回顶部</a>)</p>

<!-- CONTACT -->

## 联系方式

Jorge Rojas - [LinkedIn](https://www.linkedin.com/in/jorgerojas26/) - jorgeluisrojasb@gmail.com

<p align="right">(<a href="#readme-top">返回顶部</a>)</p>

## 备选方案

- [Mitzasql](https://github.com/vladbalmos/mitzasql)
- [Gobang](https://github.com/TaKO8Ki/gobang)

<p align="right">(<a href="#readme-top">返回顶部</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->

[contributors-shield]: https://img.shields.io/github/contributors/jorgerojas26/lazysql?style=for-the-badge
[contributors-url]: https://github.com/jorgerojas26/lazysql/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/jorgerojas26/lazysql?style=for-the-badge
[forks-url]: https://github.com/jorgerojas26/lazysql/network/members
[stars-shield]: https://img.shields.io/github/stars/jorgerojas26/lazysql?style=for-the-badge
[stars-url]: https://github.com/jorgerojas26/lazysql/stargazers
[issues-shield]: https://img.shields.io/github/issues/jorgerojas26/lazysql?style=for-the-badge
[issues-url]: https://github.com/jorgerojas26/lazysql/issues
[license-shield]: https://img.shields.io/github/license/jorgerojas26/lazysql.svg?style=for-the-badge
[license-url]: https://github.com/jorgerojas26/lazysql/blob/main/LICENSE.txt
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://linkedin.com/in/jorgerojas26
[product-screenshot1]: images/lazysql-connection-selection.png
[product-screenshot2]: images/lazysql.png
[golang-shield]: https://img.shields.io/badge/Golang-gray?style=for-the-badge&logo=go
[tview-shield]: https://img.shields.io/badge/tview-gray?style=for-the-badge&logo=go

## Star 历史

<a href="https://star-history.dera.page/#jorgerojas26/lazysql&type=date&legend=top-left">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://star-history.dera.page/svg?repos=jorgerojas26/lazysql&type=date&theme=dark&legend=top-left" />
   <source media="(prefers-color-scheme: light)" srcset="https://star-history.dera.page/svg?repos=jorgerojas26/lazysql&type=date&legend=top-left" />
   <img alt="Star History Chart" src="https://star-history.dera.page/svg?repos=jorgerojas26/lazysql&type=date&legend=top-left" />
 </picture>
</a>
