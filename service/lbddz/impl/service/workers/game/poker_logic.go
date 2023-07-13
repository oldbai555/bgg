package game

import (
	"github.com/oldbai555/bgg/client/lbddz"
	"sort"
)

var pokerLogic = &PokerLogic{}

// PokerLogic 其实现在是空的，但是万一有用呢？所以就没有写成静态方法了
type PokerLogic struct {
}

// CalcPokerType 计算牌型
func (p *PokerLogic) CalcPokerType(val []uint32) lbddz.CardType {
	var cards []int
	for i := range val {
		cards = append(cards, int(val[i]))
	}
	// 转换为点数
	points := p.CardsToPoints(cards)

	length := len(points)

	if length == 1 { // 一张牌，当然是单牌
		return lbddz.CardType_CardTypeSingleCard
	} else if length == 2 { // 两张牌，王炸或一对
		if points[0] == 16 && points[1] == 17 { // 王炸
			return lbddz.CardType_CardTypeKingBombCard
		}
		if points[0] == points[1] { // 对子
			return lbddz.CardType_CardTypeDoubleCard
		}
	} else if length == 3 { // 三张，只检查三不带
		if points[0] == points[1] && points[1] == points[2] {
			return lbddz.CardType_CardTypeThreeCard
		}
	} else if length == 4 { // 四张，炸弹或三带一
		maxSameNum := p.CalcMaxSameNum(points)
		if maxSameNum == 4 { // 四张相同，炸弹
			return lbddz.CardType_CardTypeBombCard
		}

		if maxSameNum == 3 { // 三张点数相同的，3带1
			return lbddz.CardType_CardTypeThreeOneCard
		}
	} else if length >= 5 && p.IsStraight(points) && points[length-1] < 15 { // 大于等于5张，是点数连续，且最大点数不超过2， 则是顺子
		return lbddz.CardType_CardTypeStraight
	} else if length == 5 { // 5张，只需检查3带2
		// 最多有3张相等的，且只有两种点数，则是3带2
		if p.CalcMaxSameNum(points) == 3 && p.CalcDiffPoint(points) == 2 {
			return lbddz.CardType_CardTypeThreeTwoCard
		}
	} else { // 大于6的情况，再分别判断
		maxSameNum := p.CalcMaxSameNum(points)
		diffPointNum := p.CalcDiffPoint(points)

		// length能被3整除， 最大相同数量是3， 不同点数是length/3, 且最大与最小点数相差 length/3 - 1， 则是连续三张
		if length%3 == 0 && maxSameNum == 3 && diffPointNum == length/3 && (points[length-1]-points[0] == length/3-1) && points[length-1] < 15 {
			return lbddz.CardType_CardTypeAircraft
		}

		// 与上面连续三张判断类似，连对
		if length%2 == 0 && maxSameNum == 2 && diffPointNum == length/2 && (points[length-1]-points[0] == length/2-1) && points[length-1] < 15 {
			return lbddz.CardType_CardTypeConnectCard
		}

		// 飞机三带一
		if length%4 == 0 {
			// 连续3张的数量占了length/4 以上
			threePoints := p.GetSameNumMaxStraightPoints(points, 3)
			if len(threePoints) >= length/4 && threePoints[length/4-1] < 15 {
				return lbddz.CardType_CardTypeAircraftCard
			}
		}

		// 飞机三带二
		if length%5 == 0 {
			threePoints := p.GetSameNumPoints(points, 3)
			// 三带二里面，不会出现单牌的情况
			onePoints := p.GetSameNumPoints(points, 1)
			if len(onePoints) == 0 && len(threePoints) == length/5 && p.IsStraight(threePoints) && threePoints[len(threePoints)-1] < 15 {
				return lbddz.CardType_CardTypeAircraftWing
			}
		}

		// 四带二
		if length == 6 {
			if maxSameNum == 4 {
				return lbddz.CardType_CardTypeBombTwoCard
			}
		}

		// 四带两对
		if length == 8 {
			if maxSameNum == 4 {
				// 必须没有一张和三张的出现
				onePoints := p.GetSameNumPoints(points, 1)
				threePoints := p.GetSameNumPoints(points, 3)

				// TODO： 33334444 这样的到底算连续三带一还是四带两对？
				if len(onePoints) == 0 && len(threePoints) == 0 {
					return lbddz.CardType_CardTypeBombFourCard
				}
			}
		}
		// 连续四带二
		if length%6 == 0 {
			fourPoints := p.GetSameNumMaxStraightPoints(points, 4)
			if len(fourPoints) >= length/6 && fourPoints[length/6-1] < 15 {
				return lbddz.CardType_CardTypeBombTwoStraightCard
			}
		}
		// 连续四带两对
		if length%8 == 0 {
			fourPoints := p.GetSameNumMaxStraightPoints(points, 4)

			// 其他的都必须是成对，因此检查是否没有单张和三张的出现即可
			onePoints := p.GetSameNumPoints(points, 1)
			threePoints := p.GetSameNumPoints(points, 3)

			if len(fourPoints) >= length/8 && len(onePoints) == 0 && len(threePoints) == 0 && fourPoints[length/8-1] < 15 {
				return lbddz.CardType_CardTypeBombFourStraightCard
			}
		}
	}

	// 没有这个牌型的
	return lbddz.CardType_CardTypeErrorCards
}

