import unittest
from solution import solve_a, solve_b

SAMPLE_INPUT = [
    '467..114..\n',
    '...*......\n',
    '..35..633.\n',
    '......#...\n',
    '617*......\n',
    '.....+.58.\n',
    '..592.....\n',
    '......755.\n',
    '...$.*....\n',
    '.664.598..\n',
]


class TestSolution(unittest.TestCase):
    def test_solve_a(self):
        self.assertEqual(solve_a(iter(SAMPLE_INPUT)), 4361)

    def test_solve_b(self):
        self.assertEqual(solve_b(iter(SAMPLE_INPUT)), 467835)


if __name__ == '__main__':
    unittest.main()
