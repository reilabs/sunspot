#!/usr/bin/bash
# This script builds the Noir samples

cd noir-samples

cd hello_world
nargo compile

cd ../sum_a_b
nargo compile