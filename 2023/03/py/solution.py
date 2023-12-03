# Input is a 2x2 grid. Each cell is a '.' (empty), a digit, or a symbol.
# Consecutive digits form numbers. Numbers are only horizontal. See the unit
# test for sample input.
#
# Problem A: Find the sum of all numbers that are adjacent to a symbol (even
#   diagonally).
# Problem B: Some '*' characters have exactly two adjacent numbers. Return the
#   sum of the product of the two numbers adjacent to each qualifying '*'.
#
# We make a point of processing the input as a stream. It's more efficient to
# do it this way (vs reading the entire file into memory).

from typing import Dict, Iterable, List, Tuple
from collections import namedtuple

Finding = namedtuple('Finding', ['loc', 'num'])
Lines = namedtuple('Lines', ['prev', 'cur', 'next_'])
Pos = namedtuple('Pos', ['x', 'y'])


def find_number(s: str, x: int) -> Tuple[int, int]:
    start, end = x, x
    while start > 0 and s[start-1].isdigit():
        start -= 1
    while end < len(s)-1 and s[end+1].isdigit():
        end += 1

    return start, int(s[start:end+1])


def find_numbers_at(prev: str, cur: str, next_: str, x: int, y: int) -> List[Finding]:
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


def find_numbers(prev: str, cur: str, next: str, y: int) -> List[Finding]:
    out: List[Finding] = []

    for x, c in enumerate(cur):
        if c == '.' or c.isdigit():
            continue  # We only want symbols

        out.extend(find_numbers_at(prev, cur, next, x, y))

    return out


def walker(input: Iterable) -> Tuple[str, str, str]:
    prev, next_ = "", next(input).rstrip()

    while next_:
        cur = next_
        try:
            next_ = next(input).rstrip()
        except StopIteration:
            next_ = ""

        yield Lines(prev, cur, next_)

        prev = cur


def solve_a(input: Iterable[str]):
    nums: Dict[Pos, int] = {}
    for y, lines in enumerate(walker(input)):
        for finding in find_numbers(lines.prev, lines.cur, lines.next_, y):
            nums[finding.loc] = finding.num

    return sum(nums.values())


def solve_b(input: Iterable[str]):
    out = 0
    for y, lines in enumerate(walker(input)):
        for x, c in enumerate(lines.cur):
            if c != '*':
                continue

            findings = find_numbers_at(lines.prev, lines.cur, lines.next_, x, y)
            if len(findings) != 2:
                continue

            out += findings[0].num * findings[1].num

    return out
