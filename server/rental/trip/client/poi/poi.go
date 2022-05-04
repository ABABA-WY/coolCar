package poi

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"google.golang.org/protobuf/proto"
	"hash/fnv"
)

type Manager struct {
}

var poi = []string{
	"北京",
	"上海",
	"杭州",
	"武汉",
}

func (*Manager) Resolve(c context.Context, location *rentalpb.Location) (string, error) {
	b, err := proto.Marshal(location)
	if err != nil {
		return "", err
	}
	//计算哈希值
	h := fnv.New32()
	_, _ = h.Write(b)
	return poi[int(h.Sum32())%len(poi)], nil
}
