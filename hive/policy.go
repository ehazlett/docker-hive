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
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/garyburd/redigo/redis"
)

type (
	RunPolicy interface {
		Name() string
		GetNodes(int64, string) ([]string, error)
	}

	RandomPolicy struct {
		RedisPool *redis.Pool
	}

	UniquePolicy struct {
		RedisPool *redis.Pool
	}
)

func allNodes(pool *redis.Pool) ([]string, error) {
	nodes := []string{}
	conn := pool.Get()
	keys, err := redis.Strings(conn.Do("KEYS", fmt.Sprintf("%s:*", NODE_KEY)))
	if err != nil {
		return nodes, err
	}
	for _, k := range keys {
		n := strings.Split(k, ":")
		nodes = append(nodes, n[2])
	}
	return nodes, nil
}

func getNodesByZone(pool *redis.Pool, zone string) ([]string, error) {
	nodes := []string{}
	conn := pool.Get()
	keys, err := redis.Strings(conn.Do("KEYS", fmt.Sprintf("%s:%s:*", NODE_KEY, zone)))
	if err != nil {
		return nodes, err
	}
	for _, k := range keys {
		n := strings.Replace(k, fmt.Sprintf("%s:%s:", NODE_KEY, zone), "", 1)
		nodes = append(nodes, n)
	}
	return nodes, nil

}

// Random Policy
func (p *RandomPolicy) Name() string {
	return "random"
}

func (p *RandomPolicy) GetNodes(num int64, zone string) ([]string, error) {
	// TODO: return multiple nodes based upon random
	nodes := []string{}
	zoneNodes, err := getNodesByZone(p.RedisPool, zone)
	if err != nil {
		return nodes, err
	}
	numNodes := len(zoneNodes)
	if numNodes == 0 {
		return nodes, errors.New("No nodes found in that zone")
	}
	// get as many as requested
	for i := 0; i < int(num); i++ {
		r := rand.Intn(numNodes)
		nodes = append(nodes, zoneNodes[r])
	}
	return nodes, nil
}

// Unique Policy
func (p *UniquePolicy) Name() string {
	return "unique"
}

func (p *UniquePolicy) GetNodes(num int64, zone string) ([]string, error) {
	// TODO: return unique nodes based upon criteria
	return []string{}, nil
}
