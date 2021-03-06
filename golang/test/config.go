package main


// test args
// f1.path=f1 f1.keys.map1=f1map1 f2.path=f2 f3.keys.kk.hu=f3 f8.oss.hk.76=f6 f9.2.disk.hk=f7 f5.cache=22 f5.keys.ss.ff=f9keys
// error args
// f4.cache=22


import (
	"os"
	"fmt"
	"errors"
	"strconv"
	"strings"
	"reflect"
	"encoding/json"
)

var (
	errArg 			=	errors.New("undefined args")
	errType 		= 	errors.New("undefined type")
	errValue 		= 	errors.New("undefined value type")
	errUndefined	=	errors.New("setdata use undefined interface")
)

type config struct {
	F1		File
	F2		*File
	F3		interface{}		// null
	F4		interface{}		// struct		error
	F5		interface{}		// struct pointer
	F6		map[string]File
	F7		map[string]*File
	F8		map[string]interface{}
	F9		map[int]interface{}
}

type File struct {
	Path	string
	Key		string
	Srcret	string
	Keys 	interface{}
}


type App struct {
	Cache		string
	Keys 		interface{}
}
func main() {
	c := &config{}
	c.F4 = App{}
	c.F5 = &App{}
	// c.F6 = make(map[string]File)
	// c.F6["cache"] = File{Key: "kes",}
	for _,v := range os.Args[1:] {
		Json(SetData(c,v))
	}
	Json(c)
	Json(Data(c,"f2"))
}


func Json(arg interface{}) {
	indent,_ := json.MarshalIndent(&arg, "", "\t")
	fmt.Println("\n"+string(indent))
}

func Data(p interface, arg string) (reflect.Value, error) {
	fs := strings.Split(arg,".")
	len := len(fs) - 1
	for i,_ := range fs {
		pv := reflect.ValueOf(p)
		for pv.Kind() == reflect.Ptr || pv.Kind() == reflect.Interface {
			pv= pv.Elem()
		}
		fmt.Println("\n----- ",pv.Kind(),pv.Type().Kind(),fs[i])
		switch pv.Kind() {
		case reflect.Struct:
			f,ok := pv.Type().FieldByName(strings.Title(fs[i]))
			if !ok{
				// error
				fmt.Println(errArg)
				return nil,errArg
			}
			pv = pv.Field(f.Index[0])
			if i == len {
				return pv,nil
			}
		case reflect.Map:
			if i==len {
				return pv,nil
			}
			pmk,_ := getvalue(pv.Type().Key(),fs[i])
			if pv.MapIndex(pmk).Kind() == reflect.Invalid {
				fmt.Println("=",pv.Type(),pv.Type().Kind(),pv.Type().Elem())
				fmt.Println(pv.Kind(), pv.Type().Key().Kind(), pv.Type().Elem().Kind())
				if pv.Type().Elem().Kind() == reflect.Interface {
					pv.SetMapIndex(pmk , reflect.ValueOf(make(map[string]interface{})))
					pv = pv.MapIndex(pmk)
					//pv.SetMapIndex(reflect.ValueOf(fs[i]), reflect.ValueOf(make(map[string]interface{})))
					//pv.SetMapIndex(reflect.ValueOf(fs[i]), reflect.MakeMap(pv.Type()))
				}else{
					nv := reflect.New(pv.Type().Elem())
					fmt.Println(fmt.Sprintf("%s=%s",strings.Join(fs[i+1:len+1],".") ,v))
					err := SetData(nv.Interface(),fmt.Sprintf("%s=%s",strings.Join(fs[i+1:len+1],".") ,v))
					pv.SetMapIndex(pmk, nv.Elem() )
					fmt.Printf("%p\n", nv.Interface())
					//fmt.Println( pv.MapIndex(pmk))
					s1 := nv.Elem().Interface().(File)
					s2 := pv.MapIndex(pmk).Interface().(File)
					fmt.Printf("%p\n", &s1)
					fmt.Printf("%p", &s2)
					fmt.Println()
					return nil,err
					pv = nv.Elem()
				}
			}else{
				nv := pv.MapIndex(pmk)
				//pv.SetMapIndex(pmk,nv)
				pv = nv
			}

			fmt.Println("---",pv,pv.Type().Kind())
			fmt.Println(pv.Kind(),pv.Type())
		}
		fmt.Println("?",pv.Type(),pv.Kind(),pv.Type().Kind())//,pv.IsNil())
		// is null
		switch pv.Kind() {
		case reflect.Ptr:
			if pv.IsNil() {
				fmt.Println("+new struct ptr")
				pv.Set(reflect.New(pv.Type().Elem()))	
			}
			pv=pv.Elem()
		case reflect.Invalid:
			// invalid ? ?
			fmt.Println("==",pv.Kind(),pv.Type().Kind())
			if pv.Kind() == reflect.Map {
				fmt.Println("+------------------------new Invalid")
				pv.Set(reflect.MakeMap(pv.Type()))
			}
		case reflect.Map:
			if pv.IsNil() {
				fmt.Println("+new map")
				pv.Set(reflect.MakeMap(pv.Type()))
			}
		case reflect.Interface:
			if pv.IsNil() {
				fmt.Println("+new interface")
				pv.Set(reflect.ValueOf(make(map[string]interface{})))
			}
		}
		// ? ?
		if (pv.Type().Kind() == reflect.Interface) && pv.IsNil() {
			fmt.Println("---------------------------new map",pv.Type().Elem().Kind())
			if pv.Type().Elem().Kind() == reflect.Struct {
					pv.Set(reflect.MakeMap(pv.Type()))
			}else{
				pv.Set( reflect.New(pv.Type().Elem()))		
			}
			//pv.Set(reflect.ValueOf(make(map[string]interface{})))
		}
		// next
		fmt.Println("Addr",pv.CanAddr())
		if pv.CanAddr() {
			p = pv.Addr().Interface()	
		}else {
			p= pv.Interface()
		}
	}
	fmt.Println("--end")
	return nil,nil
}


