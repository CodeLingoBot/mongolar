package configs

import (
	//"github.com/davecgh/go-spew/spew"
	"fmt"
	"github.com/jasonrichardsmith/mongolar/logger"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"log"
	"path/filepath"
	"strings"
)

const (
	SERVER_CONFIG   = "/etc/mongolar/"
	SITES_DIRECTORY = "/etc/mongolar/enabled/"
)

type Configs struct {
	Server   *Server
	SitesMap SitesMap
	Aliases  Aliases
}

func New() *Configs {
	c := new(Configs)
	c.Server = NewServer()
	c.SitesMap = NewSitesMap()
	c.Aliases = NewAliases(c.SitesMap)
	return c
}

// Simple map for the aliases
type Aliases map[string]string

// Builder function takes a SiteMap and returns the Aliases Map
func NewAliases(sm SitesMap) Aliases {
	a := make(Aliases)
	a.setAliases(&sm)
	return a
}

// Itterate over SiteMap and map all domains to their configs
func (a Aliases) setAliases(sm *SitesMap) {
	for k, s := range *sm {
		//TODO check Aliases length, if 0 fatal error
		for _, alias := range s.Aliases {
			fmt.Printf("Mapping domain  %v to sit configuration %v\n", alias, k)
			a[alias] = k
		}
	}
}

// Server Config, only port setting right now, but it will probably grow
type Server struct {
	Port string
}

// Constructor
func NewServer() *Server {
	s := new(Server)
	s.getServerConfig()
	return s
}

// Build from config file
func (s *Server) getServerConfig() {
	viper.SetConfigName("mongolar")
	viper.AddConfigPath(SERVER_CONFIG)
	viper.ReadInConfig()
	viper.Marshal(s)
}

// Individual Site Configuration Type
type SiteConfig struct {
	MongoDb           map[string]string      //Configuration for MongoDB Connection
	Directory         string                 // Directory for html and assets
	Aliases           []string               // Site Aliases/Domains
	SessionExpiration int64                  // When to expire a users Session
	TemplateEndpoint  string                 // URL where will be stored
	ForeignDomains    []string               // This will whitelist domains for loading assets from other domains
	AngularModules    []string               // A slice of angularjs modules to load
	PublicValues      map[string]string      // These values can be directly invoked from the domain controller
	Misc              map[string]interface{} // Where you can store any other value not defined here
	Logger            logger.LogChannel      // A channeel for writing Logs
	DbSession         *mgo.Session           // The master MongoDb session that gets copied
	FourOFour         string
}

// Constructor for SiteConfig
func NewSiteConfig(f string) *SiteConfig {
	s := SiteConfig{
		MongoDb: make(map[string]string),
	}
	s.getSiteConfig(f)
	s.getDbConnection(f)
	s.getLogger()
	return &s
}

// Get one site configuration and marshall it
func (s *SiteConfig) getSiteConfig(file string) {
	v := viper.New()
	v.SetConfigName(file)
	v.AddConfigPath(SITES_DIRECTORY)
	v.ReadInConfig()
	v.Marshal(s)
}

// Establish a Database connection and attach it to the site configuration
func (s *SiteConfig) getDbConnection(f string) {
	u := "mongodb://" + s.MongoDb["user"] + ":" + s.MongoDb["password"] + "@" + s.MongoDb["host"] + "/" + s.MongoDb["db"]
	var err error
	s.DbSession, err = mgo.Dial(u)
	if err != nil {
		log.Fatal(err)
	}
}

// Attach a logger channel to log errors predictably.
func (s *SiteConfig) getLogger() {
	s.Logger = logger.New(s.DbSession)
}

// The map to load site files
type SiteFiles map[int]string

// The builder for site files
func NewSiteFiles() SiteFiles {
	s := make(SiteFiles)
	s.getSiteConfigFiles()
	return s

}

//Get all enabled config file names
func (s SiteFiles) getSiteConfigFiles() {
	glob := SITES_DIRECTORY + "*.yaml"
	files, err := filepath.Glob(glob)
	if err != nil {
		log.Fatal(err)
	}
	for key, value := range files {
		var filename string
		fmt.Printf("Found configuration file %v\n", value)
		_, filename = filepath.Split(value)
		s[key] = strings.TrimSuffix(filename, ".yaml")
	}
}

// The definision of the sitemap
type SitesMap map[string]*SiteConfig

// Constructor that builds the sitemap.
func NewSitesMap() SitesMap {
	s := make(SitesMap)
	f := NewSiteFiles()
	s.getSiteConfigs(f)
	return s
}

// Builds all sites based off all found files
func (s SitesMap) getSiteConfigs(f SiteFiles) {
	for _, value := range f {
		site := NewSiteConfig(value)
		s[value] = site
	}
}
