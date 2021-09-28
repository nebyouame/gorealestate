package main

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"html/template"
	"net/http"

	orepim "trail1/allEntityActions/order/repository"
	osrvim "trail1/allEntityActions/order/service"

	prepim "trail1/allEntityActions/propertypage/repository"
	psrvim "trail1/allEntityActions/propertypage/service"
	"trail1/authFiles/csrfToken"
	"trail1/delivery/http/handler"
	"trail1/entity"

	urepimp "trail1/allEntityActions/user/repository"
	usrvimp "trail1/allEntityActions/user/service"
)

//func createTables(dbconn *gorm.DB) []error {
//	errs := dbconn.CreateTable(&entity.User{}, &entity.Bank{}, &entity.Session{}, &entity.Product{}, &entity.Category{}, &entity.Cart{}, &entity.Order{}, &entity.Role{}).GetErrors()
//	if errs != nil {
//		return errs
//	}
// }

func createTabels(dbconn *gorm.DB) []error {
	errs := dbconn.CreateTable(&entity.User{}, &entity.Session{}, &entity.Property{}, &entity.Category{}, &entity.Order{}, &entity.Role{}).GetErrors()
	if errs != nil {
		return errs
	}
	if !dbconn.HasTable(&entity.User{}) {
		errs := dbconn.CreateTable(&entity.User{}).GetErrors()
		if errs != nil {
			return errs
		}
		return nil
	}
	if !dbconn.HasTable(&entity.Role{}) {
		errs := dbconn.CreateTable(&entity.Role{}).GetErrors()
		if errs != nil {
			return errs
		}
		return nil
	}
	if !dbconn.HasTable(&entity.Session{}) {
		errs := dbconn.CreateTable(&entity.Session{}).GetErrors()
		if errs != nil {
			return errs
		}
		return nil
	}
	if !dbconn.HasTable(&entity.Property{}) {
		errs := dbconn.CreateTable(&entity.Property{}).GetErrors()
		if errs != nil {
			return errs
		}
		return nil
	}
	if !dbconn.HasTable(&entity.Category{}) {
		errs := dbconn.CreateTable(&entity.Category{}).GetErrors()
		if errs != nil {
			return errs
		}
		return nil
	}
	if !dbconn.HasTable(&entity.Order{}) {
		errs := dbconn.CreateTable(&entity.Order{}).GetErrors()
		if errs != nil {
			return errs
		}
		return nil
	}

	return nil
}

