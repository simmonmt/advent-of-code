package puzzle

var (
	map1 = []string{
		// 0000000011111111111
		// 2345678901234567890
		"         A           ", // 0
		"         A           ", // 1
		"  #######.#########  ", // 2
		"  #######.........#  ", // 3
		"  #######.#######.#  ", // 4
		"  #######.#######.#  ", // 5
		"  #######.#######.#  ", // 6
		"  #####  B    ###.#  ", // 7
		"BC...##  C    ###.#  ", // 8
		"  ##.##       ###.#  ", // 9
		"  ##...DE  F  ###.#  ", // 10
		"  #####    G  ###.#  ", // 11
		"  #########.#####.#  ", // 12
		"DE..#######...###.#  ", // 13
		"  #.#########.###.#  ", // 14
		"FG..#########.....#  ", // 15
		"  ###########.#####  ", // 16
		"             Z       ", // 17
		"             Z       ", // 18
	}

	map2 = []string{
		// 0000000011111111112222222222333333333344444
		// 2345678901234567890123456789012345678901234
		"             Z L X W       C                 ", // 0
		"             Z P Q B       K                 ", // 1
		"  ###########.#.#.#.#######.###############  ", // 2
		"  #...#.......#.#.......#.#.......#.#.#...#  ", // 3
		"  ###.#.#.#.#.#.#.#.###.#.#.#######.#.#.###  ", // 4
		"  #.#...#.#.#...#.#.#...#...#...#.#.......#  ", // 5
		"  #.###.#######.###.###.#.###.###.#.#######  ", // 6
		"  #...#.......#.#...#...#.............#...#  ", // 7
		"  #.#########.#######.#.#######.#######.###  ", // 8
		"  #...#.#    F       R I       Z    #.#.#.#  ", // 9
		"  #.###.#    D       E C       H    #.#.#.#  ", // 10
		"  #.#...#                           #...#.#  ", // 11
		"  #.###.#                           #.###.#  ", // 12
		"  #.#....OA                       WB..#.#..ZH", // 13
		"  #.###.#                           #.#.#.#  ", // 14
		"CJ......#                           #.....#  ", // 15
		"  #######                           #######  ", // 16
		"  #.#....CK                         #......IC", // 17
		"  #.###.#                           #.###.#  ", // 18
		"  #.....#                           #...#.#  ", // 19
		"  ###.###                           #.#.#.#  ", // 20
		"XF....#.#                         RF..#.#.#  ", // 21
		"  #####.#                           #######  ", // 22
		"  #......CJ                       NM..#...#  ", // 23
		"  ###.#.#                           #.###.#  ", // 24
		"RE....#.#                           #......RF", // 25
		"  ###.###        X   X       L      #.#.#.#  ", // 26
		"  #.....#        F   Q       P      #.#.#.#  ", // 27
		"  ###.###########.###.#######.#########.###  ", // 28
		"  #.....#...#.....#.......#...#.....#.#...#  ", // 29
		"  #####.#.###.#######.#######.###.###.#.#.#  ", // 30
		"  #.......#.......#.#.#.#.#...#...#...#.#.#  ", // 31
		"  #####.###.#####.#.#.#.#.###.###.#.###.###  ", // 32
		"  #.......#.....#.#...#...............#...#  ", // 33
		"  #############.#.#.###.###################  ", // 34
		"               A O F   N                     ", // 35
		"               A A D   M                     ", // 36
	}
)