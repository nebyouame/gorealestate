package handler

import (
	"html/template"
	"net/http"
	"net/url"
	"trail1/allEntityActions/user"
	"trail1/authFiles/csrfToken"
	"trail1/entity"
	"trail1/frontend/form"
)


type MenuHandler struct {
	tmp1 *template.Template
	userService    user.UserService
	csrfSignKey []byte
}

func NewMenuHandler(T *template.Template, usrServ user.UserService, csKey []byte) *MenuHandler  {
	return &MenuHandler{tmp1: T, userService: usrServ, csrfSignKey: csKey}
}

func (mh *MenuHandler) Index (w http.ResponseWriter, r *http.Request) {
	if r.URL.Path !="/"{
		http.NotFound(w, r)
		return
	}
	token, err := csrfToken.CSRFToken(mh.csrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	properties, errs := mh.userService.Properties()
	if len(errs) > 0 {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	tmp1Data := struct {
		Values url.Values
		VErrors form.ValidationErrors
		Properties []entity.Property
		CSRF string
		UserID string
	}{
		Values: nil,
		VErrors: nil,
		Properties: properties,
		CSRF: token,
		UserID: r.FormValue("userid"),
	}
	mh.tmp1.ExecuteTemplate(w, "demo.layout", tmp1Data)
}


func (mh *MenuHandler) Admin(w http.ResponseWriter, r *http.Request)  {
	token, err := csrfToken.CSRFToken(mh.csrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	tmp1Data := struct {
		Values url.Values
		VErrors form.ValidationErrors
		CSRF string
	}{
		Values: nil,
		VErrors: nil,
		CSRF: token,
	}
	mh.tmp1.ExecuteTemplate(w, "admin.index.layout", tmp1Data)
}

func (mh *MenuHandler) RegistPage(w http.ResponseWriter, req *http.Request) {
	mh.tmp1.ExecuteTemplate(w, "rform.layout", nil)
}

func (mh *MenuHandler) RegistPageAdmin(w http.ResponseWriter, req *http.Request) {
	mh.tmp1.ExecuteTemplate(w, "rformadmin.layout", nil)
}

func (mh *MenuHandler) Request(w http.ResponseWriter, req *http.Request) {
	mh.tmp1.ExecuteTemplate(w, "order.layout", nil)
}

func(mh *MenuHandler) userRequest(w http.ResponseWriter, req *http.Request) {
	mh.tmp1.ExecuteTemplate(w, "order.layout", nil)
}

func (mh *MenuHandler) LoginPage(w http.ResponseWriter, req *http.Request) {

	mh.tmp1.ExecuteTemplate(w, "login.html", nil)
}

func (mh *MenuHandler) ProductDetail(w http.ResponseWriter, req *http.Request) {
	mh.tmp1.ExecuteTemplate(w, "productdetail.html", nil)
}

func (mh *MenuHandler) PaySuccess(w http.ResponseWriter, req *http.Request) {
	mh.tmp1.ExecuteTemplate(w, "pay.success.html", nil)
}

