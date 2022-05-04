package car

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/shared/id"
)

type Manager struct {
}

//认证
func (receiver *Manager) Verify(c context.Context, aid id.CarID, location *rentalpb.Location) error {
	return nil
}

//检查是否被锁
func (receiver *Manager) Unlock(c context.Context, cID id.CarID) error {
	return nil
}
