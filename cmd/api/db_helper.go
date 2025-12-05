package main

import "github.com/harshambasta-2001/Steganography_project/internal"


func (s *APIServer) getUserbyEmail(email string) (*internal.User ,error){
	var user internal.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil

}

func (s *APIServer) getAllUsers() ([]internal.User ,error){
	var user []internal.User
	if err :=  s.db.Find(&user).Error; err!= nil{
		return nil,err
	}
	return user,nil
}

func (s *APIServer) delete_user(id int) (*internal.User ,error){
	var user internal.User

	if err:= s.db.Where("ID = ?",id).Delete(&user).Error; err !=nil{
		return nil,err
	}

	return &user,nil
}


func (s *APIServer) get_text_from_code(userId uint,code string) (*internal.Product,error){
	var product internal.Product

	if err:= s.db.Where("code = ? AND UserId = ?",code,userId).First(&product).Error; err !=nil{
		return nil,err
	}
	return &product ,nil
}

func (s *APIServer) delete_product(userId uint,code string) (*internal.Product,error){
	var product internal.Product

	if err:= s.db.Where("code = ? AND UserId= ?",code,userId).Delete(&product).Error; err !=nil{
		return nil,err
	}

	return &product,nil
}

func (s *APIServer) get_product_codes() ([]string,error){
	var codes []string
	if err := s.db.Model(&internal.Product{}).Pluck("code", &codes).Error; err != nil {
		return nil, err
	}
	return codes, nil

}