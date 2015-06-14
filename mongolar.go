// Copyright © 2014 Jason Smith <jasonrichardsmith@gmail.com>.
//
// Use of this source code is governed by the GPL-3
// license that can be found in the LICENSE file.

package main

import (
	"github.com/mongolar/mongolar/admin"
	"github.com/mongolar/mongolar/basecontrollers"
	"github.com/mongolar/mongolar/configs"
	"github.com/mongolar/mongolar/controller"
	"github.com/mongolar/mongolar/oauthlogin"
	"github.com/mongolar/mongolar/router"
	"gopkg.in/mgo.v2"
	"net/http"
	"time"
)

func main() {
	cm := controller.NewMap()
	basecontrollers.GetControllerMap(cm)
	admin.GetControllerMap(cm)
	oauthlogin.GetControllerMap(cm)
	Serve(cm)
}

func Serve(cm controller.ControllerMap) {
	c, port := configs.New()
	EnsureIndexes(c)
	HostSwitch := router.New(c.Aliases, c.SitesMap, cm)
	http.ListenAndServe(":"+port, HostSwitch)
}

func EnsureIndexes(configs *configs.Configs) {
	for _, site_config := range configs.SitesMap {
		db_session := site_config.DbSession.Copy()
		defer db_session.Close()
		duration := time.Duration(site_config.SessionExpiration * time.Hour)
		i := mgo.Index{
			Key:         []string{"updated"},
			Unique:      false,
			DropDups:    false,
			Background:  true,
			Sparse:      false,
			ExpireAfter: duration,
		}
		c := db_session.DB("").C("sessions")
		c.EnsureIndex(i)
		i = mgo.Index{
			Key:        []string{"path", "wildcard"},
			Unique:     true,
			DropDups:   true,
			Background: true,
			Sparse:     false,
		}
		c = db_session.DB("").C("paths")
		c.EnsureIndex(i)
		i = mgo.Index{
			Key:        []string{"id", "type"},
			Unique:     true,
			DropDups:   true,
			Background: true,
			Sparse:     false,
		}
		c = db_session.DB("").C("users")
		duration = time.Duration(2 * time.Hour)
		c.EnsureIndex(i)
		i = mgo.Index{
			Key:         []string{"created"},
			Unique:      false,
			DropDups:    false,
			Background:  true,
			Sparse:      false,
			ExpireAfter: duration,
		}
		c = db_session.DB("").C("forms")
		c.EnsureIndex(i)
	}
}
