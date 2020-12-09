// +build !darwin,!freebsd,!dragonfly,!linux,!openbsd,!solaris,!windows

/*
Copyright 2020 The arhat.dev Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package exechelper

import (
	"io"
	"os/exec"

	"arhat.dev/pkg/wellknownerrors"
)

func startCmdWithTty(
	cmd *exec.Cmd,
) (
	doResize resizeFunc,
	close func(),
	stdin io.WriteCloser,
	stdout io.ReadCloser,
	err error,
) {
	return nil, nil, nil, nil, wellknownerrors.ErrNotSupported
}
