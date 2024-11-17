package rest_api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ///////////////////////////////////////////////////////////////////////////////////
// // Web Interface
// type Page struct {
// 	Title string
// 	Body  []byte
// }

// func (p *Page) savePage() error {
// 	filename := "http_web/data/" + p.Title + ".txt"
// 	return os.WriteFile(filename, p.Body, 0600)
// }

// func loadPage(title string) (*Page, error) {
// 	filename := "http_web/data/" + title + ".txt"
// 	body, err := os.ReadFile(filename)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &Page{Title: title, Body: body}, nil
// }

// func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
// 	t, err := template.ParseFiles("http_web/tmpl/" + tmpl + ".html")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	err = t.Execute(w, p)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

// func view_formGenerator(w http.ResponseWriter, r *http.Request) {

// 	// view a form generator page

// 	title := r.URL.Path[len("/api/v1/form_generator_view/"):]
// 	fmt.Printf("view %v\n", title)
// 	p, err := loadPage(title)
// 	if err != nil {
// 		fmt.Printf("redirect %v\n", "/api/v1/form_generator_edit/"+title)
// 		http.Redirect(w, r, "/api/v1/form_generator_edit/"+title, http.StatusFound)
// 		return
// 	}
// 	renderTemplate(w, "view", p)

// }

// func edit_formGenerator(w http.ResponseWriter, r *http.Request) {

// 	// Edit a form generator page
// 	title := r.URL.Path[len("/api/v1/form_generator_edit/"):]
// 	p, err := loadPage(title)
// 	if err != nil {
// 		body, _ := JSON_form_DB()
// 		p = &Page{Title: title, Body: body}
// 	}
// 	renderTemplate(w, "edit", p)
// }

// func save_formGenerator(w http.ResponseWriter, r *http.Request) {

// 	// save a form generator page
// 	title := r.URL.Path[len("/api/v1/form_generator_save/"):]
// 	body := r.FormValue("body")
// 	p := &Page{Title: title, Body: []byte(body)}
// 	err := p.savePage()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	http.Redirect(w, r, "/api/v1/form_generator_view/"+title, http.StatusFound)
// }

type FormJSONGenerator_Request struct {
	Form_id uint32  `json:"form_id"`
	Version float32 `json:"version"`
}

func formGenerator(w http.ResponseWriter, r *http.Request) {

	var req FormData_Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("Failed to comprehend request: %v", err)})
	} else {
		//GenerateJSONform(req.FormId, req.Version)

	}
}

// ///////////////////////////////////////////////////////////////////////////////
// Generate JSON from DB
type JSONForm_Option struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type JSONForm_Conditional struct {
	DependsOn string `json:"dependsOn"`
	ShowIf    map[string]JSONForm_Option
}

type JSONForm_Validation struct {
	Required        bool    `json:"required"`
	MaxLen          int32   `json:"maxLength"`
	MinLen          int32   `json:"minLength"`
	MaxNumericRange float32 `json:"max"`
	MinNumericRange float32 `json:"min"`
}

type JSONForm_Field struct {
	Type             string               `json:"type"`
	Label            string               `json:"label"`
	Placeholder      string               `json:"placeholder"`
	HTMLId           string               `json:"id"`
	Value            string               `json:"value"`
	BackgroundColour string               `json:"backgroundColour"`
	SelectOptions    []JSONForm_Option    `json:"options"`
	Conditional      JSONForm_Conditional `json:"conditional"`
}

type JSONForm_Submit struct {
	Type     string `json:"type"`
	Endpoint string `json:"endpoint"`
}

type JSONForm_Main struct {
	Title   string           `json:"formTitle"`
	HTMLId  string           `json:"formId"`
	Version float32          `json:"version"`
	Fields  []JSONForm_Field `json:"fields"`
	Submit  JSONForm_Submit  `json:"submitAction"`
}

type JSONForm_CaptureApp struct {
	AppName     string                   `json:"appName"`
	AppId       uint32                   `json:"appId"`
	AppType     string                   `json:"appType"`
	AppVersion  float32                  `json:"appVersion"`
	PrimaryForm uint32                   `json:"primary_form"`
	Forms       map[uint32]JSONForm_Main `json:"forms"`
}

// func GenerateJSONform(capture_app_id, version uint32) ([]byte, error) {
// 	var err error = nil

// 	// Get Capture app
// 	if dbform, err := sqldb.GetCapTureApp(capture_app_id, version); err == nil {
// 		ca := JSONForm_CaptureApp{}

// 	if dbform, err := sqldb.GetForm(form_id, version); err == nil {
// 		frm := JSONForm_Main{Title: dbform.Name, HTMLId: strings.ReplaceAll(dbform.Name, " ", "_"), Version: dbform.Version}

// 		if dbfldversions, err := sqldb.GetFormFieldVersionsByFormId(dbform.Form_id); err == nil {
// 			for _, dbfldversion := range *dbfldversions {
// 				frmfield := JSONForm_Field{Type: dbfldversion.Field_type, Label: dbfldversion.Name, Placeholder: dbfldversion.Description, HTMLId: strings.ReplaceAll(dbfldversion.Name, " ", "_")}
// 				frm.Fields = append(frm.Fields, frmfield)
// 			}
// 		} else {
// 			fmt.Printf("GetFormFieldVersionsByFormId failed %v\n", err)
// 		}

// 		jsonstr, _ := json.Marshal(frm)
// 		sqldb.UpdateFormVersionJSON(form_id, version, string(jsonstr))
// 	}
// 	return nil, err
// }
