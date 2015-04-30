package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/jasonrichardsmith/mongolar/configs"
	"github.com/jasonrichardsmith/mongolar/controller"
	"github.com/jasonrichardsmith/mongolar/router"
	"net/http"
)

func main() {
	cm := controller.NewMap()
	cm["domian_public_value"] = controller.DomainPublicValue
	cm["path"] = controller.PathValues
	cm["content"] = controller.ContentValues
	cm["wrapper"] = controller.WrapperValues
	cm["slug"] = controller.SlugValues
	Serve(cm)
}

func Serve(cm controller.ControllerMap) {

	c := configs.New()
	//	spew.Dump(c)
	HostSwitch := router.New(c.Aliases, c.SitesMap, cm)
	err := http.ListenAndServe(":"+c.Server.Port, HostSwitch)
	spew.Dump(err)
}
