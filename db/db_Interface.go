package sqldb

import (
	"fmt"
	"time"
)

// Specific to this app
func Julianstr(t time.Time) string {
	return fmt.Sprintf("%02d%03d", t.Local().Year()-2000, t.Local().YearDay())
}
func Julian(t time.Time) uint32 {
	return uint32((t.Local().Year()-2000)*1000 + t.Local().YearDay())
}

func CheckRegToken(tok string) error {

	db := SqlDB{}
	err := db.Open("")
	if err == nil {
		defer db.Close()

		err, properties := db.Load_properties_by_name("registration_token")
		if err == nil {
			if len(*properties) != 0 {
				if (*properties)[0].Property_str == tok {
					return nil
				} else {
					return fmt.Errorf("Registration token is incorrect.")
				}

			} else {
				return fmt.Errorf("Failed to find registration token.")
			}
		} else {
			return err
		}
	} else {
		return fmt.Errorf("Could not open DB. %v", err.Error())
	}
}

func GetJWTSessionInterval() (uint32, error) {

	db := SqlDB{}
	err := db.Open("")
	if err == nil {
		defer db.Close()

		err, properties := db.Load_properties_by_name("jwt_expiration_interval")
		if err == nil {
			if len(*properties) != 0 {
				return (*properties)[0].Property_int, nil
			} else {
				return 0, fmt.Errorf("Failed to find jwt_expiration_interval.")
			}
		} else {
			return 0, err
		}
	} else {
		return 0, fmt.Errorf("Could not open DB. %v", err.Error())
	}
}

func RegisterDevice() (error, string) {
	db := SqlDB{}
	err := db.Open("")
	if err == nil {
		defer db.Close()

		if err, uuid_ := db.register_device_sp(); err == nil {
			return nil, uuid_
		} else {
			return fmt.Errorf("Could notregister device. %v", err.Error()), ""
		}
	} else {
		return fmt.Errorf("Could not open DB. %v", err.Error()), ""
	}
}

func CheckDevice(uuid_ string) (bool, error) {
	db := SqlDB{}
	err := db.Open("")
	if err == nil {
		defer db.Close()

		if err, devices := db.Load_devices_by_device_uuid(uuid_); err == nil {
			return len(*devices) != 0, nil
		} else {
			return false, fmt.Errorf("Could query for device. %v", err.Error())
		}
	} else {
		return false, fmt.Errorf("Could not open DB. %v", err.Error())
	}
}

func LoginUser(user string, password string, device_uuid string) (string, uint32, error) {
	db := SqlDB{}
	err := db.Open("")
	if err == nil {
		defer db.Close()
		return db.Login_user_sp(user, password, device_uuid)
	} else {
		return "", 0, fmt.Errorf("Could not open DB. %v", err.Error())
	}
}

func LogoutUser(userid uint32, device_uuid string) error {
	db := SqlDB{}
	err := db.Open("")
	if err == nil {
		defer db.Close()
		return db.Create_audit_log(device_uuid, userid, "LOGOUT USER", "", time.Now().Format("2006-02-01"))
	} else {
		return fmt.Errorf("Could not open DB. %v", err.Error())
	}
}

func AddUser(name string, user_name string, passwd string, card_id string, auth_type string) (uint32, error) {
	db := SqlDB{}
	err := db.Open("")
	if err == nil {
		defer db.Close()
		return db.Create_user_sp(name, user_name, passwd, card_id, auth_type)
	} else {
		return 0, fmt.Errorf("Could not open DB. %v", err.Error())
	}
}

func LoginCard(card string) (string, uint32, error) {
	db := SqlDB{}
	err := db.Open("")
	if err == nil {
		defer db.Close()
		return db.Login_card_sp(card)
	} else {
		return "", 0, fmt.Errorf("Could not open DB. %v", err.Error())
	}
}

func LogAudit(device_id string, user_id uint32, action string, details string, date_created string) error {
	db := SqlDB{}
	err := db.Open("")
	if err == nil {
		defer db.Close()
		return db.Create_audit_log(device_id, user_id, action, details, date_created)
	} else {
		return fmt.Errorf("Could not open DB. %v", err.Error())
	}
}

// func GetCapTureApp(capture_app_id, version uint32) (*FormVersion, error) {
// 	db := SqlDB{}
// 	err := db.Open("")
// 	if err == nil {
// 		defer db.Close()

// 		if err, forms := db.Load_form_by_form_id_status(form_id, version, "published"); err == nil {
// 			if len(*forms) != 0 {
// 				return (*forms)[0], nil
// 			} else {
// 				return nil, fmt.Errorf("Form %d version %d not found.", form_id, version)
// 			}
// 		} else {
// 			return nil, fmt.Errorf("Form query failed. %v", err.Error())
// 		}
// 	} else {
// 		return nil, fmt.Errorf("Could not open DB. %v", err.Error())
// 	}
// }

