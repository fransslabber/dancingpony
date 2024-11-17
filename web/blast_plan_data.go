package rest_api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	sqldb "wo-infield-service/db"
)

// ///////////////////////////////////////////////////////////////////////////////////
// // Audit_Request

type PlanData_Attribute struct {
	Id        uint32 `json:"id"`
	Name      string `json:"name"`
	Data_type string `json:"data_type"`
	Formula   string `json:"formula"`
}

type PlanData_Col struct {
	Attribute_id uint32 `json:"attribute_id"`
	Value        string `json:"value"`
}

// type PlanData_Row struct {
// 	Plan_data_id map[uint32][]PlanData_Col `json:"plan_data_id"`
// }

type PlanData_Request struct {
	MobileDeviceId string `json:"mobile_device_id"`
	Form_id        uint32 `json:"form_id"`
	User_id        uint32 `json:"user_id"`
}

type PlanData_Response struct {
	Attributes []PlanData_Attribute    `json:"attributes"`
	Plan_data  map[uint][]PlanData_Col `json:"plan_data"`
}

func planData(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		keys := r.URL.Query()
		req := PlanData_Request{}
		err := DecodeURLParameters(&req, keys, "json")
		if err != nil {
			w.WriteHeader(400)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("Failed to comprehend request: %v", err)})
		} else {
			if len(r.Header["Authorization"]) >= 1 {
				if len(strings.Split(r.Header["Authorization"][0], " ")) >= 2 {
					if userid, err := AuthenticateJWT(strings.Split(r.Header["Authorization"][0], " ")[1], req.MobileDeviceId, ecdsakey); err != nil {
						w.WriteHeader(400)
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(Error_Response{err.Error()})
					} else {

						// get all attributes
						pdres := PlanData_Response{}
						attrs, err := sqldb.GetAttributes()
						if err != nil {
							w.WriteHeader(400)
							w.Header().Set("Content-Type", "application/json")
							json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("Failed to retrieve attributes for userid %v: %v", userid, err)})
							return
						}
						pdres.Attributes = make([]PlanData_Attribute, len(*attrs))
						for indx, attr := range *attrs {
							pdres.Attributes[indx].Id = attr.Id
							pdres.Attributes[indx].Name = attr.Name
							pdres.Attributes[indx].Data_type = attr.Data_type
							if attr.Formula == nil {
								pdres.Attributes[indx].Formula = ""
							} else {
								pdres.Attributes[indx].Formula = *attr.Formula
							}
						}

						// get all plan ids
						plans, err := sqldb.GetPlans()
						if err != nil {
							w.WriteHeader(400)
							w.Header().Set("Content-Type", "application/json")
							json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("Failed to retrieve plans, err: %v", err)})
							return
						}

						if len(*plans) == 0 {
							w.WriteHeader(404)
							w.Header().Set("Content-Type", "application/json")
							json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("No plans found")})
							return
						}

						// Get all values for each plan
						pdres.Plan_data = make(map[uint][]PlanData_Col)
						for _, pln := range *plans {
							planValues, err := sqldb.GetPlanValues(pln.Id)
							if err != nil {
								w.WriteHeader(400)
								w.Header().Set("Content-Type", "application/json")
								json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("Failed to retrieve plan data values for plan %d, err: %v", pln.Id, err)})
								return
							}

							pdc := make([]PlanData_Col, len(*planValues))
							for index, pdval := range *planValues {
								pdc[index].Attribute_id = pdval.Attribute_id
								pdc[index].Value = pdval.Value
							}
							pdres.Plan_data[uint(pln.Id)] = pdc
						}

						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(pdres)
					}
				} else {
					w.WriteHeader(400)
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(Error_Response{"No authorization token found."})

				}
			} else {
				w.WriteHeader(400)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(Error_Response{"No authorization token found."})

			}
		}
	case "POST":
	default:
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{"Sorry, only GET method is supported."})
	}

}
