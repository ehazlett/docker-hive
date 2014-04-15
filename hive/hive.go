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

type (
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
