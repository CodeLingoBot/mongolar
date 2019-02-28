package admin

import (
	"fmt"
	"github.com/mongolar/mongolar/form"
	"github.com/mongolar/mongolar/models/paths"
	"github.com/mongolar/mongolar/services"
	"github.com/mongolar/mongolar/wrapper"
	"net/http"
)

// Controller to list all paths
func AdminPaths(w *wrapper.Wrapper) {
	pl, err := paths.PathList(w)
	if err != nil {
		services.AddMessage("There was an error retrieving your site paths", "Error", w)
		errmessage := fmt.Sprintf("Error getting path list: %s", err.Error())
		w.SiteConfig.Logger.Error(errmessage)
	} else {
		w.SetContent(pl)
	}
	w.Serve()
}

// Controller for editing paths
func PathEditor(w *wrapper.Wrapper) {
	if len(w.APIParams) == 0 {
		http.Error(w.Writer, "Forbidden", 403)
		w.Serve()
		return
	}
	if w.Request.Method != "POST" {
		PathEditorForm(w)
		return
	}
	PathEditorSubmit(w)
	return
}

// Controller to present path editor form
func PathEditorForm(w *wrapper.Wrapper) {
	pathid := w.APIParams[0]
	f := form.NewForm()
	f.AddText("title", "text").AddLabel("Title").Required()
	f.AddText("path", "text").AddLabel("Path").Required()
	f.AddText("template", "text").AddLabel("Template").Required()
	f.AddCheckBox("wildcard").AddLabel("Wildcard")
	ops := []map[string]string{
		map[string]string{"name": "published", "value": "published"},
		map[string]string{"name": "unpublished", "value": "unpublished"},
	}
	f.AddRadio("status", ops).AddLabel("Status").Required()
	f.AddText("id", "text").Hidden()
	var p paths.Path
	var err error
	if pathid != "new" {
		p, err = paths.LoadPath(pathid, w)
		if err != nil {
			errmessage := fmt.Sprintf("Could not retrieve path %s by %s: %s", w.APIParams[0], w.Request.Host, err.Error())
			w.SiteConfig.Logger.Error(errmessage)
			services.AddMessage("Error retrieving path information.", "Error", w)
			w.Serve()
			return
		}
	} else {
		p = paths.NewPath()
	}
	f.FormData = p
	f.Register(w)
	w.SetTemplate("admin/form.html")
	w.SetPayload("form", f)
	w.Serve()
	return
}

// Controller to handle path editor form submissions.
func PathEditorSubmit(w *wrapper.Wrapper) {
	pathid := w.APIParams[0]
	var path paths.Path
	var err error
	if pathid != "new" {
		path, err = paths.LoadPath(pathid, w)
		if err != nil {
			errmessage := fmt.Sprintf("Unable to save path %s by %s: %s", pathid,
				w.Request.Host, err.Error())
			w.SiteConfig.Logger.Error(errmessage)
			services.AddMessage("There was a problem saving your path.", "Error", w)
			w.Serve()
			return
		}
	}
	err = form.GetValidFormData(w, &path)
	if err != nil {
		return
	}
	err = path.Save(w)
	if err != nil {
		errmessage := fmt.Sprintf("Unable to save path %s by %s: %s", pathid,
			w.Request.Host, err.Error())
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("There was a problem saving your path.", "Error", w)
		w.Serve()
		return
	}
	services.AddMessage("Your path was saved.", "Success", w)
	dynamic := services.Dynamic{
		Target:     "pathbar",
		Controller: "admin/paths",
		Template:   "admin/path_list.html",
	}
	services.SetDynamic(dynamic, w)
	w.Serve()
	return

}

// PathElements retrieves a list of elements in a path for content editor.
func PathElements(w *wrapper.Wrapper) {
	if len(w.APIParams) == 0 {
		http.Error(w.Writer, "Forbidden", 403)
		w.Serve()
		return
	}
	pathid := w.APIParams[0]
	p, err := paths.LoadPath(pathid, w)
	if err != nil {
		errmessage := fmt.Sprintf("Path not found to edit for %s by %s ", pathid, w.Request.Host)
		w.SiteConfig.Logger.Error(errmessage)
		services.AddMessage("This path was not found", "Error", w)
		w.Serve()
	} else {
		w.SetPayload("path", p)
		if len(p.Elements) == 0 {
			services.AddMessage("This path has no elements.", "Info", w)
		}
		w.Serve()
	}

}
