package handler

import (
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"trail1/allEntityActions/propertypage"
	"trail1/authFiles/csrfToken"
	"trail1/frontend/form"
	"trail1/entity"
)

type AdminCategoryHandler struct {
	tmpl        *template.Template
	categorySrv propertypage.CategoryService
	csrfSignKey []byte
}

func NewAdminCategoryHandler(t *template.Template, cs propertypage.CategoryService, csKey []byte) *AdminCategoryHandler {
	return &AdminCategoryHandler{tmpl: t, categorySrv: cs, csrfSignKey: csKey}
}

func (ach *AdminCategoryHandler) AdminCategories(w http.ResponseWriter, r *http.Request) {
	categories, errs := ach.categorySrv.Categories()
	if errs != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
	token, err := csrfToken.CSRFToken(ach.csrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	tmplData := struct {
		Values     url.Values
		VErrors    form.ValidationErrors
		Categories []entity.Category
		CSRF       string
	}{
		Values:     nil,
		VErrors:    nil,
		Categories: categories,
		CSRF:       token,
	}
	ach.tmpl.ExecuteTemplate(w, "admin.categ.layout", tmplData)
}

func (ach *AdminCategoryHandler) AdminCategoriesNew(w http.ResponseWriter, r *http.Request) {
	token, err := csrfToken.CSRFToken(ach.csrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	if r.Method == http.MethodGet {
		newCatForm := struct {
			Values  url.Values
			VErrors form.ValidationErrors
			CSRF    string
		}{
			Values:  nil,
			VErrors: nil,
			CSRF:    token,
		}
		err := ach.tmpl.ExecuteTemplate(w, "admin.categ.new.layout", newCatForm)
		if err != nil {
			panic(err)
		}
	}

	if r.Method == http.MethodPost {
		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		// Validate the form contents
		newCatForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		newCatForm.Required("catname", "catdesc")
		newCatForm.MinLength("catdesc", 10)
		newCatForm.CSRF = token
		// If there are any errors, redisplay the signup form.
		if !newCatForm.Valid() {
			ach.tmpl.ExecuteTemplate(w, "admin.categ.new.layout", newCatForm)
			return
		}
		mf, fh, err := r.FormFile("catimg")
		if err != nil {
			newCatForm.VErrors.Add("catimg", "File error")
			ach.tmpl.ExecuteTemplate(w, "admin.categ.new.layout", newCatForm)
			return
		}
		defer mf.Close()
		ctg := &entity.Category{
			Name:        r.FormValue("catname"),
			Description: r.FormValue("catdesc"),
			Image:       fh.Filename,
		}
		writeFile(&mf, fh.Filename)
		_, errs := ach.categorySrv.StoreCategory(ctg)
		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
	}
}

func (ach *AdminCategoryHandler) AdminCategoriesUpdate(w http.ResponseWriter, r *http.Request) {
	token, err := csrfToken.CSRFToken(ach.csrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	if r.Method == http.MethodGet {
		idRaw := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idRaw)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
		cat, errs := ach.categorySrv.Category(uint(id))
		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		values := url.Values{}
		values.Add("catid", idRaw)
		values.Add("catname", cat.Name)
		values.Add("catdesc", cat.Description)
		values.Add("catimg", cat.Image)
		upCatForm := struct {
			Values   url.Values
			VErrors  form.ValidationErrors
			Category *entity.Category
			CSRF     string
		}{
			Values:   values,
			VErrors:  form.ValidationErrors{},
			Category: cat,
			CSRF:     token,
		}
		ach.tmpl.ExecuteTemplate(w, "admin.categ.update.layout", upCatForm)
		return
	}
	if r.Method == http.MethodPost {
		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		// Validate the form contents
		updateCatForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		updateCatForm.Required("catname", "catdesc")
		updateCatForm.MinLength("catdesc", 10)
		updateCatForm.CSRF = token

		catID, err := strconv.Atoi(r.FormValue("catid"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
		ctg := &entity.Category{
			ID:          uint(catID),
			Name:        r.FormValue("catname"),
			Description: r.FormValue("catdesc"),
			Image:       r.FormValue("imgname"),
		}
		mf, fh, err := r.FormFile("catimg")
		if err == nil {
			ctg.Image = fh.Filename
			err = writeFile(&mf, ctg.Image)
		}
		if mf != nil {
			defer mf.Close()
		}
		_, errs := ach.categorySrv.UpdateCategory(ctg)
		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
		return
	}
}

func (ach *AdminCategoryHandler) AdminCategoriesDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		idRaw := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idRaw)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
		_, errs := ach.categorySrv.DeleteCategory(uint(id))
		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	}
	http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
}

func (ach *AdminCategoryHandler) UserCateg(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		idRaw := r.URL.Query().Get("id")
		id, errs := strconv.Atoi(idRaw)
		if errs != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
		token, err := csrfToken.CSRFToken(ach.csrfSignKey)
		if err != nil {
			panic(errs)
		}
		cat, errr := ach.categorySrv.Category(uint(id))
		if len(errr) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		procat, errss :=ach.categorySrv.PropertiesInCategory(cat)

		if len(errss) > 0 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
		tmp1Data := struct {
			Values url.Values
			VErrors form.ValidationErrors
			Properties []entity.Property
			CSRF string
			UserID string
		}{
			Values:nil,
			VErrors:nil,
			Properties:procat,
			CSRF:token,
			UserID: r.FormValue("userid"),
		}
		err = ach.tmpl.ExecuteTemplate(w, "demo.layout", tmp1Data)
		if err != nil {
			panic(err.Error())
		}
	}
}
