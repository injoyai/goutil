package polygon

import (
	"github.com/injoyai/conv"
	"github.com/injoyai/logs"
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

// MaxAngle 能形成的最大角,最大180度
func (q Quadrilateral) MaxAngle(n ...int) float64 {
	i := conv.Default(0, n...)
	a := q.get(i + 0)
	b := q.get(i + 1)
	c := q.get(i + 2)
	d := q.get(i + 3)
	l := c + d
	if l > a+b {
		return 180
	}
	result := math.Acos((a*a + b*b - l*l) / (2 * a * b))
	result = result * 180 / math.Pi
	return result
}

func (q Quadrilateral) get(n int) float64 {
	return q[n%4]
}

/*
LeftAngle 已知q[0],q[1]边角度a1,求邻边q[2],q[3]角度a2

		    q[1]
		  a1------a2
		  	\     \
	    q[0] \     \ q[2]
		      \     \
		       ------
			     q[3]
*/
func (q Quadrilateral) LeftAngle(angle float64, offset ...int) (float64, float64) {

	i := conv.Default(0, offset...)

	if angle > q.MaxAngle(i) {
		angle = q.MaxAngle(i)
	}

	a := q.get(i + 0)
	b := q.get(i + 1)
	c := q.get(i + 2)
	d := q.get(i + 3)

	ll := a*a + b*b - 2*a*b*math.Cos(angle*math.Pi/180)
	l := math.Sqrt(ll)
	if l == 0 {
		r := math.Acos(0) * 2 * 180 / math.Pi
		return r, r
	}

	a1 := math.Acos((ll + b*b - a*a) / (2 * l * b))
	a1 = conv.Select(math.IsNaN(a1), 0, a1)
	a2 := math.Acos((ll + c*c - d*d) / (2 * l * c))
	a2 = conv.Select(math.IsNaN(a2), 0, a2)

	//有2个解,突出四边形,凹陷四边形
	r1 := math.Abs(a1+a2) * 180 / math.Pi
	r2 := math.Abs(a1-a2) * 180 / math.Pi

	return r1, r2
}

/*
RightAngle 已知q[0],q[1]边角度a1,求邻边q[3],q[0]角度a2

		    q[1]
		  a1------
		  	\     \
	    q[0] \     \ q[2]
		      \     \
		      a2-----
			     q[3]
*/
func (q Quadrilateral) RightAngle(angle float64, offset ...int) (float64, float64) {
	n := conv.Default(0, offset...)

	if angle > q.MaxAngle(n) {
		logs.Debug("超出最大角度", angle, q.MaxAngle(n))
		angle = q.MaxAngle(n)
	}

	a := q.get(n + 0)
	b := q.get(n + 1)
	c := q.get(n + 2)
	d := q.get(n + 3)

	ll := a*a + b*b - 2*a*b*math.Cos(angle*math.Pi/180)
	l := math.Sqrt(ll)
	if l == 0 {
		r := math.Acos(0) * 2 * 180 / math.Pi
		return r, r
	}

	a1 := math.Acos((ll + a*a - b*b) / (2 * l * a))
	a1 = conv.Select(math.IsNaN(a1), 0, a1)
	a2 := math.Acos((ll + d*d - c*c) / (2 * l * d))
	a2 = conv.Select(math.IsNaN(a2), 0, a2)

	//有2个解,突出四边形,凹陷四边形
	r1 := math.Abs(a1+a2) * 180 / math.Pi
	r2 := math.Abs(a1-a2) * 180 / math.Pi

	return r1, r2
}

// Diagonal 已知1,2边角度,求对角线(3,4)角度
func (q Quadrilateral) Diagonal(angle float64, offset ...int) float64 {
	n := conv.Default(0, offset...)
	a := q.get(n + 0)
	b := q.get(n + 1)
	c := q.get(n + 2)
	d := q.get(n + 3)
	ll := a*a + b*b - 2*a*b*math.Cos(angle*math.Pi/180)
	r := math.Acos((c*c + d*d - ll) / (2 * c * d))
	r = r * 180 / math.Pi
	return r
}
