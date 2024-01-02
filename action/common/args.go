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

package common

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
)

func ReadArgs(in io.Reader) error {
	os.Args = []string{""}

	scanner := bufio.NewScanner(in)

	var lines []string

	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), "\r\n ")
		if line[len(line)-1] == '\\' {
			line = line[:len(line)-1]
			lines = append(lines, line)
		} else {
			lines = append(lines, line)
			break
		}
	}

	argsStr := strings.Join(lines, " ")
	rawArgs := strings.Split(argsStr, "'")

	if len(rawArgs) != 1 && len(rawArgs) != 3 {
		return errors.New("read args from input error")
	}

	args := strings.Split(rawArgs[0], " ")

	if len(rawArgs) == 3 {
		args = append(args, rawArgs[1])
		args = append(args, strings.Split(rawArgs[2], " ")...)
	}

	for _, arg := range args {
		if arg != "" {
			os.Args = append(os.Args, strings.TrimSpace(arg))
		}
	}
	return nil
}

func ParseDictArg(dataStr string) (map[string]string, error) {
	var data map[string]string
	if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
		return nil, err
	}
	return data, nil
}

func ParseArrayArg(dataStr string) ([]string, error) {
	var data []string
	if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
		return nil, err
	}
	return data, nil
}
