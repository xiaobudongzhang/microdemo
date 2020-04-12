package access

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part3/basic/config"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/micro/v2/broker"
)

var (
	tokenExpiredDate  = 3600 * 24 * 30 * time.Second
	tokenIDKeyPrefix  = "token:auth:id:"
	tokenExpiredTopic = "mu.micro.book.topic.auth.tokenExpired"
)

type Subject struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

func (s *service) MakeAccessToken(subject *Subject) (ret string, err error) {
	m, err := s.createTokenClaims(subject)

	if err != nil {
		return "", fmt.Errorf("create token claim fail, %", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, m)
	ret, err = token.SignedString([]byte(config.GetJwtConfig().GetSecretKey()))
	if err != nil {
		return "", fmt.Errorf("create token fail %s", err)
	}

	err = s.saveTokenToCache(subject, ret)
	if err != nil {
		return "", fmt.Errorf("save cache fail %s", err)
	}
	return
}

func (s *service) GetCacheAccessToken(subject *Subject) (ret string, err error) {
	ret, error = s.GetTokenFromCache(subject)
	if err != nil {
		return "", fmt.Errorf("get token fail")
	}
	return
}

func (s *service) DelUserAccessToken(tk string) (err error) {
	claims, err := s.parseToken(tk)
	if err != nil {
		return fmt.Errorf("error token %s", err)
	}

	err = s.delTokenFromCache(&Subject{
		ID:claims.Subject,
	})

	if err != nil {
		return fmt.Errorf("del cache fail %s", err)
	}

	msg := &broker.Message{
		Body:[]byte(claims.Subject)
	}

	if err := broker.Publish(tokenExpiredTopic, msg) ; err != nil {
		log.Logf("发布token失败%v", err)
	} else {
		fmt.Println("发布删除消息", string(msg.Body))
	}
	return
}
