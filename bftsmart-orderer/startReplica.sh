#!/bin/bash
java -DNODE_ID=$1 -cp bin/orderingservice.jar:lib/* bft.BFTNode $@
