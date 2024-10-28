/**
 * @Author: zjj
 * @Date: 2024/10/28
 * @Desc:
**/

package server

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateOrderID 生成包含订单创建时间和随机数的订单号
func GenerateOrderID() string {
	// 获取当前时间
	now := time.Now()
	// 格式化时间为年月日时分秒，例如：20230915123456
	timestamp := now.Format("20060102150405")

	// 生成一个随机数
	rand.New(rand.NewSource(20241028)).Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(10000) // 生成一个0到9999之间的随机数

	// 将随机数格式化为四位数，不足四位的前面补0
	formattedRandomNumber := fmt.Sprintf("%04d", randomNumber)

	// 组合订单号
	orderID := fmt.Sprintf("%s%s", timestamp, formattedRandomNumber)

	return orderID
}
