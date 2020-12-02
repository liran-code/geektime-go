//借鉴https://github.com/cy22117888/Go-000/blob/main/Week02/main.go
package main

import(
	"database/sql"
	"errors"
	"fmt"
)

func main(){
	id := uint(1)
	user := UserServer{}
	user.findUserById(id)
	fmt.Println(user)
}


// Service层
type UserServer struct{
	UserDao *UserDao
}

func (u *UserServer) findUserById(id uint) *sql.Result{
	user, err := u.UserDao.findUserById(id)
	if err != nil{
		fmt.Println("error info %v", err)
		return nil
	}
	return user
}

//Dao层
type UserDao struct{
	Db *sql.DB
}

func (u *UserDao) findUserById(id uint) (*sql.Result, error){
	user, err := u.Db.Exec("SELECT * FROM user WHERE id = ?", id)
	if errors.Is(err, sql.ErrNoRows){
		return nil, errors.New("findUserById error:" + sql.ErrNoRows.Error())
	}
	if err != nil {
		return nil, errors.Wrap(err, "findUserById error...") // undefined:errors.Wrap
	}

	return &user, nil
}