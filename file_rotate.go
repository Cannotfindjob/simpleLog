package simpleLogger

import (
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Default configuration for rotation,
// included: maxlines、maxsize、max
var defaultRotate = Rotate{
	MaxLines:    10000,     // num:  10000
	MaxSize:     150 << 20, // size: 150MB
	MaxKeepDays: 7,         // days: 7days
}

var DefaultTimeFormat = "20060102150405"



type RotateOption func(*Rotate)

// WithMaxLinesOption is an optional function, the maximum
// number of rows of log files can be configured through
// the function option mode.
func WithMaxLinesOption(l int64) RotateOption {
	return func(o *Rotate) {
		o.MaxLines = l
	}
}

// WithMaxSizeOption is an optional function, the maximum
// size of the log file can be configured through the function
// option mode.
func WithMaxSizeOption(s int64) RotateOption {
	return func(o *Rotate) {
		o.MaxSize = s
	}
}

// WithMaxDaysOption is an optional function, the maximum
// number of days to save log files can be configured through
// the function option mode.
func WithMaxDaysOption(d int) RotateOption {
	return func(o *Rotate) {
		o.MaxKeepDays = d
	}
}

// DoRotate method implements the specific rotation logic:
// 1. Close the old log file handle first;
// 2. Rename the old log, for example: a.log is renamed to a.log.1;
// 3. Open and generate a new log file handle, The handle pointer is assigned to f;
// 4. If log compression is configured, the asynchronous compression function will
// be called, written to chan, and executed asynchronously.
func (f *FileObject) DoRotate() error {
	var err error

	// counter increment
	f.rotate.Count++

	// close old file handle
	f.file.Close()
	f.file = nil

	// time format
	format := time.Now().Format(DefaultTimeFormat)

	// Rename the log that will be rotated
	// For example: a.log will be renamed to a.log.1
	fName := f.path + "." + format + "_" + strconv.Itoa(f.rotate.Count)
	if err = os.Rename(f.path, fName); err != nil {
		return err
	}

	// generate a new file handle
	f.file, err = f.Open()
	if err != nil {
		return err
	}

	if f.isCompress {
		splice := "." + format + "_" + strconv.Itoa(f.rotate.Count) + ".zip"
		zipName := strings.Replace(f.path, filepath.Ext(f.path), splice, 1)

		select {
		case <-f.compress.ctx.Done():
			return nil
		case f.compress.taskQueue <- f.compress.DoCompress(zipName, path.Dir(f.path), []string{filepath.Base(fName)}):
		}
	}

	return nil
}

func (f *FileObject) RotateByLines() bool {
	return f.rotate.MaxLines > 0 && f.rotate.CurrentLine >= f.rotate.MaxLines
}

func (f *FileObject) RotateBySizes() bool {
	return f.rotate.MaxSize > 0 && f.rotate.CurrentSize >= f.rotate.MaxSize
}

func (f *FileObject) RotateByDaily() bool {
	t := f.rotate.CurrentTime
	tm := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).AddDate(0, 0, 1).Unix()

	return time.Now().Unix() > tm
}
