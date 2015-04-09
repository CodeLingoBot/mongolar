package session

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/jasonrichardsmith/mongolar/wrapper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

// Wrapper structure for Sessions
type Session struct {
	Id        string
	data      *SessionData
	dbSession *mgo.Session
	query     *mgo.Query
}

// The actual data being stored
type SessionData struct {
	MongoId   bson.ObjectId `bson:"_id,omitempty"`
	SessionId string        `bson:"session_id"`
	Data      map[string]interface{}
}

//Builds the session
func New(w *wrapper.Wrapper) {
	s := new(Session)
	var duration time.Duration = time.Duration(w.SiteConfig.SessionExpiration) * time.Hour
	expire := time.Now().Add(duration)
	c, err := w.Request.Cookie("m_session_id")
	fmt.Print(err)
	if c == nil {
		id := getSessionID()
		c = &http.Cookie{
			Name:     "m_session_id",
			Value:    id,
			Path:     "/",
			Domain:   w.Request.Host,
			MaxAge:   0,
			Secure:   true,
			HttpOnly: true,
			Raw:      "m_session_id=" + id,
			Unparsed: []string{"m_session_id=" + id},
		}
	}
	c.Expires = expire
	c.RawExpires = expire.Format(time.RFC3339)
	http.SetCookie(w.Writer, c)

	s.data = &SessionData{
		SessionId: c.Value,
		Data:      make(map[string]interface{}),
	}
	s.setDbSession()
	s.Id = c.Value
	s.dbSession = w.SiteConfig.DbSession.Copy()
	s.getQuery(duration)
}

func (s Session) getQuery(d time.Duration) {
	c := s.dbSession.DB("").C("sessions")
	setCollection(c, d)
	s.query = c.Find(bson.M{"session_id": s.Id})
}

func (s Session) setDbSession() {
	c := mgo.Change{
		Update:    s.data,
		Upsert:    true,
		ReturnNew: true,
	}

	s.query.Apply(c, &s.data)
	s.getData()
}

func (s Session) getData() {
	s.query.One(&s.data)
}

func (s Session) Close() {
	s.dbSession.Close()
}

func (s Session) Get(n string) (v interface{}, err error) {
	s.getData()
	if v, ok := s.data.Data[n]; ok {
		return v, nil
	} else {
		err := errors.New("No Value")
		return nil, err
	}
}

func (s Session) Set(n string, v interface{}) {
	s.getData()
	s.data.Data[n] = v
	s.setDbSession()
}

func getSessionID() string {
	raw := make([]byte, 30)
	_, err := rand.Read(raw)
	if err != nil {
		fmt.Print(err)
	}
	return hex.EncodeToString(raw)

}

func setCollection(c *mgo.Collection, d time.Duration) {
	i := mgo.Index{
		Key:         []string{"SessionId"},
		Unique:      true,
		DropDups:    true,
		Background:  true,
		Sparse:      false,
		ExpireAfter: d,
	}
	err := c.EnsureIndex(i)
	fmt.Print(err)
}
