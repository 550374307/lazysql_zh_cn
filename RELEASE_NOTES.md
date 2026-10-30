# LazySQL I18N 汉化版 v0.7.2-zh

> 本 fork 基于上游 [LazySQL](https://github.com/jorgerojas26/lazysql)（jorgerojas26），主线为 I18N 国际化翻译（默认简体中文、中英切换、自定义语言文件）；本次按需求新增 MySQL `--login-path` 支持（见下），未修改其他上游功能逻辑。

## 本次新增：MySQL --login-path 支持（免明文密码）

LazySQL 现在可以直接读取 `mysql_config_editor` 加密保存的登录文件（`~/.mylogin.cnf`，Windows 为 `%APPDATA%\MySQL\.mylogin.cnf`），数据库密码无需再以明文写入 `.lazysql.toml`。

- **配置引用**：`[[database]]` 条目新增 `LoginPath` 字段，留空的 `user/password/host/port` 从加密文件自动补齐；显式配置的字段优先（与 `mysql` 客户端行为一致）。
- **启动自动发现（免配置）**：启动时自动扫描登录文件，其中**每一个** login-path 条目直接显示为一条 MySQL 连接，选中即可连接，无需手写 `[[database]]`。
- **自动条目不落盘**：自动发现的连接仅存在于会话内存，不会被写回配置文件；配置中已有同名 `LoginPath` 时以显式条目为准。
- **细节**：缺 host 默认 `localhost`、缺端口默认 `3306`；带 `socket` 且无 host 时走 unix socket 连接；登录文件缺失/损坏时静默跳过，不影响启动。
- 实现为纯 Go（`internal/mylogin` 包），AES-128-ECB + PKCS#7、20 字节 key XOR 折叠算法与 MySQL 8.0 官方源码（`my_aes.cc` / `mysql_config_editor.cc`）逐行核对一致，并通过独立实现交叉验证。
- 新增 11 个测试（mylogin 解密 roundtrip / deriveKey 向量 / 补齐 / 显式优先 / 缺失报错 / MySQL DSN / socket / 自动发现 / 去重 / 不落盘等）。

## 使用示例

```bash
# 1. 生成加密凭据（密码不进 shell 历史）
mysql_config_editor set --login-path=prod --host=127.0.0.1 --port=3306 --user=root --password

# 2. 直接启动 lazysql：连接选择器中会出现 "prod"，选中即可连接
#    （也支持在 .lazysql.toml 显式引用：LoginPath = "prod"）
```

## 历史修复（v0.7.1-zh）

- 补全框中文乱码（字节索引 → `[]rune` 索引渲染）。
- 回车 / Tab 无法选中补全候选；Tab 与 Enter 行为一致。
- 表限定补全保留表名前缀（`tenant_llm.` + Enter → `tenant_llm.id`）。
