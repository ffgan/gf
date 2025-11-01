package configs

type Config struct {
	Title
	Kernel
	Distro
	Uptime
	Memory
	Packages
	Shell
	Editor
	CPU
	GPU
	Resolution
	GTK
	IPAddress
	DE
	Disk
	Song
	TextColor
	TextOptions
	ColorBlocks
	ProgressBar
	BackendSettings
	ASCIIOptions
	ImageOptions
	MiscOptions
}

type Title struct {
	TitleFqdn string `gf:"title_fqdn"`
}

type Kernel struct {
	KernelShorthand string `gf:"kernel_shorthand"`
}

type Distro struct {
	DistroShorthand string `gf:"distro_shorthand"`
	OSArch          string `gf:"os_arch"`
}

type Uptime struct {
	UptimeShorthand string `gf:"uptime_shorthand"`
}

type Memory struct {
	MemoryPercent string `gf:"memory_percent"`
	MemoryUnit    string `gf:"memory_unit"`
	MemPrecision  string `gf:"mem_precision"`
}

type Packages struct {
	PackageManagers string `gf:"package_managers"`
	PackageSeparate string `gf:"package_separate"`
	PackageMinimal  string `gf:"package_minimal"`
}

type Shell struct {
	ShellPath    string `gf:"shell_path"`
	ShellVersion string `gf:"shell_version"`
}

type Editor struct {
	EditorPath    string `gf:"editor_path"`
	EditorVersion string `gf:"editor_version"`
}

type CPU struct {
	SpeedType      string `gf:"speed_type"`
	SpeedShorthand string `gf:"speed_shorthand"`
	CPUBrand       string `gf:"cpu_brand"`
	CPUSpeed       string `gf:"cpu_speed"`
	CPUCores       string `gf:"cpu_cores"`
	CPUTemp        string `gf:"cpu_temp"`
}

type GPU struct {
	GPUBrand string `gf:"gpu_brand"`
	GPUType  string `gf:"gpu_type"`
}

type Resolution struct {
	RefreshRate string `gf:"refresh_rate"`
}

type GTK struct {
	GTKShorthand string `gf:"gtk_shorthand"`
	GTK2         string `gf:"gtk2"`
	GTK3         string `gf:"gtk3"`
	QT           string `gf:"qt"`
}

type IPAddress struct {
	PublicIPHost     string `gf:"public_ip_host"`
	PublicIPTimeout  string `gf:"public_ip_timeout"`
	LocalIPInterface string `gf:"local_ip_interface"`
}

type DE struct {
	DeVersion string `gf:"de_version"`
}

type Disk struct {
	DiskShow     string `gf:"disk_show"`
	DiskSubtitle string `gf:"disk_subtitle"`
	DiskPercent  string `gf:"disk_percent"`
}

type Song struct {
	MusicPlayer   string `gf:"music_player"`
	SongFormat    string `gf:"song_format"`
	SongShorthand string `gf:"song_shorthand"`
	MPCArgs       string `gf:"mpc_args"`
}

type TextColor struct {
	Colors           string `gf:"colors"`
	Bold             string `gf:"bold"`
	UnderlineEnabled string `gf:"underline_enabled"`
	UnderlineChar    string `gf:"underline_char"`
	Separator        string `gf:"separator"`
}

type TextOptions struct {
	Bold             string `gf:"bold"`
	UnderlineEnabled string `gf:"underline_enabled"`
	UnderlineChar    string `gf:"underline_char"`
	Separator        string `gf:"separator"`
}

type ColorBlocks struct {
	BlockRange  string `gf:"block_range"`
	ColorBlocks string `gf:"color_blocks"`
	BlockWidth  string `gf:"block_width"`
	BlockHeight string `gf:"block_height"`
	ColOffset   string `gf:"col_offset"`
}

type ProgressBar struct {
	BarCharElapsed  string `gf:"bar_char_elapsed"`
	BarCharTotal    string `gf:"bar_char_total"`
	BarBorder       string `gf:"bar_border"`
	BarLength       string `gf:"bar_length"`
	BarColorElapsed string `gf:"bar_color_elapsed"`
	BarColorTotal   string `gf:"bar_color_total"`
	MemoryDisplay   string `gf:"memory_display"`
	BatteryDisplay  string `gf:"battery_display"`
	DiskDisplay     string `gf:"disk_display"`
}

type BackendSettings struct {
	ImageBackend string `gf:"image_backend"`
	ImageSource  string `gf:"image_source"`
}

type ASCIIOptions struct {
	ASCIIDistro string `gf:"ascii_distro"`
	ASCIIColors string `gf:"ascii_colors"`
	ASCIIBold   string `gf:"ascii_bold"`
}

type ImageOptions struct {
	ImageLoop       string `gf:"image_loop"`
	ThumbnailDir    string `gf:"thumbnail_dir"`
	CropMode        string `gf:"crop_mode"`
	CropOffset      string `gf:"crop_offset"`
	ImageSize       string `gf:"image_size"`
	CatimgSize      string `gf:"catimg_size"`
	Gap             string `gf:"gap"`
	Yoffset         string `gf:"yoffset"`
	Xoffset         string `gf:"xoffset"`
	BackgroundColor string `gf:"background_color"`
}

type MiscOptions struct {
	Stdout string `gf:"stdout"`
}
