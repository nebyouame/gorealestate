package repository



import (
	"errors"
	"github.com/jinzhu/gorm"
	"log"
	"math"
	"trail1/allEntityActions/user"
	"trail1/entity"
	"golang.org/x/crypto/bcrypt"
)

type UserGormRepo struct {
	conn *gorm.DB
}

// NewUserGormRepo creates a new object of UserGormRepo
func NewUserGormRepo(db *gorm.DB) user.UserRepository {
	return &UserGormRepo{conn: db}
}

// Users return all users from the database
func (userRepo *UserGormRepo) Users() ([]entity.User, []error) {
	users := []entity.User{}
	errs := userRepo.conn.Find(&users).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return users, errs
}
func (userRepo *UserGormRepo) Login(email string) (*entity.User, []error) {
	log.Println(email)

	u := entity.User{}

	errs := userRepo.conn.First(&u, &entity.User{Email: email}).GetErrors()

	if len(errs) > 0 {
		return nil, errs
	}
	return &u, errs
}

// User retrieves a user by its id from the database
func (userRepo *UserGormRepo) User(id uint) (*entity.User, []error) {
	usr := entity.User{}
	errs := userRepo.conn.First(&usr, id).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return &usr, errs
}

// UpdateUser updates a given user in the database
func (userRepo *UserGormRepo) UpdateUser(user *entity.User) (*entity.User, []error) {
	usr := user
	errs := userRepo.conn.Model(&user).Updates(entity.User{Name: usr.Name, Email: usr.Email, Phone: usr.Phone}).GetErrors()
	//errs := userRepo.conn.Save(usr).GetErrors()

	//errs := userRepo.conn.Exec("UPDATE users SET name=$1, email=$2, phone=$3 WHERE id=$4;",
	//	usr.Name, usr.Email, usr.Phone, usr.ID).Save(usr).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return usr, errs
}

// DeleteUser deletes a given user from the database
func (userRepo *UserGormRepo) DeleteUser(id uint) (*entity.User, []error) {
	usr, errs := userRepo.User(id)
	if len(errs) > 0 {
		return nil, errs
	}
	errs = userRepo.conn.Delete(usr, id).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return usr, errs
}

// StoreUser stores a new user into the database
func (userRepo *UserGormRepo) StoreUser(user *entity.User) (*entity.User, []error) {
	usr := user
	hashedpass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err.Error())
	}
	user.Password = string(hashedpass)

	errs := userRepo.conn.Create(usr).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return usr, errs
}

func (userRepo *UserGormRepo) ChangePassword(user *entity.User) (*entity.User, []error) {
	usr := user
	//hashedpass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	//if err != nil {
	//	panic(err.Error())
	//}
	errs := userRepo.conn.Model(&user).Updates(entity.User{Password: usr.Password}).GetErrors()

	if len(errs) > 0 {
		return nil, errs
	}
	return usr, errs
}

// PhoneExists check if a given phone number is found
func (userRepo *UserGormRepo) PhoneExists(phone string) bool {
	user := entity.User{}
	errs := userRepo.conn.Find(&user, "phone=?", phone).GetErrors()
	if len(errs) > 0 {
		return false
	}
	return true
}

// EmailExists check if a given email is found
func (userRepo *UserGormRepo) EmailExists(email string) bool {
	user := entity.User{}
	errs := userRepo.conn.Find(&user, "email=?", email).GetErrors()
	if len(errs) > 0 {
		return false
	}
	return true
}

// UserRoles returns list of application roles that a given user has
func (userRepo *UserGormRepo) UserRoles(user *entity.User) ([]entity.Role, []error) {
	userRoles := []entity.Role{}
	errs := userRepo.conn.Model(user).Related(&userRoles).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return userRoles, errs
}

func (propertyRepo *UserGormRepo) Properties() ([]entity.Property, []error) {
	properties := []entity.Property{}
	errs := propertyRepo.conn.Find(&properties).GetErrors()
	//log.Println("properties", errs)
	if len(errs) > 0 {
		return nil, errs
	}
	return properties, errs
}

