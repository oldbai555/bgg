/**
 * @Author: zjj
 * @Date: 2024/9/13
 * @Desc:
**/

package marshal

import (
	"github.com/oldbai555/lbtool/pkg/jsonpb"
	"google.golang.org/protobuf/proto"
)

func JsonPbMarshal(obj proto.Message) ([]byte, error) {
	buf, err := jsonpb.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func JsonPbUnmarshal(buf []byte, obj proto.Message) error {
	err := jsonpb.Unmarshal(buf, obj)
	if err != nil {
		return err
	}
	return nil
}
