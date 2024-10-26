// Copyright 2023 Google LLC
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

package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	z3 "github.com/aclements/go-z3/z3"
	"github.com/simmonmt/aoc/2023/common/filereader"
	"github.com/simmonmt/aoc/2023/common/logger"
	"github.com/simmonmt/aoc/2023/common/pos"
)

var (
	verbose = flag.Bool("verbose", false, "verbose")
	input   = flag.String("input", "", "input file")
)

type StoneSpec struct {
	P, V pos.P3
}

func parseInput(lines []string) ([]StoneSpec, error) {
	out := []StoneSpec{}
	for i, line := range lines {
		parts := strings.Fields(line)
		for j := range parts {
			parts[j] = strings.TrimRight(parts[j], ",")
		}

		nums := []int{}
		for j, part := range parts {
			if j != 3 {
				num, err := strconv.Atoi(part)
				if err != nil {
					return nil, fmt.Errorf(`%d: num %d: "%s": %v`,
						i+1, j+1, part, err)
				}
				nums = append(nums, num)
			}
		}

		p := pos.P3{X: nums[0], Y: nums[1], Z: nums[2]}
		v := pos.P3{X: nums[3], Y: nums[4], Z: nums[5]}

		out = append(out, StoneSpec{p, v})
	}
	return out, nil
}

type Stone2D struct {
	Spec StoneSpec
	M, B float64
}

func Find2DSlope(spec StoneSpec) (m, b float64) {
	// m = spec.V.Y / spec.V.X
	// spec.P.Y = m * spec.P.X + b
	// spec.P.Y - m * spec.P.X = b

	m = float64(spec.V.Y) / float64(spec.V.X)
	b = float64(spec.P.Y) - m*float64(spec.P.X)
	return
}

func Intersects(a, b *Stone2D, lo, hi float64) bool {
	if a.M == b.M {
		return a.B == b.B
	}

	// y=m1x+b1 y=m2x+b2
	// m1x+b1=m2x+b2
	// m1x-m2x=+b2-b1
	// x=(b2-b1)/(m1-m2)

	x := (b.B - a.B) / (a.M - b.M)
	y := a.M*x + a.B

	at := (x - float64(a.Spec.P.X)) / float64(a.Spec.V.X)
	bt := (x - float64(b.Spec.P.X)) / float64(b.Spec.V.X)

	intersects := x >= lo && x <= hi && y >= lo && y <= hi
	//logger.Infof("%v and %v %v at %v bt %v pos %v,%v", *a, *b, intersects, at, bt, x, y)

	return intersects && at >= 0 && bt >= 0
}

func solveA(specs []StoneSpec, lo, hi float64) int {
	stones := make([]Stone2D, len(specs))
	for i := range specs {
		stone := &stones[i]
		stone.Spec = specs[i]
		stone.M, stone.B = Find2DSlope(specs[i])
	}

	num := 0
	for i := 0; i < len(stones); i++ {
		a := &stones[i]
		for j := i + 1; j < len(stones); j++ {
			b := &stones[j]
			if Intersects(a, b, lo, hi) {
				num++
			}
		}
	}
	return num
}

func makeInt[K int | int64](ctx *z3.Context, num K) z3.Int {
	return ctx.FromInt(int64(num), ctx.IntSort()).(z3.Int)
}

func modelInt(model *z3.Model, zi z3.Value) int64 {
	v, isLiteral, ok := model.Eval(zi, true).(z3.Int).AsInt64()
	if !ok || !isLiteral {
		logger.Fatalf("eval %v => v %v isLiteral %v ok %v",
			zi, v, isLiteral, ok)
	}
	return v
}

type Z3Pos struct {
	X, Y, Z z3.Int
}

