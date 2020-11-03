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

# native
helm-stack:
	sh scripts/build/build.sh $@

# linux
helm-stack.linux.x86:
	sh scripts/build/build.sh $@

helm-stack.linux.amd64:
	sh scripts/build/build.sh $@

helm-stack.linux.armv5:
	sh scripts/build/build.sh $@

helm-stack.linux.armv6:
	sh scripts/build/build.sh $@

helm-stack.linux.armv7:
	sh scripts/build/build.sh $@

helm-stack.linux.arm64:
	sh scripts/build/build.sh $@

helm-stack.linux.mips:
	sh scripts/build/build.sh $@

helm-stack.linux.mipshf:
	sh scripts/build/build.sh $@

helm-stack.linux.mipsle:
	sh scripts/build/build.sh $@

helm-stack.linux.mipslehf:
	sh scripts/build/build.sh $@

helm-stack.linux.mips64:
	sh scripts/build/build.sh $@

helm-stack.linux.mips64hf:
	sh scripts/build/build.sh $@

helm-stack.linux.mips64le:
	sh scripts/build/build.sh $@

helm-stack.linux.mips64lehf:
	sh scripts/build/build.sh $@

helm-stack.linux.ppc64:
	sh scripts/build/build.sh $@

helm-stack.linux.ppc64le:
	sh scripts/build/build.sh $@

helm-stack.linux.s390x:
	sh scripts/build/build.sh $@

helm-stack.linux.riscv64:
	sh scripts/build/build.sh $@

helm-stack.linux.all: \
	helm-stack.linux.x86 \
	helm-stack.linux.amd64 \
	helm-stack.linux.armv5 \
	helm-stack.linux.armv6 \
	helm-stack.linux.armv7 \
	helm-stack.linux.arm64 \
	helm-stack.linux.mips \
	helm-stack.linux.mipshf \
	helm-stack.linux.mipsle \
	helm-stack.linux.mipslehf \
	helm-stack.linux.mips64 \
	helm-stack.linux.mips64hf \
	helm-stack.linux.mips64le \
	helm-stack.linux.mips64lehf \
	helm-stack.linux.ppc64 \
	helm-stack.linux.ppc64le \
	helm-stack.linux.s390x \
	helm-stack.linux.riscv64

helm-stack.darwin.amd64:
	sh scripts/build/build.sh $@

# # currently darwin/arm64 build will fail due to golang link error
# helm-stack.darwin.arm64:
# 	sh scripts/build/build.sh $@

helm-stack.darwin.all: \
	helm-stack.darwin.amd64

helm-stack.windows.x86:
	sh scripts/build/build.sh $@

helm-stack.windows.amd64:
	sh scripts/build/build.sh $@

helm-stack.windows.armv5:
	sh scripts/build/build.sh $@

helm-stack.windows.armv6:
	sh scripts/build/build.sh $@

helm-stack.windows.armv7:
	sh scripts/build/build.sh $@

# # currently no support for windows/arm64
# helm-stack.windows.arm64:
# 	sh scripts/build/build.sh $@

helm-stack.windows.all: \
	helm-stack.windows.x86 \
	helm-stack.windows.amd64 \
	helm-stack.windows.armv5 \
	helm-stack.windows.armv6 \
	helm-stack.windows.armv7

# # android build requires android sdk
# helm-stack.android.amd64:
# 	sh scripts/build/build.sh $@

# helm-stack.android.x86:
# 	sh scripts/build/build.sh $@

# helm-stack.android.armv5:
# 	sh scripts/build/build.sh $@

# helm-stack.android.armv6:
# 	sh scripts/build/build.sh $@

# helm-stack.android.armv7:
# 	sh scripts/build/build.sh $@

# helm-stack.android.arm64:
# 	sh scripts/build/build.sh $@

# helm-stack.android.all: \
# 	helm-stack.android.amd64 \
# 	helm-stack.android.arm64 \
# 	helm-stack.android.x86 \
# 	helm-stack.android.armv7 \
# 	helm-stack.android.armv5 \
# 	helm-stack.android.armv6

helm-stack.freebsd.amd64:
	sh scripts/build/build.sh $@

helm-stack.freebsd.x86:
	sh scripts/build/build.sh $@

helm-stack.freebsd.armv5:
	sh scripts/build/build.sh $@

helm-stack.freebsd.armv6:
	sh scripts/build/build.sh $@

helm-stack.freebsd.armv7:
	sh scripts/build/build.sh $@

helm-stack.freebsd.arm64:
	sh scripts/build/build.sh $@

helm-stack.freebsd.all: \
	helm-stack.freebsd.amd64 \
	helm-stack.freebsd.arm64 \
	helm-stack.freebsd.armv7 \
	helm-stack.freebsd.x86 \
	helm-stack.freebsd.armv5 \
	helm-stack.freebsd.armv6

helm-stack.netbsd.amd64:
	sh scripts/build/build.sh $@

helm-stack.netbsd.x86:
	sh scripts/build/build.sh $@

helm-stack.netbsd.armv5:
	sh scripts/build/build.sh $@

helm-stack.netbsd.armv6:
	sh scripts/build/build.sh $@

helm-stack.netbsd.armv7:
	sh scripts/build/build.sh $@

helm-stack.netbsd.arm64:
	sh scripts/build/build.sh $@

helm-stack.netbsd.all: \
	helm-stack.netbsd.amd64 \
	helm-stack.netbsd.arm64 \
	helm-stack.netbsd.armv7 \
	helm-stack.netbsd.x86 \
	helm-stack.netbsd.armv5 \
	helm-stack.netbsd.armv6

helm-stack.openbsd.amd64:
	sh scripts/build/build.sh $@

helm-stack.openbsd.x86:
	sh scripts/build/build.sh $@

helm-stack.openbsd.armv5:
	sh scripts/build/build.sh $@

helm-stack.openbsd.armv6:
	sh scripts/build/build.sh $@

helm-stack.openbsd.armv7:
	sh scripts/build/build.sh $@

helm-stack.openbsd.arm64:
	sh scripts/build/build.sh $@

helm-stack.openbsd.all: \
	helm-stack.openbsd.amd64 \
	helm-stack.openbsd.arm64 \
	helm-stack.openbsd.armv7 \
	helm-stack.openbsd.x86 \
	helm-stack.openbsd.armv5 \
	helm-stack.openbsd.armv6

helm-stack.solaris.amd64:
	sh scripts/build/build.sh $@

helm-stack.aix.ppc64:
	sh scripts/build/build.sh $@

helm-stack.dragonfly.amd64:
	sh scripts/build/build.sh $@

helm-stack.plan9.amd64:
	sh scripts/build/build.sh $@

helm-stack.plan9.x86:
	sh scripts/build/build.sh $@

helm-stack.plan9.armv5:
	sh scripts/build/build.sh $@

helm-stack.plan9.armv6:
	sh scripts/build/build.sh $@

helm-stack.plan9.armv7:
	sh scripts/build/build.sh $@

helm-stack.plan9.all: \
	helm-stack.plan9.amd64 \
	helm-stack.plan9.armv7 \
	helm-stack.plan9.x86 \
	helm-stack.plan9.armv5 \
	helm-stack.plan9.armv6
