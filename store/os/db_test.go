package os

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func testOpenDB() *BoltDB {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	return OpenBoltDB("test.db", logger)
}

func TestBoltDB_AddFile(t *testing.T) {
	file := &DBFile{
		Name:    "nhao.jpg",
		Size:    1022,
		Path:    "./nihao.jpg",
		UUID:    "9985",
		Deleted: false,
	}
	err := testOpenDB().Add(file)
	if err != nil {
		panic(err)
	}
}

func TestBoltDB_GetFile(t *testing.T) {
	file := testOpenDB().Get("9985")
	fmt.Println(fmt.Sprintf("%#v", file))
}

func TestBoltDB_UpdateFile(t *testing.T) {
	err := testOpenDB().Update(&DBFile{UUID: "9985", Name: "shlande.py"})
	if err != nil {
		panic(err)
	}
}

func TestBoltDB_DeleteFile(t *testing.T) {
	db := testOpenDB()
	defer db.Close()
	err := db.Delete("9985")
	if err != nil {
		panic(err)
	}
	file := db.Get("9985")
	fmt.Println(file)
}
