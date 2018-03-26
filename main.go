// +build windows

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kardianos/service"
)

type program struct {
	logger service.Logger
	doneCh chan struct{}
}

func (p *program) Start(s service.Service) error {
	// this function should not block
	p.logger.Info("Start called!")
	fmt.Println("Start called!: ", s)
	p.doneCh = make(chan struct{})
	go p.run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	// this function should not block
	p.logger.Info("Stop called!")
	fmt.Println("Stop called!: ", s)
	p.doneCh <- struct{}{}
	return nil
}

func (p *program) run() error {
	for {
		select {
		case <-time.After(5 * time.Second):
		case <-p.doneCh:
			return nil
		}
		p.logger.Info("running ...")
	}
}

func main() {
	cfg := &service.Config{
		Name:        "eventlogtest",
		DisplayName: "Event Log Test",
		Description: "This is event log test from Go!",
	}

	prg := &program{}
	s, err := service.New(prg, cfg)
	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
		return
	}

	logger, err := s.Logger(nil)
	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
		return
	}
	prg.logger = logger
	prg.logger.Info("event log test running")

	if len(os.Args) > 1 {
		verb := os.Args[1]
		err = service.Control(s, verb)
		if err != nil {
			prg.logger.Error(verb, " failed: ", err.Error())
			fmt.Println(verb, " failed: ", err.Error())
			return
		}
		prg.logger.Info(verb, " succeeded")
		fmt.Println(verb, " succeeded")
		return
	}

	err = s.Run()
	if err != nil {
		prg.logger.Error("Run failed: ", err.Error())
		fmt.Println("Run failed: ", err.Error())
		return
	}
}
