package repository

import (
	"account-service/internal/models"
	"comet/utils"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	DB *gorm.DB
}

// NewRepository create new Repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

// CreateUserIfNotExist create user if not exist by email
func (r *Repository) CreateUserIfNotExist(email, firstName, lastName, password string, departmentID, roleID uint32) (*models.User, error) {
	var resultUser models.User
	b, err := utils.HashArgon(password)

	if err != nil {
		return nil, fmt.Errorf("utils.HashArgon error: %w", err)
	}

	hashPassword := string(b)

	result := r.DB.Where(models.User{
		Email: email,
	}).Attrs(models.User{
		FirstName:    firstName,
		LastName:     lastName,
		Password:     hashPassword,
		DepartmentID: departmentID,
		RoleID:       roleID,
		IsRegistered: true,
	}).FirstOrCreate(&resultUser)

	if result.Error != nil {
		return nil, fmt.Errorf("r.DB.FirstOrCreate error: %w", result.Error)
	}

	return &resultUser, nil
}

// GetUserByID get user by id
func (r *Repository) GetUserByID(id string) (*models.User, error) {
	var resultUser models.User
	result := r.DB.Where(models.User{ID: id}).First(&resultUser)

	if result.Error != nil {
		return nil, fmt.Errorf("r.DB.First error: %w", result.Error)
	}

	return &resultUser, nil
}

// GetUserByEmail get user by id
func (r *Repository) GetUserByEmail(email string) (*models.User, error) {
	var resultUser models.User
	result := r.DB.Where(models.User{Email: email}).First(&resultUser)

	if result.Error != nil {
		return nil, fmt.Errorf("r.DB.First error: %w", result.Error)
	}

	return &resultUser, nil
}

// ChangePasswordByID change password by id
func (r *Repository) ChangePasswordByID(id, password string) (*models.User, error) {
	b, err := utils.HashArgon(password)
	if err != nil {
		return nil, fmt.Errorf("utils.HashArgon error: %w", err)
	}
	hashPassword := string(b)

	var users []models.User
	result := r.DB.Model(&users).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Update("password", hashPassword)
	if result.Error != nil {
		return nil, fmt.Errorf("r.DB.Update error: %w", result.Error)
	}

	return &users[0], nil
}

// EmailExist check is email exist or not
func (r *Repository) EmailExist(email string) (bool, error) {
	count := int64(0)
	err := r.DB.Model(&models.User{}).
		Where("email = ?", email).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("r.DB.Count error: %w", err)
	}

	return count > 0, nil
}

// RemoveUserByID remove user by id
func (r *Repository) RemoveUserByID(id string) error {

	result := r.DB.Model(&models.User{}).Delete("id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("r.DB.Delete error: %w", result.Error)
	}

	return nil
}

// AddUserDepartment add new department if not exist
func (r *Repository) AddUserDepartment(name string) (*models.Department, error) {
	var resultDepartment models.Department
	result := r.DB.Where(models.Department{Name: name}).FirstOrCreate(&resultDepartment)
	if result.Error != nil {
		return nil, fmt.Errorf("r.DB.FirstOrCreate error: %w", result.Error)
	}

	return &resultDepartment, nil
}

// RemoveUserDepartment remove user department by id
func (r *Repository) RemoveUserDepartment(id uint32) error {
	result := r.DB.Delete(models.Department{ID: id})
	if result.Error != nil {
		return fmt.Errorf("r.DB.Delete error: %w", result.Error)
	}

	return nil
}

// GetUserDepartmentByID get department by id
func (r *Repository) GetUserDepartmentByID(id uint32) (*models.Department, error) {
	var resultDepartment models.Department
	result := r.DB.Where("id = ? ", id).First(&resultDepartment)
	if result.Error != nil {
		return nil, fmt.Errorf("r.DB.First error: %w", result.Error)
	}

	return &resultDepartment, nil
}

// GetUserDepartments get all departments
func (r *Repository) GetUserDepartments() (*[]models.Department, error) {
	var resultDepatments []models.Department
	result := r.DB.Find(&resultDepatments)

	if result.Error != nil {
		return nil, fmt.Errorf("r.DB.Find error: %w", result.Error)
	}

	return &resultDepatments, nil
}

// AddUserRole add new department if not exist
func (r *Repository) AddUserRole(name string) (*models.Role, error) {
	var resultRole models.Role
	result := r.DB.Where(models.Role{Name: name}).FirstOrCreate(&resultRole)
	if result.Error != nil {
		return nil, fmt.Errorf("r.DB.FirstOrCreate error: %w", result.Error)
	}

	return &resultRole, nil
}

// RemoveUserRole remove user department by id
func (r *Repository) RemoveUserRole(id uint32) error {
	result := r.DB.Delete(models.Role{ID: id})
	if result.Error != nil {
		return fmt.Errorf("r.DB.Delete error: %w", result.Error)
	}

	return nil
}

// GetUserRoleByID get department by id
func (r *Repository) GetUserRoleByID(id uint32) (*models.Role, error) {
	var resultRole models.Role
	result := r.DB.Where("id = ? ", id).First(&resultRole)
	if result.Error != nil {
		return nil, fmt.Errorf("r.DB.First error: %w", result.Error)
	}

	return &resultRole, nil
}

// GetUserRoles get all departments
func (r *Repository) GetUserRoles() (*[]models.Role, error) {
	var resultRoles []models.Role
	result := r.DB.Find(&resultRoles)

	if result.Error != nil {
		return nil, fmt.Errorf("r.DB.Find error: %w", result.Error)
	}

	return &resultRoles, nil
}
