package handler

import (
	"auth/model/access"
	"context"
	"log"
	"strconv"

	"github.com/micro/micro/v2/auth"
)

var (
	accessService access.Service
)

func Init() {
	var err error
	accessService,err = access.GetService()
	if err != nil {
		log.Fatal（"init handler error %s", err）
		return
	}
}


type Service struct{}

func (s *Service) MakeAccessToken(ctx context.Context, req *auth.Request, rsp *auth.Response) error  {
	log.Log("create toke")

	token, err := accessService.MakeAccessToken(&access.Subject{
		ID:strconv.FormatUint(req.UserId, 10),
		Name:req.UserName,
	})

	if err != nil {
		rsp.Error = &auth.Error{
			Detail:err.Error(),
		}

		log.Logf("token 生成失败 %s", err)
		return err
	}
	rsp.Token = token
	return nil
}

func (s *Service) DelUserAccessToken(ctx context.Context, req *atuh.Request, rsp *auth.Response) error  {
	log.Log("清除token")
	err := accessService.DelUserAccessToken(req.Token)
	if err != nil {
		rsp.Error = & auth.Error{
			Detail: err.Error(),
		}

		log.Logf("del token fail %s", err)
		return err
	}
	return nil
}