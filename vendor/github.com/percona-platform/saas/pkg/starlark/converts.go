package starlark

import (
	"time"

	"github.com/pkg/errors"
	"go.starlark.net/starlark"
)

func goToStarlark(v interface{}) (starlark.Value, error) {
	switch v := v.(type) {
	case nil:
		return starlark.None, nil
	case bool:
		return starlark.Bool(v), nil
	case int64:
		return starlark.MakeInt64(v), nil
	case uint64:
		return starlark.MakeUint64(v), nil
	case float64:
		return starlark.Float(v), nil
	case []byte:
		return starlark.String(v), nil
	case string:
		return starlark.String(v), nil
	case time.Time:
		return starlark.MakeInt64(v.UnixNano()), nil
	case []interface{}:
		return goToStarlarkList(v)
	case map[string]interface{}:
		return goToStarlarkDict(v)
	case map[bool]struct{}:
		return goToStarlarkSetBool(v)
	case map[int64]struct{}:
		return goToStarlarkSetInt(v)
	case map[uint64]struct{}:
		return goToStarlarkSetUint(v)
	case map[float64]struct{}:
		return goToStarlarkSetFloat(v)
	case map[string]struct{}:
		return goToStarlarkSetString(v)
	default:
		return nil, errors.Errorf("unhandled type %T", v)
	}
}

func goToStarlarkList(v []interface{}) (starlark.Value, error) {
	l := make([]starlark.Value, len(v))
	for i, o := range v {
		sv, err := goToStarlark(o)
		if err != nil {
			return nil, err
		}
		l[i] = sv
	}
	return starlark.NewList(l), nil
}

func goToStarlarkDict(v map[string]interface{}) (starlark.Value, error) {
	sd := starlark.NewDict(len(v))
	for k, o := range v {
		sv, err := goToStarlark(o)
		if err != nil {
			return nil, err
		}
		if err := sd.SetKey(starlark.String(k), sv); err != nil {
			return nil, errors.Wrap(err, "failed to set key in dict")
		}
	}
	return sd, nil
}

func goToStarlarkSetBool(v map[bool]struct{}) (starlark.Value, error) {
	ss := starlark.NewSet(len(v))
	for k := range v {
		err := ss.Insert(starlark.Bool(k))
		if err != nil {
			return nil, errors.Wrap(err, "failed to insert into set")
		}
	}
	return ss, nil
}

func goToStarlarkSetInt(v map[int64]struct{}) (starlark.Value, error) {
	ss := starlark.NewSet(len(v))
	for k := range v {
		err := ss.Insert(starlark.MakeInt64(k))
		if err != nil {
			return nil, errors.Wrap(err, "failed to insert into set")
		}
	}
	return ss, nil
}

func goToStarlarkSetUint(v map[uint64]struct{}) (starlark.Value, error) {
	ss := starlark.NewSet(len(v))
	for k := range v {
		err := ss.Insert(starlark.MakeUint64(k))
		if err != nil {
			return nil, errors.Wrap(err, "failed to insert into set")
		}
	}
	return ss, nil
}

func goToStarlarkSetFloat(v map[float64]struct{}) (starlark.Value, error) {
	ss := starlark.NewSet(len(v))
	for k := range v {
		err := ss.Insert(starlark.Float(k))
		if err != nil {
			return nil, errors.Wrap(err, "failed to insert into set")
		}
	}
	return ss, nil
}

func goToStarlarkSetString(v map[string]struct{}) (starlark.Value, error) {
	ss := starlark.NewSet(len(v))
	for k := range v {
		err := ss.Insert(starlark.String(k))
		if err != nil {
			return nil, errors.Wrap(err, "failed to insert into set")
		}
	}
	return ss, nil
}

