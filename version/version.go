package version

// 版本信息，由构建时通过 ldflags 注入
var (
	Version   = "dev"     // 版本号，默认为 dev
	BuildTime = "unknown" // 构建时间
)

// GetVersion 获取版本信息
func GetVersion() (string, string) {
	return Version, BuildTime
}
