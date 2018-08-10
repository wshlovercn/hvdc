package app

import (
	"flag"
	"github.com/golang/glog"
	"os"
	"os/signal"
	"syscall"
)

type ApplicationDelegate interface {
	Initialize() error
	WorkLoop() error
	Terminate() error
}

type Application struct {
	delegate ApplicationDelegate
	chSig chan os.Signal
}

func NewApplication(delegate ApplicationDelegate) *Application {
	return &Application{
		delegate:delegate,
		chSig: make(chan os.Signal, 1),
	}
}

func (app *Application) Run()  {
	flag.Parse()
	defer glog.Flush()

	if app.delegate == nil {
		glog.Fatalf("application with nil delegate")
		return
	}

	glog.Infof("application initialize...")
	if err := app.delegate.Initialize(); err != nil {
		glog.Fatalf("application initialize error: %", err)
		return
	}

	signal.Notify(app.chSig, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)

	glog.Infof("start application workloop...")
	go func() {
		if err := app.delegate.WorkLoop(); err != nil {
			glog.Fatalf("application work loop error: %v", err)
			app.chSig <- syscall.SIGTERM
		}
	}()

	sig := <- app.chSig
	glog.Warningf("recv signal(%v) application will terminate", sig)

	app.delegate.Terminate()
	glog.Infof("application terminated!")
}