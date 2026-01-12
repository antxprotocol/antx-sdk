#!/bin/bash

# Copy generated Python proto to python package
mkdir -p ../antx_proto
cp -r ../../proto/gen/python/* ../antx_proto/
# Create __init__.py files
find ../antx_proto -type d -exec touch {}/__init__.py \;
# Fix import paths to use antx_proto prefix
find ../antx_proto -name "*_pb2.py" -exec sed -i '' 's/^from amino import/from antx_proto.amino import/g' {} \;
find ../antx_proto -name "*_pb2.py" -exec sed -i '' 's/^from gogoproto import/from antx_proto.gogoproto import/g' {} \;
find ../antx_proto -name "*_pb2.py" -exec sed -i '' 's/^from cosmos\./from antx_proto.cosmos./g' {} \;
find ../antx_proto -name "*_pb2.py" -exec sed -i '' 's/^from cosmos_proto import/from antx_proto.cosmos_proto import/g' {} \;
find ../antx_proto -name "*_pb2.py" -exec sed -i '' 's/^from cometbft\./from antx_proto.cometbft./g' {} \;
find ../antx_proto -name "*_pb2.py" -exec sed -i '' 's/^from antx\./from antx_proto.antx./g' {} \;
echo "Python proto files copied and import paths fixed"
