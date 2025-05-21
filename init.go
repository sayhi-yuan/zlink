package zlink

import (
	"github.com/sayhi-yuan/zlink/logger"
)

type Options struct {
	Logger logger.Logger
}

func InitZLink(options *Options) {
	logger.Log = options.Logger
}
