function genProto {
    # receive params
    DOMAIN=$1
    SKIP_GATEWAY=$2

    # prepare directories
    PROTO_PATH=./${DOMAIN}/api
    GO_OUT_PATH=./${DOMAIN}/api/gen/v1
    WX_MINI_PROGIRAM_PATH=$HOME/WeChatProjects/coolcar/miniprogram
    mkdir -p $GO_OUT_PATH

    # generate proto buffer codes
    protoc -I=$PROTO_PATH --go_out=plugins=grpc,paths=source_relative:$GO_OUT_PATH ${DOMAIN}.proto

    if [ $SKIP_GATEWAY ]; then
        return
    fi

    protoc -I=$PROTO_PATH --grpc-gateway_out=paths=source_relative,grpc_api_configuration=$PROTO_PATH/${DOMAIN}.yaml:$GO_OUT_PATH ${DOMAIN}.proto

    # generate grpc-gateway codes for wx miniprogram
    PBTS_BIN_DIR=$WX_MINI_PROGIRAM_PATH/node_modules/.bin
    PBTS_OUT_DIR=$WX_MINI_PROGIRAM_PATH/service/proto_gen/${DOMAIN}
    mkdir -p $PBTS_OUT_DIR
    $PBTS_BIN_DIR/pbjs -t static -w es6 $PROTO_PATH/${DOMAIN}.proto --no-create --no-encode --no-decode --no-verify --no-delimited --force-number -o $PBTS_OUT_DIR/${DOMAIN}_pb.js
    $PBTS_BIN_DIR/pbts -o $PBTS_OUT_DIR/${DOMAIN}_pb.d.ts $PBTS_OUT_DIR/${DOMAIN}_pb.js
}

genProto auth
genProto rental