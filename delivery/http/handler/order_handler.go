package handler

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"trail1/allEntityActions/order"
	"trail1/allEntityActions/user"
	"trail1/authFiles/csrfToken"
	"trail1/entity"
	"trail1/frontend/form"
)

type OrderHandler struct {
	tmpl        *template.Template
	orderServ   order.OrderService
	usdrServ    user.UserService
	csrfSignKey []byte
}

var idr string

func NewOrderHandler(t *template.Template, os order.OrderService, us user.UserService, csKey []byte) *OrderHandler {
	return &OrderHandler{tmpl: t, orderServ: os, usdrServ: us, csrfSignKey: csKey}
}

func (oh *OrderHandler) Orders(w http.ResponseWriter, r *http.Request) {
	orders, _ := oh.orderServ.Orders()
	token, err := csrfToken.CSRFToken(oh.csrfSignKey)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
	tmplData := struct {
		Values  url.Values
		VErrors form.ValidationErrors
		Orders  []entity.Order
		CSRF    string
	}{
		Values:  nil,
		VErrors: nil,
		Orders:  orders,
		CSRF:    token,
	}
	err = oh.tmpl.ExecuteTemplate(w, "admin.order.layout", tmplData)
	if err != nil {
		panic(err.Error())
	}
}


func (oh *OrderHandler) GetSingleOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		idRaw := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idRaw)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
		ordr, errs := oh.orderServ.Order(uint(id))
		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		err = oh.tmpl.ExecuteTemplate(w, "admin.order.layout", ordr)
		if err != nil {
			panic(err.Error())
		}
		return
	}
	http.Redirect(w, r, "/admin/order", 303)
}

func (oh *OrderHandler) GetUserOrder(w http.ResponseWriter, r *http.Request) {
	user := &entity.User{}
	idRaw := r.URL.Query().Get("id")
	uuid := 0
	id, _ := strconv.Atoi(idRaw)
	products := []entity.Property{}
	log.Println("id:", id)
	user.ID = uint(uuid)
	// var id uint
	order, _ := oh.orderServ.CustomerOrders(user)
	productlist := strings.Split(order.ItemsID, ",")
	for i:=0; i< len(productlist); i++{
		productid := productlist[i]
		prodid, _ := strconv.Atoi(productid)
		pro, _ := oh.usdrServ.Property(uint(prodid))
		products = append(products, *pro)
	}
	tmplData := struct {
		Values  url.Values
		VErrors form.ValidationErrors
		Order   entity.Order
		Products []entity.Property
		CSRF    string
	}{
		Values:  nil,
		VErrors: nil,
		Order:   order,
		Products: products,
		// CSRF:    token,
	}
	err := oh.tmpl.ExecuteTemplate(w, "checkoutpage.html", tmplData)
	if err != nil {
		panic(err.Error())
	}
}

func (oh *OrderHandler) OrderNeww(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		uid, _ := strconv.Atoi(r.FormValue("userid"))
		t := time.Now()
		uuid := 0

		total, _ := strconv.ParseFloat(r.FormValue("total"), 64)
		//prodids, _ := strconv.Atoi(r.FormValue("prodids"))
		prodids := r.FormValue("prodids")
		log.Println("User order id: ", uid)
		log.Println("User order prod id: ", prodids)
		log.Println("User order total: ", total)
		idr = strconv.Itoa(uuid)
		propertyid,_ := strconv.Atoi(prodids)
		log.Println("User order property id1111: ", propertyid)
		propertyID,err := oh.usdrServ.Property(uint(propertyid))
		orderPropertyID := propertyID.UserId
		log.Println("User order property id2222: ", orderPropertyID)
		if len(err) > 0 {
			panic(err)
		}

		ord := &entity.Order{
			UserID:    uint(uuid),
			CreatedAt: t,
			ItemsID:   prodids,
			Total:     total,
			PropertyID:orderPropertyID,
		}
		// link += string(pid)
		_, errs := oh.orderServ.StoreOrder(ord)
		if len(errs) > 0 {
			panic(errs)
		}
	}
	http.Redirect(w, r, "/getorder?id="+idr, 303)
}

func (oh *OrderHandler) OrderDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		idRaw := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idRaw)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
		_, errs := oh.orderServ.DeleteOrder(uint(id))
		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
		http.Redirect(w, r, "/getusercart?id="+idr, http.StatusSeeOther)
	}
}



var UID uint
var reqid string
var TOTAL float64

func (oh *OrderHandler) OrderNew(w http.ResponseWriter, r *http.Request)   {
	token, err := csrfToken.CSRFToken(oh.csrfSignKey)
	if err!= nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	if r.Method == http.MethodGet {
		log.Println("got hereeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee get")
		idRaw := r.URL.Query().Get("id")
		id, _ := strconv.Atoi(idRaw)
		log.Println("orderinfo id:", id)
		ord, errs := oh.usdrServ.Property(uint(id))
		userID := ord.UserId
		total := ord.Price
		reqid = idRaw
		UID = userID
		TOTAL = total
		log.Println("reqid:", reqid)
		log.Println("uid:", UID)
		log.Println("total", total)
		if errs != nil {
			panic(errs)
		}

		newOrderForm := struct {
			Values url.Values
			VErrors form.ValidationErrors
			CSRF string
			Pro    *entity.Property
		}{
			Values: nil,
			VErrors: nil,
			CSRF: token,
			Pro:    ord,
		}
		err := oh.tmpl.ExecuteTemplate(w, "order.layout", newOrderForm)
		if err != nil {
			panic(err.Error())
		}
	}
	if r.Method == http.MethodPost {
		t := time.Now()
		err := r.ParseForm()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		orderForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		orderForm.Required("name", "email")
		orderForm.CSRF = token
		uuid := 0

		ord := &entity.Order{}
		ord.Name = r.FormValue("name")
		ord.Phone = r.FormValue("phone")
		ord.Email = r.FormValue("email")
		ord.UserID = uint(uuid)
		ord.CreatedAt = t
		ord.ItemsID = reqid
		ord.Total = TOTAL
		ord.PropertyID = UID
		log.Println("got hereeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee post")
		_, errs := oh.orderServ.StoreOrder(ord)
		if len(errs) > 0 {
			panic(errs)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
