package syncPool

import (
	"github.com/panhongrainbow/goCodePebblez/syncPoolUtil"
)

const stringSliceCapacity = 60

var GlobalStringSlice *syncPoolUtil.StringSlicePool

func init() {
	GlobalStringSlice = syncPoolUtil.NewStringSlicePool(stringSliceCapacity)
}