func main()  {

	csrfSignKey := []byte(csrfToken.GenerateRandomID(32))
	//tmp1 := template.Must(template.ParseGlob("../../frontend/ui/templates/*"))
	tmp1 := template.Must(template.ParseGlob("../../frontend/ui/templates/*"))
	dbconn, err := gorm.Open("postgres", "postgres://postgres:@localhost/tri?sslmode=disable")

	//createTabels(dbconn)

	if err != nil {
		panic(err)
	}

	defer dbconn.Close()

	//dbconn.Exec("Insert into users (name, email, phone, password, role_id) values ('admin', 'admin123@gmail.com', '+251911111111', 'admin123', 1);")
	//dbconn.Exec("Insert into roles (name) values ('ADMIN')")
	//dbconn.Exec("Insert into roles (name) values ('USER')")
	//dbconn.Exec("Insert into banks (account_no, balance) values ('111111', 120000.00)")
	//dbconn.Exec("Insert into banks (account_no, balance) values ('222222', 9000.00)")
	//dbconn.Exec("Insert into banks (account_no, balance) values ('333333', 30000.00)")

	categoryRepo := prepim.NewCategoryGormRepo(dbconn)
	categoryServ := psrvim.NewCategoryService(categoryRepo)



	orderRepo := orepim.NewOrderGormRepo(dbconn)
	orderServ := osrvim.NewOrderService(orderRepo)

	sessionRepo := urepimp.NewSessionGormRepo(dbconn)
	sessionSrv := usrvimp.NewSessionService(sessionRepo)

	//propertyRepo := prepim.NewPropertyGormRepo(dbconn)
	//propertyServ := psrvim.NewPropertyService(propertyRepo)

	userRepo := urepimp.NewUserGormRepo(dbconn)
	userServ := usrvimp.NewUserService(userRepo)

	roleRepo := urepimp.NewRoleGormRepo(dbconn)
	roleServ := usrvimp.NewRoleService(roleRepo)

	sess := ConfigSessions()


	ach := handler.NewAdminCategoryHandler(tmp1, categoryServ, csrfSignKey)
	uh := handler.NewUserHandler(tmp1, userServ, sessionSrv, roleServ, sess, csrfSignKey)
	oh := handler.NewOrderHandler(tmp1, orderServ, userServ, csrfSignKey)
	//sph := handler.NewSellerPropertyHandler(tmp1, propertyServ, csrfSignKey)
	mh := handler.NewMenuHandler(tmp1, userServ, csrfSignKey)
	arh := handler.NewAdminRoleHandler(roleServ)

	fs := http.FileServer(http.Dir("../../frontend/ui/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))


	http.HandleFunc("/", mh.Index)
	http.Handle("/admin", uh.Authenticated(uh.Authorized(http.HandlerFunc(mh.Admin))))
	http.HandleFunc("/Loginpage", mh.LoginPage)
	http.HandleFunc("/Registpage", mh.RegistPage)
	http.HandleFunc("/RegistpageAdmin", mh.RegistPageAdmin)
	//http.HandleFunc("/userrequest", mh.)
	http.HandleFunc("/userRequest", mh.Request)
	//http.HandleFunc("/requestorder", mh)

	http.Handle("/admin/users", uh.Authenticated(uh.Authorized(http.HandlerFunc(uh.Users))))
	http.Handle("/admin/categories", uh.Authenticated(uh.Authorized(http.HandlerFunc(ach.AdminCategories))))
	http.Handle("/admin/properties", uh.Authenticated(uh.Authorized(http.HandlerFunc(uh.SellerProperties))))
	http.Handle("/admin/categories/new", uh.Authenticated(uh.Authorized(http.HandlerFunc(ach.AdminCategoriesNew))))
	http.Handle("/admin/categories/update", uh.Authenticated(uh.Authorized(http.HandlerFunc(ach.AdminCategoriesUpdate))))
	http.Handle("/admin/categories/delete", uh.Authenticated(uh.Authorized(http.HandlerFunc(ach.AdminCategoriesDelete))))


	//http.Handle("/admin/orders", uh.Authenticated(uh.Authorized(http.HandlerFunc(uh.OrderDetail))))
	//http.Handle("/admin/order", uh.Authenticated(uh.Authorized(http.HandlerFunc(oh.GetUserOrder))))
	//http.Handle("/admin/order/delete", uh.Authenticated(uh.Authorized(http.HandlerFunc(oh.OrderDelete))))


	//http.Handle("/order/delete", http.HandlerFunc(oh.OrderDelete))
	http.Handle("/orderinfo", http.HandlerFunc(oh.OrderNew))

	http.HandleFunc("/category", ach.UserCateg)

	//http.Handle("/admin/properties", uh.Authenticated(uh.Authorized(http.HandlerFunc(uh.SellerProperties))))
	//http.Handle("/admin/properties/new", uh.Authenticated(uh.Authorized(http.HandlerFunc(uh.SellerPropertiesNew))))
	//http.Handle("/admin/properties/update", uh.Authenticated(uh.Authorized(http.HandlerFunc(sph.SellerPropertiesUpdate))) )
	//http.Handle("/admin/properties/delete", uh.Authenticated(uh.Authorized(http.HandlerFunc(sph.SellerPropertiesDelete))))


	http.Handle("/admin/roles/new", uh.Authenticated(uh.Authorized(http.HandlerFunc(arh.PostRole))))
	http.Handle("/admin/roles", uh.Authenticated(uh.Authorized(http.HandlerFunc(arh.GetRoles))))
	http.Handle("/admin/role", uh.Authenticated(uh.Authorized(http.HandlerFunc(arh.GetSingleRole))))
	http.Handle("/admin/roles/update", uh.Authenticated(uh.Authorized(http.HandlerFunc(arh.PutRole))))
	http.Handle("/admin/roles/delete", uh.Authenticated(uh.Authorized(http.HandlerFunc(arh.DeleteRole))))

	http.HandleFunc("/registrationprocess1", uh.Signup)
	http.HandleFunc("/registrationprocess11", uh.Signupadmin)

	http.HandleFunc("/Registration2", uh.Registration2)

	http.HandleFunc("/login", uh.Login)
	http.HandleFunc("/searchProperties", uh.SearchProperties)
	http.HandleFunc("/detail", uh.PropertyDetail)


	http.HandleFunc("/rate", uh.Rating)
	http.Handle("/userprof", uh.Authenticated(uh.AuthorizedUser(http.HandlerFunc(uh.User))))
	http.Handle("/user/update", uh.Authenticated(uh.AuthorizedUser(http.HandlerFunc(uh.UsersUpdate))))
	http.HandleFunc("/user/delete", uh.UsersDelete)

	http.Handle("/user/properties", uh.Authenticated(uh.AuthorizedUser(http.HandlerFunc(uh.UserPropertyDetail))))
	http.Handle("/user/propertiesNew", uh.Authenticated(uh.AuthorizedUser(http.HandlerFunc(uh.SellerPropertiesNew))))

	http.Handle("/user/properties/new", uh.Authenticated(uh.AuthorizedUser(http.HandlerFunc(uh.SellerPropertiesNew))))
	http.Handle("/user/properties/update", uh.Authenticated(uh.AuthorizedUser(http.HandlerFunc(uh.SellerPropertiesUpdate))))
	http.Handle("/user/properties/delete", uh.Authenticated(uh.AuthorizedUser(http.HandlerFunc(uh.SellerPropertiesDelete))))

	http.Handle("/user/orders", uh.Authenticated(uh.AuthorizedUser(http.HandlerFunc(uh.OrderDetail))))

	http.Handle("/logout", uh.Authenticated(http.HandlerFunc(uh.Logout)))

	http.ListenAndServe(":8080", nil)

}

func ConfigSessions() *entity.Session {
	tokenExpires := time.Now().Add(time.Minute * 30).Unix()
	sessionID := csrfToken.GenerateRandomID(32)
	signingString, err := csrfToken.GenerateRandomString(32)
	if err != nil {
		panic(err)
	}
	signingKey := []byte(signingString)

	return &entity.Session{
		Expires:    tokenExpires,
		SigningKey: signingKey,
		UUID:       sessionID,
	}
}



