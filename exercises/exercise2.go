package exercises

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

type shardConfig struct {
	id         int
	connString string
}
type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Shard string `json:"shard"`
}

func connectToShards() ([]*sql.DB, error) {
	shardConfigs := []shardConfig{
		{id: 0, connString: "host=localhost port=5432 user=admin password=admin dbname=shard0 sslmode=disable"},
		{id: 1, connString: "host=localhost port=5433 user=admin password=admin dbname=shard1 sslmode=disable"},
	}
	dbs := make([]*sql.DB, len(shardConfigs))
	for i, config := range shardConfigs {
		db, err := sql.Open("postgres", config.connString)
		if err != nil {
			return nil, fmt.Errorf("error connecting to the DB %d:%s", config.id, err)
		}
		if err = db.Ping(); err != nil {
			return nil, fmt.Errorf("failed to ping the shard %d:%s", config.id, err)
		}
		dbs[i] = db
		fmt.Printf("Successfully connected to shard %d\n", config.id)
	}
	return dbs, nil
}

func handleQuery(w http.ResponseWriter, r *http.Request, dbConns []*sql.DB) {
	params := r.URL.Query().Get("user_id")
	userId, err := strconv.Atoi(params)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
	}
	shardIndex := userId % len(dbConns)
	db := dbConns[shardIndex]

	row := db.QueryRow("SELECT user_id, name, email FROM users WHERE user_id=$1", userId)

	var user User

	err = row.Scan(&user.Id, &user.Name, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		return
	}
	user.Shard = fmt.Sprintf("shard-%d", shardIndex)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
	}

}

func ShardExercise() {
	dbConns, err := connectToShards()
	if err != nil {
		fmt.Printf("Error connecting to the shards: %s", err)
		return
	}
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		handleQuery(w, r, dbConns)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
