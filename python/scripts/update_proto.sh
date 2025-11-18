#!/bin/bash

# Copy generated Python proto to python package
mkdir -p ../antex_proto
cp -r ../../proto/gen/python/* ../antex_proto/
# Create __init__.py files
find ../antex_proto -type d -exec touch {}/__init__.py \;
# Fix import paths to use antex_proto prefix
find ../antex_proto -name "*_pb2.py" -exec sed -i '' 's/^from amino import/from antex_proto.amino import/g' {} \;
find ../antex_proto -name "*_pb2.py" -exec sed -i '' 's/^from gogoproto import/from antex_proto.gogoproto import/g' {} \;
find ../antex_proto -name "*_pb2.py" -exec sed -i '' 's/^from cosmos\./from antex_proto.cosmos./g' {} \;
find ../antex_proto -name "*_pb2.py" -exec sed -i '' 's/^from cosmos_proto import/from antex_proto.cosmos_proto import/g' {} \;
find ../antex_proto -name "*_pb2.py" -exec sed -i '' 's/^from cometbft\./from antex_proto.cometbft./g' {} \;
find ../antex_proto -name "*_pb2.py" -exec sed -i '' 's/^from antex\./from antex_proto.antex./g' {} \;
echo "Python proto files copied and import paths fixed"
