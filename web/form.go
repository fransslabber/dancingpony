package rest_api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	sqldb "wo-infield-service/db"
)

// ///////////////////////////////////////////////////////////////////////////////////
// //
type FormData_Request struct {
	FormId         uint32 `json:"form_id"`
	Version        uint32 `json:"version"`
	MobileDeviceId string `json:"mobile_device_id"`
}

type FormData_Response struct {
	Table_headers []FormData_ResponseTableHeaders `json:"table_headers"`
	Form          FormData_ResponseData           `json:"form"`
}

type FormData_ResponseTableHeaders struct {
	Name             string `json:"name"`
	Field_id_version int    `json:"field_version_id"`
	Field_type       string `json:"field_type"`
}

type FormData_ResponseData struct {
	Id          int                                  `json:"id"`
	Version     int                                  `json:"version"`
	Name        string                               `json:"name"`
	Description string                               `json:"description"`
	Data        string                               `json:"data"`
	Options     map[uint32][]FormOptions_FieldOption `json:"options"`
}

func formData(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":

		keys := r.URL.Query()
		req := FormData_Request{}
		err := DecodeURLParameters(&req, keys, "json")
		//err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(400)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("Failed to comprehend request: %v", err)})
		} else {
			if len(strings.Split(r.Header["Authorization"][0], " ")) > 1 {
				if userid, err := AuthenticateJWT(strings.Split(r.Header["Authorization"][0], " ")[1], req.MobileDeviceId, ecdsakey); err != nil {
					w.WriteHeader(401)
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(Error_Response{err.Error()})
				} else {

					if form, err := sqldb.GetForm(req.FormId, req.Version); err == nil {
						if form != nil {
							// Get table headers for form fields
							res := FormData_Response{}
							if fieldversions, err := sqldb.GetFormFieldVersionsByFormId(req.FormId); err == nil {
								for _, fv := range *fieldversions {
									res.Table_headers = append(res.Table_headers, FormData_ResponseTableHeaders{Name: fv.Name, Field_id_version: int(fv.Id), Field_type: fv.Field_type})
								}
							} else {
								w.WriteHeader(404)
								w.Header().Set("Content-Type", "application/json")
								json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("DB field versions query failed %v.", err.Error())})
							}

							// Get Form Field options
							if field_options, err := sqldb.GetFormFieldOptionsByUser(userid, req.FormId, req.Version); err == nil {
								res.Form = FormData_ResponseData{Id: int(form.Id), Version: int(form.Version), Name: form.Name, Description: form.Description, Data: form.Form_JSON}
								res.Form.Options = make(map[uint32][]FormOptions_FieldOption)
								for _, field_option := range *field_options {
									fo := FormOptions_FieldOption{field_option.Option_filter, field_option.Field_id, field_option.Option_value}
									res.Form.Options[field_option.Id] = append(res.Form.Options[field_option.Id], fo) // this query returns the field_version_id in the field_option.id variable
								}
								w.Header().Set("Content-Type", "application/json")
								json.NewEncoder(w).Encode(res)
							} else {
								w.WriteHeader(404)
								w.Header().Set("Content-Type", "application/json")
								json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("DB field options query failed %v.", err.Error())})
							}

						} else {
							w.WriteHeader(404)
							w.Header().Set("Content-Type", "application/json")
							json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("Form not found id=%d version=%d.", req.FormId, req.Version)})
						}
					} else {
						w.WriteHeader(404)
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("DB Form query failed %v.", err.Error())})
					}
				}
			} else {
				w.WriteHeader(400)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("Failed to comprehend request: %v", err)})
			}
		}
	case "POST":
	default:
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{"Sorry, only GET and POST(not implemented yet) methods are supported."})
	}
}

// ////////////////////////////////////////////////////////////////////////////////////////////////
// Form Options
type FormOptions_Request struct {
	FormId         uint32 `json:"form_id"`
	Version        uint32 `json:"version"`
	MobileDeviceId string `json:"mobile_device_id"`
}

type FormOptions_FieldOption struct {
	MO          string `json:"filter_option"`
	Field_id    uint32 `json:"field_id"`
	Field_value string `json:"field_value"`
}

type FormOptions_ResponseData struct {
	Id      uint32                               `json:"id"`
	Options map[uint32][]FormOptions_FieldOption `json:"options"`
}

type FormOptions_Response struct {
	Form FormOptions_ResponseData `json:"form"`
}

func formOptions(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		keys := r.URL.Query()
		req := FormOptions_Request{}
		err := DecodeURLParameters(&req, keys, "json")
		if err != nil {
			w.WriteHeader(400)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("Failed to comprehend request: %v", err)})
		} else {
			if len(strings.Split(r.Header["Authorization"][0], " ")) > 1 {
				if userid, err := AuthenticateJWT(strings.Split(r.Header["Authorization"][0], " ")[1], req.MobileDeviceId, ecdsakey); err != nil {
					w.WriteHeader(401)
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(Error_Response{err.Error()})
				} else {

					if field_options, err := sqldb.GetFormFieldOptionsByUser(userid, req.FormId, req.Version); err == nil {
						res := FormOptions_Response{}
						res.Form.Id = req.FormId
						res.Form.Options = make(map[uint32][]FormOptions_FieldOption)
						for _, field_option := range *field_options {
							fo := FormOptions_FieldOption{field_option.Option_filter, field_option.Field_id, field_option.Option_value}
							res.Form.Options[field_option.Id] = append(res.Form.Options[field_option.Id], fo) // this query returns the field_version_id in the field_option.id variable
						}
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(res)

					} else {
						w.WriteHeader(404)
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("DB field options query failed %v.", err.Error())})
					}
				}
			} else {
				w.WriteHeader(400)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("Failed to comprehend request: %v", err)})
			}
		}
	case "POST":
	default:
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{"Sorry, only GET method is supported."})
	}
}
