package users

type User struct {
	Id       string `json:"id" sql:"id"`
	Login    string `json:"login" sql:"login"`
	Password string `json:"password" sql:"password"`
}
