/**
 * @Author: zjj
 * @Date: 2024/9/13
 * @Desc:
**/

package marshal

import "google.golang.org/protobuf/proto"

func PbMarshal(obj proto.Message) ([]byte, error) {
	buf, err := proto.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func PbUnmarshal(buf []byte, obj proto.Message) error {
	err := proto.Unmarshal(buf, obj)
	if err != nil {
		return err
	}
	return nil
}
