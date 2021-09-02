package atena

import (
	"context"
	"crypto/tls"
	"errors"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/mchl-coder/atenadb-go-driver/atena"
)

const (
	address = "localhost:5001"
)

// Ops implements atena Db operations
type Ops struct {

	// User info
	user string

	// DB connection info
	DB  string
	url string

	// Auth
	token string

	// gRPC Connection
	conn *grpc.ClientConn
	c    pb.AtenaDBClient
}

// Build creates the DB client
func Build(user string, password string, url string, db string) (*Ops, error) {
	client := new(Ops)
	log.Println("done")
	// Set up a connection to the server.
	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(credentials.NewTLS(config)), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
		return client, err
	}

	c := pb.NewAtenaDBClient(conn)

	client.conn = conn
	client.c = c

	token, err := client.auth(user, password, url, db)
	if err != nil || token == "" {
		return client, err
	}

	client.token = token
	client.user = user
	client.url = url
	client.DB = db

	return client, nil
}

// BuildDefault creates the DB client with default url
func BuildDefault(user string, password string, db string) (*Ops, error) {
	client := new(Ops)
	// Set up a connection to the server.
	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(credentials.NewTLS(config)), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
		return client, err
	}

	c := pb.NewAtenaDBClient(conn)

	client.conn = conn
	client.c = c

	token, err := client.auth(user, password, address, db)
	if err != nil || token == "" {
		return client, err
	}

	client.token = token
	client.user = user
	client.url = address
	client.DB = db

	return client, nil
}

// Auth
func (client *Ops) auth(user string, password string, url string, db string) (string, error) {
	if user == "" || password == "" || db == "" {
		return "", errors.New("Error")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.c.Auth(ctx, &pb.AuthLookupModel{User: user, Password: password, Db: db})
	if err != nil || !r.Successful {
		return "", errors.New("Error")
	}
	return r.GetToken(), nil
}

// Set operation
func (client *Ops) Set(key string, value string) bool {
	// if string is null or empty or whitespace => Exception
	if key == "" || value == "" {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.c.SetRecord(ctx, &pb.AtenaSet{Token: client.token, Key: key, Value: value})
	if err != nil {
		return false
	}
	return r.GetSuccessful()
}

// Get operation
func (client *Ops) Get(key string) (string, error) {
	if key == "" {
		return "", errors.New("Error: Unvalid key")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.c.GetRecord(ctx, &pb.AtenaGet{Token: client.token, Query: key})
	if err != nil {
		return "", err
	}
	return r.GetValue(), nil
}

// Del operation
func (client *Ops) Del(key string) bool {
	if key == "" {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.c.DeleteRecord(ctx, &pb.AtenaDel{Token: client.token, Key: key})
	if err != nil {
		return false
	}
	return r.GetSuccessful()
}

// Incr operation
func (client *Ops) Incr(key string, inc int32) (string, error) {
	if inc == 0 {
		return "", errors.New("Error: Increment must be different than 0")
	}
	if key == "" {
		return "", errors.New("Error: Unvalid Key")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.c.IncrRecord(ctx, &pb.AtenaIncr{Token: client.token, Key: key, Inc: inc})
	if err != nil {
		return "", err
	}
	return r.GetValue(), nil
}

// RemoveAll removes all the K/V records stored
func (client *Ops) RemoveAll() bool {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.c.RemoveAll(ctx, &pb.RemoveAllRecords{Token: client.token})
	if err != nil {
		return false
	}
	return r.GetSuccessful()
}

// Dispose connection
func (client *Ops) Dispose() {
	client.conn.Close()
}
