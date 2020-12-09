package conf

import (
	"crypto/rand"
	"encoding/hex"
	"strings"

	"arhat.dev/pkg/hashhelper"
)

func generateRandomName(fallback string) string {
	buf := make([]byte, 128)
	n, err := rand.Read(buf)
	if err != nil || n != 128 {
		return hex.EncodeToString(hashhelper.Sha256Sum([]byte(fallback)))
	}

	return hex.EncodeToString(hashhelper.Sha256Sum(buf))
}

func getChartRepoNameChartNameChartVersion(name string) (repoName, chartName, chartVersion string) {
	parts := strings.SplitN(name, "@", 2)
	chartName, chartVersion = parts[0], parts[1]

	nameParts := strings.SplitN(chartName, "/", 2)
	switch len(nameParts) {
	case 1:
		return
	case 2:
		repoName = nameParts[0]
		chartName = nameParts[1]
		return
	}

	return
}

func assembleCommandWithoutEmptyString(prefix []string, args ...string) []string {
	return append(append([]string{}, removeEmptyString(prefix)...), removeEmptyString(args)...)
}

func removeEmptyString(s []string) []string {
	var ret []string
	for _, str := range s {
		if str == "" {
			continue
		}

		ret = append(ret, str)
	}

	return ret
}
