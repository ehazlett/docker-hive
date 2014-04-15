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
)

const (
	CONTAINER_JOB_KEY         = "jobs:containers"
	IMAGE_JOB_KEY             = "jobs:images"
	JOB_KEY                   = "jobs"
	JOB_NODE_KEY              = "nodes:jobs"
	JOB_INTERVAL              = 10
	MASTER_HEARTBEAT_INTERVAL = 2
	MASTER_KEY                = "master"
	NODE_HEARTBEAT_INTERVAL   = 1
	NODE_KEY                  = "nodes"
)

// Returns node key
func getNodeKey(node string, zone string) string {
	return fmt.Sprintf("%s:%s:%s", NODE_KEY, zone, node)
}
