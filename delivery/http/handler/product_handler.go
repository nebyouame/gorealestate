package handler

import (
	"html/template"
	"trail1/allEntityActions/propertypage"
)

var outid, categoid, userid uint
// ProductHandler handles product handler admin requests
//type SellerProductHandler struct {
//	tmpl        *template.Template
//	productSrv  productpage.ItemService
//	csrfSignKey []byte
//}
//
type SellerPropertyHandler struct {
	tmp1 *template.Template
	propertySrv propertypage.PropertyService
	csrfSignKey []byte
}


func NewSellerPropertyHandler(t *template.Template, ps propertypage.PropertyService, csKey []byte) *SellerPropertyHandler {
	return &SellerPropertyHandler{tmp1:t, propertySrv: ps, csrfSignKey:csKey}
}


//func (sph *SellerPropertyHandler) SellerProperties(w http.ResponseWriter, r *http.Request) {
//	properties, errs := sph.propertySrv.Properties()
//	token, err := csrfToken.CSRFToken(sph.csrfSignKey)
//	if err != nil {
//		panic(errs)
//	}
//	tmp1Data := struct {
//		Values url.Values
//		VErrors form.ValidationErrors
//		Properties []entity.Property
//		CSRF string
//	}{
//		Values:nil,
//		VErrors:nil,
//		Properties:properties,
//		CSRF:token,
//	}
//
//	err = sph.tmp1.ExecuteTemplate(w, "seller.properties.layout", tmp1Data)
//	if err != nil {
//		panic(err.Error())
//	}
//}

//func (sph *SellerPropertyHandler) SellerPropertiesNew(w http.ResponseWriter, r *http.Request)  {
//	token, err := csrfToken.CSRFToken(sph.csrfSignKey)
//	if err != nil {
//		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
//	}
//	if r.Method == http.MethodGet {
//		newProForm := struct {
//			Values url.Values
//			VErrors form.ValidationErrors
//			CSRF string
//		}{
//			Values: nil,
//			VErrors: nil,
//			CSRF: token,
//		}
//		err := sph.tmp1.ExecuteTemplate(w, "seller.property.new.layout", newProForm)
//		if err != nil {
//			panic(err.Error())
//		}
//	}
//
//	if r.Method == http.MethodPost {
//		err := r.ParseForm()
//		if err != nil {
//			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
//			return
//		}
//		newProForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
//		newProForm.MinLength("description", 10)
//		newProForm.CSRF = token
//
//		mf, fh, err := r.FormFile("catimg")
//		if err != nil {
//			newProForm.VErrors.Add("catimg", "File error")
//			err := sph.tmp1.ExecuteTemplate(w, "seller.property.new.layout", newProForm)
//			if err != nil {
//				panic(err.Error())
//			}
//			return
//		}
//
//		defer mf.Close()
//		pro := &entity.Property{}
//		pro.Name = r.FormValue("name")
//		pro.Quantity, _ = strconv.Atoi(r.FormValue("quantity"))
//		pro.Description = r.FormValue("description")
//		pro.Price, _ = strconv.ParseFloat(r.FormValue("price"), 64)
//		categ, _ :=strconv.Atoi(r.FormValue("type"))
//		pro.CategoryID = uint(categ)
//		pro.Image = fh.Filename
//
//		writeFile(&mf, fh.Filename)
//
//		_, errs := sph.propertySrv.StoreProperty(pro)
//
//		if len(errs) > 0 {
//			panic(errs)
//		}
//
//		pro1 := &entity.Property{CategoryID: pro.CategoryID, ID: pro.ID}
//		errs = sph.propertySrv.StorePropertyCateg(pro1)
//		if len(errs) > 0 {
//			panic(errs)
//		} else {
//			http.Redirect(w, r, "/admin/properties", http.StatusSeeOther)
//		}
//	}
//
//}