func GetForm(form_id, version uint32) (*FormVersion, error) {
	db := SqlDB{}
	err := db.Open("")
	if err == nil {
		defer db.Close()

		if err, forms := db.Load_form_by_form_id_status(form_id, version, "published"); err == nil {
			if len(*forms) != 0 {
				return (*forms)[0], nil
			} else {
				return nil, fmt.Errorf("Form %d version %d not found.", form_id, version)
			}
		} else {
			return nil, fmt.Errorf("Form query failed. %v", err.Error())
		}
	} else {
		return nil, fmt.Errorf("Could not open DB. %v", err.Error())
	}
}

func GetFormFieldOptions(form_id, version uint32) (*Array_FieldOptions, error) {
	db := SqlDB{}
	err := db.Open("")
	if err == nil {
		defer db.Close()

		if fieldoptions, err := db.Load_form_field_options(form_id, version); err == nil {
			return fieldoptions, nil
		} else {
			return nil, fmt.Errorf("Form query failed. %v", err.Error())
		}
	} else {
		return nil, fmt.Errorf("Could not open DB. %v", err.Error())
	}
}
func GetFormFieldOptionsByUser(userid, form_id, version uint32) (*Array_FieldOptions, error) {
	db := SqlDB{}
	err := db.Open("")
	if err == nil {
		defer db.Close()

		if fieldoptions, err := db.Load_form_field_options_by_user(userid, form_id, version); err == nil {
			return fieldoptions, nil
		} else {
			return nil, fmt.Errorf("Form query failed. %v", err.Error())
		}
	} else {
		return nil, fmt.Errorf("Could not open DB. %v", err.Error())
	}
}

func AddFormFieldOptions(field_id uint32, value, filter string) (uint32, error) {
	db := SqlDB{}
	err := db.Open("")
	if err == nil {
		defer db.Close()

		if fieldoptions_id, err := db.Create_field_option(&FieldOption{Field_id: field_id, Option_value: value, Option_filter: filter}); err == nil {
			return fieldoptions_id, nil
		} else {
			return 0, fmt.Errorf("Form query failed. %v", err.Error())
		}
	} else {
		return 0, fmt.Errorf("Could not open DB. %v", err.Error())
	}
}

func GetFormFieldVersionsByFormId(form_id uint32) (*Array_FieldVersions, error) {
	db := SqlDB{}
	err := db.Open("")
	if err == nil {
		defer db.Close()

		if fieldversions, err := db.Load_field_versions_by_form_id(form_id); err == nil {
			return fieldversions, nil
		} else {
			return nil, fmt.Errorf("Form query failed. %v", err.Error())
		}
	} else {
		return nil, fmt.Errorf("Could not open DB. %v", err.Error())
	}
}

func UpdateFormVersionJSON(form_id, version uint32, strjson string) (uint32, error) {
	db := SqlDB{}
	err := db.Open("")
	if err == nil {
		defer db.Close()

		if rowsaffected, err := db.Update_form_version_JSON(form_id, version, strjson); err == nil {
			return rowsaffected, nil
		} else {
			return 0, fmt.Errorf("Form query failed. %v", err.Error())
		}
	} else {
		return 0, fmt.Errorf("Could not open DB. %v", err.Error())
	}
}

func GetForms(user_id uint32) (*Array_Forms, error) {
	db := SqlDB{}
	err := db.Open("")
	if err == nil {
		defer db.Close()

		if err, forms := db.Load_forms_by_user_id_status(user_id, "published"); err == nil {
			return forms, nil
		} else {
			return nil, fmt.Errorf("Get form for user query failed. %v", err.Error())
		}
	} else {
		return nil, fmt.Errorf("Could not open DB. %v", err.Error())
	}
}

func GetAllForms() (*Array_Forms, error) {
	db := SqlDB{}
	err := db.Open("")
	if err == nil {
		defer db.Close()

		if err, forms := db.Load_all_forms(); err == nil {
			return forms, nil
		} else {
			return nil, fmt.Errorf("Get form for user query failed. %v", err.Error())
		}
	} else {
		return nil, fmt.Errorf("Could not open DB. %v", err.Error())
	}
}

func GetDataSources() (*Array_Data_Sources, error) {
	db := SqlDB{}
	err := db.Open("")
	if err == nil {
		defer db.Close()

		if dss, err := db.Load_all_distinct_data_sources(); err == nil {
			return dss, nil
		} else {
			return nil, fmt.Errorf("Get data sources query failed. %v", err.Error())
		}
	} else {
		return nil, fmt.Errorf("Could not open DB. %v", err.Error())
	}
}

