package main

import (
	"context"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/shaj13/raft"
	"github.com/shaj13/raft/transport"
	"github.com/shaj13/raft/transport/raftgrpc"
	"google.golang.org/grpc"
)

// Plan -
// Packages
// 	| - FSM 		=> This holds our data storage backend
// 	| - Server 		=> This is the HTTP server interface to our store
// 	| - Raft 		=> This holds the raft layer to communicate with other members on the cluster

type Entry struct {
	Key   string
	Value string
}

type StateMachine struct {
	mu sync.Mutex
	kv map[string][]string
}

var (
	node *raft.Node
	fsm  *StateMachine
)

/*
   STATE MACHINE METHODS
   * TODO see if it's possible to add support for random BLOBs or JSON Objects instead of just text
   * TODO Add Timestamp on each message
*/

func NewStateMachine() *StateMachine {
	return &StateMachine{
		kv: make(map[string][]string),
	}
}

// Function to apply a new data element to the state machine
func (s *StateMachine) Apply(data []byte) {

	var e Entry
	if err := json.Unmarshal(data, &e); err != nil {
		log.Println("Unable to unmarshal entry ->", err)
		return
	}
	log.Println("Unmarshalled entry successfully")

	s.mu.Lock()
	defer s.mu.Unlock()
	// TODO Refactor this to support writing the latest message to the first part of the list
	s.kv[e.Key] = append([]string{e.Value}, s.kv[e.Key]...)
	// s.kv[e.Key] = e.Value
}

// Function to write the state to Disk
func (s *StateMachine) Snapshot() (io.ReadCloser, error) {
	// TODO Refactor this to write to directories based on topic
	s.mu.Lock()
	defer s.mu.Unlock()
	buf, err := json.Marshal(&s.kv)
	if err != nil {
		return nil, err
	}
	return io.NopCloser(strings.NewReader(string(buf))), nil
}

// Function to read the state from disk
func (s *StateMachine) Restore(r io.ReadCloser) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, &s.kv)
	if err != nil {
		return err
	}

	return r.Close()
}

// Function to read from the FSM
func (s *StateMachine) Read(key string) []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	// TODO Adjust logic to read latest/all messages from this block
	return s.kv[key]
}

func (s *StateMachine) ReadLatest(key string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.kv[key][0]
}

/*
   Request Handler Methods
*/

func history(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()

	if err := node.LinearizableRead(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	value := fsm.Read(key)
	val, err := json.Marshal(value)
	if err != nil {
		w.Write([]byte("Cannot Marshal!"))
	} else {
		w.Write(val)
	}
	w.Write([]byte{'\n'})
}

func get(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()

	if err := node.LinearizableRead(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	value := fsm.ReadLatest(key)
	// Return the Latest Message in the message queue for the specific topic
	w.Write([]byte(value))
	w.Write([]byte{'\n'})
}

func nodes(w http.ResponseWriter, r *http.Request) {
	raws := []raft.RawMember{}
	membs := node.Members()
	for _, m := range membs {
		raws = append(raws, m.Raw())
	}

	buf, err := json.Marshal(raws)
	if err != nil {
		panic(err)
	}

	w.Write(buf)
	w.Write([]byte{'\n'})
}

func removeNode(w http.ResponseWriter, r *http.Request) {
	sid := mux.Vars(r)["id"]
	id, err := strconv.ParseUint(sid, 0, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()

	if err := node.RemoveMember(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func save(w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buf, new(Entry)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()

	if err := node.Replicate(ctx, buf); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {

	// Add the command line arguments to parse
	// https://betterprogramming.pub/writing-distributed-and-replicated-state-machine-in-golang-using-raft-dad79b58cd3a

	//  Various command-line flags
	addr := flag.String("raft", "", "raft server address")
	join := flag.String("join", "", "join cluster address")
	api := flag.String("api", "", "api server address")
	state := flag.String("state_dir", "", "raft state directory (WAL, Snapshots)")
	log.Println(addr, join, api, state)
	flag.Parse()

	router := mux.NewRouter()
	router.HandleFunc("/", http.HandlerFunc(save)).Methods("PUT", "POST")
	router.HandleFunc("/{key}", http.HandlerFunc(get)).Methods("GET")
	router.HandleFunc("/history/{key}", http.HandlerFunc(history)).Methods("GET")
	router.HandleFunc("/mgmt/nodes", http.HandlerFunc(nodes)).Methods("GET")
	router.HandleFunc("/mgmt/nodes/{id}", http.HandlerFunc(removeNode)).Methods("DELETE")

	var (
		opts      []raft.Option
		startOpts []raft.StartOption
	)

	startOpts = append(startOpts, raft.WithAddress(*addr))
	opts = append(opts, raft.WithStateDIR(*state))

	if *join != "" {
		opt := raft.WithFallback(
			raft.WithJoin(*join, time.Second),
			raft.WithRestart(),
		)
		startOpts = append(startOpts, opt)
	} else {
		opt := raft.WithFallback(
			raft.WithInitCluster(),
			raft.WithRestart(),
		)
		startOpts = append(startOpts, opt)
	}
	log.Println("Raft Options -> ", startOpts, opts)

	raftgrpc.Register(
		raftgrpc.WithDialOptions(grpc.WithInsecure()),
	)
	fsm = NewStateMachine()
	node = raft.NewNode(fsm, transport.GRPC, opts...)
	raftServer := grpc.NewServer()
	raftgrpc.RegisterHandler(raftServer, node.Handler())

	go func() {
		lis, err := net.Listen("tcp", *addr)
		if err != nil {
			log.Fatal(err)
		}

		err = raftServer.Serve(lis)
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		err := node.Start(startOpts...)
		if err != nil && err != raft.ErrNodeStopped {
			log.Fatal(err)
		}
	}()

	if err := http.ListenAndServe(*api, router); err != nil {
		log.Fatal(err)
	}

}
