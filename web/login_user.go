package rest_api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	sqldb "wo-infield-service/db"
)

var tenant_id string = "c568dcc7-4f04-4544-8f8f-4a39ab542d2b"
var client_id string = "827283aa-822b-4fee-8c9c-6d5e762eae1b"
var client_secret string = "uYG8Q~9TG32~UVHqU4QTDvk2QAZ2jyX4ukroWbIr"

// ///////////////////////////////////////////////////////////////////////////////////
// // Login User
type LoginUser_Request struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	MobileDeviceId string `json:"mobile_device_id"`
}

type LoginUser_Response struct {
	User  LoginUser_ResponseUser   `json:"user"`
	Forms []LoginUser_ResponseForm `json:"forms"`
}

type LoginUser_ResponseUser struct {
	Id       uint32 `json:"id"`
	Name     string `json:"name"`
	Position string `json:"position"`
	JWT      string `json:"jwt"`
}

type LoginUser_ResponseForm struct {
	Id          uint32  `json:"id"`
	Version     float32 `json:"version"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Form_type   string  `json:"form_type"`
}

// ////////////////////////////////////////////////////////////////////////////////////
// Oauth2
type Oauth2ROPCResponseOK struct {
	Token_type     string `json:"token_type"`
	Scope          string `json:"scope"`
	Expires_in     uint32 `json:"expires_in"`
	Ext_expires_in uint32 `json:"ext_expires_in"`
	Access_token   string `json:"access_token"`
	Refresh_token  string `json:"refresh_token"`
	Id_token       string `json:"id_token"`
}

type Oauth2ROPCResponseError struct {
	Error             string   `json:"error"`
	Error_description string   `json:"error_description"`
	Error_codes       []uint32 `json:"error_codes"`
	Timestamp         string   `json:"timestamp"`
	Trace_id          string   `json:"trace_id"`
	Correlation_id    string   `json:"correlation_id"`
	Error_uri         string   `json:"error_uri"`
}

func AddAndBuildUserLoginResponse(user_name string, passwd string, card_id string, auth_type string, device_uuid string, w http.ResponseWriter) {
	userid, err := sqldb.AddUser(user_name, passwd, card_id, auth_type, device_uuid)
	if err != nil {
		w.WriteHeader(401)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("Failed to add new user to DB: %v", err)})
	} else {
		BuildUserLoginResponse(userid, user_name, device_uuid, w)
	}
}

func BuildUserLoginResponse(userid uint32, username string, device_uuid string, w http.ResponseWriter) {
	// Generate JWT token
	jwtToken, err := GenerateJWT(userid, device_uuid, ecdsakey)
	if err != nil {
		w.WriteHeader(401)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("Failed to generate JWT: %v", err)})
	} else {

		// Get Form Data
		frmArray := make([]LoginUser_ResponseForm, 0)
		forms, err := sqldb.GetForms(userid)
		if err == nil {
			for _, form := range *forms {
				frm := LoginUser_ResponseForm{form.Id, form.Version, form.Name, form.Description, "Generated"}
				frmArray = append(frmArray, frm)
			}
		}

		// Get User Data TODO
		usr := LoginUser_ResponseUser{userid, username, "Exobiologist", jwtToken}
		lresp := LoginUser_Response{usr, frmArray}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(lresp)
	}

}

func DisplayResponseInWebBrowser(res *http.Response) {
	resbod := []byte{}
	tmpbuf := make([]byte, 256)
	offset := 0
	for {
		numred, _ := res.Body.Read(tmpbuf)
		if numred == 0 {
			break
		}
		resbod = append(resbod, tmpbuf...)
		offset += numred
	}
	os.WriteFile("oathres.html", resbod, 0644)
	args := []string{"/c", "start", "oathres.html"}
	exec.Command("cmd", args...).Start()
}

func QueryOauth2ROPC(username string, password string, device_uuid string, userid uint32, w http.ResponseWriter) {

	apiUrl := "https://login.microsoftonline.com"
	resource := fmt.Sprintf("/%s/oauth2/v2.0/token", tenant_id)
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("scope", "openid")
	data.Set("client_id", client_id)
	data.Set("client_secret", client_secret)
	data.Set("username", username)
	data.Set("password", password)

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()

	r, _ := http.NewRequest(http.MethodGet, urlStr, strings.NewReader(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		w.WriteHeader(401)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("Oauth2 http request failed. %v", err)})
		fmt.Printf("Oauth2 request response: %v, %v http request failed.", res.StatusCode, err)
	} else {
		if res.StatusCode == http.StatusOK { // Oauth2ROPCResponseOK
			var oauthres Oauth2ROPCResponseOK
			err := json.NewDecoder(res.Body).Decode(&oauthres)
			if err != nil {
				var oauthreserror Oauth2ROPCResponseError
				err := json.NewDecoder(res.Body).Decode(&oauthreserror)
				if err != nil {
					DisplayResponseInWebBrowser(res)
					w.WriteHeader(401)
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("Failed to decode ROPC response JSON. %v", err)})
				} else {
					w.WriteHeader(401)
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("Oauth2 request failed. %+v", oauthreserror)})

				}

			} else {
				// all good user found via SSO and login credentials good
				// User NOT in WO DB
				// Create user in WO DB and return all good to mobile app
				if userid == 0 {
					AddAndBuildUserLoginResponse(username, password, "000", "sso", device_uuid, w)
				} else {
					BuildUserLoginResponse(userid, username, device_uuid, w)
				}
				// Else user in DB
				// return all good to mobile app
				Sugar.Infof("Oauth2 request response: %v. Success.", res.StatusCode)
			}
		} else {

			var oauthreserror Oauth2ROPCResponseError
			err := json.NewDecoder(res.Body).Decode(&oauthreserror)
			if err != nil {
				DisplayResponseInWebBrowser(res)
				w.WriteHeader(401)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("Failed to decode ROPC response JSON. Status %v\n%v", res.StatusCode, err)})
			} else {
				w.WriteHeader(401)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("Oauth2 request failed. Status %v\n%+v", res.StatusCode, oauthreserror)})

			}
		}
	}
}

func QueryOauth2AuthUserCredentials(username string, password string, device_uuid string, userid uint32, w http.ResponseWriter) {

	apiUrl := "https://login.microsoftonline.com"
	resource := fmt.Sprintf("/%s/oauth2/v2.0/authorize", "c568dcc7-4f04-4544-8f8f-4a39ab542d2b")
	data := url.Values{}
	data.Set("response_type", "id_token")
	//data.Set("redirect_uri", "http://ec2-13-53-42-29.eu-north-1.compute.amazonaws.com:80/api/v1/oauthcallback")
	data.Set("client_id", "827283aa-822b-4fee-8c9c-6d5e762eae1b")
	data.Set("response_mode", "code")
	data.Set("scope", "openid")
	data.Set("state", "12345")
	data.Set("nonce", "678910")
	//data.Set("prompt", "none")
	data.Set("login_hint", "lcsassays@kamoacopper.com")

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()

	r, _ := http.NewRequest(http.MethodGet, urlStr, strings.NewReader(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		w.WriteHeader(401)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("Oauth2 http request failed. %v", err)})
		fmt.Printf("Oauth2 request response: %v, %v http request failed.", res.StatusCode, err)
	} else {
		if res.StatusCode == http.StatusOK { // Oauth2ROPCResponseOK
			var oauthres Oauth2ROPCResponseOK
			err := json.NewDecoder(res.Body).Decode(&oauthres)
			if err != nil {
				// Write reponse body to fiel and open in web browser
				resbod := []byte{}
				tmpbuf := make([]byte, 256)
				offset := 0
				for {
					numred, _ := res.Body.Read(tmpbuf)
					if numred == 0 {
						break
					}
					resbod = append(resbod, tmpbuf...)
					offset += numred
				}
				os.WriteFile("oathres.html", resbod, 0644)
				args := []string{"/c", "start", "oathres.html"}
				exec.Command("cmd", args...).Start()
				Sugar.Infof("Oauth2 request response: %v. Failed %v", res.StatusCode, string(resbod))

				w.WriteHeader(401)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(Error_Response{string(resbod)})
			} else {
				// all good user found via SSO and login credentials good
				// User NOT in WO DB
				// Create user in WO DB and return all good to mobile app
				if userid == 0 {
					AddAndBuildUserLoginResponse(username, password, "000", "sso", device_uuid, w)
				} else {
					BuildUserLoginResponse(userid, username, device_uuid, w)
				}
				// Else user in DB
				// return all good to mobile app
				Sugar.Infof("Oauth2 request response: %v. Success.", res.StatusCode)
			}
		} else {
			w.WriteHeader(401)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("Oauth2 http request failed. Status=%v", res.StatusCode)})
		}
	}
}

func OauthCallback(w http.ResponseWriter, r *http.Request) {
	resbod := make([]byte, 100000)
	numred, _ := r.Body.Read(resbod)
	resbod[numred] = 0
	Sugar.Infof("Oauth2 callback request: \n%s\n", string(resbod))

}

func loginUser(w http.ResponseWriter, r *http.Request) {
	var req LoginUser_Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("Failed to comprehend request: %v", err)})
	} else {

		// Check Mobile device uuid
		if bAllGood, err := sqldb.CheckDevice(req.MobileDeviceId); !bAllGood || err != nil {
			w.WriteHeader(401)
			w.Header().Set("Content-Type", "application/json")
			if err == nil {
				json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("Device %s not found.", req.MobileDeviceId)})
			} else {
				json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("Device %s not found or DB error: %v", req.MobileDeviceId, err)})
			}
		} else {

			// if user exists AND user auth type == manual
			if users, err := sqldb.GetUser(req.Username); err == nil {
				if len(*users) > 0 {
					if (*users)[0].Auth_type == "manual" {
						// Authenticate User/Password against DB
						responsemsg, userid, err1 := sqldb.LoginUser(req.Username, req.Password, req.MobileDeviceId)
						if err1 != nil {
							w.WriteHeader(401)
							w.Header().Set("Content-Type", "application/json")
							json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("Failed to login user(DB Error): %v", responsemsg)})
						} else {
							if userid == 0 {
								w.WriteHeader(401)
								w.Header().Set("Content-Type", "application/json")
								json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("Failed to login user(Auth Error): %v", responsemsg)})
							} else {

								// if DB login fails => status 401
								// if mobile device UUID not found => status 400
								// else status 200
								BuildUserLoginResponse(userid, req.Username, req.MobileDeviceId, w)
							}
						}
					} else if (*users)[0].Auth_type == "sso" {
						// Do MS365 Oauth SSO
						// Line breaks for legibility only
						QueryOauth2ROPC(req.Username, req.Password, req.MobileDeviceId, (*users)[0].User_id, w) // QueryOauth2AuthUserCredentials
					}
				} else {
					// user not found
					// Do MS365 Oauth SSO
					// ROPC
					QueryOauth2ROPC(req.Username, req.Password, req.MobileDeviceId, 0, w) // QueryOauth2AuthUserCredentials
				}
			} else {
				w.WriteHeader(401)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("User not found(DB Error): %v", err.Error())})
			}
		}
	}
}
