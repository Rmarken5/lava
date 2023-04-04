package main

import (
	"database/sql"
	"fmt"
	"github.com/rmarken5/lava/inspect/data-access"
	"github.com/rmarken5/lava/inspect/logic"
	"io"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {

	// Define the connection parameters
	connStr := "user=user password=password dbname=mlb host=localhost port=5432 sslmode=disable"

	// Create the connection pool
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
	}
	defer db.Close()

	inspector := data_access.NewInspector(db, log.New(NewStdoutWriter(), "data-access: ", 0))

	logic := logic.New(log.New(NewStdoutWriter(), "logic: ", 0), inspector)

	str, err := logic.BuildStructsForQuery(`select * from team;`)
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
	}
	fmt.Printf("%s", str)

}

func NewStdoutWriter() io.Writer {
	return os.Stdout
}
