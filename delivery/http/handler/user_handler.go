package handler

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"trail1/allEntityActions/user"
	"trail1/authFiles/csrfToken"
	"trail1/authFiles/permission"
	"trail1/authFiles/session"
	"trail1/entity"
	"trail1/frontend/form"
)


type UserHandler struct {
	tmpl           *template.Template
	userService    user.UserService
	sessionService user.SessionService
	userSess       *entity.Session
	loggedInUser   *entity.User
	userRole       user.RoleService
	csrfSignKey    []byte
}

type contextKey string

var ctxUserSessionKey = contextKey("signed_in_user_session")
var name, email, phone, pass string
var id, roleid uint

var cid string

// NewUserHandler returns new UserHandler object
func NewUserHandler(t *template.Template, usrServ user.UserService,
	sessServ user.SessionService, uRole user.RoleService,
	usrSess *entity.Session, csKey []byte) *UserHandler {
	return &UserHandler{tmpl: t, userService: usrServ, sessionService: sessServ,
		userRole: uRole, userSess: usrSess, csrfSignKey: csKey}
}

func (uh *UserHandler) Authenticated(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ok := uh.loggedIn(r)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), ctxUserSessionKey, uh.userSess)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// Authorized checks if a user has proper authority to access a give route
func (uh *UserHandler) Authorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if uh.loggedInUser == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		roles, errs := uh.userService.UserRoles(uh.loggedInUser)
		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		for _, role := range roles {
			permitted := permission.HasPermission(r.URL.Path, role.Name, r.Method)
			if !permitted {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
		}
		if r.Method == http.MethodPost {
			ok, err := csrfToken.ValidCSRF(r.FormValue("_csrf"), uh.csrfSignKey)
			if !ok || (err != nil) {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (uh *UserHandler) AuthorizedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if uh.loggedInUser == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		roles, errs := uh.userService.UserRoles(uh.loggedInUser)
		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		for _, role := range roles {
			permitted := permission.HasPermissionUser(r.URL.Path, role.Name, r.Method)
			if !permitted {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
		}
		if r.Method == http.MethodPost {
			ok, err := csrfToken.ValidCSRF(r.FormValue("_csrf"), uh.csrfSignKey)
			if !ok || (err != nil) {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// Login hanldes the GET/POST /login requests
func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	token, err := csrfToken.CSRFToken(uh.csrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	if r.Method == http.MethodGet {
		loginForm := struct {
			Values  url.Values
			VErrors form.ValidationErrors
			CSRF    string
		}{
			Values:  nil,
			VErrors: nil,
			CSRF:    token,
		}
		uh.tmpl.ExecuteTemplate(w, "login.html", loginForm)
		return
	}

	if r.Method == http.MethodPost {
		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		loginForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		usr, errs := uh.userService.Login(r.FormValue("email"))
		if len(errs) > 0 {
			loginForm.VErrors.Add("generic", "Your Email Address and/or Password is Wrong")
			uh.tmpl.ExecuteTemplate(w, "login.html", loginForm)
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(r.FormValue("password")))
		if err == bcrypt.ErrMismatchedHashAndPassword {
			loginForm.VErrors.Add("generic", "Your Email Address and/or Password is Wrong")
			uh.tmpl.ExecuteTemplate(w, "login.html", loginForm)
			return
		}
		uh.loggedInUser = usr
		log.Println("Usersess login:", uh.userSess)
		claims := csrfToken.Claims(usr.Email, uh.userSess.Expires)
		session.Create(claims, uh.userSess.UUID, uh.userSess.SigningKey, w)
		newSess, errs := uh.sessionService.StoreSession(uh.userSess)
		log.Println("Usersess login:", uh.userSess)
		if len(errs) > 0 {
			loginForm.VErrors.Add("generic", "Failed to Store Session")
			uh.tmpl.ExecuteTemplate(w, "login.layout", loginForm)
			return
		}
		uh.userSess = newSess
		roles, _ := uh.userService.UserRoles(usr)
		if uh.checkAdmin(roles) {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		}
		cid = fmt.Sprint(usr.ID)
		link := "/?userid=" + fmt.Sprint(usr.ID)
		http.Redirect(w, r, link, http.StatusSeeOther)
	}
}

// Logout hanldes the POST /logout requests
func (uh *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userSess, _ := r.Context().Value(ctxUserSessionKey).(*entity.Session)
	session.Remove(userSess.UUID, w)
	uh.sessionService.DeleteSession(userSess.UUID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Signup hanldes the GET/POST /registrationprocess1 requests
func (uh *UserHandler) Signupadmin(w http.ResponseWriter, r *http.Request) {
	token, err := csrfToken.CSRFToken(uh.csrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	if r.Method == http.MethodGet {
		signUpForm := struct {
			Values  url.Values
			VErrors form.ValidationErrors
			CSRF    string
		}{
			Values:  nil,
			VErrors: nil,
			CSRF:    token,
		}
		err := uh.tmpl.ExecuteTemplate(w, "Registrationform.html", signUpForm)
		if err != nil {
			panic(err.Error())
		}
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
		//log.Println("Got here 190")
		signUpForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		signUpForm.Required("name", "email", "password", "confpass")
		signUpForm.MatchesPattern("email", form.EmailRX)
		signUpForm.MatchesPhonePattern("phone", form.PhoneRX)
		signUpForm.MinLength("password", 8)
		signUpForm.PasswordMatches("password", "confpass")
		signUpForm.CSRF = token
		//If there are any errors, redisplay the signup form.
		if !signUpForm.Valid() {
			err = uh.tmpl.ExecuteTemplate(w, "rform.layout", signUpForm)
			if err != nil {
				panic(err.Error())
			}
			return
		}

		pExists := uh.userService.PhoneExists(r.FormValue("phone"))
		if pExists {
			signUpForm.VErrors.Add("phone", "Phone Already Exists")
			err := uh.tmpl.ExecuteTemplate(w, "rform.layout", signUpForm)
			if err != nil {
				panic(err.Error())
			}
			return
		}

		eExists := uh.userService.EmailExists(r.FormValue("email"))
		if eExists {
			signUpForm.VErrors.Add("email", "Email Already Exists")
			err := uh.tmpl.ExecuteTemplate(w, "rform.layout", signUpForm)
			if err != nil {
				panic(err.Error())
			}
			return
		}
		role, errs := uh.userRole.RoleByName("ADMIN")

		if len(errs) > 0 {
			signUpForm.VErrors.Add("role", "could not assign role to the user")
			err := uh.tmpl.ExecuteTemplate(w, "rform.layout", signUpForm)
			if err != nil {
				panic(err.Error())
			}
			return
		}
		user := &entity.User{
			Name:     r.FormValue("name"),
			Email:    r.FormValue("email"),
			Phone:    r.FormValue("phone"),
			Password: r.FormValue("password"),
			RoleID:   role.ID,
		}

		name = user.Name
		email = user.Email
		phone = user.Phone
		pass = user.Password
		id = user.ID
		roleid = user.RoleID

		hostURL := "smtp.gmail.com"
		hostPort := "587"
		emailSender := "jermayah7@gmail.com"
		password := "qnzfgwbnaxykglvu"
		emailReceiver := user.Email

		emailAuth := smtp.PlainAuth(
			"",
			emailSender,
			password,
			hostURL,
		)

		msg := []byte("To: " + emailReceiver + "\r\n" +
			"Subject: " + "Hello " + user.Name + "\r\n" +
			"This is your OTP. 123456789")

		err = smtp.SendMail(
			hostURL+":"+hostPort,
			emailAuth,
			emailSender,
			[]string{emailReceiver},
			msg,
		)

		if err != nil {
			fmt.Print("Error: ", err)
		}
		fmt.Print("Email Sent")

		_ = uh.tmpl.ExecuteTemplate(w, "Registrationformpart2.html", user)
	}
}

func (uh *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	token, err := csrfToken.CSRFToken(uh.csrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	if r.Method == http.MethodGet {
		signUpForm := struct {
			Values  url.Values
			VErrors form.ValidationErrors
			CSRF    string
		}{
			Values:  nil,
			VErrors: nil,
			CSRF:    token,
		}
		err := uh.tmpl.ExecuteTemplate(w, "rform.html", signUpForm)
		if err != nil {
			panic(err.Error())
		}
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
		//log.Println("Got here 190")
		signUpForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		signUpForm.Required("name", "email", "password", "confpass")
		signUpForm.MatchesPattern("email", form.EmailRX)
		signUpForm.MatchesPhonePattern("phone", form.PhoneRX)
		signUpForm.MinLength("password", 8)
		signUpForm.PasswordMatches("password", "confpass")
		signUpForm.CSRF = token
		//If there are any errors, redisplay the signup form.
		if !signUpForm.Valid() {
			err = uh.tmpl.ExecuteTemplate(w, "rform.layout", signUpForm)
			if err != nil {
				panic(err.Error())
			}
			return
		}

		pExists := uh.userService.PhoneExists(r.FormValue("phone"))
		if pExists {
			signUpForm.VErrors.Add("phone", "Phone Already Exists")
			err := uh.tmpl.ExecuteTemplate(w, "rform.layout", signUpForm)
			if err != nil {
				panic(err.Error())
			}
			return
		}

		eExists := uh.userService.EmailExists(r.FormValue("email"))
		if eExists {
			signUpForm.VErrors.Add("email", "Email Already Exists")
			err := uh.tmpl.ExecuteTemplate(w, "rform.layout", signUpForm)
			if err != nil {
				panic(err.Error())
			}
			return
		}
		role, errs := uh.userRole.RoleByName("USER")

		if len(errs) > 0 {
			signUpForm.VErrors.Add("role", "could not assign role to the user")
			err := uh.tmpl.ExecuteTemplate(w, "rform.layout", signUpForm)
			if err != nil {
				panic(err.Error())
			}
			return
		}
		user := &entity.User{
			Name:     r.FormValue("name"),
			Email:    r.FormValue("email"),
			Phone:    r.FormValue("phone"),
			Password: r.FormValue("password"),
			RoleID:   role.ID,
		}

		name = user.Name
		email = user.Email
		phone = user.Phone
		pass = user.Password
		id = user.ID
		roleid = user.RoleID
		log.Println("roleid :", roleid)
		hostURL := "smtp.gmail.com"
		hostPort := "587"
		emailSender := "jermayah7@gmail.com"
		password := "ba23atalwarad"
		emailReceiver := user.Email

		emailAuth := smtp.PlainAuth(
			"",
			emailSender,
			password,
			hostURL,
		)

		msg := []byte("To: " + emailReceiver + "\r\n" +
			"Subject: " + "Hello " + user.Name + "\r\n" +
			"This is your OTP. 123456789")

		err = smtp.SendMail(
			hostURL+":"+hostPort,
			emailAuth,
			emailSender,
			[]string{emailReceiver},
			msg,
		)

		if err != nil {
			fmt.Print("Error: ", err)
		}
		fmt.Print("Email Sent")

		_ = uh.tmpl.ExecuteTemplate(w, "Registrationformpart2.html", user)
	}
}

//Second stage registration /Registration2
func (uh *UserHandler) Registration2(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Redirect(w, req, "/Registpage", http.StatusSeeOther)
		return
	}
	otp := req.FormValue("otpfield")

	usrinfo := &entity.User{ID: id, Name: name, Email: email, Phone: phone, Password: pass, RoleID: roleid}

	if otp == "123456789" {
		_, err := uh.userService.StoreUser(usrinfo)
		if err != nil {
			http.Redirect(w, req, "/Registration2", http.StatusSeeOther)
		}
		http.Redirect(w, req, "/Loginpage", http.StatusSeeOther)
	} else {
		fmt.Print("Wrong otp")
		http.Redirect(w, req, "/Registpage", http.StatusSeeOther)
	}
	http.Redirect(w, req, "/Registpage", http.StatusSeeOther)
	return

}
func (uh *UserHandler) loggedIn(r *http.Request) bool {
	if uh.userSess == nil {
		return false
	}
	userSess := uh.userSess
	c, err := r.Cookie(userSess.UUID)
	if err != nil {
		return false
	}
	ok, err := session.Valid(c.Value, userSess.SigningKey)

	if !ok || (err != nil) {
		return false
	}
	return true
}

func (uh *UserHandler) checkAdmin(rs []entity.Role) bool {
	for _, r := range rs {
		if strings.ToUpper(r.Name) == strings.ToUpper("Admin") {
			return true
		}
	}
	return false
}

func (uh *UserHandler) Users(w http.ResponseWriter, r *http.Request) {
	users, errs := uh.userService.Users()
	if errs != nil {
		panic(errs)
	}
	uh.tmpl.ExecuteTemplate(w, "admin.users.layout", users)
}

func (uh *UserHandler) User(w http.ResponseWriter, r *http.Request) {
	token, err := csrfToken.CSRFToken(uh.csrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	if r.Method == http.MethodGet {
		userr :=uh.loggedInUser
		uuid := userr.ID
		uidd := int(uuid)
		log.Println("UUIDddddddddddddddddddddddddddddddd", uuid)
		idraw := r.URL.Query().Get("id")
		id, _ := strconv.Atoi(idraw)
		log.Println("IDraw loginnnnnnnnnnnnnnnnnn", id)
		if(uidd != id) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		usr, errs := uh.userService.User(uint(id))
		if errs != nil {
			panic(errs)
		}
		userProf := struct {
			Values  url.Values
			VErrors form.ValidationErrors
			CSRF    string
			User    *entity.User
		}{
			Values:  nil,
			VErrors: nil,
			CSRF:    token,
			User:    usr,
		}

		uh.tmpl.ExecuteTemplate(w, "user.index.layout", userProf)
	}
	if r.Method == http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

// UsersUpdate handles GET/POST /users/update?id={id} request
func (uh *UserHandler) UsersUpdate(w http.ResponseWriter, r *http.Request) {
	token, err := csrfToken.CSRFToken(uh.csrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	if r.Method == http.MethodGet {
		userr :=uh.loggedInUser
		uuid := userr.ID
		uidd := int(uuid)
		log.Println("UUIDddddddddddddddddddddddddddddddd", uuid)
		idRaw := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idRaw)
		log.Println("IDraw loginnnnnnnnnnnnnnnnnn", id)
		if(uidd != id) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		user, errs := uh.userService.User(uint(id))
		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		values := url.Values{}
		values.Add("userid", idRaw)
		values.Add("name", user.Name)
		values.Add("email", user.Email)
		values.Add("phone", user.Phone)

		upAccForm := struct {
			Values  url.Values
			VErrors form.ValidationErrors
			User    *entity.User
			CSRF    string
		}{
			Values:  values,
			VErrors: form.ValidationErrors{},
			User:    user,
			CSRF:    token,
		}
		err = uh.tmpl.ExecuteTemplate(w, "user.update.html", upAccForm)
		if err != nil {
			panic(err)
		}
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
		upAccForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		upAccForm.Required("name", "email", "phone")
		upAccForm.MatchesPattern("email", form.EmailRX)
		upAccForm.MatchesPhonePattern("phone", form.PhoneRX)
		upAccForm.CSRF = token
		// If there are any errors, redisplay the signup form.
		if !upAccForm.Valid() {
			uh.tmpl.ExecuteTemplate(w, "user.update.layout", upAccForm)
			return
		}
		userID := r.FormValue("userid")
		uid, err := strconv.Atoi(userID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		user, errs := uh.userService.User(uint(uid))
		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		eExists := uh.userService.EmailExists(r.FormValue("email"))
		if (user.Email != r.FormValue("email")) && eExists {
			upAccForm.VErrors.Add("email", "Email Already Exists")
			uh.tmpl.ExecuteTemplate(w, "user.update.layout", upAccForm)
			return
		}

		pExists := uh.userService.PhoneExists(r.FormValue("phone"))

		if (user.Phone != r.FormValue("phone")) && pExists {
			upAccForm.VErrors.Add("phone", "Phone Already Exists")
			uh.tmpl.ExecuteTemplate(w, "user.update.layout", upAccForm)
			return
		}
		if err != nil {
			upAccForm.VErrors.Add("role", "could not retrieve role id")
			uh.tmpl.ExecuteTemplate(w, "admin.user.update.layout", upAccForm)
			return
		}
		usr := &entity.User{
			ID:       user.ID,
			Name:     r.FormValue("name"),
			Email:    r.FormValue("email"),
			Phone:    r.FormValue("phone"),
			Password: user.Password,
			RoleID:   user.RoleID,
		}
		_, errs = uh.userService.UpdateUser(usr)
		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/userprof?id="+cid, http.StatusSeeOther)
	}
}

// UsersDelete handles Delete /users/delete?id={id} request
func (uh *UserHandler) UsersDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		idRaw := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idRaw)
		log.Println("Delete id:", id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
		_, errs := uh.userService.DeleteUser(uint(id))

		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}

		userSess, _ := r.Context().Value(ctxUserSessionKey).(*entity.Session)
		session.Remove(userSess.UUID, w)
		uh.sessionService.DeleteSession(userSess.UUID)

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// AdminUsersDelete handles Delete /admin/users/delete?id={id} request
func (uh *UserHandler) AdminUsersDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		idRaw := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idRaw)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
		_, errs := uh.userService.DeleteUser(uint(id))
		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	}
	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (uh *UserHandler) SellerProperties(w http.ResponseWriter, r *http.Request) {
	properties, errs := uh.userService.Properties()
	token, err := csrfToken.CSRFToken(uh.csrfSignKey)
	if err != nil {
		panic(errs)
	}
	tmp1Data := struct {
		Values url.Values
		VErrors form.ValidationErrors
		Properties []entity.Property
		CSRF string
	}{
		Values:nil,
		VErrors:nil,
		Properties:properties,
		CSRF:token,
	}

	err = uh.tmpl.ExecuteTemplate(w, "seller.properties.layout", tmp1Data)
	if err != nil {
		panic(err.Error())
	}
}

func (uh *UserHandler) UserPropertyDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		userr :=uh.loggedInUser
		uuid := userr.ID
		properties, errs := uh.userService.UserProperty(uuid)
		usr, errs := uh.userService.User(uuid)
		token, err := csrfToken.CSRFToken(uh.csrfSignKey)
		if err != nil {
			panic(errs)
		}
		tmp1Data := struct {
			Values url.Values
			VErrors form.ValidationErrors
			Properties []entity.Property
			CSRF string
			User    *entity.User
		}{
			Values:nil,
			VErrors:nil,
			Properties:properties,
			CSRF:token,
			User: usr,
		}
		err = uh.tmpl.ExecuteTemplate(w, "seller.properties.layout", tmp1Data)
		if err != nil {
			panic(err.Error())
		}

	}


}

func (uh *UserHandler) OrderDetail(w http.ResponseWriter, r *http.Request) {
	userr :=uh.loggedInUser
	uuid := userr.ID

	orders, errs := uh.userService.UserOrder(uuid)
	usr, errs := uh.userService.User(uuid)
	token, err := csrfToken.CSRFToken(uh.csrfSignKey)
	if err != nil {
		panic(errs)
	}
	tmp1Data := struct {
		Values url.Values
		VErrors form.ValidationErrors
		Orders []entity.Order
		CSRF string
		User    *entity.User
	}{
		Values:nil,
		VErrors:nil,
		Orders:orders,
		CSRF:token,
		User:usr,
	}
	err = uh.tmpl.ExecuteTemplate(w, "admin.order.layout", tmp1Data)
	if err != nil {
		panic(err.Error())
	}

}

func (uh *UserHandler) SellerPropertiesNew(w http.ResponseWriter, r *http.Request)  {
	token, err := csrfToken.CSRFToken(uh.csrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	userr :=uh.loggedInUser
	uuid := userr.ID

	if r.Method == http.MethodGet {
		usr, errs := uh.userService.User(uuid)
		if err != nil {
			panic(errs)
		}
		newProForm := struct {
			Values url.Values
			VErrors form.ValidationErrors
			CSRF string
			User    *entity.User
		}{
			Values: nil,
			VErrors: nil,
			CSRF: token,
			User:    usr,
		}
		err := uh.tmpl.ExecuteTemplate(w, "seller.property.new.layout", newProForm)
		if err != nil {
			panic(err.Error())
		}

	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		newProForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
		newProForm.MinLength("description", 10)
		newProForm.CSRF = token

		mf, fh, err := r.FormFile("catimg")
		mf2, fh2, err2 := r.FormFile("catimg2")
		log.Println(err2)
		mf3, fh3, err3 := r.FormFile("catimg3")
		log.Println(err3)
		mf4, fh4, err4 := r.FormFile("catimg4")
		log.Println(err4)
		if err != nil {
			newProForm.VErrors.Add("catimg", "File error")
			err := uh.tmpl.ExecuteTemplate(w, "seller.property.new.layout", newProForm)
			if err != nil {
				panic(err.Error())
			}
			return
		}

		defer mf.Close()
		userr :=uh.loggedInUser
		uuid := userr.ID
		log.Println("userid", userr)
		log.Println("userid num", uuid)

		pro := &entity.Property{}
		pro.Name = r.FormValue("name")
		pro.Quantity, _ = strconv.Atoi(r.FormValue("quantity"))
		pro.Description = r.FormValue("description")
		pro.Price, _ = strconv.ParseFloat(r.FormValue("price"), 64)
		categ, _ :=strconv.Atoi(r.FormValue("type"))
		pro.CategoryID = uint(categ)
		pro.Image = fh.Filename
		pro.Image2 = fh2.Filename
		pro.Image3 = fh3.Filename
		pro.Image4 = fh4.Filename
		pro.UserId = uint(uuid)
		writeFile(&mf, fh.Filename)
		writeFile(&mf2, fh2.Filename)
		writeFile(&mf3, fh3.Filename)
		writeFile(&mf4, fh4.Filename)

		_, errs := uh.userService.StoreProperty(pro)

		if len(errs) > 0 {
			panic(errs)
		}

		pro1 := &entity.Property{CategoryID: pro.CategoryID, ID: pro.ID}
		errs = uh.userService.StorePropertyCateg(pro1)
		if len(errs) > 0 {
			panic(errs)
		} else {
			http.Redirect(w, r, "/user/properties", http.StatusSeeOther)
		}
	}

}

var updateId uint
func (uh *UserHandler) SellerPropertiesUpdate(w http.ResponseWriter, r *http.Request)  {
	token, err := csrfToken.CSRFToken(uh.csrfSignKey)
	if err!= nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	if r.Method == http.MethodGet {
		idRaw := r.URL.Query().Get("id")
		//log.Println("Sellerpropertiesupdate id:", idRaw)
		id, err := strconv.Atoi(idRaw)

		if err != nil {
			panic(err)
		}

		pro, errs := uh.userService.Property(uint(id))
		outid = uint(id)
		log.Println("outid", outid)
		if len(errs) > 0 {
			panic(errs)
		}

		price := strconv.FormatFloat(pro.Price, 'f', 2 , 64)
		rating := strconv.FormatFloat(pro.Rating, 'f', 2, 64)
		ratercount := strconv.FormatFloat(pro.RatersCount, 'f', 2, 64)
		quan := strconv.Itoa(pro.Quantity)
		catid := pro.CategoryID
		usrId := pro.UserId
		categoid = catid
		updateId = usrId
		//uid := int(usrId)
		//updateid := strconv.Itoa(uid)
		values := url.Values{}
		values.Add("proid", idRaw)
		values.Add("name", pro.Name)
		values.Add("description", pro.Description)
		values.Add("price", price)
		values.Add("quantity", quan)
		values.Add("catimg", pro.Image)
		values.Add("catimg2", pro.Image2)
		values.Add("catimg3", pro.Image3)
		values.Add("catimg4", pro.Image4)
		values.Add("ratcount", ratercount)
		values.Add("rate", rating)
		//values.Add("userId", updateid)
		//.Add("userId", us)
		upProForm := struct {
			Values url.Values
			VErrors form.ValidationErrors
			Property *entity.Property
			CSRF string
		}{
			Values: values,
			VErrors: form.ValidationErrors{},
			Property: pro,
			CSRF: token,
		}

		err = uh.tmpl.ExecuteTemplate(w, "seller.properties.update.layout", upProForm)
		if err != nil {
			err.Error()
		}
		return
	}

	if r.Method == http.MethodPost {

		log.Println("ID", outid)
		if err!= nil {
			panic(err.Error())
		}

		quan, _ := strconv.Atoi(r.FormValue("quantity"))
		price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
		rating, _ := strconv.ParseFloat(r.FormValue("rate"), 64)
		ratercount, _ := strconv.ParseFloat(r.FormValue("ratcount"), 64)
		prop := &entity.Property{
			ID: outid,
			Name: r.FormValue("name"),
			CategoryID:categoid,
			Description:r.FormValue("description"),
			Quantity: quan,
			Price: price,
			RatersCount: ratercount,
			Rating: rating,
			Image: r.FormValue("imgname"),
			Image2: r.FormValue("imgname2"),
			Image3: r.FormValue("imgname3"),
			Image4: r.FormValue("imgname4"),
			UserId: updateId,
		}
		log.Println("Name", prop.Name)
		log.Println("Price", prop.Price)
		log.Println("Descr", prop.Description)
		log.Println("Quan", prop.Quantity)
		log.Println("Image", prop.Image)
		log.Println("Image", prop.Image2)
		log.Println("Image", prop.Image3)
		log.Println("Image", prop.Image4)
		log.Println("rate", prop.Rating)
		log.Println("count", prop.RatersCount)
		log.Println("userid", prop.UserId)

		mf, fh, err := r.FormFile("catimg")
		if err == nil {
			prop.Image = fh.Filename
			err = writeFile(&mf, prop.Image)
		}
		if mf != nil {
			defer mf.Close()
		}

		mf, fh2, err2 := r.FormFile("catimg2")
		if err2 == nil {
			prop.Image2 = fh2.Filename
			err2 = writeFile(&mf, prop.Image2)
		}
		if mf != nil {
			defer mf.Close()
		}
		mf, fh3, err3 := r.FormFile("catimg3")
		if err3 == nil {
			prop.Image3 = fh3.Filename
			err3 = writeFile(&mf, prop.Image3)
		}
		if mf != nil {
			defer mf.Close()
		}
		mf, fh4, err4 := r.FormFile("catimg4")
		if err4 == nil {
			prop.Image4 = fh4.Filename
			err4 = writeFile(&mf, prop.Image4)
		}
		if mf != nil {
			defer mf.Close()
		}

		_, errs := uh.userService.UpdateProperty(prop)

		if len(errs) > 0 {
			panic(errs)
		}
		http.Redirect(w, r, "/user/properties", http.StatusSeeOther)
		return
	}
}


func (uh *UserHandler) SellerPropertiesDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		idRaw := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idRaw)
		if err!= nil {
			panic(err)
		}
		_, errs := uh.userService.DeleteProperty(uint(id))
		if len(errs) > 0 {
			panic(err)
		}
	}
	http.Redirect(w, r, "/user/properties", http.StatusSeeOther)
}

func (uh *UserHandler) SearchProperties(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		res := r.URL.Query().Get("search")
		log.Println("Ressssss:", res)
		if len(res) == 0 {
			http.Redirect(w, r, "/", 303)
		}
		results, err := uh.userService.SearchProperty(res)
		if err != nil {
			panic(err)
		}
		uh.tmpl.ExecuteTemplate(w, "searchresults.layout", results)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (uh *UserHandler) PropertyDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		idRaw := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idRaw)
		if err != nil {
			panic(err)
		}
		pro, errs := uh.userService.Property(uint(id))

		if len(errs) > 0 {
			panic(errs)
		}
		_ = uh.tmpl.ExecuteTemplate(w, "detail.layout", pro)
	}
}


func (uh *UserHandler) Rating(w http.ResponseWriter, req *http.Request)  {
	if req.Method == http.MethodGet {
		idRaw := req.URL.Query().Get("id")
		id, _ := strconv.Atoi(idRaw)
		pro, errs := uh.userService.Property(uint(id))
		if len(errs) > 0 {
			panic(errs)
		}

		_= uh.tmpl.ExecuteTemplate(w, "ratings.html", pro)
	} else if req.Method == http.MethodPost {
		prop := &entity.Property{}
		idRaw, _ := strconv.Atoi(req.FormValue("id"))
		prop.ID = uint(idRaw)
		prop.Rating, _ = strconv.ParseFloat(req.FormValue("star"), 64)
		log.Println("prop.rating", prop.Rating)
		log.Println("prop.id", prop.ID)
		_, err := uh.userService.RateProperty(prop)
		if err != nil {
			panic(err)
		}
		http.Redirect(w, req, "/", http.StatusSeeOther)

	} else {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
}




func writeFile(mf *multipart.File, fname string) error {
	wd, err := os.Getwd()
	log.Println("Working dir", wd)
	if err != nil {
		return err
	}
	path := filepath.Join(wd, "../", "../", "frontend", "ui", "assets", "img", fname)
	image, err := os.Create(path)
	if err != nil {
		return err
	}
	defer image.Close()
	io.Copy(image, *mf)
	return nil
}