// GetSameNumPoints 取出所有点数数量等于num的点数
// 例如，现在牌中有3个3，3个4，2个5，1个6， 取出数量等于3的点数，则返回[3, 4]，取出数量等于2的点数，则返回[5]，取出数量等于1的点数，则返回[6]，其他都返回空数组
func (p *PokerLogic) GetSameNumPoints(points []int, num int) []int {
	length := len(points)
	newPoints := make([]int, length)
	pointIndex := 0

	nowNum := 1

	for i := 1; i < length; i++ {
		if points[i] == points[i-1] { // 与前一张相同
			nowNum++
		} else { // 与前一张不同，若前一张出现num次，加入数组
			if nowNum == num {
				newPoints[pointIndex] = points[i-1]
				pointIndex++
			}
			nowNum = 1
		}
	}

	if nowNum == num {
		newPoints[pointIndex] = points[length-1]
		pointIndex++
	}

	return newPoints[0:pointIndex]
}

// GetGeNumPoints 取出所有点数数量大于等于num的点数
// 例如，现在牌中有3个3，3个4，2个5，1个6， 取出数量大于等于3的点数，则返回[3, 4]，取出数量大于等于2的点数，则返回[3, 4, 5]，取出数量大于等于1的点数，则返回[3, 4, 5, 6]
func (p *PokerLogic) GetGeNumPoints(points []int, num int) []int {
	length := len(points)
	newPoints := make([]int, length)
	pointIndex := 0

	nowNum := 1

	for i := 1; i < length; i++ {
		if points[i] == points[i-1] { // 与前一张相同
			nowNum++
		} else { // 与前一张不同，若前一张出现大于等于num次，加入数组
			if nowNum >= num {
				newPoints[pointIndex] = points[i-1]
				pointIndex++
			}
			nowNum = 1
		}
	}

	if nowNum >= num {
		newPoints[pointIndex] = points[length-1]
		pointIndex++
	}

	return newPoints[0:pointIndex]
}

// GetSameNumMaxStraightPoints 从所有点数数量大于等于num的列表中，取出最长连续递增子列表
func (p *PokerLogic) GetSameNumMaxStraightPoints(points []int, num int) []int {
	geNumPoints := p.GetGeNumPoints(points, num)

	// 没有，直接返回
	length := len(geNumPoints)
	if length == 0 {
		return geNumPoints
	}

	maxStartPoint := geNumPoints[0]
	maxNum := 1

	nowStartPoint := geNumPoints[0]
	nowNum := 1

	for i := 1; i < length; i++ {
		if geNumPoints[i] == geNumPoints[i-1]+1 { // 比上一张多1
			nowNum++
		} else { // 重新开始计算
			if nowNum > maxNum {
				maxNum = nowNum
				maxStartPoint = nowStartPoint
			}
			nowNum = 1
			nowStartPoint = geNumPoints[i]
		}
	}

	if nowNum > maxNum {
		maxNum = nowNum
		maxStartPoint = nowStartPoint
	}

	newPoints := make([]int, maxNum)
	for i := 0; i < maxNum; i++ {
		newPoints[i] = maxStartPoint + i
	}

	return newPoints
}

// IsStraight 是否是顺子
func (p *PokerLogic) IsStraight(points []int) bool {
	length := len(points)
	for i := 1; i < length; i++ {
		if points[i] != points[i-1]+1 { // 与前一张相同
			return false
		}
	}

	return true
}

// CalcDiffPoint 有多少种不同的点数
func (p *PokerLogic) CalcDiffPoint(points []int) int {
	diffNum := 1

	length := len(points)
	for i := 1; i < length; i++ {
		if points[i] != points[i-1] { // 与前一张不同，则出现了新的点数
			diffNum++
		}
	}

	return diffNum
}

