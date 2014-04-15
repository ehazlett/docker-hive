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
	_ "log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ehazlett/docker-hive/utils"
)

var (
	listenPort  = 45000
	nodeName    = "testNode"
	testAddress = fmt.Sprintf("http://localhost:%d", listenPort)
)

func newTestEngine() *Engine {
	dockerPath := os.Getenv("DOCKER_PATH")
	if dockerPath == "" {
		dockerPath = "/var/run/docker.sock"
	}
	pool := utils.NewRedisPool("127.0.0.1", 6379, "")
	testEngine := NewEngine("", listenPort, dockerPath, "test", nodeName, "default", pool, "default")
	testEngine.Start()
	return testEngine
}

func getTestUrl(path string) string {
	return fmt.Sprintf("%s%s", testAddress, path)
}

func TestHandleIndexReturnsWithStatusOK(t *testing.T) {
	request, _ := http.NewRequest("GET", getTestUrl("/"), nil)
	response := httptest.NewRecorder()

	testEngine := newTestEngine()
	testEngine.indexHandler(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("Non-expected status code %v: expected %v\nbody: %v", response.Code, "200", response.Body)
	}
}

func TestHandlePingReturnsWithStatusOK(t *testing.T) {
	request, _ := http.NewRequest("GET", getTestUrl("/ping"), nil)
	response := httptest.NewRecorder()

	testEngine := newTestEngine()
	testEngine.indexHandler(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("Non-expected status code %v: expected %v\nbody: %v", response.Code, "200", response.Body)
	}
}
