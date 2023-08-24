package alg

import "testing"

func TestPointInRectangle(t *testing.T) {
	{
		q := map[Point]bool{
			{-7, 0}:  false,
			{-7, -7}: false,
			{0, 0}:   true,
			{7, 7}:   false,
			{5, 6}:   true,
			{5, 5.5}: true,
			{5, 6.1}: false,
		}
		for i, v := range q {
			if PointInRectangle(i, [2]Point{
				{-6, -6},
				{6, 6},
			}) != v {
				t.Log("错误")
			}
		}
	}
	{
		q := map[Point]bool{
			{-7, 0}:  false,
			{-7, -7}: false,
			{0, 0}:   true,
			{0, 7}:   false,
			{5, 6}:   false,
			{0, 5.5}: true,
			{5, 6.1}: false,
		}
		for i, v := range q {
			if PointInRectangle(i, [2]Point{
				{0, 0},
				{0, 6},
			}) != v {
				t.Log("错误")
			}
		}
	}

}

func TestPointInLine(t *testing.T) {
	{
		q := map[Point]bool{
			{-7, 0}:  false,
			{-7, -7}: false,
			{0, 0}:   true,
			{0, 7}:   false,
			{5, 6}:   false,
			{3, 1}:   true,
			{5, 6.1}: false,
		}
		for i, v := range q {
			if PointInLine(i, [2]Point{
				{-6, -2},
				{6, 2},
			}) != v {
				t.Log("错误")
			}
		}
	}
	{
		q := map[Point]bool{
			{-7, 0}:  false,
			{0, -2}:  false,
			{0, 0}:   true,
			{0, 7}:   false,
			{5, 6}:   false,
			{0, 5.5}: true,
			{5, 6.1}: false,
		}
		for i, v := range q {
			if PointInLine(i, [2]Point{
				{0, -1},
				{0, 6},
			}) != v {
				t.Log("错误")
			}
		}
	}
}

func TestLineIntersectLine(t *testing.T) {
	//t.Log(LineIntersectLine(Line{}, Line{}))
	t.Log(LineIntersect(Line{
		{0, 0},
		{6, 6},
	}, Line{
		{0, 6},
		{2, 0},
	}))
}
