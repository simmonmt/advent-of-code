import argparse
import importlib

solution = importlib.import_module('03.py.solution')

parser = argparse.ArgumentParser(description='AoC 2023 Day 03')
parser.add_argument('--input', type=str, required=True, help='input file')


def main():
    args = parser.parse_args()

    with open(args.input) as input_file:
        print("A: %d" % (solution.solve_a(input_file)))

    with open(args.input) as input_file:
        print("B: %d" % (solution.solve_b(input_file)))


if __name__ == '__main__':
    main()
