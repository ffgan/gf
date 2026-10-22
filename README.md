一个用于快速展示当前系统信息的命令行工具，对标 neofetch / hyfetch。纯 Go 实现，单个二进制文件，默认在左侧绘制发行版 ASCII logo，右侧列出主机信息。

hyfetch 的旗帜配色、预设、`hyfetch.json` 和命令行参数可以直接用。没有 hyfetch 配置、也没传相关参数时，输出仍是原来的 neofetch 风格。

## 安装

[Releases](https://github.com/ffgan/gf/releases) 里按平台放了压缩包。Windows 是 `.zip`，其他系统是 `.tar.gz`。包里是可执行文件、`LICENSE` 和这份说明。同一个 release 里的 `gf_<version>_checksums.txt` 是 SHA-256。

这些包用当前还在维护的最老 Go 版本构建，版本号写在 [`.go-version`](.go-version)（现在是 Go 1.26）。二进制沿用这个工具链的系统要求，例如 macOS 12，以及 Linux ppc64 的 ELFv1。

目标来自该版本的 `go tool dist list`。32 位 ARM 会另外打出 ARMv5、ARMv6、ARMv7；MIPS 会同时打出 hardfloat 和 softfloat。文件名里的系统、架构和变体能直接对上：

- `linux_armv6` 给树莓派 Zero 这种 ARMv6，`linux_armv5` 给更老的设备，一般的 32 位 ARM 用 `linux_armv7`。FreeBSD、NetBSD、OpenBSD、Plan 9 也是同样三份。
- 没有硬件浮点的 MIPS 机器用 `*_softfloat`，有 FPU 的用 `*_hardfloat`。
- `android_arm64` 不依赖 NDK，可以放进 Termux。`android_386`、`android_amd64`、`android_armv7` 链接 Android 系统 libc（API 21）。
- `ios_arm64` 面向真机（最低 iOS 12），`ios_amd64` 面向 x86_64 模拟器。
- `js_wasm` 和 `wasip1_wasm` 里的文件名是 `gf.wasm`。

Go 1.26 已经去掉 `windows/arm`，也不再提供 `freebsd/riscv64`，所以发布里没有这两项。

## 构建

需要 Go 1.26 或更新的版本。

```bash
go build -o gf .
```

打出和 release 相同的一套包：

```bash
# android/386、android/amd64、android/arm 需要 ANDROID_NDK_HOME。
# ios/arm64 只能在装了 Xcode 的 macOS 上编。
VERSION=0.3.2 ANDROID_NDK_HOME=/opt/android-ndk bash scripts/build-release.sh
```

## 用法

```bash
gf
gf --stdout
gf --logo Arch
gf --disable disk network
gf --print_config > ~/.config/gf/config.conf

gf -p transgender          # 用跨性别旗给 logo 上色
gf -p '#ff0000,#00ff00,#0000ff' -m rgb
gf --distro Arch           # 指定发行版 logo，并按 hyfetch 配置上色
gf -c                      # 写入 ~/.config/hyfetch.json
gf -b fastfetch            # 把上好色的 logo 交给 fastfetch
gf --print-font-logo
gf --june
```

已安装 hyfetch 时，`~/.config/hyfetch.json`（或 `-C` 指定的文件）会被直接读取：预设、`8bit`/`rgb`、水平/垂直/随机上色、明度、`--args` 里的 neofetch 参数都会生效。预设表与 hyfetch 附带的 `presets.json` 相同。`-b` 的取值如果是 `ascii` 这类图片后端，仍然交给原来的 logo 逻辑；只有 `neofetch`、`fastfetch`、`macchina` 会切到 hyfetch 的后端。

不带参数时，`gf` 会依次查找这些配置：

- `$XDG_CONFIG_HOME/gf/config.conf`
- `~/.config/gf/config.conf`
- `~/.config/neofetch/config.conf`
- `/etc/gf/config.conf`
- `/etc/neofetch/config.conf`

配置格式和 neofetch 相同：`print_info()` 决定显示哪些行，其余是 `key="value"`。用 `--print_config` 可以打印一份当前默认配置。

信息行里暂时检测不到的项目会跳过，不会占一行空输出。`--v` 会把跳过的项目写到标准错误。

## 默认会显示的内容

用户名和主机名、操作系统、机型、内核、开机时间、软件包、Shell、分辨率、桌面环境和窗口管理器、主题、CPU、内存、磁盘、本地 IP、网卡，以及两行色块。桌面、电池、蓝牙这些只在检测得到时出现。

图片 logo 会在安装了 `chafa`、`catimg` 或 `viu` 时使用，例如 `gf --source ./logo.png`。没有这些工具时回退到发行版 ASCII。
