#!/bin/sh

read a b
echo $(( ($a + $b) / 0 ))
