// Package types: the file 'http_code_range.go' was copied from:
// https://github.com/traefik/traefik/blob/master/pkg/types/http_code_range.go
// under the MIT license. It was copied, since traefik plugins do not support package dependencies
// but not vendored to minimize the size of the plugin
package types

import (
	"github.com/m-riedel/traefik-plugin-redirect-on-status/pkg/util"
	"strconv"
	"strings"
)

// HTTPCodeRanges holds HTTP code ranges.
type HTTPCodeRanges [][2]int

// NewHTTPCodeRanges creates HTTPCodeRanges from a given []string.
// Break out the http status code ranges into a low int and high int
// for ease of use at runtime.
func NewHTTPCodeRanges(strBlocks []string) (HTTPCodeRanges, error) {
	var blocks HTTPCodeRanges
	for _, block := range strBlocks {
		codes := strings.Split(block, "-")
		// if only a single HTTP code was configured, assume the best and create the correct configuration on the user's behalf
		if len(codes) == 1 {
			codes = append(codes, codes[0])
		}
		lowCode, err := strconv.Atoi(codes[0])
		if err != nil {
			return nil, err
		}
		highCode, err := strconv.Atoi(codes[1])
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, [2]int{lowCode, highCode})
	}
	return blocks, nil
}

// Contains tests whether the passed status code is within one of its HTTP code ranges.
func (h HTTPCodeRanges) Contains(statusCode int) bool {
	return util.ContainsFunc(h, func(block [2]int) bool {
		return statusCode >= block[0] && statusCode <= block[1]
	})
}
