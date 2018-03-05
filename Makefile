TARGET = anycast-operator
VERSION = latest

REPO = as65342
TAG = ${REPO}/k8s-${TARGET}:${VERSION}

BUILD_DIR = ./build
GOSRC_DIR = ${BUILD_DIR}/src/github.com/r3boot/anycast-operator

all: ${BUILD_DIR}/${TARGET}

${BUILD_DIR}:
	mkdir -p $@

${GOSRC_DIR}:
	mkdir -p $@

${BUILD_DIR}/${TARGET}: ${GOSRC_DIR}
	cp -Rvp cmd pkg ${GOSRC_DIR}/
	docker run --rm -it \
		-v $(shell pwd)/build:/build \
		-v $(shell pwd)/scripts/build.sh:/build.sh \
		as65342/alpine-builder:3.7 /build.sh

container:
	docker build -t ${TAG} .
	docker push ${TAG}

clean:
	[[ -d "${BUILD_DIR}" ]] && rm -rf "${BUILD_DIR}"