package vo

type ReqCreateUser struct {
	UserName string `json:"user_name" binding:"required"`
	NickName string `json:"nick_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=128"`
	Mobile   string `json:"mobile" binding:"required,len=11"`
}

type ReqLogin struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required,min=6,max=128"`
}
