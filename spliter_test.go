package spliter_test

import (
	"fmt"
	"testing"
	"time"

	spliter "github.com/klarkxy/go-logfile-spliter"
)

func TestWrite(t *testing.T) {
	names := []string{}
	s := spliter.NewSpliter("@every 1s", func() string {
		now := time.Now()
		name := fmt.Sprintf("test/%d-%d-%d-%d-%d-%d.log", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
		names = append(names, name)
		return name
	})
	for i := 1; i <= 10; i++ {
		s.Write([]byte(fmt.Sprintf("%d\n", i)))
		time.Sleep(time.Second)
	}
}

func TestDoubleSplit(t *testing.T) {
	s := spliter.NewSpliter("0 1 0 0 0 0", func() string {
		name := "test/TestDoubleSplit.log"
		return name
	})
	for i := 1; i <= 10; i++ {
		s.Write([]byte(fmt.Sprintf("%d\n", i)))
		s.Split()
	}
}
