# Distributed-Data-Store-with-Queue-Support
Implementation of a simple data store API server that accepts HTTP POST requests with JSON payloads to handle operations like set, get, delete, exists, queue push (QPush), queue pop (QPop), and blocking queue pop (BQPop.

This is a Go program that implements a simple key-value store with some additional functionality like a queue. The program listens for incoming HTTP requests and handles them by performing the specified operation on the data store. 

The supported operations are:
      GET: Get the value of a key

      SET: Set the value of a key

      DELETE: Delete a key and its value

      EXISTS: Check if a key exists

      QPUSH: Push one or more values onto the end of a queue

      QPOP: Pop a value from the end of a queue

      BQPOP: Pop a value from the beginning of a queue (blocking)



The program is structured as follows:
      The main function creates a new instance of DataStore, which is a struct that contains a map of strings to strings for storing key-value pairs.

      The handleAPI function creates an HTTP handler function that parses incoming requests, performs the requested operation on the data store, and returns a JSON-  encoded response.

      The handleSet, handleGet, handleDelete, handleExists, handleQPush, handleQPop, and handleBQPop functions are called by handleAPI to perform the requested operation on the data store.


The program also defines the following types:
      Request: A struct that represents an incoming API request. It has fields for the command, key, value, and options.

      Response: A struct that represents an outgoing API response. It has fields for whether the operation was successful, a message (if applicable), the key, and the value.

      DataStore: A struct that represents the data store. It has a map of strings to strings for storing key-value pairs and a RWMutex for synchronizing access to the map. It also has methods for getting, setting, deleting, and locking queues.

The program uses the following standard Go packages:

      encoding/json: for encoding and decoding JSON data
      log: for logging errors
      net/http: for handling HTTP requests and responses
      sync: for synchronization primitives (RWMutex)