func SetData(p interface{} , arg string) error {
	kv := append(strings.SplitN(arg,"=",2),"")
	k,v := kv[0],kv[1]
	fs := strings.Split(k,".")
	len := len(fs) - 1
	for i,_ := range fs {
		pv := reflect.ValueOf(p)
		for pv.Kind() == reflect.Ptr || pv.Kind() == reflect.Interface {
			pv= pv.Elem()
		}
		fmt.Println("\n----- ",pv.Kind(),pv.Type().Kind(),fs[i])
		switch pv.Kind() {
		case reflect.Struct:
			f,ok := pv.Type().FieldByName(strings.Title(fs[i]))
			if !ok{
				// error
				fmt.Println(errArg)
				return errArg
			}
			pv = pv.Field(f.Index[0])
			if i == len {
				return setvalue(pv,v)
			}
		case reflect.Map:
			if i==len {
				pv.SetMapIndex(reflect.ValueOf(fs[i]), reflect.ValueOf(v))
				return nil
			}
			pmk,_ := getvalue(pv.Type().Key(),fs[i])
			if pv.MapIndex(pmk).Kind() == reflect.Invalid {
				fmt.Println("=",pv.Type(),pv.Type().Kind(),pv.Type().Elem())
				fmt.Println(pv.Kind(), pv.Type().Key().Kind(), pv.Type().Elem().Kind())
				if pv.Type().Elem().Kind() == reflect.Interface {
					pv.SetMapIndex(pmk , reflect.ValueOf(make(map[string]interface{})))
					pv = pv.MapIndex(pmk)
					//pv.SetMapIndex(reflect.ValueOf(fs[i]), reflect.ValueOf(make(map[string]interface{})))
					//pv.SetMapIndex(reflect.ValueOf(fs[i]), reflect.MakeMap(pv.Type()))
				}else{
					nv := reflect.New(pv.Type().Elem())
					fmt.Println(fmt.Sprintf("%s=%s",strings.Join(fs[i+1:len+1],".") ,v))
					err := SetData(nv.Interface(),fmt.Sprintf("%s=%s",strings.Join(fs[i+1:len+1],".") ,v))
					pv.SetMapIndex(pmk, nv.Elem() )
					fmt.Printf("%p\n", nv.Interface())
					//fmt.Println( pv.MapIndex(pmk))
					s1 := nv.Elem().Interface().(File)
					s2 := pv.MapIndex(pmk).Interface().(File)
					fmt.Printf("%p\n", &s1)
					fmt.Printf("%p", &s2)
					fmt.Println()
					return err
					pv = nv.Elem()
				}
			}else{
				nv := pv.MapIndex(pmk)
				//pv.SetMapIndex(pmk,nv)
				pv = nv
			}

			fmt.Println("---",pv,pv.Type().Kind())
			fmt.Println(pv.Kind(),pv.Type())
		}
		fmt.Println("?",pv.Type(),pv.Kind(),pv.Type().Kind())//,pv.IsNil())
		// is null
		switch pv.Kind() {
		case reflect.Ptr:
			if pv.IsNil() {
				fmt.Println("+new struct ptr")
				pv.Set(reflect.New(pv.Type().Elem()))	
			}
			pv=pv.Elem()
		case reflect.Invalid:
			// invalid ? ?
			fmt.Println("==",pv.Kind(),pv.Type().Kind())
			if pv.Kind() == reflect.Map {
				fmt.Println("+------------------------new Invalid")
				pv.Set(reflect.MakeMap(pv.Type()))
			}
		case reflect.Map:
			if pv.IsNil() {
				fmt.Println("+new map")
				pv.Set(reflect.MakeMap(pv.Type()))
			}
		case reflect.Interface:
			if pv.IsNil() {
				fmt.Println("+new interface")
				pv.Set(reflect.ValueOf(make(map[string]interface{})))
			}
		}
		// ? ?
		if (pv.Type().Kind() == reflect.Interface) && pv.IsNil() {
			fmt.Println("---------------------------new map",pv.Type().Elem().Kind())
			if pv.Type().Elem().Kind() == reflect.Struct {
					pv.Set(reflect.MakeMap(pv.Type()))
			}else{
				pv.Set( reflect.New(pv.Type().Elem()))		
			}
			//pv.Set(reflect.ValueOf(make(map[string]interface{})))
		}
		// next
		fmt.Println("Addr",pv.CanAddr())
		if pv.CanAddr() {
			p = pv.Addr().Interface()	
		}else {
			p= pv.Interface()
		}
	}
	fmt.Println("--end")
	return nil
}