// CalcMaxSameNum 最多有几张点数相等的牌
func (p *PokerLogic) CalcMaxSameNum(points []int) int {
	length := len(points)
	nowNum := 1
	maxNum := 1

	for i := 1; i < length; i++ {
		if points[i] == points[i-1] { // 与前一张相同
			nowNum++
		} else { // 与前一张不同，重新开始计数
			if nowNum > maxNum {
				maxNum = nowNum
			}
			nowNum = 1
		}
	}

	if nowNum > maxNum {
		maxNum = nowNum
	}

	return maxNum
}

// CardsToPoints 牌id转点数
func (p *PokerLogic) CardsToPoints(cards []int) []int {
	length := len(cards)
	points := make([]int, length)

	var point int
	for i := 0; i < length; i++ {
		value := cards[i]
		if value < 53 { // id 1-4 对应的是3， 5-8对应4， 依此类推 45 - 48对应14(A)， 49-52对应15(2)
			if value%4 == 0 {
				point = value/4 + 2
			} else {
				point = value/4 + 3
			}
		} else { // 小王和大王
			point = value/4 + 2 + value%4
		}

		points[i] = point
	}

	// 按点数升序排序
	sort.Ints(points)

	return points
}

// CalcPokerHeader 计算头牌
func (p *PokerLogic) CalcPokerHeader(val []uint32, cardType lbddz.CardType) int {
	var cards []int
	for i := range val {
		cards = append(cards, int(val[i]))
	}
	points := p.CardsToPoints(cards)

	switch cardType {
	case lbddz.CardType_CardTypeSingleCard, lbddz.CardType_CardTypeDoubleCard, lbddz.CardType_CardTypeThreeCard, lbddz.CardType_CardTypeStraight, lbddz.CardType_CardTypeConnectCard, lbddz.CardType_CardTypeAircraft, lbddz.CardType_CardTypeBombCard:
		return points[0]
	case lbddz.CardType_CardTypeThreeOneCard, lbddz.CardType_CardTypeThreeTwoCard, lbddz.CardType_CardTypeBombTwoCard:
		return points[2]
	case lbddz.CardType_CardTypeAircraftCard: // 连续三带一
		threePoints := p.GetSameNumMaxStraightPoints(points, 3)
		return threePoints[0]
	case lbddz.CardType_CardTypeAircraftWing: // 连续三带二
		return p.FirstPoint(points, 3)
	case lbddz.CardType_CardTypeBombFourCard: // 四带两对
		fourPoints := p.GetSameNumPoints(points, 4)
		return fourPoints[len(fourPoints)-1]
	case lbddz.CardType_CardTypeBombTwoStraightCard, lbddz.CardType_CardTypeBombFourStraightCard:
		fourPoints := p.GetSameNumMaxStraightPoints(points, 4)
		return fourPoints[0]
	}

	return 0
}

// FirstPoint 获得首个出现num次的点数
func (p *PokerLogic) FirstPoint(points []int, num int) int {
	nowNum := 1
	length := len(points)

	for i := 1; i < length; i++ {
		if points[i] == points[i-1] { //与上一张相同，数量加1
			nowNum++
		} else { //重新开始计算
			if nowNum == num {
				return points[i-1]
			}
			nowNum = 1
		}
	}

	if nowNum == num {
		return points[length-1]
	}

	return 0
}

// CanOut 是否可以出牌
func (p *PokerLogic) CanOut(newCardSet *lbddz.CardSet, nowCardSet *lbddz.CardSet) bool {
	// 当前是第一次出牌，牌型正确即可
	if lbddz.CardType(nowCardSet.CardType) == lbddz.CardType_CardTypeNoCards && lbddz.CardType(newCardSet.CardType) != lbddz.CardType_CardTypeErrorCards {
		return true
	}

	// 王炸，天下第一
	if lbddz.CardType(newCardSet.CardType) == lbddz.CardType_CardTypeKingBombCard {
		return true
	}

	// 炸弹，检查前面是不是也是炸弹
	if lbddz.CardType(newCardSet.CardType) == lbddz.CardType_CardTypeBombCard {
		if lbddz.CardType(nowCardSet.CardType) == lbddz.CardType_CardTypeBombCard {
			return newCardSet.Header > nowCardSet.Header
		} else {
			return true
		}
	}

	// 同类型，张数相同，头牌更大
	if newCardSet.CardType == nowCardSet.CardType && len(newCardSet.Cards) == len(nowCardSet.Cards) && newCardSet.Header > nowCardSet.Header {
		return true
	}

	return false
}
