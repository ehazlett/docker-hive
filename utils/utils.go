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
package utils

import (
	"net"
        "net/http"
	"net/http/httputil"
        "fmt"

	"github.com/garyburd/redigo/redis"
)

const (
	DEFAULT_POOL_SIZE = 10
)

func NewRedisPool(addr string, port int, password string) *redis.Pool {
	return redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
		if err != nil {
			return nil, err
		}
		if password != "" {
			if _, err := c.Do("AUTH", password); err != nil {
				return nil, err
			}
		}
		return c, nil
	}, DEFAULT_POOL_SIZE)
}

// Creates a new Docker client using the Docker unix socket.
func NewDockerClient(dockerSocketPath string) (*httputil.ClientConn, error) {
	conn, err := net.Dial("unix", dockerSocketPath)
	if err != nil {
		return nil, err
	}
	return httputil.NewClientConn(conn, nil), nil
}

// Utility function for copying HTTP Headers.
func CopyHeaders(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
