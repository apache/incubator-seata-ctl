/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package model

type Config struct {
	Kubernetes Kubernetes `yaml:"kubernetes"`
	Prometheus Prometheus `yaml:"prometheus"`
	Log        Log        `yaml:"log"`
	Context    Context    `yaml:"context"`
}

type Kubernetes struct {
	Cluster []KubernetesCluster `yaml:"clusters"`
}

type Prometheus struct {
	Servers []Server `yaml:"servers,omitempty"`
}

type Log struct {
	Clusters []Cluster `yaml:"clusters"`
}

// Context Select the appropriate configuration based on the Context field
type Context struct {
	Kubernetes string `yaml:"kubernetes"`
	Prometheus string `yaml:"prometheus"`
	Log        string `yaml:"log"`
}

type KubernetesCluster struct {
	Name           string `yaml:"name"`
	KubeConfigPath string `yaml:"kubeconfigpath"`
	YmlPath        string `yaml:"ymlpath"`
}

type Server struct {
	Name    string `yaml:"name"`
	Address string `yaml:"address"`
	Auth    string `yaml:"auth"`
}

type Cluster struct {
	Name     string `yaml:"name"`
	Types    string `yaml:"types"`
	Address  string `yaml:"address"`
	Source   string `yaml:"source"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Index    string `yaml:"index"`
}
