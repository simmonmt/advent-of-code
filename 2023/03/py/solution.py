"""The AOC 2023 day 03 solution code.

Input is a 2x2 grid. Each cell is a '.' (empty), a digit, or a symbol.
Consecutive digits form numbers. Numbers are only horizontal. See the unit
test for sample input.

Problem A: Find the sum of all numbers that are adjacent to a symbol (even
  diagonally).
Problem B: Some '*' characters have exactly two adjacent numbers. Return the
  sum of the product of the two numbers adjacent to each qualifying '*'.

We make a point of processing the input as a stream. It's more efficient to
do it this way (vs reading the entire file into memory).
"""

from typing import Dict, Generator, Iterator, List, Tuple
from collections import namedtuple

Finding = namedtuple('Finding', ['loc', 'num'])
Lines = namedtuple('Lines', ['prev', 'cur', 'next_'])
Pos = namedtuple('Pos', ['x', 'y'])


def find_number(s: str, x: int) -> Tuple[int, int]:
    """Find a number in a string.

    Given a string like '...123...' we can be called with x set to any index
    within the number (x=3..5). Returns the index of the first character of the
    number as well as the number itself.
    """

    start, end = x, x
    while start > 0 and s[start-1].isdigit():
        start -= 1
    while end < len(s)-1 and s[end+1].isdigit():
        end += 1

    return start, int(s[start:end+1])


def find_numbers_at(prev: str, cur: str, next_: str, x: int, y: int) -> List[Finding]:
    """Find the numbers adjacent to a given character in the current line.

    x is an index into cur. This function finds all numbers that have digits in
    any of the cells adjacent to cur[x]. Adjacency includes diagonals. y gives the
    y coordinate for cur, with prev and next_ containing the lines before and after,
    respectively, cur.
    """

    out: List[Finding] = []

    def save_number(s: str, x: int, y: int):
        start, num = find_number(s, x)
        out.append(Finding(Pos(start, y), num))

    def check_other(s: str, x: int, y: int):
        if s[x].isdigit():
            save_number(s, x, y)

            # We don't need to check x-1 or x+1 because save_number did it
            # for us (it looked at adjacent characters).
        else:
            if x > 0 and s[x-1].isdigit():
                save_number(s, x-1, y)
            if x < len(cur)-1 and s[x+1].isdigit():
                save_number(s, x+1, y)

    if prev:
        check_other(prev, x, y-1)
    if next_:
        check_other(next_, x, y+1)

    if x > 0 and cur[x-1].isdigit():
        save_number(cur, x-1, y)
    if x < len(cur)-1 and cur[x+1].isdigit():
        save_number(cur, x+1, y)

    return out


def find_numbers(prev: str, cur: str, next_: str, y: int) -> List[Finding]:
    """Find numbers adjacent to any symbol in the current line.

    Returns any numbers (in the form of Finding instances) adjacent to each
    symbols (non-digit, non-'.' character) in cur. y contains the y coordinate
    for cur, used to calculate the location in returned findings.
    """

    out: List[Finding] = []

    for x, c in enumerate(cur):
        if c == '.' or c.isdigit():
            continue  # We only want symbols

        out.extend(find_numbers_at(prev, cur, next_, x, y))

    return out


def walker(input_: Iterator[str]) -> Generator[Lines, None, None]:
    """Walks an interator, returning previous, current, and next lines.

    This generator is used by callers that process a stream, but need the
    previous and next lines in addition to the current one. prev and next_
    will be empty at the beginning and end of the iteration, respectively.
    """

    prev = ""
    try:
        next_ = next(input_).rstrip()
    except StopIteration:
        return

    while next_:
        cur = next_
        try:
            next_ = next(input_).rstrip()
        except StopIteration:
            next_ = ""

        yield Lines(prev, cur, next_)

        prev = cur


def solve_a(input_: Iterator[str]):
    """Solve part A."""

    # A number could be adjacent to multiple symbols (maybe?). This dict
    # ensures we process each number once.
    nums: Dict[Pos, int] = {}

    for y, lines in enumerate(walker(input_)):
        for finding in find_numbers(lines.prev, lines.cur, lines.next_, y):
            nums[finding.loc] = finding.num

    return sum(nums.values())


def solve_b(input_: Iterator[str]):
    """Solve part B."""

    out = 0
    for y, lines in enumerate(walker(input_)):
        for x, c in enumerate(lines.cur):
            if c != '*':
                continue

            findings = find_numbers_at(lines.prev, lines.cur, lines.next_, x, y)
            if len(findings) != 2:
                continue

            out += findings[0].num * findings[1].num

    return out
