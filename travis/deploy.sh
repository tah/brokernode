#!/bin/bash

echo "HEREEEEEEEEEEEEEEEE"
pwd

ssh ubuntu@52.14.218.135 -i ./travis/id_rsa <<EOF
  echo "PRINT SOMETHINGGGGGGGGGGGGGGGGG"
EOF
