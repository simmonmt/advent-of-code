package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Firewall struct {
	layers   map[int]int
	maxLayer int

	userPosition     int
	scannerPositions map[int]int
	scannerForward   map[int]bool
}

func NewFirewall(layers map[int]int) *Firewall {
	f := &Firewall{
		layers:           map[int]int{},
		maxLayer:         0,
		userPosition:     -1,
		scannerPositions: map[int]int{},
		scannerForward:   map[int]bool{},
	}

	for num, size := range layers {
		f.addLayer(num, size)
	}

	return f
}

func (f *Firewall) addLayer(num, size int) {
	f.layers[num] = size
	if num > f.maxLayer {
		f.maxLayer = num
	}
	f.scannerPositions[num] = 0
	f.scannerForward[num] = true
}

func (f *Firewall) Dump(out io.Writer) {
	for i := 0; i <= f.maxLayer; i++ {
		userInLayer := f.userPosition == i

		if size, found := f.layers[i]; found {
			fmt.Fprintf(out, " %2d", i)

			for j := 0; j < size; j++ {
				contents := ""
				if f.scannerPositions[i] == j {
					contents = "S"
				} else {
					contents = " "
				}

				bounds := "[]"
				if userInLayer && j == 0 {
					bounds = "()"
				}

				fmt.Fprintf(out, " %c%v%c", bounds[0], contents, bounds[1])
			}

			if f.scannerForward[i] {
				fmt.Fprintf(out, " ->")
			} else {
				fmt.Fprintf(out, " <-")
			}
		} else {
			fmt.Fprintf(out, "  : ")
			if userInLayer {
				fmt.Fprintf(out, "(.)")
			}
			fmt.Fprintf(out, "...")
		}
		fmt.Fprintf(out, "\n")
	}
}

func (f *Firewall) AdvanceUser(userMovement int) int {
	severity := 0

	f.userPosition += userMovement
	//fmt.Printf("user position now %d\n", f.userPosition)
	if scannerPosition, found := f.scannerPositions[f.userPosition]; found {
		//fmt.Printf("found scanner in %d, scanner in %d\n", f.userPosition, scannerPosition)
		if scannerPosition == 0 {
			severity = (f.userPosition + 1) * f.layers[f.userPosition]
			//fmt.Printf("collision; severity %d\n", severity)
		}
	}

	return severity
}

func (f *Firewall) AdvanceScanners(amt int) {
	for num, position := range f.scannerPositions {
		scannerAmt := amt
		for scannerAmt%((f.layers[num]-1)*2) != 0 {
			scannerAmt--

			if f.scannerForward[num] {
				position++
				if position == f.layers[num] {
					f.scannerForward[num] = false
					position -= 2
				}
			} else {
				position--
				if position < 0 {
					f.scannerForward[num] = true
					position += 2
				}
			}
			f.scannerPositions[num] = position
		}
	}
}

func (f *Firewall) UserEscaped() bool {
	return f.userPosition > f.maxLayer
}

func (f *Firewall) IsInitial() bool {
	for _, position := range f.scannerPositions {
		if position != 0 {
			return false
		}
	}
	return true
}

func readLayers(in io.Reader) (map[int]int, error) {
	layers := map[int]int{}

	reader := bufio.NewReader(in)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)

		pieces := strings.Split(line, ": ")

		layerNum, err := strconv.Atoi(pieces[0])
		if err != nil {
			return nil, fmt.Errorf("failed to parse layer number %v\n",
				pieces[0])
		}
		layerSize, err := strconv.Atoi(pieces[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse layer size %v\n",
				pieces[1])
		}

		layers[layerNum] = layerSize
	}

	return layers, nil
}

func runFirewall(layers map[int]int, initialDelay int) int {
	firewall := NewFirewall(layers)

	firewall.AdvanceScanners(initialDelay)

	//fmt.Printf("after initial delay:\n")
	//firewall.Dump(os.Stdout)

	if initialDelay != 0 && firewall.IsInitial() {
		panic(fmt.Sprintf("initial at %d", initialDelay))
	}

	round := initialDelay
	for !firewall.UserEscaped() {
		round++
		//fmt.Printf("--- beginning %d\n", round)

		if severity := firewall.AdvanceUser(1); severity > 0 {
			return severity
		}
		//fmt.Printf("severity now %d\n\n", severity)
		//firewall.Dump(os.Stdout)

		firewall.AdvanceScanners(1)
		//fmt.Printf("\n")
		//firewall.Dump(os.Stdout)
		//fmt.Printf("\n")
	}

	return 0
}

func main() {
	layers, err := readLayers(os.Stdin)
	if err != nil {
		log.Fatalf("failed to read input: %v", err)
	}

	for i := 0; ; i++ {
		if i%1000 == 0 {
			fmt.Printf("%d\n", i)
		}
		if severity := runFirewall(layers, i); severity == 0 {
			fmt.Printf("delay: %d\n", i)
			break
		}
	}
}
