package service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ostafen/kronos/internal/config"
	"github.com/ostafen/kronos/internal/model"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type AlertService interface {
	Send(sched *model.Schedule) error
}

type alertService struct {
	emailConf config.Email
}

func NewAlertService(conf config.Email) *alertService {
	return &alertService{
		emailConf: conf,
	}
}

func extractHostPort(addr string) (string, int) {
	i := strings.Index(addr, ":")
	log.Println(i)
	if i < 0 {
		return "", -1
	}

	port, err := strconv.Atoi(addr[i+1:])
	if err != nil {
		return "", -1
	}
	return addr[:i], port
}

func (s *alertService) Send(sched *model.Schedule) error {
	host, port := extractHostPort(s.emailConf.Server)

	log.WithFields(log.Fields{
		"from": s.emailConf.Address,
		"to":   sched.Email,
	}).Info("sending email alert")

	m := gomail.NewMessage()
	m.SetHeader("From", s.emailConf.Address)
	m.SetHeader("To", sched.Email)
	m.SetBody("text/html", fmt.Sprintf("Kronos notification for schedule <b>%s</b>", sched.Title))

	dialer := gomail.NewDialer(host, port, s.emailConf.Address, s.emailConf.Password)
	err := dialer.DialAndSend(m)
	log.Println(err)
	return err
}
