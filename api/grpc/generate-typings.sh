#!/bin/bash

# ./api/grpc/generate-typings.sh

# Path to this plugin
PROJECT_DIR=$(pwd)
GOPATH=$(go env GOPATH)
GO_BIN=$(go env GOBIN)
PROTOC_PATH=$(which protoc)               # brew install protobuf && brew link --overwrite protobuf
PROTO_DIR="${PROJECT_DIR}/api/grpc/proto" # Path to your proto files
# PB_OUTPUT_DIR="${PROJECT_DIR}/generated/grpc/pb"      # Path to generated PB files
PB_OUTPUT_DIR="${PROJECT_DIR}/generated/grpc/pb"     # Path to generated PB files
OPENAPI_OUTPUT_DIR="${PROJECT_DIR}/doc/openapi"      # Path to generated PB files
GORM_OUTPUT_DIR="${PROJECT_DIR}/generated/grpc/gorm" # Path to generated PB files

PROTOS=(
    $(
        cd ${PROTO_DIR}
        ls *.proto | sort
    )
)

if [ "$#" -gt "0" ]; then
    PROTOS=("$@")
fi

# PROTOS=(
#     hello.proto
#     health.proto
#     common.proto
#     user.proto
#     auth.proto
#     team.proto
# )

rm -rf ${PB_OUTPUT_DIR} ${OPENAPI_OUTPUT_DIR} ${GORM_OUTPUT_DIR}
mkdir -p ${PB_OUTPUT_DIR} ${OPENAPI_OUTPUT_DIR} ${GORM_OUTPUT_DIR}

echo "Protos:"
printf "%s\n" "${PROTOS[@]}"
echo "Generating proto typings..."

for PROTO in ${PROTOS[@]}; do
    # OUTPUT_DIR="${PB_OUTPUT_DIR}/${PROTO%.*}"

    echo -e "\n>>> Processing protos [${PROTO_DIR}/${PROTO}]...\n"
    # mkdir -p ${OUTPUT_DIR}

    ${PROTOC_PATH} \
        --plugin=protoc-gen-go="${GO_BIN}/protoc-gen-go" \
        --plugin=protoc-gen-go-grpc="${GO_BIN}/protoc-gen-go-grpc" \
        --plugin=protoc-gen-openapiv2="${GO_BIN}/protoc-gen-openapiv2" \
        --proto_path="${PROTO_DIR}" \
        --go_out="${PB_OUTPUT_DIR}" \
        --go-grpc_out="${PB_OUTPUT_DIR}" \
        --openapiv2_out="${OPENAPI_OUTPUT_DIR}" \
        "${PROTO_DIR}/${PROTO}"
    # --gorm_out="${GORM_OUTPUT_DIR}" \
    # --go_opt=paths=source_relative
    # --go-grpc_opt=paths=source_relative

    echo "Generated typings for ${PROTO} in ${PB_OUTPUT_DIR}"
    echo "Generated OpenAPI specs for ${PROTO} in ${OPENAPI_OUTPUT_DIR}"

    # Define the replacement patterns
    patterns=(
        '"\.\/common"' '"github.com\/cclhsu\/gin-grpc-gorm\/generated\/grpc\/pb\/common"'
        '"\.\/user"' '"github.com\/cclhsu\/gin-grpc-gorm\/generated\/grpc\/pb\/user"'
        '"\.\/team"' '"github.com\/cclhsu\/gin-grpc-gorm\/generated\/grpc\/pb\/team"'
    )

    # Iterate over patterns and perform replacements
    for ((i = 0; i < ${#patterns[@]}; i += 2)); do
        # find "${PB_OUTPUT_DIR}" -type f -name '*.pb.go' -exec sed -i.bak "s/${patterns[i]}/${patterns[i + 1]}/g" {} +
        find "${PB_OUTPUT_DIR}" -type f -name '*.pb.go' -exec sed -i "s/${patterns[i]}/${patterns[i + 1]}/g" {} +
    done

done

echo "gRPC typings generated successfully!"
