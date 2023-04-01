# Distributed-Data-Store-with-Queue-Support
Implementation of a simple data store API server that accepts HTTP POST requests with JSON payloads to handle operations like set, get, delete, exists, queue push (QPush), queue pop (QPop), and blocking queue pop (BQPop.

The  Distributed-Data-Store-with-Queue-Support implementation uses a mutex lock to synchronize access to the data store across concurrent goroutines. 
