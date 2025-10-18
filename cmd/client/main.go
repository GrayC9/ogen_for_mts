//go:build ogen

package main

import (
	"context"
	"log"
	"os"
	"time"

	api "github.com/example/ogen_for_mts/internal/api_1"
)

func main() {
	addr := "http://localhost:8080"
	if v := os.Getenv("BASE_URL"); v != "" {
		addr = v
	}

	cl, err := api.NewClient(addr)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cres, err := cl.CreateUser(ctx, &api.UserCreate{Name: "Alice"})
	if err != nil {
		log.Fatal("create:", err)
	}
	var id int64
	switch v := cres.(type) {
	case *api.User:
		log.Printf("Created: %+v", v)
		id = v.ID
	default:
		log.Fatal("create: unexpected response type")
	}

	lst, err := cl.ListUsers(ctx)
	if err != nil {
		log.Fatal("list:", err)
	}
	log.Printf("List: %+v", lst)

	gres, err := cl.GetUser(ctx, api.GetUserParams{ID: id})
	if err != nil {
		log.Fatal("get:", err)
	}
	switch v := gres.(type) {
	case *api.User:
		log.Printf("Get: %+v", v)
	default:
		log.Printf("Get: not found")
	}

	newName := api.NewOptString("Alice Updated")
	newDesc := api.NewOptNilString("Updated user")
	ures, err := cl.UpdateUser(ctx, &api.UserUpdate{Name: newName, Description: newDesc}, api.UpdateUserParams{ID: id})
	if err != nil {
		log.Fatal("update:", err)
	}
	switch v := ures.(type) {
	case *api.User:
		log.Printf("Updated: %+v", v)
	default:
		log.Printf("Update: not found")
	}

	if _, err := cl.DeleteUser(ctx, api.DeleteUserParams{ID: id}); err != nil {
		log.Fatal("delete:", err)
	}
	log.Printf("Deleted: %d", id)
}
