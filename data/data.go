package data

import (
	"crypto/rand"
	"crypto/sha1"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("postgres", "user=tomokokawase password=zhy677097 dbname=chitchat sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return
}

func createUUID() (uuid string) {
	u := new([16]byte)
	// rand实现了read方法，所以它是个io.reader
	// Read 应该会返回读取的字节数和一个错误(如果有的话)
	_, err := rand.Read(u[:])
	if err != nil {
		log.Fatal("Cannot generate UUID", err)
	}
	u[8] = (u[8] | 0x40) & 0x7F
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	// 由于函数签名制定了返回值的变量名，这里可以省略，编译器会自己去找
	return
}

func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return
}
