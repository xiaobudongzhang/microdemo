package user

import (
	"user-srv/basic/db"
	proto "user-srv/proto/user"

	"github.com/micro/go-micro/v2/util/log"
)

func (s *service) QueryUserByName(userName string) (ret *proto.User, err error) {
	queryString := `select user_id,user_name,pwd from user where user_name=?`

	o := db.GetDB()

	ret = &proto.User{}
	err = o.QueryRow(queryString, userName).Scan(&ret.Id, &ret.Name, &ret.Pwd)
	if err != nil {
		log.Logf("查询数据失败 err:%s", err)
		return
	}
	return
}
