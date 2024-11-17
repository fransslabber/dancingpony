package sqldb

import (
	"database/sql"
	"fmt"
	"time"
)

type FormVersion struct {
	Id                 uint32
	Form_id            uint32
	Name               string
	Description        string
	Version            float32
	Form_JSON          string
	Status             string
	Created_by_user_id uint32
	Date_created       time.Time
}

type Array_Forms []*FormVersion

func (d *SqlDB) Load_forms_by_user_id_status(user_id uint32, status string) (error, *Array_Forms) {
	return d.load_forms_sql(fmt.Sprintf("select * from [dbo].[form_version] where form_id in (select form_id from user_form where user_id = %d ) and status = '%s'", user_id, status))
}

func (d *SqlDB) Load_all_forms() (error, *Array_Forms) {
	return d.load_forms_sql("select * from [dbo].[form_version]")
}

func (d *SqlDB) Load_form_by_form_id_status(form_id uint32, version uint32, status string) (error, *Array_Forms) {
	return d.load_forms_sql(fmt.Sprintf("select top 1 * from [dbo].[form_version] where form_id = %d and status = '%s' and version = %d;", form_id, status, version))
}

func (d *SqlDB) load_forms_sql(sql string) (error, *Array_Forms) {
	rows, err := d.db.Query(sql)
	if err != nil {
		return err, nil
	}

	forms := make(Array_Forms, 0)
	for rows.Next() {
		frm := FormVersion{}
		err = rows.Scan(&frm.Id, &frm.Form_id, &frm.Name, &frm.Description, &frm.Version, &frm.Form_JSON, &frm.Status, &frm.Created_by_user_id, &frm.Date_created)
		if err != nil {
			rows.Close()
			return err, nil
		}
		forms = append(forms, &frm)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return err, nil
	}
	return nil, &forms
}

func (d *SqlDB) Update_form_version_JSON(form_id, version uint32, jsonstr string) (uint32, error) {
	res, err := d.db.Exec("Update form_version set form_JSON = @p1 where form_id = @p2 AND version = @p3;", jsonstr, form_id, version)
	if err != nil {
		return 0, err
	}

	id, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return uint32(id), nil
}

type FieldOption struct {
	Id              uint32
	Field_id        uint32
	Option_value    string
	Option_filter   string
	Option_metadata string
	Date_created    time.Time
}

type Array_FieldOptions []*FieldOption

// ////////////////////////////////////////////////////////////////////////////////////////
// form_id and version and status == 'published' from form_version => all form_version_id
// form_version_id from form_field => all field_version_id
// field_version_id and status == published from field_version => all field_id
// field_id from field_options => all field options required
func (d *SqlDB) Load_form_field_options(form_id uint32, version uint32) (*Array_FieldOptions, error) {
	return d.load_field_options_sql(fmt.Sprintf("SELECT form_field.field_version_id ,field_options.field_id, field_options.option_value, field_options.option_filter, field_options.date_created "+
		"FROM form_version "+
		"INNER JOIN form_field on form_version.id = form_field.form_version_id "+
		"INNER JOIN field_version on form_field.field_version_id = field_version.id "+
		"INNER JOIN field_options on field_options.field_id = field_version.field_id "+
		"WHERE  form_version.form_id = %d AND form_version.status = 'published' AND form_version.version = %d AND field_version.status = 'published'", form_id, version))
}

func (d *SqlDB) Load_form_field_options_by_user(userid, form_id, version uint32) (*Array_FieldOptions, error) {
	return d.load_field_options_sql(fmt.Sprintf("SELECT form_field.field_version_id ,field_options.field_id, field_options.option_value, field_options.option_filter,field_options.option_metadata, field_options.date_created "+
		"FROM form_version "+
		"INNER JOIN form_field on form_version.id = form_field.form_version_id "+
		"INNER JOIN field_version on form_field.field_version_id = field_version.id "+
		"INNER JOIN field_options on field_options.field_id = field_version.field_id "+
		"WHERE  form_version.form_id = %d AND form_version.status = 'published' AND form_version.version = %d AND field_version.status = 'published' "+
		"AND field_options.option_filter IN ((SELECT attribute_ref_value from user_app where user_id = %d), 'ALL')", form_id, version, userid))
}

func (d *SqlDB) Load_field_options(user_id uint32, version uint32) (*Array_FieldOptions, error) {
	return d.load_field_options_sql("select * from [dbo].[field_options]")
}

