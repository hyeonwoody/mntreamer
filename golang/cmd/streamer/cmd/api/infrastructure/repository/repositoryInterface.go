// SELECT p.id AS platform_id, s.id AS streamer_id
// FROM platform p
// INNER JOIN streamer s ON p.id = s.platform_id
// WHERE p.name = ?
// AND s.nickname = ?
// LIMIT 1;

package repository

import (
	model "mntreamer/streamer/cmd/model"
)

type IRepository interface {
	Save(*model.Streamer) (*model.Streamer, error)
}