//func (sph *SellerPropertyHandler) SellerPropertiesUpdate(w http.ResponseWriter, r *http.Request)  {
//	token, err := csrfToken.CSRFToken(sph.csrfSignKey)
//	if err!= nil {
//		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
//	}
//	if r.Method == http.MethodGet {
//		idRaw := r.URL.Query().Get("id")
//		id, err := strconv.Atoi(idRaw)
//
//		if err != nil {
//			panic(err)
//		}
//
//		pro, errs := sph.propertySrv.Property(uint(id))
//		outid = uint(id)
//		log.Println("outid", outid)
//		if len(errs) > 0 {
//			panic(errs)
//		}
//
//		price := strconv.FormatFloat(pro.Price, 'f', 2 , 64)
//		rating := strconv.FormatFloat(pro.Rating, 'f', 2, 64)
//		ratercount := strconv.FormatFloat(pro.RatersCount, 'f', 2, 64)
//		quan := strconv.Itoa(pro.Quantity)
//		catid := pro.CategoryID
//		categoid = catid
//		values := url.Values{}
//		values.Add("proid", idRaw)
//		values.Add("name", pro.Name)
//		values.Add("description", pro.Description)
//		values.Add("price", price)
//		values.Add("quantity", quan)
//		values.Add("catimg", pro.Image)
//		values.Add("ratcount", ratercount)
//		values.Add("rate", rating)
//		upProForm := struct {
//			Values url.Values
//			VErrors form.ValidationErrors
//			Property *entity.Property
//			CSRF string
//		}{
//			Values: values,
//			VErrors: form.ValidationErrors{},
//			Property: pro,
//			CSRF: token,
//		}
//
//		err = sph.tmp1.ExecuteTemplate(w, "seller.properties.update.layout", upProForm)
//		if err != nil {
//			err.Error()
//		}
//		return
//	}
//
//	if r.Method == http.MethodPost {
//
//		log.Println("ID", outid)
//		if err!= nil {
//			panic(err.Error())
//		}
//
//		quan, _ := strconv.Atoi(r.FormValue("quantity"))
//		price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
//		rating, _ := strconv.ParseFloat(r.FormValue("rate"), 64)
//		ratercount, _ := strconv.ParseFloat(r.FormValue("ratcount"), 64)
//		prop := &entity.Property{
//			ID: outid,
//			Name: r.FormValue("name"),
//			CategoryID:categoid,
//			Description:r.FormValue("description"),
//			Quantity: quan,
//			Price: price,
//			RatersCount: ratercount,
//			Rating: rating,
//			Image: r.FormValue("imgname"),
//		}
//		log.Println("Name", prop.Name)
//		log.Println("Price", prop.Price)
//		log.Println("Descr", prop.Description)
//		log.Println("Quan", prop.Quantity)
//		log.Println("Image", prop.Image)
//		log.Println("rate", prop.Rating)
//		log.Println("count", prop.RatersCount)
//
//		mf, fh, err := r.FormFile("catimg")
//		if err == nil {
//			prop.Image = fh.Filename
//			err = writeFile(&mf, prop.Image)
//		}
//		if mf != nil {
//			defer mf.Close()
//		}
//
//		_, errs := sph.propertySrv.UpdateProperty(prop)
//
//		if len(errs) > 0 {
//			panic(errs)
//		}
//		http.Redirect(w, r, "/admin/properties", http.StatusSeeOther)
//		return
//	}
//	}
//
//
//
//func (sph *SellerPropertyHandler) SellerPropertiesDelete(w http.ResponseWriter, r *http.Request) {
//	if r.Method == http.MethodGet {
//		idRaw := r.URL.Query().Get("id")
//		id, err := strconv.Atoi(idRaw)
//		if err!= nil {
//			panic(err)
//		}
//		_, errs := sph.propertySrv.DeleteProperty(uint(id))
//		if len(errs) > 0 {
//			panic(err)
//		}
//	}
//	http.Redirect(w, r, "/admin/properties", http.StatusSeeOther)
//}
//
//
//func (sph *SellerPropertyHandler) SearchProperties(w http.ResponseWriter, r *http.Request) {
//	if r.Method == http.MethodGet {
//		res := r.URL.Query().Get("search")
//		if len(res) == 0 {
//			http.Redirect(w, r, "/", 303)
//		}
//		results, err := sph.propertySrv.SearchProperty(res)
//		if err != nil {
//			panic(err)
//		}
//		sph.tmp1.ExecuteTemplate(w, "searchresults.layout", results)
//	} else {
//		http.Redirect(w, r, "/", http.StatusSeeOther)
//	}
//}
//


//func (sph *SellerPropertyHandler) PropertyDetail(w http.ResponseWriter, r *http.Request) {
//	if r.Method == http.MethodGet {
//		idRaw := r.URL.Query().Get("id")
//		id, err := strconv.Atoi(idRaw)
//		if err != nil {
//			panic(err)
//		}
//		pro, errs := sph.propertySrv.Property(uint(id))
//
//		if len(errs) > 0 {
//			panic(errs)
//		}
//		_ = sph.tmp1.ExecuteTemplate(w, "detail.layout", pro)
//	}
//}
//
//
//func (sph *SellerPropertyHandler) Rating(w http.ResponseWriter, req *http.Request)  {
//	if req.Method == http.MethodGet {
//		idRaw := req.URL.Query().Get("id")
//		id, _ := strconv.Atoi(idRaw)
//		pro, errs := sph.propertySrv.Property(uint(id))
//		if len(errs) > 0 {
//			panic(errs)
//		}
//
//		_= sph.tmp1.ExecuteTemplate(w, "ratings.html", pro)
//		} else if req.Method == http.MethodPost {
//			prop := &entity.Property{}
//			idRaw, _ := strconv.Atoi(req.FormValue("id"))
//			prop.ID = uint(idRaw)
//			prop.Rating, _ = strconv.ParseFloat(req.FormValue("star"), 64)
//			log.Println("prop.rating", prop.Rating)
//			log.Println("prop.id", prop.ID)
//			_, err := sph.propertySrv.RateProperty(prop)
//			if err != nil {
//				panic(err)
//			}
//			http.Redirect(w, req, "/", http.StatusSeeOther)
//
//	} else {
//		http.Redirect(w, req, "/", http.StatusSeeOther)
//	}
//}
//
//func writeFile(mf *multipart.File, fname string) error {
//	wd, err := os.Getwd()
//	log.Println("Working dir", wd)
//	if err != nil {
//		return err
//	}
//	path := filepath.Join(wd, "../", "../", "frontend", "ui", "assets", "img", fname)
//	image, err := os.Create(path)
//	if err != nil {
//		return err
//	}
//	defer image.Close()
//	io.Copy(image, *mf)
//	return nil
//}



//
