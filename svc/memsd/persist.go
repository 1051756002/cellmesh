package main

import (
	"github.com/davyxu/cellmesh/svc/memsd/model"
	"github.com/davyxu/ulog"
	"os"
	"time"
)

func loadPersistFile(fileName string) {

	fileHandle, err := os.OpenFile(fileName, os.O_RDONLY, 0666)

	// 可能文件不存在，忽略
	if err != nil {
		return
	}

	ulog.Infoln("Load values...")

	err = model.LoadValue(fileHandle)
	if err != nil {
		ulog.Errorf("load values failed: %s %s", fileName, err.Error())
		return
	}

	ulog.Infof("Load %d values", model.ValueCount())
}

func startPersistCheck(fileName string) {

	ticker := time.NewTicker(time.Second * 20)

	for {

		<-ticker.C

		// 与收发在一个队列中，保证无锁
		model.Queue.Post(func() {

			if model.ValueDirty {

				ulog.Infoln("Save values...")

				fileHandle, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0666)
				if err != nil {
					ulog.Errorf("save persist file failed: %s %s", fileName, err.Error())
					return
				}

				valuesSaved, err := model.SaveValue(fileHandle)

				if err != nil {
					ulog.Errorf("save values failed: %s %s", fileName, err.Error())
					return
				}

				if valuesSaved > 0 {
					ulog.Infof("Save %d values", valuesSaved)
				}

				model.ValueDirty = false

			}

		})

	}

}
