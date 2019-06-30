package simpleLogger

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

// DeleteOld method implements the function of deleting the history
// log. First, it traverses all the log files in the log directory
// to obtain two types of log files:
// 1.zip archive log;
// 2.program.log.timeunix_num type non-archive log,
// and then judge whether the log file is satisfied. Delete the
// condition (has the maximum number of days saved), and delete the
// history log when it is satisfied. This logic is executed every
// 24 after the first execution.
func (f *FileObject) DeleteOld() {
	for {
		f := func() {
			// traverse the archived log files in the log directory and
			// delete the log files that exceed the maximum number of days
			// reserved.
			dir := filepath.Dir(f.path)
			os.Chdir(dir)

			filepath.Walk(dir, func(path string, file os.FileInfo, err error) error {
				if file == nil {
					return nil
				}

				if !file.IsDir() {
					name := file.Name()
					if strings.HasSuffix(name, ".zip") {
						timestamp := name[2:16]
						if f.isDelete(timestamp) {
							os.Remove(name)
						}
					} else if strings.Index(name, ".log.") == 1 {
						timestamp := name[6:20]
						if f.isDelete(timestamp) {
							os.Remove(name)
						}
					}
				}

				return nil
			})
		}
		f()

		now := time.Now()
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
	}
}

// Determine whether the current timestamp is greater than or equal to
// the log file creation time + maximum retention days * 86400.
func (f *FileObject) isDelete(timestamp string) bool {
	timeFormat, _ := time.ParseInLocation(DefaultTimeFormat, timestamp, time.Local)
	return time.Now().Unix() >= (timeFormat.Unix() + int64(f.rotate.MaxKeepDays * 86400))
}
