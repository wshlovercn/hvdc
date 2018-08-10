package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"fmt"
	"github.com/golang/glog"
	"hvdc/rtuserver"
)

var conf_file = flag.String("conf", "config.toml", "config file path")

func init()  {
	flag.Set("log_dir", "/Users/fengchen/go/src/hvdc/hvdc/logs")

}

type HdvcServer struct {
	config *Config

	rtuServer *rtuserver.RtuServer
}

/**
 ApplicationDelegate
 */
func (server *HdvcServer) Initialize() error {
	server.config = &Config{}
	if _, err := toml.DecodeFile(*conf_file, server.config); err != nil {
		glog.Fatalf("load config file %s error: %v", *conf_file, err)
		return fmt.Errorf("load config file %s error: %v", *conf_file, err)
	}

	if server.rtuServer = rtuserver.NewRtuServer(server.config.RtuConf); server.rtuServer == nil {
		glog.Fatalf("create rtu server failed")
		return fmt.Errorf("create rtu server failed")
	}

	return nil
}

func (server *HdvcServer) WorkLoop() error {
	//TODO:服务异常停止监控
	go server.rtuServer.Serve()

	return nil
}

func (server *HdvcServer) Terminate() error {
	server.rtuServer.Stop()

	return nil
}
