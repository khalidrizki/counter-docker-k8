package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	conn, err := pgx.Connect(
		context.Background(),
		dbURL,
	)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close(context.Background())

	log.Println("Connected to PostgreSQL")

	_, err = conn.Exec(
		context.Background(),
		`
        CREATE TABLE IF NOT EXISTS counter (
            id SERIAL PRIMARY KEY,
            value INTEGER NOT NULL
        )
        `,
	)

	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Exec(
		context.Background(),
		`
        INSERT INTO counter (id, value)
        VALUES (1, 0)
        ON CONFLICT (id) DO NOTHING
        `,
	)

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/api/count", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "POST" {

			//   jika backend belum terhubung ke database
			// 		count = count + 1

			_, err := conn.Exec(
				context.Background(),
				"UPDATE counter SET value = value + 1 WHERE id = 1",
			)

			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
		}

		var count int

		err := conn.QueryRow(
			context.Background(),
			"SELECT value FROM counter WHERE id = 1",
		).Scan(&count)

		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		w.Header().Set("Content-Type", "application/json")
		response := map[string]int{
			"count": count,
		}
		json.NewEncoder(w).Encode(response)
	})

	http.ListenAndServe(":8080", nil)
}
