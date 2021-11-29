// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package node

import (
	"fmt"
	"strconv"
	"strings"
)

type Node struct {
	X, Y     int
	Passcode string
	Path     string
}

func New(x, y int, passcode string) *Node {
	return &Node{
		X:        x,
		Y:        y,
		Passcode: passcode,
		Path:     "",
	}
}

func Deserialize(ser string) (*Node, error) {
	parts := strings.SplitN(ser, ",", 4)
	if len(parts) != 4 {
		return nil, fmt.Errorf("illegal node name %v", ser)
	}

	x, xErr := strconv.ParseInt(parts[0], 10, 32)
	y, yErr := strconv.ParseInt(parts[1], 10, 32)
	if xErr != nil || yErr != nil {
		return nil, fmt.Errorf("failed to parse x,y from %v", ser)
	}

	return &Node{
		X:        int(x),
		Y:        int(y),
		Passcode: parts[2],
		Path:     parts[3],
	}, nil
}

func (n *Node) Serialize() string {
	return fmt.Sprintf("%d,%d,%v,%v", n.X, n.Y, n.Passcode, n.Path)
}
