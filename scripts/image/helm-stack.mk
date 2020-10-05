# Copyright 2020 The arhat.dev Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# build
image.build.helm-stack.linux.x86:
	sh scripts/image/build.sh $@

image.build.helm-stack.linux.amd64:
	sh scripts/image/build.sh $@

image.build.helm-stack.linux.armv5:
	sh scripts/image/build.sh $@

image.build.helm-stack.linux.armv6:
	sh scripts/image/build.sh $@

image.build.helm-stack.linux.armv7:
	sh scripts/image/build.sh $@

image.build.helm-stack.linux.arm64:
	sh scripts/image/build.sh $@

image.build.helm-stack.linux.ppc64le:
	sh scripts/image/build.sh $@

image.build.helm-stack.linux.mips64le:
	sh scripts/image/build.sh $@

image.build.helm-stack.linux.s390x:
	sh scripts/image/build.sh $@

image.build.helm-stack.linux.all: \
	image.build.helm-stack.linux.amd64 \
	image.build.helm-stack.linux.arm64 \
	image.build.helm-stack.linux.armv7 \
	image.build.helm-stack.linux.armv6 \
	image.build.helm-stack.linux.armv5 \
	image.build.helm-stack.linux.x86 \
	image.build.helm-stack.linux.s390x \
	image.build.helm-stack.linux.ppc64le \
	image.build.helm-stack.linux.mips64le

image.build.helm-stack.windows.amd64:
	sh scripts/image/build.sh $@

image.build.helm-stack.windows.armv7:
	sh scripts/image/build.sh $@

image.build.helm-stack.windows.all: \
	image.build.helm-stack.windows.amd64 \
	image.build.helm-stack.windows.armv7

# push
image.push.helm-stack.linux.x86:
	sh scripts/image/push.sh $@

image.push.helm-stack.linux.amd64:
	sh scripts/image/push.sh $@

image.push.helm-stack.linux.armv5:
	sh scripts/image/push.sh $@

image.push.helm-stack.linux.armv6:
	sh scripts/image/push.sh $@

image.push.helm-stack.linux.armv7:
	sh scripts/image/push.sh $@

image.push.helm-stack.linux.arm64:
	sh scripts/image/push.sh $@

image.push.helm-stack.linux.ppc64le:
	sh scripts/image/push.sh $@

image.push.helm-stack.linux.mips64le:
	sh scripts/image/push.sh $@

image.push.helm-stack.linux.s390x:
	sh scripts/image/push.sh $@

image.push.helm-stack.linux.all: \
	image.push.helm-stack.linux.amd64 \
	image.push.helm-stack.linux.arm64 \
	image.push.helm-stack.linux.armv7 \
	image.push.helm-stack.linux.armv6 \
	image.push.helm-stack.linux.armv5 \
	image.push.helm-stack.linux.x86 \
	image.push.helm-stack.linux.s390x \
	image.push.helm-stack.linux.ppc64le \
	image.push.helm-stack.linux.mips64le

image.push.helm-stack.windows.amd64:
	sh scripts/image/push.sh $@

image.push.helm-stack.windows.armv7:
	sh scripts/image/push.sh $@

image.push.helm-stack.windows.all: \
	image.push.helm-stack.windows.amd64 \
	image.push.helm-stack.windows.armv7
