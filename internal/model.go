package internal

type User struct {
	ID       uint   `gorm:"primaryKey;column:ID" json:"ID" validate:"required"`
	Name     string `gorm:"column:Name" json:"Name" validate:"required"`
	Email    string `gorm:"unique;column:Email" json:"Email" validate:"required"`
	Password string `gorm:"column:Password" json:"Password"`
}

func (User) TableName() string {
	return "users"
}

type Product struct {
	ID     uint   `gorm:"primaryKey;column:ID" json:"id"`
	Text   string `gorm:"type:text;column:Text" json:"text"`
	UserID uint   `gorm:"column:UserId" json:"userid"`
	Code   string `gorm:"column:code" json:"code"`
}

func (Product) TableName() string {
	return "products"
}

type RegisterUser struct{
	Name string `json:"Name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8,max=130"`
}

type PayloadProduct struct{
	Text string `json:"Text" binding:"required"`
}

type LoginUser struct{
	Email string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8,max=130"`

}

// type CreateProduct struct {
// 	Text string `json:"Text" binding: "required"`
// 	UserID uint `json: "UserID" binding: "required"`
// }