func starlarkToGo(v starlark.Value) (interface{}, error) {
	switch v := v.(type) {
	case starlark.NoneType:
		return nil, nil
	case starlark.Bool:
		return bool(v), nil
	case starlark.Int:
		if i, ok := v.Int64(); ok {
			return i, nil
		}
		if u, ok := v.Uint64(); ok {
			return u, nil
		}
		return nil, errors.Errorf("unhandled type %T", v)
	case starlark.Float:
		return float64(v), nil
	case starlark.String:
		return string(v), nil
	case starlark.Tuple:
		res := make([]interface{}, len(v))
		for i, o := range v {
			no, err := starlarkToGo(o)
			if err != nil {
				return nil, err
			}
			res[i] = no
		}
		return res, nil
	case *starlark.List:
		res := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			no, err := starlarkToGo(v.Index(i))
			if err != nil {
				return nil, err
			}
			res[i] = no
		}
		return res, nil
	case *starlark.Dict:
		res := make(map[string]interface{}, v.Len())
		for _, tu := range v.Items() {
			var err error
			k := tu[0].(starlark.String).GoString()
			res[k], err = starlarkToGo(tu[1])
			if err != nil {
				return nil, err
			}
		}
		return res, nil
	case *starlark.Set:
		return starlarkSetToGo(v)
	default:
		return nil, errors.Errorf("unhandled type %T", v)
	}
}

func starlarkSetToGo(v *starlark.Set) (interface{}, error) {
	iter := v.Iterate()
	var x starlark.Value
	var tp string
	var notValid bool
	for iter.Next(&x) {
		if tp != "" && tp != x.Type() {
			notValid = true
			iter.Done()
		}
		tp = x.Type()
	}
	if notValid {
		return nil, errors.New("more types in starlark.Set")
	}
	iter = v.Iterate()
	defer iter.Done()

	switch tp {
	case "bool":
		return starlarkSetBoolToGo(v, iter, x)
	case "int":
		vl := x.(starlark.Int)
		if _, ok := vl.Int64(); ok {
			return starlarkSetIntToGo(v, iter, x)
		}
		if _, ok := vl.Uint64(); ok {
			return starlarkSetUintToGo(v, iter, x)
		}
		return nil, errors.Errorf("unhandled type %s", tp)
	case "float":
		return starlarkSetFloatToGo(v, iter, x)
	case "string":
		res := make(map[string]struct{}, v.Len())
		for iter.Next(&x) {
			res[x.(starlark.String).GoString()] = struct{}{}
		}
		return res, nil
	default:
		return nil, errors.Errorf("unhandled type %s", tp)
	}
}

func starlarkSetBoolToGo(v *starlark.Set, iter starlark.Iterator, x starlark.Value) (map[bool]struct{}, error) {
	res := make(map[bool]struct{}, v.Len())
	for iter.Next(&x) {
		nv, err := starlarkToGo(x.(starlark.Bool))
		if err != nil {
			return nil, err
		}
		res[nv.(bool)] = struct{}{}
	}

	return res, nil
}

func starlarkSetIntToGo(v *starlark.Set, iter starlark.Iterator, x starlark.Value) (map[int64]struct{}, error) {
	res := make(map[int64]struct{}, v.Len())
	for iter.Next(&x) {
		nv, err := starlarkToGo(x.(starlark.Int))
		if err != nil {
			return nil, err
		}
		res[nv.(int64)] = struct{}{}
	}
	return res, nil
}

func starlarkSetUintToGo(v *starlark.Set, iter starlark.Iterator, x starlark.Value) (map[uint64]struct{}, error) {
	res := make(map[uint64]struct{}, v.Len())
	for iter.Next(&x) {
		nv, err := starlarkToGo(x.(starlark.Int))
		if err != nil {
			return nil, err
		}
		res[nv.(uint64)] = struct{}{}
	}
	return res, nil
}

func starlarkSetFloatToGo(v *starlark.Set, iter starlark.Iterator, x starlark.Value) (map[float64]struct{}, error) {
	res := make(map[float64]struct{}, v.Len())
	for iter.Next(&x) {
		nv, err := starlarkToGo(x.(starlark.Float))
		if err != nil {
			return nil, err
		}
		res[nv.(float64)] = struct{}{}
	}
	return res, nil
}
