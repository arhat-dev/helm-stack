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

#
# linux
#
package.helm-stack.deb.amd64:
	sh scripts/package/package.sh $@

package.helm-stack.deb.armv6:
	sh scripts/package/package.sh $@

package.helm-stack.deb.armv7:
	sh scripts/package/package.sh $@

package.helm-stack.deb.arm64:
	sh scripts/package/package.sh $@

package.helm-stack.deb.all: \
	package.helm-stack.deb.amd64 \
	package.helm-stack.deb.armv6 \
	package.helm-stack.deb.armv7 \
	package.helm-stack.deb.arm64

package.helm-stack.rpm.amd64:
	sh scripts/package/package.sh $@

package.helm-stack.rpm.armv7:
	sh scripts/package/package.sh $@

package.helm-stack.rpm.arm64:
	sh scripts/package/package.sh $@

package.helm-stack.rpm.all: \
	package.helm-stack.rpm.amd64 \
	package.helm-stack.rpm.armv7 \
	package.helm-stack.rpm.arm64

package.helm-stack.linux.all: \
	package.helm-stack.deb.all \
	package.helm-stack.rpm.all

#
# windows
#

package.helm-stack.msi.amd64:
	sh scripts/package/package.sh $@

package.helm-stack.msi.arm64:
	sh scripts/package/package.sh $@

package.helm-stack.msi.all: \
	package.helm-stack.msi.amd64 \
	package.helm-stack.msi.arm64

package.helm-stack.windows.all: \
	package.helm-stack.msi.all

#
# darwin
#

package.helm-stack.pkg.amd64:
	sh scripts/package/package.sh $@

package.helm-stack.pkg.arm64:
	sh scripts/package/package.sh $@

package.helm-stack.pkg.all: \
	package.helm-stack.pkg.amd64 \
	package.helm-stack.pkg.arm64

package.helm-stack.darwin.all: \
	package.helm-stack.pkg.all
