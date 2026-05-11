package external

import (
	"fmt"
	"os"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/wiliehidayat87/rmqp"
	"gorm.io/gorm"
)

type SigCH struct {
	Logs    *logrus.Logger
	Sigch   chan os.Signal
	Forever chan bool
	DB      *gorm.DB
	Rabbit  rmqp.AMQP
}

func SigHandler(obj SigCH) {

	signalType := <-obj.Sigch

	obj.Logs.Info("Gracefully shutdown the application ...")

	switch signalType {
	default:

		obj.Logs.Info(fmt.Sprintf("got with default signal = %v", signalType.String()))

	case syscall.SIGHUP:
		obj.Logs.Info("got Hangup/SIGHUP - portable number 1")

	case syscall.SIGINT:
		obj.Logs.Info("got Terminal interrupt signal/SIGINT - portable number 2")

	case syscall.SIGQUIT:
		obj.Logs.Info("got Terminal quit signal/SIGQUIT - portable number 3 - will core dump")

	case syscall.SIGABRT:
		obj.Logs.Info("got Process abort signal/SIGABRT - portable number 6 - will core dump")

	case syscall.SIGKILL:
		obj.Logs.Info("got Kill signal/SIGKILL - portable number 9")

	case syscall.SIGALRM:
		obj.Logs.Info("got Alarm clock signal/SIGALRM - portable number 14")

	case syscall.SIGTERM:
		obj.Logs.Info("got Termination signal/SIGTERM - portable number 15")

	case syscall.SIGUSR1:
		obj.Logs.Info("got User-defined signal 1/SIGUSR1")

	case syscall.SIGUSR2:
		obj.Logs.Info("got User-defined signal 2/SIGUSR2")
	}

	// Cleanup
	sqlDB, _ := obj.DB.DB()
	defer sqlDB.Close()

	obj.Logs.Info(fmt.Sprintf("Database Postgre (%#v connection) stopped successfully", sqlDB.Stats().InUse))

	defer obj.Rabbit.Channel.Close()

	obj.Logs.Info(fmt.Sprintf("AMQP Channel RabbitMQ stopped (%#v) successfully", obj.Rabbit.Channel))

	defer obj.Rabbit.Connection.Close()

	obj.Logs.Info(fmt.Sprintf("AMQP Connection RabbitMQ stopped (%#v) successfully", obj.Rabbit.Connection))

	// Stopping goroutine
	obj.Forever <- true
}
