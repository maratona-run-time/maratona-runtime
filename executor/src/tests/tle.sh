#!/bin/sh
read a b
c=0
while [ $a -gt 0 ]
do
    c=$((c+1)) 
    a=$((a - 1))
done
while [ $b -gt 0 ]
do
    c=$((c+1)) 
    b=$((b - 1))
done
echo $c
