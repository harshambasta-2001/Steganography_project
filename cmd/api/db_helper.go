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