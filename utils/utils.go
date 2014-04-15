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
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"

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

// Proxies request to local Docker instance.
func ProxyLocalDockerRequest(w http.ResponseWriter, req *http.Request, dockerPath string) {
	req.ParseForm()
	params := req.Form
	path := fmt.Sprintf("%s?%s", req.URL.Path, params.Encode())
	log.Printf("Proxying Docker request: %s", path)
	c, err := NewDockerClient(dockerPath)
	defer c.Close()
	if err != nil {
		msg := fmt.Sprintf("Error connecting to Docker: %s", err)
		log.Println(msg)
		w.Write([]byte(msg))
		return
	}
	r, err := http.NewRequest(req.Method, path, req.Body)
	if err != nil {
		msg := fmt.Sprintf("Error connecting to Docker: %s", err)
		log.Println(msg)
		w.Write([]byte(msg))
		return
	}
	resp, err := c.Do(r)
	if err != nil {
		msg := fmt.Sprintf("Error connecting to Docker: %s", err)
		w.WriteHeader(resp.StatusCode)
		w.Write([]byte(msg))
		return
	}
	contents, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		msg := fmt.Sprintf("Error connecting to Docker: %s", err)
		log.Println(msg)
		w.WriteHeader(resp.StatusCode)
		w.Write([]byte(msg))
		return
	}
	w.WriteHeader(resp.StatusCode)
	io.WriteString(w, string(contents))
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}
