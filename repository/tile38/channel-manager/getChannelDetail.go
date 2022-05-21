package channelmanager

import (
	"github.com/xjem/t38c"
)

func (cm *chanManager) GetChannelDetail(pattern string) ([]t38c.Chan, error) {

	var result []t38c.Chan
	detail, err := cm.tile.Channels.Chans(pattern)
	if err != nil {
		return result, err
	}
	result = append(result, detail...)
	return result, nil
}
