package portfinger

import (
	"regexp"
	"strconv"
	"strings"
)

// 预编译正则表达式
var (
	whitespaceRegex = regexp.MustCompile(`\s+`)

	// 版本信息字段解析正则 - 支持斜线和竖线两种分隔符
	fieldRegexes = map[string][]*regexp.Regexp{
		" p": {regexp.MustCompile(` p/([^/]*)/`), regexp.MustCompile(` p\|([^|]*)\|`)},
		" v": {regexp.MustCompile(` v/([^/]*)/`), regexp.MustCompile(` v\|([^|]*)\|`)},
		" i": {regexp.MustCompile(` i/([^/]*)/`), regexp.MustCompile(` i\|([^|]*)\|`)},
		" h": {regexp.MustCompile(` h/([^/]*)/`), regexp.MustCompile(` h\|([^|]*)\|`)},
		" o": {regexp.MustCompile(` o/([^/]*)/`), regexp.MustCompile(` o\|([^|]*)\|`)},
		" d": {regexp.MustCompile(` d/([^/]*)/`), regexp.MustCompile(` d\|([^|]*)\|`)},
	}

	// CPE解析正则
	cpeRegexSlash = regexp.MustCompile(`cpe:/([^/]*)`)
	cpeRegexPipe  = regexp.MustCompile(`cpe:\|([^|]*)`)
)

// ParseVersionInfo 解析版本信息并返回额外信息结构
func (m *Match) ParseVersionInfo(response []byte) Extras {
	var extras = Extras{}

	// 替换版本信息中的占位符（如 $1, $2 等）
	versionInfo := m.VersionInfo
	if len(m.FoundItems) > 0 {
		replacements := make([]string, 0, len(m.FoundItems)*2)
		for i, value := range m.FoundItems {
			replacements = append(replacements, "$"+strconv.Itoa(i+1), value)
		}
		versionInfo = strings.NewReplacer(replacements...).Replace(versionInfo)
	}

	// 定义解析函数 - 使用预编译正则
	parseField := func(field string) string {
		regexes, ok := fieldRegexes[field]
		if !ok || !strings.Contains(versionInfo, field) {
			return ""
		}
		for _, regex := range regexes {
			if matches := regex.FindStringSubmatch(versionInfo); len(matches) > 1 {
				return matches[1]
			}
		}
		return ""
	}

	// 解析各个字段
	extras.VendorProduct = parseField(" p")
	extras.Version = parseField(" v")
	extras.Info = parseField(" i")
	extras.Hostname = parseField(" h")
	extras.OperatingSystem = parseField(" o")
	extras.DeviceType = parseField(" d")

	// 特殊处理CPE - 使用预编译正则
	if strings.Contains(versionInfo, " cpe:/") || strings.Contains(versionInfo, " cpe:|") {
		for _, regex := range []*regexp.Regexp{cpeRegexSlash, cpeRegexPipe} {
			if matches := regex.FindStringSubmatch(versionInfo); len(matches) > 1 {
				extras.CPE = matches[1]
				break
			}
		}
	}

	return extras
}

// ToMap 将 Extras 转换为 map[string]string
func (e *Extras) ToMap() map[string]string {
	result := make(map[string]string)

	// 定义字段映射
	fields := map[string]string{
		"vendor_product": e.VendorProduct,
		"version":        e.Version,
		"info":           e.Info,
		"hostname":       e.Hostname,
		"os":             e.OperatingSystem,
		"device_type":    e.DeviceType,
		"cpe":            e.CPE,
	}

	// 添加非空字段到结果map
	for key, value := range fields {
		if value != "" {
			result[key] = value
		}
	}

	return result
}

// TrimBanner 清理横幅数据，移除不可打印字符
func TrimBanner(banner string) string {
	// 移除开头和结尾的空白字符
	banner = strings.TrimSpace(banner)

	// 移除控制字符，但保留换行符和制表符
	var result strings.Builder
	for _, r := range banner {
		if r >= 32 && r <= 126 { // 可打印ASCII字符
			result.WriteRune(r)
		} else if r == '\n' || r == '\t' { // 保留换行符和制表符
			result.WriteRune(r)
		} else {
			result.WriteRune(' ') // 其他控制字符替换为空格
		}
	}

	// 压缩多个连续空格为单个空格
	resultStr := result.String()
	resultStr = whitespaceRegex.ReplaceAllString(resultStr, " ")

	return strings.TrimSpace(resultStr)
}
