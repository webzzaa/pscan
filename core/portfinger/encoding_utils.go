package portfinger

import (
	"encoding/hex"
	"strconv"
)

// DecodePattern 解码匹配模式
func DecodePattern(s string) ([]byte, error) {
	b := []byte(s)
	var result []byte

	for i := 0; i < len(b); {
		if b[i] == '\\' && i+1 < len(b) {
			// 处理转义序列
			switch b[i+1] {
			case 'x':
				// 十六进制编码 \xNN
				if i+3 < len(b) {
					if hexStr := string(b[i+2 : i+4]); isValidHex(hexStr) {
						if decoded, err := hex.DecodeString(hexStr); err == nil {
							result = append(result, decoded...)
							i += 4
							continue
						}
					}
				}
			case 'a':
				result = append(result, '\a')
				i += 2
				continue
			case 'f':
				result = append(result, '\f')
				i += 2
				continue
			case 't':
				result = append(result, '\t')
				i += 2
				continue
			case 'n':
				result = append(result, '\n')
				i += 2
				continue
			case 'r':
				result = append(result, '\r')
				i += 2
				continue
			case 'v':
				result = append(result, '\v')
				i += 2
				continue
			case '\\':
				result = append(result, '\\')
				i += 2
				continue
			default:
				// 八进制编码 \NNN
				if i+1 < len(b) && b[i+1] >= '0' && b[i+1] <= '7' {
					octalStr := ""
					j := i + 1
					for j < len(b) && j < i+4 && b[j] >= '0' && b[j] <= '7' {
						octalStr += string(b[j])
						j++
					}
					// 使用16位解析避免int8溢出（\377=255超出int8范围）
					if octal, err := strconv.ParseInt(octalStr, 8, 16); err == nil && octal <= 255 {
						result = append(result, byte(octal))
						i = j
						continue
					}
				}
			}
		}

		// 普通字符
		result = append(result, b[i])
		i++
	}

	return result, nil
}

// DecodeData 解码探测数据
func DecodeData(s string) ([]byte, error) {
	// 移除首尾的分隔符
	if len(s) > 0 && (s[0] == '"' || s[0] == '\'') {
		s = s[1:]
	}
	if len(s) > 0 && (s[len(s)-1] == '"' || s[len(s)-1] == '\'') {
		s = s[:len(s)-1]
	}

	return DecodePattern(s)
}

// isValidHex 检查字符串是否为有效的十六进制
func isValidHex(s string) bool {
	for _, c := range s {
		if (c < '0' || c > '9') && (c < 'A' || c > 'F') && (c < 'a' || c > 'f') {
			return false
		}
	}
	return len(s) == 2
}
