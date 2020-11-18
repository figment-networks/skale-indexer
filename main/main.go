package main

import (
	"../client"
	"../handler"
	"../store"
	"../store/postgres"
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx := context.Background()
	db, err := sql.Open("postgres", "postgresql://localhost:5432/skale?user=postgres&password=admin&sslmode=disable")
	if err != nil {
		return
	}

	if err := db.PingContext(ctx); err != nil {
		return
	}

	defer db.Close()

	pgsqlDriver := postgres.NewDriver(ctx, db)
	storee := store.New(pgsqlDriver)

	hClient := client.NewClientContractor(storee)
	HTTPTransport := handler.NewClientConnector(*hClient)

	mux := http.NewServeMux()
	HTTPTransport.AttachToHandler(mux)

	srv := &http.Server{
		Handler: mux,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
	/*delegation, err := delegationStore.GetDelegationById(ctx, 3434)
	fmt.Println(delegation.ID)

	delegations, err := delegationStore.GetDelegationsByValidatorId(ctx, 0)
	if err == nil {
		fmt.Println(len(delegations))
	}	*/
}
