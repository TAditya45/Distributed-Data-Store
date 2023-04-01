package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type Request struct {
	Command string            `json:"command"`
	Key     string            `json:"key,omitempty"`
	Value   string            `json:"value,omitempty"`
	Options map[string]string `json:"options,omitempty"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Key     string `json:"key,omitempty"`
	Value   string `json:"value,omitempty"`
}

type DataStore struct {
	sync.RWMutex
	data map[string]string
}

func NewDataStore() *DataStore {
    return &DataStore{
        data: make(map[string]string),
        lock: sync.RWMutex{},
    }
}

func (ds *DataStore) Get(key string) (string, bool) {
	ds.RLock()
	defer ds.RUnlock()
	value, ok := ds.data[key]
	return value, ok
}

func (ds *DataStore) Set(key, value string) {
	ds.Lock()
	defer ds.Unlock()
	ds.data[key] = value
}

func (ds *DataStore) QPush(key string, values []string) {
	ds.Lock()
	defer ds.Unlock()
	queue, ok := ds.data[key]
	if !ok {
		queue = ""
	}
	for _, value := range values {
		queue += value + " "
	}
	ds.data[key] = queue
}

func (ds *DataStore) QPop(key string) (string, bool) {
	ds.Lock()
	defer ds.Unlock()
	queue, ok := ds.data[key]
	if !ok {
		return "", false
	}
	values := strings.Fields(queue)
	if len(values) == 0 {
		return "", false
	}
	value := values[len(values)-1]
	ds.data[key] = strings.Join(values[:len(values)-1], " ")
	return value, true
}

func (ds *DataStore) Delete(key string) bool {
	ds.Lock()
	defer ds.Unlock()
	_, ok := ds.data[key]
	if ok {
		delete(ds.data, key)
	}
	return ok
}

func (ds *DataStore) isQueueLocked(key string) bool {
	ds.lockMutex.RLock()
	defer ds.lockMutex.RUnlock()

	_, ok := ds.locks[key]
	return ok
}

func (ds *DataStore) lockQueue(key string) {
	ds.lockMutex.Lock()
	defer ds.lockMutex.Unlock()

	ds.locks[key] = struct{}{}
}

func (ds *DataStore) unlockQueue(key string) {
	ds.lockMutex.Lock()
	defer ds.lockMutex.Unlock()

	delete(ds.locks, key)
}


func handleAPI(ds *DataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			response := Response{
				Success: false,
				Message: "Invalid request method",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			response := Response{
				Success: false,
				Message: "Invalid request format",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		switch req.Command {
		case "GET":
			handleGet(w, ds, req)
		case "SET":
			handleSet(w, ds, req)
		case "DELETE":
			handleDelete(w, ds, req)
		case "EXISTS":
			handleExists(w, ds, req)
		case "QPUSH":
			handleQPush(w, ds, req)
		case "QPOP":
			handleQPop(w, ds, req)
		case "BQPOP":
			handleBQPop(w, ds, req)
		default:
			response := Response{
				Success: false,
				Message: "Invalid command",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
	}
}

func main() {
	ds := NewDataStore()

	http.HandleFunc("/", handleAPI(ds))

	log.Fatal(http.ListenAndServe(":8080", nil))
}


func handleSet(w http.ResponseWriter, ds *DataStore, req Request) {
	key := req.Key
	value := req.Value
	if key == "" {
		response := Response{
			Success: false,
			Message: "Key not provided",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	if value == "" {
		response := Response{
			Success: false,
			Message: "Value not provided",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	ds.Set(key, value)
	response := Response{
		Success: true,
		Key:     key,
		Value:   value,
	}
	json.NewEncoder(w).Encode(response)
}

func handleGet(w http.ResponseWriter, ds *DataStore, req Request) {
	key := req.Key
	value, ok := ds.Get(key)
	if !ok {
		response := Response{
			Success: false,
			Message: "Key not found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}
	response := Response{
		Success: true,
		Key:     key,
		Value:   value,
	}
	json.NewEncoder(w).Encode(response)
}

func handleQPush(w http.ResponseWriter, ds *DataStore, req Request) {
	key := req.Key
	values, ok := req.Options["values"]
	if !ok {
		response := Response{
			Success: false,
			Message: "Values not provided",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	ds.QPush(key, strings.Fields(values))
	response := Response{
		Success: true,
		Key:     key,
		Message: "Values added to queue",
	}
	json.NewEncoder(w).Encode(response)
}

func handleQPop(w http.ResponseWriter, ds *DataStore, req Request) {
	key := req.Key
	value, ok := ds.QPop(key)
	if !ok {
		response := Response{
			Success: false,
			Message: "Queue is empty or does not exist",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}
	response := Response{
		Success: true,
		Key:     key,
		Value:   value,
	}
	json.NewEncoder(w).Encode(response)
}

func handleBQPop(w http.ResponseWriter, ds *DataStore, req Request) {
	key := req.Key
	timeout, err := strconv.ParseFloat(req.Args[0], 64)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Invalid timeout value",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	queue, ok := ds.store.Load(key)
	if !ok {
		response := Response{
			Success: false,
			Message: "Key not found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Lock the queue to prevent concurrent reads
	queueLock, _ := ds.locks.LoadOrStore(key, &sync.Mutex{})
	queueLock.(*sync.Mutex).Lock()
	defer queueLock.(*sync.Mutex).Unlock()

	// If the queue is empty and timeout is 0, return immediately
	if queue.Len() == 0 && timeout == 0 {
		response := Response{
			Success: false,
			Message: "Queue is empty",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Wait for a value to be added to the queue
	startTime := time.Now()
	for queue.Len() == 0 {
		elapsedTime := time.Since(startTime).Seconds()
		if elapsedTime >= timeout {
			response := Response{
				Success: false,
				Message: "Timeout",
			}
			w.WriteHeader(http.StatusRequestTimeout)
			json.NewEncoder(w).Encode(response)
			return
		}

		time.Sleep(100 * time.Millisecond)
	}

	// Pop the last value from the queue and return it
	val := queue.Remove(queue.Back())
	response := Response{
		Success: true,
		Value:   val,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}