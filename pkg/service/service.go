package service

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/JokeTrue/image-previewer/pkg/logging"
)

type StartStopper interface {
	Start()
	Stop()
}

func Run(service StartStopper, appName string) {
	logging.Infof("starting %s", appName)
	service.Start()

	SetReady(true)
	logging.Infof("%s ready", appName)
	logging.Infof("received %s", Wait([]os.Signal{syscall.SIGTERM, syscall.SIGINT}))
	logging.Infof("stopping %s", appName)
	SetReady(false)

	service.Stop()
}

func Wait(signals []os.Signal) os.Signal {
	sig := make(chan os.Signal, len(signals))
	signal.Notify(sig, signals...)
	return <-sig
}
