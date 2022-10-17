package main

import (
	"context"
	"flag"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	mysqlNS  string
	mysqlSvc string
	rootPW   string
)

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	flag.StringVar(&mysqlSvc, "mysqlSvc", "mysql", "mysql svc name.")
	flag.StringVar(&mysqlNS, "mysqlNS", "default", "mysql namespace name.")
	flag.StringVar(&rootPW, "rootPW", "root", "root password.")

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	svc, err := clientset.CoreV1().Services(mysqlNS).Get(context.Background(), mysqlSvc, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}

	conn, err := getConn(svc.Spec.ClusterIP, svc.Spec.Ports[0].Port)
	if err != nil {
		panic(err.Error())
	}
	// create database
	err = conn.Exec("CREATE DATABASE IF NOT EXISTS test").Error
	if err != nil {
		panic(err.Error())
	}
}

func getConn(svc string, port int32) (*gorm.DB, error) {
	conn, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/", "root", rootPW, svc, port)))
	return conn, err
}
