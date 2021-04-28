#!/usr/bin/env awk

{
  print $0
  if (match($0, /^total:/)) {
    sub(/%/, "", $NF);
    if (strtonum($NF) < target) {
      exit 1;
    }
  }
}
