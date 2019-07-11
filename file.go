package simpleLogger

import (
	"archive/zip"
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type FileObject struct {
	sync.RWMutex
	file         *os.File
	Filepath     string     `json:"filepath"`
	Perm         string     `json:"perm"`
	Rotate       bool       `json:"rotate"`
	Compress     bool       `json:"compress"`

	Count         int
	MaxLines      int64     `json:"max_lines"`
	CurrentLine   int64
	MaxSize       int64     `json:"max_size"`
	CurrentSize   int64
	MaxKeepDays   int       `json:"max_keep_days"`
	CurrentTime   time.Time

	taskQueue     chan func() error
	ctx           context.Context
	cancel        context.CancelFunc
}

func NewFileObject() Logger {
	obj := new(FileObject)
	obj.Perm = "0660"
	obj.Filepath = "./Simple.log"
	obj.Rotate = true
	obj.Compress = true

	obj.MaxLines = 100000
	obj.MaxSize = 500 << 20
	obj.MaxKeepDays = 7

	obj.taskQueue = make(chan func() error, 20)
	obj.ctx, obj.cancel = context.WithCancel(context.Background())

	go obj.TaskListen()
	go obj.DeleteOld()


	return obj
}

func (f *FileObject) Init(jsonConfig string) error {

	err := json.Unmarshal([]byte(jsonConfig), f)
	if err != nil {
		return err
	}

	if len(f.Filepath) == 0 {
		return errors.New("SimpleLog: jsonconfig must have filepath")
	}

	var file *os.File
	if file, err = f.Open(); err != nil {
		return err
	}
	if f.file != nil {
		_:f.file.Close()
	}
	f.file = file

	if err = f.initStat(); err != nil {
		return err
	}

	return nil
}

func (f *FileObject)Write(p []byte) (n int, err error) {

	if len(p) == 0 {
		return len(p), nil
	}

	if f.Rotate {
		if f.rotateByLines() {
			f.Lock()
			f.CurrentLine = 0
			f.Unlock()
			f.DoRotate()
		}

		if f.rotateBySizes() {
			f.Lock()
			f.CurrentSize = 0
			f.Unlock()
			f.DoRotate()
		}

		if f.rotateByDaily() {
			f.Lock()
			f.CurrentTime = time.Now()
			f.Unlock()
			f.DoRotate()
		}
	}

	f.Lock()
	_, err = f.file.Write(p)
	f.Unlock()

	if err != nil {
		 _:fmt.Fprintln(os.Stderr, "SimpleLog: file write", err)
	} else {
		atomic.AddInt64(&f.CurrentLine, 1)
		atomic.AddInt64(&f.CurrentSize, int64(len(p)))
	}

	return len(p),nil
}

func (f *FileObject) Flush() {
	if err := f.file.Sync(); err != nil {
		_:fmt.Fprintln(os.Stderr, "SimpleLog: file flush", err)
		return
	}
}

func (f *FileObject) Close() {
	_:f.file.Close()
	f.cancel()
}

func (f *FileObject) Open() (*os.File, error) {

	if f.file == nil {
		perm, err := strconv.ParseInt(f.Perm, 8, 64)
		if err != nil {
			return nil, err
		}
		pathSplit := path.Dir(f.Filepath)
		_, err = os.Stat(pathSplit)
		if err != nil && os.IsNotExist(err) {
			if err := f.Create(pathSplit, os.FileMode(perm)); err != nil {
				return nil, err
			}
		}
		fd, err := os.OpenFile(f.Filepath, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.FileMode(perm))
		if err != nil {
			return nil, err
		}
		return fd, nil
	}
	return f.file, nil
}

func (f *FileObject) Create(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (f *FileObject) initStat() error {
	info, err := f.file.Stat()
	if err != nil {
		return err
	}
	f.CurrentSize = info.Size()
	f.CurrentLine = f.initLine()
	f.CurrentTime = time.Now()

	return nil
}

func (f *FileObject) initLine() int64 {
	var line int64

	scanner := bufio.NewScanner(f.file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line++
	}

	if err := scanner.Err(); err != nil {
		_:fmt.Fprintln(os.Stderr, "SimpleLog: ile read err", err)
	}

	return line
}

func (f *FileObject) DoRotate() {
	var err error
	f.Count++

	format := time.Now().Format("20060102150405")

	fName := f.Filepath + "." + format + "_" + strconv.Itoa(f.Count)
	if err = os.Rename(f.Filepath, fName); err != nil {
		_:fmt.Fprintln(os.Stderr, "SimpleLog: file rename failed: ", err)
		return
	}

	tmpfile := f.file
	f.file = nil

	f.file, err = f.Open()
	if err != nil {
		_:fmt.Fprintln(os.Stderr, "SimpleLog: new file handle failed to open: ", err)
		f.file = tmpfile
		if err = os.Rename(fName, f.Filepath); err != nil {
			_:fmt.Fprintln(os.Stderr, "SimpleLog: file rename failed: ", err)
			return
		}
		f.Count--
		return
	}
	_:tmpfile.Close()

	if f.Compress {
		splice := "." + format + "_" + strconv.Itoa(f.Count) + ".zip"
		zipName := strings.Replace(f.Filepath, filepath.Ext(f.Filepath), splice, 1)
		f.taskQueue <- f.DoCompress(zipName, path.Dir(f.Filepath), []string{filepath.Base(fName)})
	}
}

func (f *FileObject) rotateByLines() bool {
	return f.MaxLines > 0 && f.CurrentLine >= f.MaxLines
}

func (f *FileObject) rotateBySizes() bool {
	return f.MaxSize > 0 && f.CurrentSize >= f.MaxSize
}

func (f *FileObject) rotateByDaily() bool {
	t := f.CurrentTime
	tm := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).AddDate(0, 0, 1).Unix()

	return time.Now().Unix() > tm
}

func (f *FileObject) DoCompress(zipName string, path string, sources []string) func() error {
	return func() error {
		err := Compress(zipName, sources)
		if err == nil {
			return os.Remove(path + "/" + sources[0])
		}

		return err
	}
}

func (f *FileObject) TaskListen() {
	for {
		select {
		case fn, ok := <-f.taskQueue:
			if !ok {
				return
			}

			if err := fn(); err != nil {
				_:fmt.Fprintln(os.Stderr, "SimpleLog: log compression error: ", err)
			}
		case <-f.ctx.Done():
			close(f.taskQueue)
			return
		}
	}
}


func Compress(destination string, sources []string) error {
	pwd := path.Dir(destination)
	if err := os.Chdir(pwd); err != nil {
		return err
	}

	out, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("SimpleLog: error creating %s: %v", destination, err)
	}
	defer out.Close()

	w := zip.NewWriter(out)
	err = zipFile(w, sources[0])
	if err != nil {
		fmt.Println(err)
	    _:w.Close()
		return err
	}

	return w.Close()
}

