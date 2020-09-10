package test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func Test_All_Default(t *testing.T) {

	var (
		user_name = "root"
		user_pass = "root"
		db_host   = "127.0.0.1:3306"
		db_name   = "database_test"
		dsn       = fmt.Sprintf(
			"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true",
			user_name,
			user_pass,
			db_host,
			db_name,
		)
		wg = new(sync.WaitGroup)
	)
	t1 := time.Now()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("sql.Open error %v", err)
	}
	defer db.Close()
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if err := db.Ping(); err != nil {
				log.Fatal(err)
			}
			tx, err := db.BeginTx(context.Background(), &sql.TxOptions{})
			if err != nil {
				log.Fatal(err)
			}
			if err := tx.Commit(); err != nil {
				log.Fatal(err)
			}
			log.Printf("%+v", db.Stats())
		}(i)
	}
	wg.Wait()
	t2 := time.Now()
	diff := t2.Sub(t1)
	log.Println(diff)
	for {
		log.Printf("%+v", db.Stats())
		time.Sleep(time.Second * 1)
	}
	os.Exit(0)
}
