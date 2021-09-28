package user

import "trail1/entity"

type UserRepository interface {
	Users() ([]entity.User, []error)
	User(id uint) (*entity.User, []error)
	Login(email string) (*entity.User, []error)
	UpdateUser(user *entity.User) (*entity.User, []error)
	DeleteUser(id uint) (*entity.User, []error)
	StoreUser(user *entity.User) (*entity.User, []error)
	ChangePassword(user *entity.User) (*entity.User, []error)
	PhoneExists(phone string) bool
	EmailExists(email string) bool
	UserRoles(*entity.User) ([]entity.Role, []error)

	Properties() ([]entity.Property, []error)
	Property(id uint) (*entity.Property, []error)
	UpdateProperty(property *entity.Property) (*entity.Property, []error)
	DeleteProperty(id uint) (*entity.Property, []error)
	StoreProperty(property *entity.Property) (*entity.Property, []error)
	RateProperty(property *entity.Property) (*entity.Property, []error)
	SearchProperty(index string) ([]entity.Property, error)
	StorePropertyCateg(property *entity.Property) []error

	//UserProperties(user *entity.User) ([]entity.Property, []error)
	UserProperty(id uint) ([]entity.Property, []error)
	UserOrder(id uint) ([]entity.Order, []error)
}

// RoleRepository speifies application user role related database operations
type RoleRepository interface {
	Roles() ([]entity.Role, []error)
	Role(id uint) (*entity.Role, []error)
	RoleByName(name string) (*entity.Role, []error)
	UpdateRole(role *entity.Role) (*entity.Role, []error)
	DeleteRole(id uint) (*entity.Role, []error)
	StoreRole(role *entity.Role) (*entity.Role, []error)
}

// SessionService specifies logged in user session related service
type SessionRepository interface {
	Session(sessionID string) (*entity.Session, []error)
	StoreSession(session *entity.Session) (*entity.Session, []error)
	DeleteSession(sessionID string) (*entity.Session, []error)
}

