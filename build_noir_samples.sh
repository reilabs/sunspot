#!/usr/bin/bash
# This script builds the Noir samples

git submodule update --init --recursive

cd noir-samples/black_box_functions/aes128encrypt
nargo build
nargo execute
cd ../and
nargo build
nargo execute
cd ../keccakf1600
nargo build
nargo execute
cd ../range
nargo build
nargo execute
cd ../sha256
nargo build
nargo execute
cd ../xor
nargo build
nargo execute

cd ../../expressions/hello_world
nargo build
nargo execute
cd ../lcchecker
nargo build
nargo execute
cd ../linear_equation
nargo build
nargo execute
cd ../polynomial
nargo build
nargo execute
cd ../rock_paper_scissors
nargo build
nargo execute
cd ../square_equation
nargo build
nargo execute
cd ../sum_a_b
nargo build
nargo execute

cd ../../real_world/ProveKit/noir-examples/basic
nargo build
nargo execute
cd ../basic-2
nargo build
nargo execute
cd ../basic-3
nargo build
nargo execute
cd ../noir-native-sha256
nargo build
nargo execute

cd ../../../zk-noir-voting/circuits
nargo build
nargo execute