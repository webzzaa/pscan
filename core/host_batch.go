package core

import (
	"scanner/core/common"
	"scanner/core/common/parsers"
)

const maxHostBatchSize = 65536

func targetHostBatchSize(config *common.Config) int {
	size := parsers.DefaultHostBatchSize
	if config != nil && config.ThreadNum > 0 {
		threadWindow := config.ThreadNum * 8
		if threadWindow > size {
			size = threadWindow
		}
	}
	if size > maxHostBatchSize {
		return maxHostBatchSize
	}
	return size
}
