import random
import math

import yaml


TRIALS = 10_000
MAX_EXP = 17


def simulate(n: int) -> bool:
    """Simulates to see if there's a self match for n size"""
    randomized = list(range(n))
    random.shuffle(randomized)
    for i, x in enumerate(randomized):
        if (i == x):
            return True
    return False


def main() -> None:
    """Runs line_calc"""

    print(f'running {TRIALS} with a max of {2**MAX_EXP} ({MAX_EXP})')

    results = {}
    for ex in range(MAX_EXP):
        i = int(2 ** ex)
        print(f'simulating for {i} ({ex})')
        results[ex] = 1 - (sum([simulate(i) for _ in range(TRIALS)]) / TRIALS)
        print(f'{i} ({ex}) result {results[ex]}')

    print(yaml.dump(results))


if __name__ == "__main__":
    main()