func setvalue(v reflect.Value,s string) error {
	switch v.Kind() {
	case reflect.Bool,
	reflect.Float32, reflect.Float64,
	reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
	reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
	reflect.String:
		k ,_ := getvalue(v.Type(),s)
		v.Set(k)
/*	case reflect.Bool:
		if s == "" {
			v.SetBool(true)
		}else {
			rb,_ := strconv.ParseBool(s)
			v.SetBool(rb)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil || v.OverflowInt(n) {
			return errType
		}
		v.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil || v.OverflowUint(n) {
			return errType
		}
		v.SetUint(n)
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(s, v.Type().Bits())
		if err != nil || v.OverflowFloat(n) {
			return errType
		}
		v.SetFloat(n)
	case reflect.String:
		v.SetString(s)*/
	case reflect.Ptr:
		if v.Elem().Kind() == reflect.Invalid {
			v.Set(reflect.New(v.Type().Elem()))
		}
		return setvalue(v.Elem(),s)
	case reflect.Array:
	case reflect.Slice:
		vs := strings.Split(s,",")
		v.Set(reflect.MakeSlice(v.Type(),len(vs),len(vs)))
		for i,n := range vs {
			setvalue(v.Index(i),n)
		}
	case reflect.Interface:
		fmt.Println("Interface")
	case reflect.Map:
		fmt.Println("Interface Map")
		return json.Unmarshal([]byte(s),v.Interface())
	default:
		fmt.Println("default")
		return errValue
	}
	return nil
}

