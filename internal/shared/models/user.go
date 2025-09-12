package models

import (
	"gorm.io/gorm"
)

// User 定义了用户的基本信息
type User struct {
	BaseModel

	// 基本信息
	Nickname     string  `gorm:"unique;size:128;not null"`
	Email        string  `gorm:"unique;size:128"`
	QQNumber     *uint64 `gorm:"unique;index;default:null"`
	PasswordHash string  `gorm:"size:256;not null"`

	// 控制信息
	IsAdmin bool `gorm:"default:false"`
	// 补充备注
	Remark *string `gorm:"type:text;default:null"`

	// 尨译相关信息
	MoetranId  string `gorm:"type:text"`
	MoetranJwt string `gorm:"type:text"`
}

// 外部调用用的指针，避免每一次都创建
var emptyUser = &User{}

// GetUser 获取一个空的 User 实例
func GetUser() *User {
	return emptyUser
}

// TableName 返回 User 的表名
func (*User) TableName() string {
	return "users"
}

// UserSpec 定义了用户的查询条件
type UserSpec struct {
	Id *PKey

	Nickname *string
	Email    *string
	QQNumber *int

	IsAdmin *bool

	MoetranId  *string
	MoetranJwt *string
}

// Apply 将 UserSpec 应用到 WHERE 子句
func (s *UserSpec) Apply(query *gorm.DB) {
	if s.Id != nil {
		query = query.Where("id = ?", *s.Id)
	}

	if s.Nickname != nil {
		query = query.Where("nickname = ?", *s.Nickname)
	}
	if s.Email != nil {
		query = query.Where("email = ?", *s.Email)
	}
	if s.QQNumber != nil {
		query = query.Where("qq_number = ?", *s.QQNumber)
	}

	if s.IsAdmin != nil {
		query = query.Where("is_admin = ?", *s.IsAdmin)
	}
	if s.MoetranId != nil {
		query = query.Where("moetran_id = ?", *s.MoetranId)
	}
	if s.MoetranJwt != nil {
		query = query.Where("moetran_jwt = ?", *s.MoetranJwt)
	}
}

// UserFields 定义了用户的字段，用于查询时选择特定字段返回
type UserFields struct {
	// BaseModel 字段
	Id   bool
	Time bool

	// User 字段
	Nickname     bool
	Email        bool
	QQNumber     bool
	PasswordHash bool

	IsAdmin bool
	Remark  bool

	MoetranId  bool
	MoetranJwt bool
}

// Apply 将 UserFields 转换为字符串切片并应用到 SELECT 子句
func (u *UserFields) Apply(query *gorm.DB) {
	var fields []string

	if u.Id {
		fields = append(fields, "id")
	}
	if u.Time {
		fields = append(fields, "created_at", "updated_at", "deleted_at")
	}

	if u.Nickname {
		fields = append(fields, "nickname")
	}
	if u.Email {
		fields = append(fields, "email")
	}
	if u.QQNumber {
		fields = append(fields, "qq_number")
	}
	if u.PasswordHash {
		fields = append(fields, "password_hash")
	}

	if u.IsAdmin {
		fields = append(fields, "is_admin")
	}
	if u.Remark {
		fields = append(fields, "remark")
	}

	if u.MoetranId {
		fields = append(fields, "moetran_id")
	}
	if u.MoetranJwt {
		fields = append(fields, "moetran_jwt")
	}

	query = query.Select(fields)
}

// Insert 创建一个新的用户
func (*User) Insert(hdl *gorm.DB, u *User) error {
	if u == nil {
		return &InvalidParameterError{}
	}

	if u.Nickname == "" ||
		u.Email == "" ||
		u.PasswordHash == "" {
		return &LackOfRequiredFieldError{}
	}

	return hdl.
		Model(&User{}).
		Create(u).Error
}

// Update 更新一个指定用户
func (*User) Update(hdl *gorm.DB, u *User) error {
	if u == nil {
		return &InvalidParameterError{}
	}

	if u.BaseModel.Id == 0 {
		return &LackOfPrimaryKeyError{}
	}

	return hdl.Updates(u).Error
}

// SelectFirst 查询一个用户
func (*User) SelectFirst(
	hdl *gorm.DB, cnd *UserSpec, fields *UserFields,
) (*User, error) {
	var u User

	query := hdl.Model(&User{})

	if cnd != nil {
		cnd.Apply(query)
	}

	if fields != nil {
		fields.Apply(query)
	}

	if err := query.
		First(&u).Error; err != nil {
		return nil, err
	}

	return &u, nil
}

// SelectMany 查询多个用户
func (*User) SelectMany(
	hdl *gorm.DB, cnd *UserSpec, fields *UserFields,
	offset *int, limit *int,
) ([]*User, error) {
	var users []*User

	query := hdl.Model(&User{})

	if cnd != nil {
		cnd.Apply(query)
	}

	if fields != nil {
		fields.Apply(query)
	}

	if offset != nil {
		query = query.Offset(*offset)
	}
	if limit != nil {
		query = query.Limit(*limit)
	}

	if err := query.
		Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// Delete 删除一个用户
func (*User) Delete(hdl *gorm.DB, id PKey) error {
	if id == 0 {
		return &LackOfPrimaryKeyError{}
	}

	return hdl.
		Model(&User{}).
		Where("id = ?", id).
		Delete(&User{}).Error
}
