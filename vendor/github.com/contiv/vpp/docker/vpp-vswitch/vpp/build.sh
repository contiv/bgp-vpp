#!/bin/bash
# Copyright (c) 2017 Cisco and/or its affiliates.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at:
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# fail in case of error
set -e

DOCKERFILETAG=""

BUILDARCH=`uname -m`

if [ ${BUILDARCH} = "aarch64" ] ; then
  DOCKERFILETAG=".arm64"
  BUILDARCH="arm64"
fi

if [ ${BUILDARCH} = "x86_64" ] ; then
  BUILDARCH="amd64"
fi

DOCKERFILE=Dockerfile${DOCKERFILETAG}

# determine extra vpp version based args
VPP_COMMIT_ARGS=""
VPP_COMMIT_VERSION="latest"
if [ -n "${VPP_COMMIT_ID}" ]
then
  VPP_COMMIT_ARGS="--build-arg VPP_COMMIT_ID=${VPP_COMMIT_ID}"
  VPP_COMMIT_VERSION="${VPP_COMMIT_ID}"
fi

# execute the build
docker build -t contivvpp/vpp-${BUILDARCH}:${VPP_COMMIT_VERSION} \
	-f ${DOCKERFILE} \
	--build-arg VPP_REPO_URL=${VPP_REPO_URL} \
	--build-arg VPP_BRANCH_NAME=${VPP_BRANCH_NAME} \
	--build-arg SKIP_DEBUG_BUILD=${SKIP_DEBUG_BUILD} \
	${VPP_COMMIT_ARGS} \
	${DOCKER_BUILD_ARGS} .

if [ ${BUILDARCH} = "amd64" ] ; then
   docker tag contivvpp/vpp-${BUILDARCH}:${VPP_COMMIT_VERSION} contivvpp/vpp:${VPP_COMMIT_VERSION}
fi
