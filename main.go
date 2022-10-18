package main

import (
	"context"
	"flag"
	"fmt"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
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

	db, err := sql.Open("mysql", fmt.Sprintf("root:%s@tcp(%s:%d)/", rootPW, svc.Spec.ClusterIP, svc.Spec.Ports[0].Port))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// create database
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS test")
	if err != nil {
		panic(err.Error())
	}
}
