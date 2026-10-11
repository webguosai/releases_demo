# hello-cli

一个用来练习 GitHub Releases 打包的极简命令行程序。推送 `v*` 标签后，Actions 会交叉编译 7 个平台包，并上传到 [Releases](https://github.com/webguosai/releases_demo/releases)。

无参数时打印 `hello`。`hello-cli version` 打印版本号和 `系统/架构`，用来确认下载到的是哪一份包。

## 产物命名

标签 `v1.0.0` 会去掉前缀 `v`，得到版本号 `1.0.0`。文件名是：

```text
hello-cli-{版本号}-{系统}-{架构}.tar.gz
hello-cli-{版本号}-windows-{架构}.zip
```

一次发布包含：

- `hello-cli-1.0.0-darwin-amd64.tar.gz`
- `hello-cli-1.0.0-darwin-arm64.tar.gz`
- `hello-cli-1.0.0-linux-amd64.tar.gz`
- `hello-cli-1.0.0-linux-arm64.tar.gz`
- `hello-cli-1.0.0-linux-riscv64.tar.gz`
- `hello-cli-1.0.0-windows-amd64.zip`
- `hello-cli-1.0.0-windows-arm64.zip`
- `checksums.txt`（每个包一行 sha256）

Unix 用 `.tar.gz`，Windows 用 `.zip`。压缩包里只有一个文件：`hello-cli` 或 `hello-cli.exe`。

## 本地试跑

```powershell
go run .
go run . version
```

第二行在还没注入版本时会显示 `hello-cli dev windows/amd64`（系统以本机为准）。

## 本地交叉编译

Go 用 `GOOS` 和 `GOARCH` 决定目标系统。下面在 Windows 上编一份 Linux amd64，并按发布时的规则命名。

```powershell
$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"
go build -ldflags "-s -w -X main.version=1.0.0" -o hello-cli .
tar -czf hello-cli-1.0.0-linux-amd64.tar.gz hello-cli
```

`-X main.version=1.0.0` 把版本写进程序。文件名由四段拼出来：程序名、版本号、`GOOS`、`GOARCH`。Windows 目标要把输出改成 `hello-cli.exe`，再用 `Compress-Archive` 打成 `.zip`。

编完后清掉环境变量，避免后面的 `go run` 仍指向 Linux：

```powershell
Remove-Item Env:GOOS, Env:GOARCH, Env:CGO_ENABLED
```

## 发布到 GitHub Releases

远程仓库是 `https://github.com/webguosai/releases_demo.git`。工作流文件必须在被打标签的那一次提交里，否则推送标签时不会运行。

1. 提交并推送到 `main`（含 `.github/workflows/release.yml`）。
2. 在这个提交上打标签并推送标签：

```powershell
git tag v1.0.0
git push origin v1.0.0
```

推送 `v1.0.0` 后，Actions 会：

1. 按 7 组 `GOOS` / `GOARCH` 在 Ubuntu 上交叉编译（`CGO_ENABLED=0`）。
2. Linux、macOS 打成 `.tar.gz`，Windows 打成 `.zip`。
3. 汇总 `checksums.txt`。
4. 用该标签创建 Release，标题为 `v1.0.0`，并挂上全部附件。

到 [Releases](https://github.com/webguosai/releases_demo/releases) 能看到这 7 个包和 `checksums.txt`。页面上每个附件也会显示 sha256。

发下一个版本就再打一个新标签，例如 `v1.0.1`，然后 `git push origin v1.0.1`。不要移动已经发出去的旧标签。

## 下载后校验并运行

在 Release 页面下载对应系统的包，用 `checksums.txt` 里的 sha256 对照。

Windows：

```powershell
Get-FileHash .\hello-cli-1.0.0-windows-amd64.zip -Algorithm SHA256
Expand-Archive .\hello-cli-1.0.0-windows-amd64.zip -DestinationPath .
.\hello-cli.exe version
```

Linux 或 macOS：

```bash
sha256sum hello-cli-1.0.0-linux-amd64.tar.gz
tar -xzf hello-cli-1.0.0-linux-amd64.tar.gz
./hello-cli version
```

输出应类似 `hello-cli 1.0.0 linux/amd64`。版本号和 `系统/架构` 与文件名一致，就说明下到了对应的包。
