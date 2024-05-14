package zygo

import (
	"errors"
	"fmt"
)

var NotAList = errors.New("not a list")

func ListToArray(expr Sexp) ([]Sexp, error) {
	if !IsList(expr) {
		return nil, NotAList
	}
	arr := make([]Sexp, 0)

	for expr != SexpNull {
		list := expr.(*SexpPair)
		arr = append(arr, list.Head)
		expr = list.Tail
	}

	return arr, nil
}

func MakeList(expressions []Sexp) Sexp {
	if len(expressions) == 0 {
		return SexpNull
	}

	return Cons(expressions[0], MakeList(expressions[1:]))
}

func MapList(env *Zlisp, fun *SexpFunction, expr Sexp) (Sexp, error) {
	if expr == SexpNull {
		return SexpNull, nil
	}

	var list = &SexpPair{}
	switch e := expr.(type) {
	case *SexpPair:
		list.Head = e.Head
		list.Tail = e.Tail
	default:
		return SexpNull, NotAList
	}

	var err error

	list.Head, err = env.Apply(fun, []Sexp{list.Head})

	if err != nil {
		return SexpNull, err
	}

	list.Tail, err = MapList(env, fun, list.Tail)

	if err != nil {
		return SexpNull, err
	}

	return list, nil
}

// O(n^2) for n total nodes in all lists. So this is
// not super efficient. We have to
// find the tail of each list in turn by
// linear search. Avoid lists if possible in favor
// of arrays.
func ConcatLists(a *SexpPair, bs []Sexp) (Sexp, error) {
	result := a
	for _, b := range bs {
		res, err := ConcatTwoLists(result, b)
		if err != nil {
			return SexpNull, err
		}
		x, ok := res.(*SexpPair)
		if !ok {
			return SexpNull, NotAList
		}
		result = x
	}
	return result, nil
}

func ConcatTwoLists(a *SexpPair, b Sexp) (Sexp, error) {
	if !IsList(b) {
		return SexpNull, NotAList
	}

	if a.Tail == SexpNull {
		return Cons(a.Head, b), nil
	}

	switch t := a.Tail.(type) {
	case *SexpPair:
		newtail, err := ConcatTwoLists(t, b)
		if err != nil {
			return SexpNull, err
		}
		return Cons(a.Head, newtail), nil
	}

	return SexpNull, NotAList
}

func ListLen(expr Sexp) (int, error) {
	sz := 0
	var list *SexpPair
	ok := false
	for expr != SexpNull {
		list, ok = expr.(*SexpPair)
		if !ok {
			return 0, fmt.Errorf("ListLen() called on non-list")
		}
		sz++
		expr = list.Tail
	}
	return sz, nil
}