func zipFile(w *zip.Writer, file string) error {
	info, err := os.Stat(file)
	if err != nil {
		return fmt.Errorf("SimpleLog: %s: stat: %v", file, err)
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("SimpleLog: %s: getting header: %v", file, err)
	}
	header.Name = file
	header.Method = zip.Deflate

	writer, err := w.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("SimpleLog: %s: making header: %v", file, err)
	}

	if info.IsDir() {
		return nil
	}

	if header.Mode().IsRegular() {
		filefd, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("SimpleLog: %s: opening: %v", file, err)
		}
		defer filefd.Close()

		_, err = io.CopyN(writer, filefd, info.Size())
		if err != nil && err != io.EOF {
			return fmt.Errorf("SimpleLog: %s: copying contents: %v", file, err)
		}
	}

	return nil
}

func (f *FileObject) DeleteOld() {
	for {

		dir := filepath.Dir(f.Filepath)
		_:os.Chdir(dir)
		_:filepath.Walk(dir, func(path string, file os.FileInfo, err error) error {
			if file == nil {
				return nil
			}
			if !file.IsDir() {
				name := file.Name()
				if strings.HasSuffix(name, ".zip") {
					timestamp := name[2:16]
					if f.isDelete(timestamp) {
						_:os.Remove(name)
					}
				} else if strings.Index(name, ".log.") == 1 {
					timestamp := name[6:20]
					if f.isDelete(timestamp) {
						_:os.Remove(name)
					}
				}
			}

			return nil
			})

		now := time.Now()
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		fmt.Println(next.Sub(now))
		t := time.NewTimer(next.Sub(now))
		<-t.C
	}
}

func (f *FileObject) isDelete(timestamp string) bool {
	timeFormat, _ := time.ParseInLocation("20060102150405", timestamp, time.Local)
	return time.Now().Unix() >= (timeFormat.Unix() + int64(f.MaxKeepDays * 86400))
}


func init() {
	Register(AdapterFile, NewFileObject)
}
