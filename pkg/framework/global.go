// Utils
package framework

import (
	"log"

	def "github.com/s3studio/cloud-bench-checker/pkg/definition"
)

var (
	_logger = log.Default()
	_opt    = def.ConfOption{
		PageSize: 50,
	}
)

// SetLogger: Set global logger
// @param: newLogger: New logger instance
func SetLogger(newLogger *log.Logger) {
	_logger = newLogger
}

func glog() *log.Logger {
	return _logger
}

// SetPageSize: Set PageSize option
// @param: pageSize: New value
func SetPageSize(pageSize int) {
	_opt.PageSize = pageSize
}
