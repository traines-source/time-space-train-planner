SRC_DIR=$1
DST_DIR=./internal/stost/
PATH=$PATH:/go/bin
cp $SRC_DIR/src/wire/wire.proto $DST_DIR

protoc --version
protoc -I=$DST_DIR --go_opt=paths=source_relative --go_opt=Mwire.proto=traines.eu/time-space-train-planner/providers/stost --go_out=$DST_DIR $DST_DIR/wire.proto 