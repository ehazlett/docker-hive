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
package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/ehazlett/docker-hive/hive"
	"github.com/ehazlett/docker-hive/utils"
)

const (
	VERSION = "0.3.0"
)

var (
	dockerPath string
	version    bool
	nodeName   string
	port       int
	host       string
	redisHost  string
	redisPort  int
	redisPass  string
	zone       string
	runPolicy  string
)

func init() {
	flag.StringVar(&dockerPath, "docker", "/var/run/docker.sock", "Path to Docker socket")
	flag.BoolVar(&version, "version", false, "Shows version")
	flag.StringVar(&nodeName, "n", "", "Node name (default: hostname)")
	flag.StringVar(&host, "l", "", "Listen address (also used for communication with ndoes)")
	flag.IntVar(&port, "p", 4500, "Listen port")
	flag.StringVar(&zone, "z", "default", "Zone for node")
	flag.StringVar(&runPolicy, "r", "default", "Run Policy")
	flag.StringVar(&redisHost, "redis-host", "localhost", "Redis hostname")
	flag.IntVar(&redisPort, "redis-port", 6379, "Redis port")
	flag.StringVar(&redisPass, "redis-password", "", "Redis password")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [arguments]\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	log.SetFlags(0)
	flag.Parse()
	if version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.LstdFlags)
	log.Printf("Docker Hive %s\n", VERSION)

	// connect to redis
	pool := utils.NewRedisPool(redisHost, redisPort, redisPass)
	// set node name
	if nodeName == "" {
		name, err := os.Hostname()
		if err != nil {
			log.Printf("Error getting hostname: %s", err)
			nodeName = "localhost"
		}
		nodeName = name
	}
	// start node
	engine := hive.NewEngine(host, port, dockerPath, VERSION, nodeName, zone, pool, runPolicy)

	waiter, err := engine.Start()
	if err != nil {
		log.Fatal(err)
		return
	}
	waiter.Wait()
}
