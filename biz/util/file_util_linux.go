package util

func GetFileCreatedTime(file os.FileInfo) time.Time {
	stat := file.Sys().(*syscall.Stat_t)
	ctim := stat.Ctim
	return time.Unix(ctim.Sec, ctim.Nsec)
}
