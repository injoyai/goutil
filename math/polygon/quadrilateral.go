package polygon

import (
	"errors"
	"math"
)

// Quadrilateral 四边形
type Quadrilateral [4]float64

// Invalid 校验四边形是否有效
func (q Quadrilateral) Invalid() bool {
	return q[0] < q[1]+q[2]+q[3] &&
		q[1] < q[0]+q[2]+q[3] &&
		q[2] < q[0]+q[1]+q[3] &&
		q[3] < q[0]+q[1]+q[2]
}

// Perimeter 四边形周长
func (q Quadrilateral) Perimeter() float64 {
	return q[0] + q[1] + q[2] + q[3]
}

/*
AdjacentAngle 已知q[0],q[1]边角度a1,求邻边q[2],q[3]角度a2

			    q[1]
	          a1------a2
			  	\     \
		    q[0] \     \ q[2]
			      \     \
			       ------
				     q[3]
*/
func (q Quadrilateral) AdjacentAngle(angle float64) (float64, error) {
	ll := q[0]*q[0] + q[1]*q[1] - 2*q[0]*q[1]*math.Cos(angle*math.Pi/180)
	if ll < 0 {
		return 0, errors.New("数据有误")
	}
	l := math.Sqrt(ll)
	a1 := math.Acos((ll + q[1]*q[1] - q[0]*q[0]) / (2 * l * q[1]))
	a2 := math.Acos((ll + q[2]*q[2] - q[3]*q[3]) / (2 * l * q[2]))
	a := (a1 + a2) * 180 / math.Pi
	return a, nil
}

/*
Diagonal 已知q[0],q[1]边角度a1,求对角线q[2],q[3]角度a2

			    q[1]
	          a1------
			  	\     \
		    q[0] \     \ q[2]
			      \     \
			       ------a2
				     q[3]
*/
func (q Quadrilateral) Diagonal(angle float64) (float64, error) {
	ll := q[0]*q[0] + q[1]*q[1] - 2*q[0]*q[1]*math.Cos(angle*math.Pi/180)
	if ll < 0 {
		return 0, errors.New("数据有误")
	}
	a := math.Acos((q[2]*q[2] + q[3]*q[3] - ll) / (2 * q[2] * q[3]))
	a = a * 180 / math.Pi
	return a, nil
}
