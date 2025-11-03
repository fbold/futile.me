package main

import (
	// "context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	// "github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	newCmd := flag.NewFlagSet("new", flag.ExitOnError)
	newName := newCmd.String("name", "untitled", "Name of migration")

	switch os.Args[1] {
	case "new":
		newCmd.Parse(os.Args[2:])
		fmt.Println("new migration:", *newName)
		os.MkdirAll("db/migrations", 0777)
		if err != nil {
			fmt.Println("new migration:", err)
		}
		filename := time.Now().Format("03040502012006")
		filename = strings.Join([]string{filename, *newName}, "_")
		filenameUp := strings.Join([]string{"db/migrations/", filename, ".up.sql"}, "")
		filenameDown := strings.Join([]string{"db/migrations/", filename, ".down.sql"}, "")
		os.WriteFile(filenameUp, []byte{}, 0666)
		os.WriteFile(filenameDown, []byte{}, 0666)
	}

	// dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
	// 	os.Exit(1)
	// }
	//
	// defer dbpool.Close()

}
