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
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
)

type (
	ContainerJob struct {
		Config *ContainerConfig
		Zone   string
	}
	ImageJob struct {
		Image string
		Zone  string
	}
	Scheduler interface {
		AddContainerJob(j *ContainerJob) (string, error)
		RemoveContainerJob(id string) (bool, error)
		AddImageJob(j *ImageJob) (bool, error)
		RemoveImageJob(id string) (bool, error)
	}
	DefaultScheduler struct {
		RedisPool *redis.Pool
	}
)

func (s *DefaultScheduler) AddContainerJob(j *ContainerJob) (string, error) {
	// config
	buf := bytes.NewBufferString("")
	if err := json.NewEncoder(buf).Encode(j.Config); err != nil {
		return "", err
	}
	// add to redis
	conn := s.RedisPool.Get()
	k := fmt.Sprintf("%s:%s", CONTAINER_JOB_KEY, j.Config.Name)
	conn.Do("SET", k, buf)
	log.Printf("Added container job %s for zone %s", j.Config.Name, j.Zone)
	return j.Config.Name, nil
}
func (s *DefaultScheduler) RemoveContainerJob(id string) (bool, error) {
	log.Printf("TODO: Removed job %s", id)
	return true, nil
}
func (s *DefaultScheduler) AddImageJob(j *ImageJob) (bool, error) {
	log.Printf("TODO: Added job %s", j.Image)
	return true, nil
}
func (s *DefaultScheduler) RemoveImageJob(id string) (bool, error) {
	log.Printf("TODO: Removed job %s", id)
	return true, nil
}
