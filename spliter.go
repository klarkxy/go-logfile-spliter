package spliter

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/robfig/cron/v3"
)

type Spliter struct {
	// 文件名生成函数
	FilenameFunc func() string
	// 文件分离失败处理函数
	OnSplitError func(err error)
	// 文件句柄
	file *os.File
	// 文件锁
	lock sync.Mutex
}

var globalCron *cron.Cron
var once sync.Once

// 新建一个Spliter
//
// cronString: 见https://pkg.go.dev/github.com/robfig/cron/v3@v3.0.1
//
// filenameFunc: 文件名生成函数
func NewSpliter(cronString string, FilenameFunc func() string) *Spliter {
	s := &Spliter{
		FilenameFunc: FilenameFunc,
	}
	once.Do(func() {
		globalCron = cron.New(cron.WithSeconds())
		globalCron.Start()
	})
	globalCron.AddFunc(cronString, s.Split)
	s.Split()
	return s
}

// 默认文件分离失败处理函数
func defaultOnSplitError(err error) {
	fmt.Println(err)
}

// 设置文件分离失败处理函数
func (s *Spliter) SetOnSplitError(fn func(err error)) {
	s.OnSplitError = fn
}

// 文件分离
func (s *Spliter) Split() {
	s.lock.Lock()
	defer s.lock.Unlock()
	// 打开新文件
	filename := s.FilenameFunc()
	os.MkdirAll(filepath.Dir(filename), 0755)
	file, err := os.OpenFile(s.FilenameFunc(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		if s.OnSplitError != nil {
			s.OnSplitError(err)
		} else {
			defaultOnSplitError(err)
		}
		return
	}
	// 关闭旧文件
	if s.file != nil {
		s.file.Close()
	}
	s.file = file
}

// 写文件
func (s *Spliter) Write(data []byte) (n int, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.file.Write(data)
}
