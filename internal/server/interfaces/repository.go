package interfaces

import "account-service/internal/models"

type Repository interface {
	GetUserByID(id string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	ChangePasswordByID(id, password string) (*models.User, error)
	EmailExist(email string) (bool, error)
	RemoveUserByID(id string) error
	AddUserDepartment(name string) (*models.Department, error)
	RemoveUserDepartment(id uint32) error
	GetUserDepartmentByID(id uint32) (*models.Department, error)
	GetUserDepartments() (*[]models.Department, error)
	AddUserRole(name string) (*models.Role, error)
	RemoveUserRole(id uint32) error
	GetUserRoleByID(id uint32) (*models.Role, error)
	GetUserRoles() (*[]models.Role, error)
	//UserAuth(email, password string) models.User
	//UserVerify(email, token string) bool
	//UserTokenRemove(email, token string)
}
