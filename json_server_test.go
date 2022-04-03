package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"testing"
)



func TestSimpleServer(t *testing.T) {

	logger = logrus.New()

	bindAddress := "0.0.0.0"
	listenPort := 24556
	serverRoot := "D:/Go_projects/src/github.com/bingerambo/go-file-json-server/tmp"
	corsEnabled := false


	jserver := NewJsonServer(serverRoot, corsEnabled)
	jserver.Start()
	//http.Handle("/xxx", jserver)
	http.HandleFunc("/mock", jserver.HandleGet)

	errors := make(chan error)

	go func() {
		logger.WithFields(logrus.Fields{
			"ip":               bindAddress,
			"port":             listenPort,
			"root":             serverRoot,
			"cors":             corsEnabled,
		}).Info("start listening")

		if err := http.ListenAndServe(fmt.Sprintf("%s:%d", bindAddress, listenPort), nil); err != nil {
			errors <- err
		}
	}()

	select{

	}
}
