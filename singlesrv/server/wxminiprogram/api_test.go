/**
 * @Author: zjj
 * @Date: 2024/10/22
 * @Desc:
**/

package wxminiprogram

import (
	"github.com/oldbai555/bgg/singlesrv/client"
	"testing"
)

func TestCode2Session(t *testing.T) {
	Code2Session(&client.JsCodeToSessionReq{})
}
