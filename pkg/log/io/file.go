package io

import (
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

type DailyFileHandler struct {
	name string
	dir  string
	mtx  sync.Mutex
}

func NewDailyFileHandler() *DailyFileHandler {
	return &DailyFileHandler{
		name: "logs-%y-%m-%d.log",
		dir:  "./",
	}
}

func (h *DailyFileHandler) SetName(name string) *DailyFileHandler {
	h.name = name
	return h
}

func (h *DailyFileHandler) SetDirectory(dir string) *DailyFileHandler {
	h.dir = dir
	return h
}

func (h *DailyFileHandler) Write(p []byte) (n int, err error) {
	file, err := os.OpenFile(path.Join(h.dir, h.dailyName()), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	defer func() {
		file.Close()
	}()

	return file.Write(p)
}

func (h *DailyFileHandler) dailyName() string {
	now := time.Now()
	y, m, d := strconv.Itoa(now.Year()), strconv.Itoa(int(now.Month())), strconv.Itoa(now.Day())
	name := strings.Replace(h.name, "%y", y, -1)
	name = strings.Replace(name, "%m", m, -1)
	name = strings.Replace(name, "%d", d, -1)

	return name
}
