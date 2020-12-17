#!/bin/sh

read a b
v = $(seq 1 300000)
echo $(( $a + $b ))