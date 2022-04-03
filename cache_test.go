package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"testing"
)

func TestCache(t *testing.T) {
	//fullFilename := "/Users/itfanr/Documents/test.txt"
	//fmt.Println("fullFilename =", fullFilename)
	//var filenameWithSuffix string
	//filenameWithSuffix = path.Base(fullFilename)
	//fmt.Println("filenameWithSuffix =", filenameWithSuffix)
	//var fileSuffix string
	//fileSuffix = path.Ext(filenameWithSuffix)
	//fmt.Println("fileSuffix =", fileSuffix)
	//
	//var filenameOnly string
	//filenameOnly = strings.TrimSuffix(filenameWithSuffix, fileSuffix)
	//fmt.Println("filenameOnly =", filenameOnly)

	serverRoot := "D:/Go_projects/src/github.com/bingerambo/go-file-json-server/tmp"
	cc := NewCache(serverRoot)
	cc.Boot()

	fmt.Println(cc.data.Set())

	select {}
}

func TestCacheServer(t *testing.T) {
	logger = logrus.New()

	bindAddress := "0.0.0.0"
	listenPort := 24556
	serverRoot := "D:/Go_projects/src/github.com/bingerambo/go-file-json-server/tmp"
	corsEnabled := false

	cc := NewCache(serverRoot)

	cc.Boot()

	jserver := NewJsonServer(serverRoot, corsEnabled)
	jserver.Start()

	fmt.Println(cc.data.Set())

	errors := make(chan error)

	go func() {
		logger.WithFields(logrus.Fields{
			"ip":   bindAddress,
			"port": listenPort,
			"root": serverRoot,
			"cors": corsEnabled,
		}).Info("start listening")

		/*
		http://127.0.0.1:24556/mock?name=sample
		http://127.0.0.1:24556/mock?name=hello
		*/
		if err := http.ListenAndServe(fmt.Sprintf("%s:%d", bindAddress, listenPort), nil); err != nil {
			errors <- err
		}
	}()

	select {}
}
