/*
   Copyright Evan Hazlett

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
package hive

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/ehazlett/docker-hive/utils"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
)

type (
	Engine struct {
		Name       string
		Host       string
		Port       int
		httpServer *http.Server
		waiter     *sync.WaitGroup
		redisPool  *redis.Pool
		Router     *mux.Router
		DockerPath string
		Version    string
		Zone       string
		RunPolicy  RunPolicy
		Scheduler  Scheduler
		Master     bool
	}
	Image struct {
		Id          string
		Created     int64
		RepoTags    []string
		Size        int
		VirtualSize int
	}

	InfoPort struct {
		IP          string
		PrivatePort int
		PublicPort  int
		Type        string
	}

	APIContainer struct {
		Id      string
		Created int
		Image   string
		Status  string
		Command string
		Ports   []InfoPort
		Names   []string
		Node    string
	}
	ContainerInfo struct {
		Container  Container
		ServerName string
	}
)

// Creates a new Engine
func NewEngine(host string, port int, dockerPath string, version string, nodeName string, zone string, redisPool *redis.Pool, runPolicy string) *Engine {
	// select launch policy
	var rp RunPolicy
	switch runPolicy {
	default:
		rp = &RandomPolicy{RedisPool: redisPool}
	case "unique":
		rp = &UniquePolicy{RedisPool: redisPool}
	}
	// scheduler
	scheduler := &DefaultScheduler{RedisPool: redisPool}

	e := &Engine{
		Name:       nodeName,
		Host:       host,
		Port:       port,
		DockerPath: dockerPath,
		waiter:     new(sync.WaitGroup),
		redisPool:  redisPool,
		Router:     mux.NewRouter(),
		Version:    version,
		Zone:       zone,
		RunPolicy:  rp,
		Scheduler:  scheduler,
		Master:     false,
	}

	// check for empty host
	if e.Host == "" {
		addrs, err := net.LookupHost(nodeName)
		// if unable to lookup host, use localhost
		if err != nil {
			e.Host = "localhost"
		}
		localAddrs := map[string]bool{
			"127.0.0.1": true,
		}
		for _, a := range addrs {
			if !localAddrs[a] {
				e.Host = a
				break
			}
		}
	}
	return e
}

// Starts the Engine
func (e *Engine) Start() (*sync.WaitGroup, error) {
	log.Println("Initializing HTTP API")

	// Initialize and start HTTP server.
	e.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", e.Port),
		Handler: e.Router,
	}

	// docker router
	dockerRouter := NewDockerSubrouter(e.Router, e.DockerPath)

	// setup router
	e.Router.HandleFunc("/ping", e.pingHandler).Methods("GET").Name("ping")
	// addon docker router
	e.Router.Handle("/{apiVersion:v1.*}", dockerRouter.Subrouter).Methods("GET", "PUT", "POST", "DELETE")
	// index
	e.Router.HandleFunc("/", e.indexHandler).Methods("GET")

	log.Printf("Server name: %s", e.Name)
	log.Printf("Zone: %s", e.Zone)
	log.Printf("Run Policy: %s", e.RunPolicy.Name())
	log.Printf("Listening at: %s", e.ConnectionString())

	// serve
	go e.listenAndServe()
	e.waiter.Add(1)
	go e.run()
	return e.waiter, nil
}

// Returns the connection string
func (e *Engine) ConnectionString() string {
	return fmt.Sprintf("http://%s:%d", e.Host, e.Port)
}

// Stops the Engine
func (e *Engine) Stop() {
	log.Println("Stopping server")
	e.waiter.Done()
}

// Checks for master node ; self-elects if missing
func (e *Engine) checkMasterStatus() {
	conn := e.redisPool.Get()
	master, err := redis.String(conn.Do("GET", MASTER_KEY))
	if err != nil {
		log.Printf("Assuming master role")
		conn.Do("SET", MASTER_KEY, e.Name)
		// set expiration to avoid the race condition between set
		// and expiring if server crashes
		conn.Do("EXPIRE", MASTER_KEY, MASTER_HEARTBEAT_INTERVAL+1)
	}
	// update master heartbeat
	if master == e.Name {
		conn.Do("EXPIRE", MASTER_KEY, MASTER_HEARTBEAT_INTERVAL+1)
		e.Master = true
	} else {
		e.Master = false
	}
}

// Updates node heartbeat ttl
func (e *Engine) nodeHeartbeat() {
	key := getNodeKey(e.Name, e.Zone)
	conn := e.redisPool.Get()
	conn.Do("SET", key, e.ConnectionString())
	conn.Do("EXPIRE", key, 5)
}

// ---- Handlers ----

// Index handler
func (e *Engine) indexHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(fmt.Sprintf("Docker Hive %s", e.Version)))
}

func (e *Engine) pingHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("pong"))
}

// Generic error handler
func handlerError(msg string, status int, w http.ResponseWriter) {
	w.WriteHeader(status)
	w.Write([]byte(msg))
}

// Proxies requests to the local Docker daemon
func (e *Engine) dockerHandler(w http.ResponseWriter, req *http.Request) {
	utils.ProxyLocalDockerRequest(w, req, e.DockerPath)
}

func (e *Engine) listenAndServe() {
	go func() {
		e.httpServer.ListenAndServe()
	}()
}

func (e *Engine) run() {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)

	heartbeatTick := time.Tick(NODE_HEARTBEAT_INTERVAL * time.Second)
	masterTick := time.Tick(MASTER_HEARTBEAT_INTERVAL * time.Second)

run:
	for {
		select {
		case <-masterTick:
			go e.checkMasterStatus()
		case <-heartbeatTick:
			go e.nodeHeartbeat()
		case <-sig:
			break run
		}
	}
	e.Stop()
}
