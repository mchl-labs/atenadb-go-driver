# Atena Go Driver

[![docs-source](https://img.shields.io/badge/source-atenadb-9cf?e&logo=go&labelColor=2a2a2a)](https://github.com/mchl-labs/atenadb-go-driver) [![Go Report Card](https://goreportcard.com/badge/github.com/mchl-labs/atenadb-go-driver)](https://goreportcard.com/report/github.com/mchl-labs/atenadb-go-driver)

# Getting started

Atena Go driver is a fully fledged tool for a fast and easy management of everything about Atena Db.
This library hosts all the necessary to handle requests to Atena Server. Thanks to this package you will be able to completely administer, manage and use Atena DB. We are making updates as a new core version comes out and tag the release according to the core release.

#### Just a few steps and you are ready to go.

So let's start developing!!!

## Step 1. Import the `atena-go-driver` module

  ```bash
   go get github.com/mchl-labs/atenadb-go-driver
  ```

## Step 2. Import the Required Namespace

  ```go
   import (
      "fmt"
      "context"
      "github.com/mchl-labs/atenadb-go-driver"
   )
  ```

### What provides Atena golang driver?

Atena golang driver provides 2 different classes to manage Atena db:

- **[AtenaAdmin]()** - `AtenaAdmin` provides everything you need for dbs and users management.
- **[AtenaDb]()** - In the other hand, `AtenaDb` features a complete set of features that enable all kind of queries to your db making easier your development experience.

## Step 3. Initialize connection

Connection is set up when an AtenaDb or AtenaAdmin object is instantiated.
The AtenaDb and the AtenaAdmin are the main arbiters of the connection to Atena, your application should maintain a single instance of these classes throughout its runtime. 

In order to initialize the connection you need the following info:

 `USERNAME, PASSWORD, HOST_NAME:PORT_NUMBER, DBNAME` 

where [USERNAME]() is your atena username (e.g. default user "[Admin]()"), 

[PASSWORD]() is your atena server's password (e.g. default psw user Admin is "[Admin]()", you can change it whenever you want), 

[HOST_NAME]() is the host name of your server (e.g. [localhost]()), 

[PORT_NUMBER]() is the port number Atena is listening on (usually e.g. [5001]()) 

and [DBNAME]() the DB you want to connect to (e.g. "[Atena]()" our default db which is always available, you don't need to create it. When you turn on Atena, you can just start using it).

`HOST_NAME:PORT_NUMBER` is optional in both the classes. 
If you don't need advanced configuration you can just omit it. Atena will connect to localhost:5001 the default AtenaDb endpoint.

```golang

// `HOST_NAME:PORT_NUMBER` is optional in both the classes. 
// If you don't need advanced configuration you can just omit it. Atena will connect to localhost:5001 the default AtenaDb endpoint.

atenaAdmin, err := atena.BuildAdmin(username, password)
db, err := atena.BuildDefault(username, password,"Atena")

// OR

atenaAdmin, err := atena.BuildAdmin(username, password, "HOST_NAME:PORT_NUMBER")
db, err := atena.Build(username, password,"HOST_NAME:PORT_NUMBER", "Atena")

```

## AtenaAdmin

### Users and dbs management 💼

Now that you've retreived the connection to the database, all that's left is to use it. Here are some simple operations:

  ### Create New User

  #### Naming rules:
  Usernames and passwords can contain letters, numbers and the underscore (_).

  ```golang
  result := atenaAdmin.CreateUser(username, password)
  ```

  Returns a `bool` representing the result of the operation.

  ### Change Password

  #### Naming rules:
  Usernames and passwords can contain letters, numbers and the underscore (_).
  
  ```golang
   result := atenaAdmin.ChangePassword(newpassword)
  ```

  Returns a `bool` representing the result of the operation.

  ### Create New DB ( using the default db engine ) ⚙️

  #### This method allow you to create a new db using the default Db engine of Atena Db.

  This db is saved [persistently on-disk by default]() and it automatically optimize the usage of your resources to provide the best performance using as few resources as possible.

  #### Stop care about save resources and limit waste. Atena takes care for you.

  Unlike all the other K/V stores on the market who work entirely in-memory or on-disk this db optimize the usage of your memory.

  It supports data [larger than your memory](), by leveraging fast external storage. So it is also consistent and you won't lose any data.
  It uses consistent recovery using a [fast non-blocking checkpointing technique](), that lets applications trade-off performance for commit latency.

  Atena is strongly [clound-oriented]() and [memory is very expensive in the cloud](), so run a K/V store completely in-memory can be very expensive. Atena doesn't need a lot of memory to works and at the same time doesn't have capacity limits due to the limited memory resources. It's designed for `heavy updates and read/write loads` as well as [top-class performance]().

  ```golang
  result := atenaAdmin.CreateDB("mynewdbname")
  ```

  Returns a `bool` representing the result of the operation.

  ### Create New DB ( using the Atena RBTree Engine ) ⚙️

> This db engine is still in [beta]() and we don't recommend use it. However it has already been thoroughly tested and is already stable 🟢

  #### This method allow you to create a new db using the Atena RBTree Engine.

  This db is a [completely custom db engine]() realized by Atena. CRUD operations and transactions on DB are made in log time (O(log n)). The engine is based on our RedBlack Tree implementation in C#.

  It is meant as an advanced and more complex db engine option. In fact it has got two option in addition to the usual parameters. 
  
  These option are:

  - **[LFU]()** - LFU is an `eviction policy` which autonomously removes the oldest and least used record from your store ( as cache ), optimizing your db size and resources usage
  - **[on-disk persistence]()** - Instead on-disk persistence enable the saving on disk of the db to avoid data loss.


  > 💡: We recommend enable this two feature together to get the best experience possible.

  > ⚠️ NOTE: Enabling LFU without on-disk persistence enabled can cause data loss.

  ```golang
  result := atenaAdmin.CreateDBRBT("mynewdbname", bool: lfu, bool: on-disk-persistence)
  ```

  Returns a `bool` representing the result of the operation. 

  ### Create New DB ( using the Atena HT Engine ) ⚙️

> This db engine is still in [beta]() and we don't recommend use it. However it has already been thoroughly tested and is already stable 🟢

  #### This method allow you to create a new db using the Atena RBTree Engine.

  This db is a [completely custom db engine]() realized by Atena. The engine is based on an hashtable.

  It is meant as an advanced and more complex db engine option. In fact it has got two option in addition to the usual parameters. 
  
  These option are:

  - **[LFU]()** - LFU is an `eviction policy` which autonomously removes the oldest and least used records from your store ( as cache ), optimizing your db size and resources usage
  - **[on-disk persistence]()** - Instead on-disk persistence enable the saving on disk of the db to avoid data loss.


  > 💡: We recommend enable this two feature together to get the best experience possible.

  > ⚠️ NOTE: Enabling LFU without on-disk persistence enabled can cause data loss.

  ```golang
  result := atenaAdmin.CreateDBHT("mynewdbname", bool: lfu, bool: on-disk-persistence);
  ```

  Returns a `bool` representing the result of the operation.

  ### Delete DB

  ```golang
  result := atenaAdmin.DeleteDB(dbname)
  ```

  Returns a `bool` representing the result of the operation.
  ### Logout

  > 💡: We recommend logging out if you finished using AtenaAdmin. This will boost AtenaDb performance.

  > ⚠️ NOTE: One time you Logout you need do reinitialize the connection and Build anothe Admin object ( Step 3.)

  ```golang
  result := atenaAdmin.Logout();
  ```

  Returns a `bool` representing the result of the operation.

## AtenaDb

### Db operations 👷

Let's make some operations on our db.
Now that you've retreived the connection to the database, all that's left is to use it. Here are some simple operations:

  ### Set K/V record

  This command allow you to set a key and a value of type string. 
  You can put practically anything inside this K/V record because everything, from complex object (JSON, Protobuf) to byte[], can be serialized into a string.

  ```golang
  result := db.Set("hello", "hi")
  ```

  Returns a `bool` representing the result of the operation.

  ### Get Value ( relative to the key we ask for)

  ```golang
  value, err := db.Get("hello")
  if err != nil {
			log.Fatal("Error")
    }
	log.Println("GET OP")
	log.Println("KEY: hello VALUE: " + value)
  ```

  Returns a `bool` with the status of the current operation and if the operation was successful the value related to the key we asked for.

  ### Modify the value

  To modify the value you can simply use the `Set` method.

  ### Increment

  Specifically designed for analytics and counters storage this method allows you to modify the value of your counter adding or substracting directly through one single call. 

  This feature is designed to optimize performance, reduce latency and minimize pointless DB calls.

  > To use this feature the value must be an integer

  ```golang

  // SET a counter.

  result := db.Set("clic", "1")

  if result
  {
      log.Println("SET clic : 1")
  }

  // Increment the counter by 11 units
  newvalue, err := db.Incr("clic",11)
  if err != nil {
			log.Fatal("Error")
    }
	log.Println("KEY: clic; NEWVALUE: " + newvalue)
  ```

  Returns a `bool` with the status of the current operation and, if the operation was successful, the new value related to the counter we wanted to modify.

  ### Delete K/V record

  ```golang
  result := atena.Del("hello")
  ```

  Returns a `bool` representing the result of the operation.

  ### RemoveAll records

  Thanks to this method you are able to completely remove every K/V record from the db.

  ```golang
  result := atena.RemoveAll()
  ```

  Returns a `bool` representing the result of the operation.
