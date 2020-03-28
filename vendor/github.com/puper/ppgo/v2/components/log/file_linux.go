package log

import (
	"os"
	"syscall"
	"time"
)

func FileCreateTime(fi os.FileInfo) time.Time {
	stat := fi.Sys().(*syscall.Stat_t)
	return time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec))
}