func (propertyRepo *UserGormRepo) Property(id uint) (*entity.Property, []error) {
	property := entity.Property{}
	errs := propertyRepo.conn.First(&property, id).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return &property, errs
}

func (propertyRepo *UserGormRepo) UserProperty(id uint) ([]entity.Property, []error) {
	property := []entity.Property{}
	log.Println("gorm property", property)
	errs := propertyRepo.conn.Where("user_id = ?", id).Find(&property).GetErrors()
	log.Println("id uint", id)
	if len(errs) > 0 {
		return nil, errs
	}
	return property, errs
}

func (propertyRepo *UserGormRepo) UserOrder(id uint) ([]entity.Order, []error) {
	order := []entity.Order{}
	log.Println("gorm order", order)
	errs := propertyRepo.conn.Where("property_id = ?", id).Find(&order).GetErrors()
	log.Println("id uint", id)
	if len(errs) > 0 {
		return nil, errs
	}
	return order, errs
}


func (propertyRepo *UserGormRepo) UpdateProperty(property *entity.Property) (*entity.Property, []error)  {
	pro := property
	errs := propertyRepo.conn.Save(pro).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return pro, errs
}

func (propertyRepo *UserGormRepo) DeleteProperty(id uint) (*entity.Property, []error) {
	pro, errs := propertyRepo.Property(id)

	if len(errs) > 0 {
		return nil, errs
	}
	errs = propertyRepo.conn.Delete(pro, id).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return pro, errs
}

func (propertyRepo *UserGormRepo) StoreProperty(property *entity.Property) (*entity.Property, []error)  {
	pro := property
	errs := propertyRepo.conn.Create(pro).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return pro, errs
}

func (propertyRepo *UserGormRepo) SearchProperty(index string) ([]entity.Property, error) {
	properties := []entity.Property{}
	log.Println("Properties", properties)
	log.Println("index in repooooooo:", index)
	//err := propertyRepo.conn.Where("name ILIKE ?", "%"+index+"&").Find(&properties).GetErrors()
	err := propertyRepo.conn.Where("name ILIKE ?", "%"+index+"%").Find(&properties).GetErrors()
	if len(err) !=0 {
		errors.New("Search Property not working")
	}
	return properties, nil
}

func (propertyRepo *UserGormRepo) RateProperty(pro *entity.Property) (*entity.Property, []error) {
	u := entity.Property{}
	property := entity.Property{}
	row := propertyRepo.conn.Select("rating").First(&property).Where("id = ?", pro.ID).Scan(&u)
	log.Println("Old rate", u.Rating)
	if row.RecordNotFound() {
		panic(row.Error)
	}

	row = propertyRepo.conn.Select("raters_count").First(&property).Where("id = ?", pro.ID).Scan(&u)
	log.Println("Old count", u.RatersCount)
	if row.RecordNotFound() {
		panic(row.Error)
	}

	newratings := ((u.Rating * u.RatersCount) + pro.Rating) / (u.RatersCount + 1)
	log.Println(newratings)
	log.Println("Pro ", pro.Rating)

	row = propertyRepo.conn.Model(&pro).Updates(entity.Property{Rating: float64((math.Round((newratings * 2 )))) / 2, RatersCount:u.RatersCount + 1 })
	if row.RowsAffected < 1 {
		return &property, []error{errors.New("Error")}
	}
	return &property, nil
}

func (propertyRepo *UserGormRepo) StorePropertyCateg(property *entity.Property) []error {
	pro := property

	err := propertyRepo.conn.Exec("Insert into property_categories (property_id, category_id) values (?, ?)", pro.ID, pro.CategoryID).GetErrors()
	if err != nil {
		return err
	}
	return nil
}

//func (propertyRepo *UserGormRepo) UserProperties(user *entity.User) ([]entity.Property, []error) {
//	properties := []entity.Property{}
//	userProp := user.ID
//
//	log.Println("userProp", user)
//	log.Println("userPropID", userProp)
//
//	errs := propertyRepo.conn.Find(&properties).GetErrors()
//
//
//
//	if len(errs) > 0 {
//		return nil, errs
//	}
//	return properties, errs
//
//}


