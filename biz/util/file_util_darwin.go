package util

import (
	"os"
	"syscall"
	"time"
)

func GetFileCreatedTime(file os.FileInfo) time.Time {
	stat := file.Sys().(*syscall.Stat_t)
	birth := stat.Birthtimespec
	return time.Unix(birth.Sec, birth.Nsec)
}
