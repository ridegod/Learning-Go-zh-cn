package main

import (
	"fmt"
	"strconv"
	"container/vector"
	"flag"
)

const (
	_ = 1000 * iota
	ADD
	SUB
	MUL
	DIV
	MAXPOS = 11
)

var mop = map[int]string{
	ADD: "+",
	SUB: "-",
	MUL: "*",
	DIV: "/",
}

var (
	ok bool
	value int
)

type Stack struct {
	i    int
	data [MAXPOS]int
}

func (s *Stack) Reset() {
	s.i = 0
}

func (s *Stack) Len() int {
	return s.i
}

func (s *Stack) Push(k int) {
	s.data[s.i] = k
	s.i++
}

func (s *Stack) Pop() int {
	s.i--
	return s.data[s.i]
}

var found int
var stack = new(Stack)

func main() {
	flag.Parse()
	list := []int{1, 6, 7, 8, 8, 75, ADD, SUB, MUL, DIV}
	//	list := []int{1, 6, 7, ADD, SUB, MUL, DIV}
	magic, ok := strconv.Atoi(flag.Arg(0))
	if ok != nil {
		return
	}
	f := make([]int, MAXPOS)
	solve(f, list, 0, magic)
}

func solve(form, numberop []int, index, magic int) {
	var tmp int
	for i, v := range numberop {
		if v == 0 {
			goto NEXT
		}

		if v < ADD {
			// 是一个数字，保存起来
			tmp = numberop[i]
			numberop[i] = 0
		}
		form[index] = v
		value, ok = rpncalc(form[0 : index+1])

		if ok && value == magic {
			if v < ADD {
				numberop[i] = tmp // 重置并继续
			}
			found++
			fmt.Printf("%s = %d  #%d\n", rpnstr(form[0:index+1]), value, found)
			//goto NEXT
		}

		if index == MAXPOS-1 {
			if v < ADD {
				numberop[i] = tmp // 重置并继续
			}
			goto NEXT
		}
		solve(form, numberop, index+1, magic)
		if v < ADD {
			numberop[i] = tmp // 重置并继续
		}
	NEXT:
	}
}

// convert rpn to nice infix notation and string
// the r must be valid rpn form
func rpnstr(r []int) (ret string) {
	s := new(vector.StringVector)
	for k, t := range r {
		switch t {
		case ADD, SUB, MUL, DIV:
			a := s.Pop()
			b := s.Pop()
			if k == len(r)-1 {
				s.Push(b + mop[t] + a)
			} else {
				s.Push("(" + b + mop[t] + a + ")")
			}
		default:
			s.Push(strconv.Itoa(t))
		}
	}
	for _, v := range *s {
		ret += v
	}
	return
}

// return result from the rpn form.
// if the expression is not valid, ok is false
func rpncalc(r []int) (int, bool) {
	stack.Reset()
	for _, t := range r {
		switch t {
		case ADD, SUB, MUL, DIV:
			if stack.Len() < 2 {
				return 0, false
			}
			a := stack.Pop()
			b := stack.Pop()
			if t == ADD {
				stack.Push(b + a)
			}
			if t == SUB {
				// disallow negative subresults
				if b-a < 0 {
					return 0, false
				}
				stack.Push(b - a)
			}
			if t == MUL {
				stack.Push(b * a)
			}
			if t == DIV {
				if a == 0 {
					return 0, false
				}
				// disallow fractions
				if b%a != 0 {
					return 0, false
				}
				stack.Push(b / a)
			}
		default:
			stack.Push(t)
		}
	}
	if stack.Len() == 1 { // there is only one!
		return stack.Pop(), true
	}
	return 0, false
}
