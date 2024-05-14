package zygo

import (
	"fmt"
	"reflect"
	"runtime"
)

// call Go methods

// Using reflection, invoke a Go method on a struct or interface.
// args[0] is a hash with an an attached GoStruct
// args[1] is a hash representing a method call on that struct.
// The returned Sexp is a hash that represents the result of that call.
func CallGoMethodFunction(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	Q("_method user func running!\n")

	// protect against bad calls/bad reflection
	var wasPanic bool
	var recovered interface{}
	tr := make([]byte, 16384)
	trace := &tr
	sx, err := func() (Sexp, error) {
		defer func() {
			recovered = recover()
			if recovered != nil {
				wasPanic = true
				nbyte := runtime.Stack(*trace, false)
				*trace = (*trace)[:nbyte]
			}
		}()

		narg := len(args)
		if narg < 2 {
			return SexpNull, WrongNargs
		}
		obj, isHash := args[0].(*SexpHash)
		if !isHash {
			return SexpNull, fmt.Errorf("_method error: first argument must be a hash or defmap (a record) with an attached GoObject")
		}

		var methodname string
		switch m := args[1].(type) {
		case *SexpSymbol:
			methodname = m.name
		case *SexpStr:
			methodname = m.S
		default:
			return SexpNull, fmt.Errorf("_method error: second argument must be a method name in symbol or string form (got %T)", args[1])
		}

		// get the method list, verify the method exists and get its type
		if obj.NumMethod == -1 {
			err := obj.SetMethodList(env)
			if err != nil {
				return SexpNull, fmt.Errorf("could not get method list for object: %s", err)
			}
		}

		var method reflect.Method
		found := false
		for _, me := range obj.GoMethods {
			if me.Name == methodname {
				method = me
				found = true
				break
			}
		}
		if !found {
			return SexpNull, fmt.Errorf("no such method '%s' on %s. choices are: %s",
				methodname, obj.TypeName,
				(obj.GoMethSx).SexpString(nil))
		}
		// INVAR: var method holds our call target

		// try always expecting this to be already done... test crashes
		//P("in CallGoMethod '%s' obj.GoShadowStructVa = '%#v'", methodname, obj.GoShadowStructVa)
		if obj.GoShadowStructVa.Kind() == reflect.Invalid {
			// ready the struct... but only because there isn't already a shadow struct there!!
			if !obj.ShadowSet {
				_, err := ToGoFunction(env, "togo", []Sexp{obj})
				if err != nil {
					return SexpNull, fmt.Errorf("error converting object to Go struct: '%s'", err)
				}
			}
		}
		inputVa := []reflect.Value{(obj.GoShadowStructVa)}

		// prep args.
		needed := method.Type.NumIn() - 1 // one for the receiver
		avail := narg - 2
		if needed != avail {
			// TODO: support varargs eventually
			return SexpNull, fmt.Errorf("method %s needs %d arguments, but we have %d", method.Name, needed, avail)
		}

		var va reflect.Value
		for i := 2; i < narg; i++ {
			typ := method.Type.In(i - 1)
			pdepth := PointerDepth(typ)
			// we only handle 0 and 1 for now
			Q("pdepth = %v\n", pdepth)
			switch pdepth {
			case 0:
				va = reflect.New(typ)
			case 1:
				// handle the common single pointer to struct case
				va = reflect.New(typ.Elem())
			default:
				return SexpNull, fmt.Errorf("error converting %d-th argument to "+
					"Go: we don't handle double pointers", i-2)
			}
			Q("converting to go '%#v' into -> %#v\n", args[i], va.Interface())
			iface, err := SexpToGoStructs(args[i], va.Interface(), env, nil)
			if err != nil {
				return SexpNull, fmt.Errorf("error converting %d-th "+
					"argument to Go: '%s'", i-2, err)
			}
			switch pdepth {
			case 0:
				inputVa = append(inputVa, reflect.ValueOf(iface).Elem())
			case 1:
				inputVa = append(inputVa, reflect.ValueOf(iface))
			}
			Q("\n allocated new %T/val=%#v /i=%#v\n", va, va, va.Interface())
		}

		//P("_method: about to .Call by reflection!\n")

		out := method.Func.Call(inputVa)

		var iout []interface{}
		for _, o := range out {
			iout = append(iout, o.Interface())
		}
		Q("done with _method call, iout = %#v\n", iout)
		Q("done with _method call, iout[0] = %#v\n", iout[0])

		nout := len(out)
		r := make([]Sexp, 0)
		for i := 0; i < nout; i++ {
			f := out[i].Interface()
			switch e := f.(type) {
			case nil:
				r = append(r, SexpNull)
			case int64:
				r = append(r, &SexpInt{Val: e})
			case int:
				r = append(r, &SexpInt{Val: int64(e)})
			case error:
				r = append(r, &SexpError{e})
			case string:
				r = append(r, &SexpStr{S: e})
			case float64:
				r = append(r, &SexpFloat{Val: e})
			case []byte:
				r = append(r, &SexpRaw{Val: e})
			case rune:
				r = append(r, &SexpChar{Val: e})
			default:
				// go through the type registry
				found := false
				for hashName, factory := range GoStructRegistry.Registry {
					st, err := factory.Factory(env, nil)
					if err != nil {
						return SexpNull, fmt.Errorf("MakeHash '%s' problem on Factory call: %s",
							hashName, err)
					}
					Q("got st from Factory, checking if types match")
					if reflect.ValueOf(st).Type() == out[i].Type() {
						Q("types match")
						retHash, err := MakeHash([]Sexp{}, factory.RegisteredName, env)
						if err != nil {
							return SexpNull, fmt.Errorf("MakeHash '%s' problem: %s",
								hashName, err)
						}

						Q("filling from shadow")
						err = retHash.FillHashFromShadow(env, f)
						if err != nil {
							return SexpNull, err
						}
						r = append(r, retHash)
						found = true
						break
					}
				}
				if !found {
					r = append(r, &SexpReflect{Val: out[i]})
				}
			}
		}
		return env.NewSexpArray(r), nil
	}()
	if wasPanic {
		return SexpNull, fmt.Errorf("\n recovered from panic "+
			"during CallGo. panic on = '%v'\n"+
			"stack trace:\n%s\n", recovered, string(*trace))
	}
	return sx, err
}

// detect if inteface is holding anything
func NilOrHoldsNil(iface interface{}) bool {
	if iface == nil {
		return true
	}
	return reflect.ValueOf(iface).IsNil()
}