func GetUser(user_name string) (*Array_Users, error) {
	db := SqlDB{}
	err := db.Open("")
	if err == nil {
		defer db.Close()

		if err, users := db.Load_users_by_user_name(user_name); err == nil {
			return users, nil
		} else {
			return nil, fmt.Errorf("Get user query failed. %v", err.Error())
		}
	} else {
		return nil, fmt.Errorf("Could not open DB. %v", err.Error())
	}
}

func AddPlanDataRow(form_id uint32, record []string) error {
	db := SqlDB{}
	err := db.Open("")
	if err == nil {
		defer db.Close()
		return db.Add_Plan_Data_sp(form_id, record)
	} else {
		return fmt.Errorf("Could not open DB. %v", err.Error())
	}
}

func AddDownloadDataFieldHeaders(record string, data_source string, user string) error {
	db := SqlDB{}
	err := db.Open("")
	if err == nil {
		defer db.Close()
		return db.Add_download_data_field_headers_sp(record, data_source, user)
	} else {
		return fmt.Errorf("Could not open DB. %v", err.Error())
	}
}

func DeleteAllDownloadDataFieldOptions() (uint32, error) {
	db := SqlDB{}
	err := db.Open("")
	if err == nil {
		defer db.Close()
		return db.Delete_All_Field_Options()
	} else {
		return 0, fmt.Errorf("Could not open DB. %v", err.Error())
	}
}

func DeleteDownloadDataFieldOptionsByName(data_source_name string) (uint32, error) {
	db := SqlDB{}
	err := db.Open("")
	if err == nil {
		defer db.Close()
		return db.Delete_Field_Options_by_name(data_source_name)
	} else {
		return 0, fmt.Errorf("Could not open DB. %v", err.Error())
	}
}
func DeleteDownloadDataGroupsByName(name string) (uint32, error) {
	db := SqlDB{}
	err := db.Open("")
	if err == nil {
		defer db.Close()
		return db.Delete_groups_by_name(name)
	} else {
		return 0, fmt.Errorf("Could not open DB. %v", err.Error())
	}
}
func AddDownloadDataRowItem(headers, record, metadata string, data_source string) error {
	db := SqlDB{}
	err := db.Open("")
	if err == nil {
		defer db.Close()
		return db.Add_download_data_row_item_sp(headers, record, data_source, metadata)
	} else {
		return fmt.Errorf("Could not open DB. %v", err.Error())
	}
}

func AddDownloadDataGroupItem(headers string, record string, parentid uint32) (uint32, error) {
	db := SqlDB{}
	err := db.Open("")
	if err == nil {
		defer db.Close()
		//grp := Group{Name: headers, Filter: record, Parent_id: parentid}
		//return db.Create_group(&grp)
		return db.Add_download_data_group_item_sp(headers, record, parentid)
	} else {
		return 0, fmt.Errorf("Could not open DB. %v", err.Error())
	}
}

func AddDownloadDataJSONArray(itemsJSON, filter, metadata string, data_source string) error {
	db := SqlDB{}
	err := db.Open("")
	if err == nil {
		defer db.Close()
		return db.Add_download_data_JSON_array_sp(itemsJSON, filter, data_source, metadata)
	} else {
		return fmt.Errorf("Could not open DB. %v", err.Error())
	}
}

func GetAttributes() (*Array_Attributes, error) {
	db := SqlDB{}
	err := db.Open("")
	if err == nil {
		defer db.Close()

		if err, attrs := db.Load_attributes(); err == nil {
			return attrs, nil
		} else {
			return nil, fmt.Errorf("Get attributes query failed. %v", err.Error())
		}
	} else {
		return nil, fmt.Errorf("Could not open DB. %v", err.Error())
	}
}

func GetPlans() (*Array_BlastPlanData, error) {
	db := SqlDB{}
	err := db.Open("")
	if err == nil {
		defer db.Close()

		if err, plndatas := db.Load_plandata(); err == nil {
			return plndatas, nil
		} else {
			return nil, fmt.Errorf("Get plan data rows query failed. %v", err.Error())
		}
	} else {
		return nil, fmt.Errorf("Could not open DB. %v", err.Error())
	}
}

func GetPlanValues(planid uint32) (*Array_BlastPlanDataValue, error) {
	db := SqlDB{}
	err := db.Open("")
	if err == nil {
		defer db.Close()

		if err, plndatavals := db.Load_plandatavalues(planid); err == nil {
			return plndatavals, nil
		} else {
			return nil, fmt.Errorf("Get plan data rows query failed. %v", err.Error())
		}
	} else {
		return nil, fmt.Errorf("Could not open DB. %v", err.Error())
	}
}
