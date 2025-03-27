// SELECT p.id AS platform_id, s.id AS streamer_id
// FROM platform p
// INNER JOIN streamer s ON p.id = s.platform_id
// WHERE p.name = ?
// AND s.nickname = ?
// LIMIT 1;

package repository

import (
	mntreamerModel "mntreamer/shared/model"
)

type IRepository interface {
	Create(streamer *mntreamerModel.Streamer) (*mntreamerModel.Streamer, error)
	Save(streamer *mntreamerModel.Streamer) (*mntreamerModel.Streamer, error)
	FindByPlatformIdAndStreamerId(platformId uint16, streamerId uint32) (*mntreamerModel.Streamer, error)
	FindByPlatformIdAndChannelName(platformId uint16, channelName string) (*mntreamerModel.Streamer, error)
}
