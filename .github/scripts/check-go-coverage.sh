#!/usr/bin/env bash
set -euo pipefail

# Read raw statement counts so rounding cannot turn 94.99% into a passing 95%.
awk '
  NR == 1 { next }
  NF == 3 {
    package = $1
    sub("/[^/]*$", "", package)
    statements[package] += $2
    total += $2
    if ($3 > 0) {
      covered[package] += $2
      totalCovered += $2
    }
  }
  END {
    if (total == 0) {
      print "No coverage statements found" > "/dev/stderr"
      exit 1
    }
    for (package in statements) {
      if (statements[package] == 0) continue
      printf "%s: %.2f%%\n", package, 100 * covered[package] / statements[package]
      if (100 * covered[package] < 95 * statements[package]) {
        printf "Coverage below 95%% for %s\n", package > "/dev/stderr"
        failed = 1
      }
    }
    printf "Total: %.2f%%\n", 100 * totalCovered / total
    if (100 * totalCovered < 95 * total) failed = 1
    exit failed
  }
' "${1:-coverage.out}"
