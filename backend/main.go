package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	dbURL := os.Getenv("DATABASE_URL")

	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	pool, err := pgxpool.New(
		context.Background(),
		os.Getenv("DATABASE_URL"),
	)

	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to PostgreSQL")

	_, err = pool.Exec(
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

	_, err = pool.Exec(
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

			_, err := pool.Exec(
				context.Background(),
				"UPDATE counter SET value = value + 1 WHERE id = 1",
			)

			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
		}

		var count int

		err := pool.QueryRow(
			context.Background(),
			"SELECT value FROM counter WHERE id = 1",
		).Scan(&count)

		if err != nil {
			log.Printf("GET database error: %v", err)
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
