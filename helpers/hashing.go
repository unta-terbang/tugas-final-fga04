package helpers

import "golang.org/x/crypto/bcrypt"

func HashPwd(password string) (res string, err error){

	salt := 10
	arrByte := []byte(password)
	hash,err := bcrypt.GenerateFromPassword(arrByte,salt)

	return string(hash), err

}

func PasswordValid(h, p string) bool {

	hass, pass := []byte(h), []byte(p)

	err := bcrypt.CompareHashAndPassword(hass, pass)


	return err == nil

}