func verifyPath(stone StoneSpec, projPos, projVel [3]int64) error {
	ctx := z3.NewContext(nil)
	solver := z3.NewSolver(ctx)

	t := ctx.IntConst("t")

	solver.Assert(
		makeInt(ctx, stone.P.X).Add(makeInt(ctx, stone.V.X).Mul(t)).Eq(
			makeInt(ctx, projPos[0]).Add(makeInt(ctx, projVel[0]).Mul(t))))
	solver.Assert(
		makeInt(ctx, stone.P.Y).Add(makeInt(ctx, stone.V.Y).Mul(t)).Eq(
			makeInt(ctx, projPos[1]).Add(makeInt(ctx, projVel[1]).Mul(t))))
	solver.Assert(
		makeInt(ctx, stone.P.Z).Add(makeInt(ctx, stone.V.Z).Mul(t)).Eq(
			makeInt(ctx, projPos[2]).Add(makeInt(ctx, projVel[2]).Mul(t))))

	if sat, err := solver.Check(); !sat || err != nil {
		return fmt.Errorf("check failed: sat %v err %v", sat, err)
	}
	tVal := modelInt(solver.Model(), t)

	logger.Infof("stone %v at t=%v", stone, tVal)
	return nil
}

func verifyPaths(stones []StoneSpec, projPos, projVel [3]int64) error {
	for i := 0; i < len(stones); i++ {
		if err := verifyPath(stones[i], projPos, projVel); err != nil {
			return fmt.Errorf("failed to verify stone %d %v: %v",
				i, stones[i], err)
		}
	}
	return nil
}

func solveB(stones []StoneSpec) int64 {
	ctx := z3.NewContext(nil)

	solver := z3.NewSolver(ctx)

	projZ3Pos := Z3Pos{
		X: ctx.IntConst("projPosX"),
		Y: ctx.IntConst("projPosY"),
		Z: ctx.IntConst("projPosZ"),
	}

	projZ3Vel := Z3Pos{
		X: ctx.IntConst("projVelX"),
		Y: ctx.IntConst("projVelY"),
		Z: ctx.IntConst("projVelZ"),
	}

	tZ3 := [3]z3.Int{
		ctx.IntConst("ta"),
		ctx.IntConst("tb"),
		ctx.IntConst("tc"),
	}

	for i := 0; i < 3; i++ {
		stonePos := [3]z3.Int{
			makeInt(ctx, stones[i].P.X),
			makeInt(ctx, stones[i].P.Y),
			makeInt(ctx, stones[i].P.Z),
		}

		stoneVel := [3]z3.Int{
			makeInt(ctx, stones[i].V.X),
			makeInt(ctx, stones[i].V.Y),
			makeInt(ctx, stones[i].V.Z),
		}

		solver.Assert(
			stonePos[0].Add(stoneVel[0].Mul(tZ3[i])).Eq(
				projZ3Pos.X.Add(projZ3Vel.X.Mul(tZ3[i]))))
		solver.Assert(
			stonePos[1].Add(stoneVel[1].Mul(tZ3[i])).Eq(
				projZ3Pos.Y.Add(projZ3Vel.Y.Mul(tZ3[i]))))
		solver.Assert(
			stonePos[2].Add(stoneVel[2].Mul(tZ3[i])).Eq(
				projZ3Pos.Z.Add(projZ3Vel.Z.Mul(tZ3[i]))))
	}

	if sat, err := solver.Check(); !sat || err != nil {
		logger.Fatalf("check failed: sat %v err %v", sat, err)
	}
	model := solver.Model()

	projPos := [3]int64{
		modelInt(model, projZ3Pos.X),
		modelInt(model, projZ3Pos.Y),
		modelInt(model, projZ3Pos.Z),
	}

	projVel := [3]int64{
		modelInt(model, projZ3Vel.X),
		modelInt(model, projZ3Vel.Y),
		modelInt(model, projZ3Vel.Z),
	}

	if err := verifyPaths(stones, projPos, projVel); err != nil {
		logger.Fatalf("bad verify: %v", err)
	}

	logger.Infof("result: pos %v, vel %v", projPos, projVel)

	return projPos[0] + projPos[1] + projPos[2]
}

func main() {
	flag.Parse()
	logger.Init(*verbose)

	if *input == "" {
		logger.Fatalf("--input is required")
	}

	lines, err := filereader.Lines(*input)
	if err != nil {
		logger.Fatalf("failed to read input: %v", err)
	}

	input, err := parseInput(lines)
	if err != nil {
		logger.Fatalf("failed to parse input: %v", err)
	}

	fmt.Println("A", solveA(input, float64(200000000000000), float64(400000000000000)))
	fmt.Println("B", solveB(input))
}
