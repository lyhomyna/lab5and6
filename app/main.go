package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatal("DATABASE_URL variable is not set")
	}

	var err error
	db, err = sql.Open("postgres", connString)
	if err != nil {
		log.Fatal("Engine failure (DB connect):", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Database is unreachable:", err)
	}

	query := `CREATE TABLE IF NOT EXISTS jdm_parts (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		car_model TEXT NOT NULL
	);`
	if _, err := db.Exec(query); err != nil {
		log.Fatal("Failed to create table:", err)
	}

	// Prefill database
	prefilDatabase()

	http.HandleFunc("/", loggingMiddleware(homeHandler))
	http.HandleFunc("/parts", loggingMiddleware(partsHandler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	log.Printf("JDM Registry starting on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func prefilDatabase() {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM jdm_parts").Scan(&count)

	if count == 0 {
		log.Println("Database is empty. Tuning the engine with 11 items...")
		items := []struct {
			name  string
			model string
		}{
			{"RB26DETT Engine", "Nissan Skyline"},
			{"2JZ-GTE Engine", "Toyota Supra"},
			{"TE37 Wheels", "Nissan Silvia"},
			{"Brembo Brakes", "Mitsubishi Evo"},
			{"Momo Steering Wheel", "Honda NSX"},
			{"Recaro Seats", "Mazda RX-7"},
			{"Tomei Expreme Exhaust", "Subaru Impreza"},
			{"HKS Turbo Kit", "Toyota AE86"},
			{"Ohlins Suspension", "Nissan 350Z"},
			{"Nismo Body Kit", "Honda Civic Type R"},
			{"Greddy Intercooler", "Nissan GT-R"},
		}
		for _, item := range items {
			_, err := db.Exec("INSERT INTO jdm_parts (name, car_model) VALUES ($1, $2)", item.name, item.model)
			if err != nil {
				log.Printf("Failed to create item %s: %v", item.name, err)
			}
		}
		log.Println("Complete. 11 items added")
	}
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", req.Method, req.URL.Path)

		next.ServeHTTP(w, req)

		log.Printf("Completed %s in %v", req.URL.Path, time.Since(start))
	}
}

func homeHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Welcome to the JDM Parts Registry API!")
	fmt.Fprintln(w, "Use /parts to see inventory.")
}

func partsHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		getParts(w, req)
	case "POST":
		addPart(w, req)
	case "PUT":
		updatePart(w, req)
	case "DELETE":
		deletePart(w, req)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "Unknown method")
	}
}

func getParts(w http.ResponseWriter, _ *http.Request) {
	rows, err := db.Query("SELECT id, name, car_model FROM jdm_parts")
	if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	fmt.Fprintln(w, "Current Inventory:")
	for rows.Next() {
		var id int
		var name, model string
		if err := rows.Scan(&id, &name, &model); err != nil {
			continue
		}
		fmt.Fprintf(w, "[ID: %d] %s for %s\n", id, name, model)
	}
}

func addPart(w http.ResponseWriter, req *http.Request) {
	name := req.URL.Query().Get("name")
	model := req.URL.Query().Get("model")
	if name == "" || model == "" {
		http.Error(w, "Missing params", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO jdm_parts (name, car_model) VALUES ($1, $2)", name, model)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Item added successfully")
}

func updatePart(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	newName := req.URL.Query().Get("name")

	result, err := db.Exec("UPDATE jdm_parts SET name = $1 WHERE id = $2", newName, id)
	if err != nil {
		http.Error(w, "Update failed", http.StatusInternalServerError)
		return
	}

	count, _ := result.RowsAffected()
	if count == 0 {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}
	fmt.Fprintln(w, "Item updated")
}

func deletePart(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")

	result, err := db.Exec("DELETE FROM jdm_parts WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Delete failed", http.StatusInternalServerError)
		return
	}

	count, _ := result.RowsAffected()
	if count == 0 {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}
	fmt.Fprintln(w, "Item deleted")
}
