package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
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
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("sql.Open error %v", err)
	}
	defer db.Close()
	// ↓ DB Option ↓
	// DBへの最大接続数（デフォルトは無制限）
	db.SetMaxOpenConns(100)
	// DBへの接続の持続時間（time.Duration型）
	db.SetConnMaxLifetime(time.Duration(360))
	// コネクションプールの待機時間
	db.SetConnMaxIdleTime(time.Duration(360))
	// コネクションプール数（デフォルトは2）
	db.SetMaxIdleConns(2)
	for i := 0; i < 120; i++ {
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
			if i % 2 == 0 {
				if err := tx.Commit(); err != nil {
					log.Fatal(err)
				}
			}
			log.Printf("%+v", db.Stats())
		}(i)
	}
	wg.Wait()
	for {
		log.Printf("%+v", db.Stats())
		time.Sleep(time.Second * 1)
	}
	os.Exit(0)
}
