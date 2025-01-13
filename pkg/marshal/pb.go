/**
 * @Author: zjj
 * @Date: 2024/9/13
 * @Desc:
**/

package marshal

import (
	"github.com/oldbai555/lbtool/pkg/lberr"
	"google.golang.org/protobuf/proto"
)

func PbMarshal(obj proto.Message) ([]byte, error) {
	buf, err := proto.Marshal(obj)
	if err != nil {
		return nil, lberr.Wrap(err)
	}
	return buf, nil
}

func PbUnmarshal(buf []byte, obj proto.Message) error {
	err := proto.Unmarshal(buf, obj)
	if err != nil {
		return lberr.Wrap(err)
	}
	return nil
}
