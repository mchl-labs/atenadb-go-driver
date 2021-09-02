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

// Manage implements atena Db operations
type Manage struct{

	// User info
	user string

	// atena connection info
	url string

	// Auth
	token string

	// gRPC Connection
	conn *grpc.ClientConn
	c pb.AtenaDBClient
}


// BuildAdmin creates the admin DB client
func BuildAdmin(user string, password string, url string) (*Manage, error) {
	client := new(Manage)
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
	
	token, err := client.auth(user, password, url)
	if err != nil || token == "" {
		return client, err
	}

	client.token = token
	client.user = user
	client.url = url

	return client,nil
}

// BuildAdminDefault creates the admin DB client with default url
func BuildAdminDefault(user string, password string) (*Manage, error) {
	client := new(Manage)
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
	
	token, err := client.auth(user, password, address)
	if err != nil || token == "" {
		return client, err
	}

	client.token = token
	client.user = user
	client.url = address

	return client, nil
}

// Auth
func (client *Manage) auth(user string, password string, url string) (string, error) {
	if user == "" || password == "" {
		return "", errors.New("Error")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.c.AuthUser(ctx, &pb.AuthUserLookupModel{User: user,Password: password})
	if err != nil || !r.Successful {
		return "", errors.New("Error")
	}
	return r.GetToken(), nil
}

// CreateUser creates a new user
func (client *Manage) CreateUser(user string, password string) bool {
	// if string is null or empty or whitespace => Exception
	if user == "" || password == "" {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.c.CreateUser(ctx, &pb.CreateUserModel{Token: client.token, Name: user, Password: password})
	if err != nil {
		return false
	}
	return r.GetSuccessful()
}

// ChangePassword changes user password
func (client *Manage) ChangePassword(password string) bool {
	// if string is null or empty or whitespace => Exception
	if password == "" {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.c.CreateUser(ctx, &pb.CreateUserModel{Token: client.token, Name: client.user, Password: password})
	if err != nil {
		return false
	}
	return r.GetSuccessful()
}

// CreateDB creates a new DB
func (client *Manage) CreateDB(name string) bool {
	// if string is null or empty or whitespace => Exception
	if name == "" {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.c.CreateDB(ctx, &pb.DBInfo{Token: client.token, Name: name})
	if err != nil {
		return false
	}
	return r.GetSuccessful()
}

// CreateDBRBT creates a new DB using RBT db engine
func (client *Manage) CreateDBRBT(name string, lfu bool, save bool) bool {
	// if string is null or empty or whitespace => Exception
	if name == "" {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.c.CreateDBRBT(ctx, &pb.DBModel{Token: client.token, Name: name, LFU: lfu, Save: save})
	if err != nil {
		return false
	}
	return r.GetSuccessful()
}

// CreateDBHT creates a new DB using HT db engine
func (client *Manage) CreateDBHT(name string, lfu bool, save bool) bool {
	// if string is null or empty or whitespace => Exception
	if name == "" {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.c.CreateDBHT(ctx, &pb.DBModel{Token: client.token, Name: name, LFU: lfu, Save: save})
	if err != nil {
		return false
	}
	return r.GetSuccessful()
}

// DeleteDB deletes the selected db
func (client *Manage) DeleteDB(name string) bool {
	// if string is null or empty or whitespace => Exception
	if name == "" {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.c.DeleteDB(ctx, &pb.DBInfo{Token: client.token, Name: name})
	if err != nil {
		return false
	}
	return r.GetSuccessful()
}

// Logout
func (client *Manage) Logout() bool {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := client.c.Logout(ctx, &pb.LogoutUser{Token: client.token})
	if err != nil {
		return false
	}
	if r.GetSuccessful() {
		client.token = ""
	}
	return r.GetSuccessful()
}

// Dispose connection
func (client *Manage) Dispose(){
	client.conn.Close()
}