package main

import (
	// "context"
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	// "github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	err := godotenv.Load(".env")
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
			log.Println("new migration:", err)
		}
		filename := time.Now().Format("03040502012006")
		filename = strings.Join([]string{filename, *newName}, "_")
		filenameUp := strings.Join([]string{"db/migrations/", filename, ".up.sql"}, "")
		filenameDown := strings.Join([]string{"db/migrations/", filename, ".down.sql"}, "")
		os.WriteFile(filenameUp, []byte{}, 0666)
		os.WriteFile(filenameDown, []byte{}, 0666)

	case "run":
		migrationsToRun := []string{}
		direction := os.Args[2]
		if direction == "up" {
			migrationsToRun = findMigrationsToRun(direction)
		} else if direction == "down" {
			migrationsToRun = findMigrationsToRun(direction)
		} else {
			log.Fatal("Must specify migration direction (up / down)")
		}

		log.Println("MIGRATIONS TO RUN:")
		for _, m := range migrationsToRun {
			log.Println(m)
		}

		db := connectDB()
		for _, migration := range migrationsToRun {
			log.Printf("RUNNING MIGRATION: %s", migration)
			dat, err := os.ReadFile(fmt.Sprintf("db/migrations/%s", migration))
			check(err)
			_, errdb := db.Exec(context.Background(), string(dat))
			check(errdb)
			migrationName := strings.TrimSuffix(migration, fmt.Sprintf(".%s.sql", direction))
			os.WriteFile("db/migrations/tally", []byte(migrationName), 0666)
			log.Println("DONE ✅️")
		}
	}
}

func findMigrationsToRun(direction string) []string {
	ls, err := os.ReadDir("db/migrations")
	check(err)

	// holds all the migration filenames (up and down)
	// look like this: 07575503112025_test.up.sql
	migrations := []string{}
	for _, line := range ls {
		if strings.HasSuffix(line.Name(), fmt.Sprintf("%s.sql", direction)) {
			migrations = append(migrations, line.Name())
		}
	}

	// tally file keeps track of which migrations have been run
	// when migrate run is run, any from above that aren't here
	// are run. each line looks like: 07575503112025_test
	f, err := os.OpenFile("db/migrations/tally", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	migrationsAlreadyRun := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		migrationsAlreadyRun = append(migrationsAlreadyRun, line)
	}

	migrationsToRun := []string{}
	re := regexp.MustCompile(`\..*\.sql`)
	for _, mig := range migrations {
		if len(migrationsAlreadyRun) == 0 {
			if direction == "up" {
				migrationsToRun = migrations
			}
			break
		}
		for _, doneMig := range migrationsAlreadyRun {
			migName := re.ReplaceAllLiteralString(mig, "")
			log.Println("comparing", mig, migName, doneMig)
			if (migName != doneMig && direction == "up") ||
				(migName == doneMig && direction == "down") {
				migrationsToRun = append(migrationsToRun, mig)
			}
		}
	}

	return migrationsToRun
}

func connectDB() *pgxpool.Pool {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	return dbpool
}
