# CronPanel — Crontab 管理面板

现代化的 Linux Crontab 定时任务 Web 管理工具，单二进制文件部署，无需任何依赖。

---

## 功能特性

- 📋 **查看任务** — 实时列出所有 crontab 任务，支持启用/禁用状态显示
- ➕ **添加任务** — 引导式多步骤表单，支持多种时间模式
- ✏️ **编辑任务** — 在线修改已有任务的时间、频率与命令
- 🗑️ **删除任务** — 一键删除定时任务
- ⏸️ **启用/停用** — 无需删除即可暂时禁用任务（通过注释行实现）
- 📝 **三种命令模式** — 直接命令、指定脚本路径、在线编写 Shell 脚本
- 🔐 **基本认证** — 可选的用户名/密码访问保护
- 🌐 **现代化 UI** — 深色主题，响应式设计，无需安装前端依赖

---

## 时间模式说明

添加/编辑任务时支持以下调度模式：

| 模式 | 说明 | 示例 Cron |
|------|------|-----------|
| 每天 | 指定时分，每天执行 | `30 9 * * *` |
| 每 N 天 | 每隔 N 天执行一次 | `0 8 */3 * *` |
| 每周 | 指定星期几执行 | `0 10 * * 1` |
| 每月 | 指定每月几号执行 | `0 6 15 * *` |
| 自定义 | 手动输入任意 Cron 表达式 | `*/5 * * * *` |

---

## 命令模式说明

| 模式 | 说明 |
|------|------|
| 直接命令 | 直接输入 shell 命令，例如 `/usr/bin/python3 /opt/job.py` |
| 脚本路径 | 指定服务器上已有的 `.sh` 文件路径，将以 `/bin/bash` 执行 |
| 编写脚本 | 在线编写 Shell 脚本内容，保存时自动生成 `.sh` 文件并写入 crontab |

> **关于「编写脚本」模式的当前行为：**
>
> 选择「编写脚本」模式后，脚本内容会被保存到 `/tmp/crontab-manager-scripts/` 目录下，文件名格式为 `script_<时间戳纳秒>.sh`，crontab 中记录的命令为该文件的绝对路径（例如 `/bin/bash /tmp/crontab-manager-scripts/script_1771532433851647800.sh`）。
>
> ⚠️ **已知限制与待改进项（供后续开发参考）：**
>
> 1. **脚本存储路径**：脚本目前保存在系统临时目录 `/tmp`，该目录在重启后通常会被清空，导致定时任务执行失败。建议改为保存在程序二进制文件所在目录下的专用子目录（如 `./cronpanel-scripts/`），以保证持久性。
>
> 2. **编辑时不显示脚本内容**：编辑一个由「编写脚本」模式创建的任务时，界面上仅能看到脚本的文件路径，无法查看或修改脚本的实际内容。建议在编辑时自动读取对应 `.sh` 文件的内容并回填到脚本编辑器中。
>
> 3. **添加任务分步骤表单**：当前添加任务采用四步向导（时间 → 频率 → 命令 → 确认），信息分散。建议将时间、频率、命令合并在同一页面从上至下展示，右下角改为「确认」按钮，提升操作效率。

---

## 编译

需要 Go 1.18 或更高版本。

### 快速编译

```bash
chmod +x build.sh
./build.sh
```

编译产物输出至 `./dist/` 目录。

### 手动编译

```bash
# Linux amd64
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/cronpanel-linux-amd64 .

# Linux arm64
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o dist/cronpanel-linux-arm64 .
```

---

## 运行

```bash
# 赋予执行权限
chmod +x cronpanel-linux-amd64

# 默认端口 8899，无认证
./cronpanel-linux-amd64

# 自定义端口
PORT=9090 ./cronpanel-linux-amd64

# 启用基本认证（用户名:密码）
./cronpanel-linux-amd64 --auth admin:yourpassword

# 组合使用
PORT=9090 ./cronpanel-linux-amd64 --auth admin:yourpassword
```

打开浏览器访问：`http://localhost:8899`（或自定义端口）

---

## 注意事项

- 程序需要有权限执行 `crontab` 命令，建议以实际需要管理 crontab 的用户身份运行。
- 「编写脚本」模式生成的 `.sh` 文件当前保存在 `/tmp/crontab-manager-scripts/`，**重启后可能丢失**，请注意备份或参考上文改进建议。
- 启用认证后，Session 有效期为 24 小时，过期后需重新登录。
- 建议在内网或通过 SSH 端口转发使用，避免直接暴露在公网（尤其是未启用认证时）。

---

## 文件说明

```
CronPanel-main/
├── main.go       # 主程序：HTTP 服务器、Crontab 读写、脚本管理
├── html.go       # 前端：HTML / CSS / JavaScript（嵌入在 Go 中，单文件部署）
├── go.mod        # Go 模块定义
├── build.sh      # 一键交叉编译脚本
└── README.md     # 本说明文档
```

---

## License

MIT
