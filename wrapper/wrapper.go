// Wrapper defines the struct to be passed back to the controller.
// It contains the entirety of the response and performs all marshalling and write operations
package wrapper

import (
	"encoding/json"
	"errors"
	"github.com/jasonrichardsmith/mongolar/configs"
	"github.com/jasonrichardsmith/mongolar/session"
	"net/http"
)

// Wrapper structure required to be passed back to the Controller
type Wrapper struct {
	Writer     http.ResponseWriter    // The response writer
	Request    *http.Request          // The request
	SiteConfig *configs.SiteConfig    // The configuration for the site being accessed
	Session    *session.Session       // Session for user
	Payload    map[string]interface{} // This is the sum of the payload that will be returned to the user
}

//Constructor for the Wrapper
func New(w http.ResponseWriter, r *http.Request, s *configs.SiteConfig) *Wrapper {
	wr := Wrapper{Writer: w, Request: r, SiteConfig: s}
	// Get session
	//	wr.Session = session.New(w, r, s)
	// Define payload
	wr.Payload = make(map[string]interface{})
	return &wr
}

// Helper function for the controller to easily add its final content to the Payload
func (w *Wrapper) SetContent(c interface{}) {
	w.SetPayload("content", c)
}

// Helper function for the controller to easily add its final content to the Payload
func (w *Wrapper) SetTemplate(t string) {
	w.SetPayload("template", t)
}

// Helper function for the controller to easily add its final content to the Payload
func (w *Wrapper) SetDynamicId(i string) {
	w.SetPayload("dynamic_id", i)
}

// Sets payload based on a keyvalue
func (w *Wrapper) SetPayload(n string, v interface{}) {
	w.Payload[n] = v
}

// Gets payload based on a keyvalue
func (w *Wrapper) GetAPayload(n string) (interface{}, error) {
	if v, ok := w.Payload[n]; ok {
		return v, nil
	}
	err := errors.New("Payload value not set")
	return nil, err
}

// Gets entire payload
func (w *Wrapper) GetPayload() map[string]interface{} {
	return w.Payload
}

// The final serve function.  This will marshall the payload and serve it to the user.
func (w *Wrapper) Serve() {
	js, err := json.Marshal(w.Payload)
	if err != nil {
		http.Error(w.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Writer.Write(js)
}