func getvalue(t reflect.Type,s string) (reflect.Value, error) {
	switch t.Kind() {
	case reflect.Bool:
		if s == "" {
			return reflect.ValueOf(true), nil
		}else {
			rb,_ := strconv.ParseBool(s)
			return reflect.ValueOf(rb), nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return reflect.Zero(t), errType
		}
		return reflect.ValueOf(n).Convert(t),nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return reflect.Zero(t), errType
		}
		return reflect.ValueOf(n).Convert(t),nil
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(s, t.Bits())
		if err != nil {
			return reflect.Zero(t), errType
		}
		return reflect.ValueOf(n).Convert(t),nil
	case reflect.String:
		return reflect.ValueOf(s),nil
	}
	return reflect.Zero(t), errType
}
/*

func SetData(p interface{} , arg string) error {
	kv := append(strings.SplitN(arg,"=",2),"")
	k,v := kv[0],kv[1]
	fs := strings.Split(k,".")
	len := len(fs) - 1
	for i,_ := range fs {
		pv := reflect.ValueOf(p)
		for pv.Kind() == reflect.Ptr || pv.Kind() == reflect.Interface {
			pv= pv.Elem()
		}
		fmt.Println("\n----- ",pv.Kind(),pv.Type().Kind(),fs[i])
		switch pv.Kind() {
		case reflect.Struct:
			f,ok := pv.Type().FieldByName(strings.Title(fs[i]))
			if !ok{
				// error
				fmt.Println(errArg)
				return errArg
			}
			pv = pv.Field(f.Index[0])
			if i == len {
				return setvalue(pv,v)
			}
		case reflect.Map:
			if i==len {
				pv.SetMapIndex(reflect.ValueOf(fs[i]), reflect.ValueOf(v))
				return nil
			}
			if pv.MapIndex(reflect.ValueOf(fs[i])).Kind() == reflect.Invalid {
				fmt.Println("=",pv.Type(),pv.Type().Kind(),pv.Type().Elem())
				fmt.Println(pv.Kind(), pv.Type().Key().Kind(), pv.Type().Elem().Kind())
				if pv.Type().Elem().Kind() == reflect.Interface {
					pv.SetMapIndex(reflect.ValueOf(fs[i]), reflect.MakeMap(pv.Type()))
				}else{
					pv.SetMapIndex(reflect.ValueOf(fs[i]), reflect.New(pv.Type().Elem()))		
				}
			}
			pv = pv.MapIndex(reflect.ValueOf(fs[i]))	
			fmt.Println("---",pv,pv.Type().Kind())
			fmt.Println(pv.Kind(),pv.Type())
		}
		fmt.Println("?",pv.Type(),pv.Kind(),pv.Type().Kind())//,pv.IsNil())
		// is pointer
		switch pv.Kind() {
		case reflect.Ptr:
			// is null
			//pv.IsNil()
			if pv.IsNil() {
				fmt.Println("+new struct ptr")
				pv.Set(reflect.New(pv.Type().Elem()))	
			}
			pv=pv.Elem()
		case reflect.Invalid:
			// invalid ? ?
			fmt.Println("==",pv.Kind(),pv.Type().Kind())
			if pv.Kind() == reflect.Map {
				fmt.Println("+new Invalid")
				pv.Set(reflect.MakeMap(pv.Type()))
			}
		case reflect.Map:
			if pv.IsNil() {
				fmt.Println("+new map")
				pv.Set(reflect.MakeMap(pv.Type()))
			}
		case reflect.Interface:
			if pv.IsNil() {
				fmt.Println("+new interface")
				pv.Set(reflect.ValueOf(make(map[string]interface{})))
			}
		}
		if (pv.Type().Kind() == reflect.Interface) && pv.IsNil() {
			fmt.Println("new map",pv.Type().Elem().Kind())
			if pv.Type().Elem().Kind() == reflect.Struct {
					pv.Set(reflect.MakeMap(pv.Type()))
			}else{
				pv.Set( reflect.New(pv.Type().Elem()))		
			}
			//pv.Set(reflect.ValueOf(make(map[string]interface{})))
		}
		// next
		fmt.Println("Addr",pv.CanAddr())
		if pv.CanAddr() {
			p = pv.Addr().Interface()	
		}else {
			p= pv.Interface()
		}
	}
	fmt.Println("--end")
	return nil
}
*/
