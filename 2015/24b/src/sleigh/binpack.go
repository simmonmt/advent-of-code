package sleigh

func OneBinPack(values []int, cap int) []int {
	// The easy way, with lots of allocations

	//fmt.Printf("BinPack(%v, %v)\n", values, cap)

	if cap == 0 {
		//fmt.Printf("  easy out; cap==0")
		return []int{}
	} else if len(values) == 0 {
		//fmt.Printf("  easy fail; len==0, cap>0 %v\n", cap)
		return nil
	}

	for i := 0; i < len(values); i++ {
		cand := values[i]
		//fmt.Printf("  trying cand %v\n", cand)
		if cand > cap {
			//fmt.Printf("    doesn't fit\n")
			continue
		}

		rest := values[i+1 : len(values)]
		if found := OneBinPack(rest, cap-cand); found != nil {
			out := make([]int, len(found)+1)
			out[0] = cand
			for j, val := range found {
				out[j+1] = val
			}
			//fmt.Printf("  worked; returning %v\n", out)
			return out
		}
	}

	//fmt.Printf("  end of list; return nil\n")
	return nil
}

func AllBinPacks(values []int, cap int) [][]int {
	// The easy way, with lots of allocations

	//fmt.Printf("BinPack(%v, %v)\n", values, cap)

	out := [][]int{}

	if cap == 0 {
		//fmt.Printf("  easy out; cap==0")
		out = append(out, []int{})
		return out
	} else if len(values) == 0 {
		//fmt.Printf("  easy fail; len==0, cap>0 %v\n", cap)
		return out
	}

	for i := 0; i < len(values); i++ {
		cand := values[i]
		//fmt.Printf("  trying cand %v\n", cand)
		if cand > cap {
			//fmt.Printf("    doesn't fit\n")
			continue
		}

		rest := values[i+1 : len(values)]
		foundPacks := AllBinPacks(rest, cap-cand)
		for _, foundPack := range foundPacks {
			withCand := make([]int, len(foundPack)+1)
			withCand[0] = cand
			for j, val := range foundPack {
				withCand[j+1] = val
			}
			out = append(out, withCand)
		}
	}

	//fmt.Printf("  end of list; return nil\n")
	return out
}