func (d *SqlDB) load_field_options_sql(sql string) (*Array_FieldOptions, error) {
	rows, err := d.db.Query(sql)
	if err != nil {
		return nil, err
	}

	fieldoptionsarray := make(Array_FieldOptions, 0)
	for rows.Next() {
		fo := FieldOption{}
		err = rows.Scan(&fo.Id, &fo.Field_id, &fo.Option_value, &fo.Option_filter, &fo.Option_metadata, &fo.Date_created)
		if err != nil {
			rows.Close()
			return nil, err
		}
		fieldoptionsarray = append(fieldoptionsarray, &fo)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &fieldoptionsarray, nil
}

func (d *SqlDB) Create_field_option(fo *FieldOption) (uint32, error) {

	insertStr := fmt.Sprintf("INSERT INTO [dbo].[field_options] ([field_id]	,[option_value]	,[option_filter],[option_metadata],[date_created]) OUTPUT INSERTED.[id] values( %d, '%s', '%s','%s', CURRENT_TIMESTAMP);",
		fo.Field_id, fo.Option_value, fo.Option_filter, fo.Option_metadata)
	rows, err := d.db.Query(insertStr)
	if err != nil {
		return 0, err
	}

	if rows.Next() {
		defer rows.Close()
		field_options_id := uint32(0)
		err = rows.Scan(&field_options_id)
		if err != nil {
			return 0, err
		}
		return field_options_id, nil
	}
	return 0, err
}

// Add CSV file column headers to field table
func (d *SqlDB) Add_download_data_field_headers_sp(record string, data_source string, user string) error {
	_, err := d.db.Exec("spAddCSVHeaderFields",
		sql.Named("pColumnHeaderArray", record),
		sql.Named("pDataSourceName", data_source),
		sql.Named("pUser", user))
	if err != nil {
		return err
	}
	return nil
}

// Add CSV file rows to field table EXLUDING 'MO' column, 'MO' value stored as field filter option
func (d *SqlDB) Add_download_data_row_item_sp(headers, record, metadata string, data_source string) error {
	_, err := d.db.Exec("spAddCSVRow",
		sql.Named("pColumnHeaderArray", headers),
		sql.Named("pColumnArray", record),
		sql.Named("pDataSource", data_source),
		sql.Named("pMetaData", metadata),
	)
	if err != nil {
		return err
	}
	return nil
}

func (d *SqlDB) Add_download_data_group_item_sp(headers string, record string, parent_id uint32) (uint32, error) {
	group_id := uint32(0)
	if parent_id != 0 {
		_, err := d.db.Exec("spInsertGroup",
			sql.Named("pName", headers),
			sql.Named("pFilter", record),
			sql.Named("pParent_group_id", parent_id),
			sql.Named("GroupID", sql.Out{Dest: &group_id}))

		return group_id, err
	} else {
		_, err := d.db.Exec("spInsertGroup",
			sql.Named("pName", headers),
			sql.Named("pFilter", record),
			sql.Named("pParent_group_id", sql.NullInt32{}),
			sql.Named("GroupID", sql.Out{Dest: &group_id}))

		return group_id, err
	}
}

// Add JSON array items to field table setting field filter option
func (d *SqlDB) Add_download_data_JSON_array_sp(JSONitems, filter, metadata string, data_source string) error {
	_, err := d.db.Exec("spAddJSONArray",
		sql.Named("pJSONItemsArray", JSONitems),
		sql.Named("pFilter", filter),
		sql.Named("pMetaData", metadata),
		sql.Named("pDataSource", data_source))
	if err != nil {
		return err
	}
	return nil
}

func (d *SqlDB) Delete_All_Field_Options() (uint32, error) {

	res, err := d.db.Exec("delete from [dbo].[field_options];")
	if err != nil {
		return 0, err
	}

	id, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return uint32(id), nil
}

func (d *SqlDB) Delete_Field_Options_by_name(data_source_name string) (uint32, error) {

	res, err := d.db.Exec(fmt.Sprintf("DELETE field_options FROM field_options INNER JOIN field ON field_options.field_id = field.id WHERE field.data_source_name = '%s'", data_source_name))
	if err != nil {
		return 0, err
	}

	id, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return uint32(id), nil
}

type FieldVersion struct {
	Id                 uint32
	Field_id           uint32
	Name               string
	Description        string
	Field_type         string
	Data_type          string
	Validation_rules   string
	Is_required        bool
	Table_display      bool
	Is_hidden          bool
	Formula            string
	Target             string
	Version            float32
	Screen_order       uint32
	Position           string
	Status             string
	Created_by_user_id uint32
	Date_created       time.Time
}

type Array_FieldVersions []*FieldVersion

func (d *SqlDB) Load_field_versions_by_form_id(form_id uint32) (*Array_FieldVersions, error) {
	return d.load_field_versions_sql(fmt.Sprintf("SELECT [dbo].[field_version].[id]	, [dbo].[field_version].[field_id]	, [dbo].[field_version].[name]	, [dbo].[field_version].[description]	, [dbo].[field_version].[field_type]"+
		", [dbo].[field_version].[data_type] , [dbo].[field_version].[validation_rules], [dbo].[field_version].[is_required], [dbo].[field_version].[table_display], [dbo].[field_version].[is_hidden], [dbo].[field_version].[formula]"+
		", [dbo].[field_version].[target], [dbo].[field_version].[version], [dbo].[field_version].[screen_order], [dbo].[field_version].[position], [dbo].[field_version].[status], [dbo].[field_version].[created_by_user_id], [dbo].[field_version].[date_created]"+
		"FROM [dbo].[field_version] INNER JOIN form_field on field_version.id = form_field.field_version_id WHERE form_field.form_version_id = (Select id from form_version where form_id = %d)", form_id))
}

func (d *SqlDB) Load_field_version_by_field_id(field_id uint32, version uint32) (*Array_FieldVersions, error) {
	return d.load_field_versions_sql(fmt.Sprintf("SELECT * FROM field_version WHERE field_id = %d AND version = %d", field_id, version))
}

func (d *SqlDB) Load_field_versions() (*Array_FieldVersions, error) {
	return d.load_field_versions_sql("select * from [dbo].[field_version]")
}

func (d *SqlDB) load_field_versions_sql(sql string) (*Array_FieldVersions, error) {
	rows, err := d.db.Query(sql)
	if err != nil {
		return nil, err
	}

	fieldversionsarray := make(Array_FieldVersions, 0)
	for rows.Next() {
		fv := FieldVersion{}
		err = rows.Scan(&fv.Id, &fv.Field_id, &fv.Name, &fv.Description, &fv.Field_type, &fv.Data_type, &fv.Validation_rules, &fv.Is_required, &fv.Table_display, &fv.Is_hidden,
			&fv.Formula, &fv.Target, &fv.Version, &fv.Screen_order, &fv.Position, &fv.Status, &fv.Created_by_user_id, &fv.Date_created)
		if err != nil {
			rows.Close()
			return nil, err
		}
		fieldversionsarray = append(fieldversionsarray, &fv)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &fieldversionsarray, nil
}

type CaptureApp struct {
	Id                 uint32
	App_type           string
	Name               string
	Published_version  float32
	Created_by_user_id uint32
	Date_created       time.Time
	Date_updated       time.Time
}

type Array_CaptureApps []*CaptureApp

func (d *SqlDB) Load_capture_app_by_id_version(form_id uint32, version uint32, status string) (error, *Array_CaptureApps) {
	return d.load_capture_app_sql(fmt.Sprintf("select top 1 * from [dbo].[capture_app] where id = %d and published_version = %d;", form_id, status, version))
}

func (d *SqlDB) load_capture_app_sql(sql string) (error, *Array_CaptureApps) {
	rows, err := d.db.Query(sql)
	if err != nil {
		return err, nil
	}

	capture_apps := make(Array_CaptureApps, 0)
	for rows.Next() {
		ca := CaptureApp{}
		err = rows.Scan(&ca.Id, &ca.App_type, &ca.Name, &ca.Published_version, &ca.Created_by_user_id, &ca.Date_created, &ca.Date_updated)
		if err != nil {
			rows.Close()
			return err, nil
		}
		capture_apps = append(capture_apps, &ca)
	}
	rows.Close()
	if err = rows.Err(); err != nil {
		return err, nil
	}
	return nil, &capture_apps
}

// func (d *SqlDB) Update_form_version_JSON(form_id, version uint32, jsonstr string) (uint32, error) {
// 	res, err := d.db.Exec("Update form_version set form_JSON = @p1 where form_id = @p2 AND version = @p3;", jsonstr, form_id, version)
// 	if err != nil {
// 		return 0, err
// 	}

// 	id, err := res.RowsAffected()
// 	if err != nil {
// 		return 0, err
// 	}
// 	return uint32(id), nil
// }
