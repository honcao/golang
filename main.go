package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	//"github.com/honcao/golang/pkg/apiv2"
	"math/rand"
	"reflect"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2017-03-30/compute"
	azcompute "github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-04-01/compute"
)

func nextint() func() int {
	i := 0
	return func() int {
		i++
		return i
	}
}

func basicSendRecv() {
	ch := make(chan string)
	var txt = "123"
	go func() {
		txt = "abc"
		close(ch)
		txt = "abc1"
	}()
	fmt.Println(txt)
	<-ch
	fmt.Println(txt)

	// can not close chanel twice
	// make(chan int)
}

func signalAck() {
	ch := make(chan string)

	go func() {
		fmt.Println(<-ch)
		ch <- "pong"
	}()
	ch <- "ping"
	fmt.Println(<-ch)
}

func closeRange() {
	ch := make(chan int, 5)
	go func() {
		for i := 0; i < 15; i++ {
			ch <- i
		}
		close(ch)
	}()

	for k := range ch {
		fmt.Println(k)
	}

}

func selectRece() {
	ch := make(chan string, 1)
	go func() {
		time.Sleep(15 * time.Millisecond)
		ch <- "do work complete"
	}()

	select {
	case v := <-ch:
		fmt.Println(v)
	case <-time.After(10 * time.Millisecond):
		fmt.Println("Time out")
	}
}

func waitGroup() {
	w := &sync.WaitGroup{}
	w.Add(2)
	go func() {
		time.Sleep(5000 * time.Millisecond)
		fmt.Println("this is one")
		w.Done()
	}()

	go func() {
		time.Sleep(15 * time.Millisecond)
		fmt.Println("this is two")
		w.Done()
	}()
	w.Wait()
	fmt.Println("done")
}

// Play dkdkd
func Play(name string, hit chan int) {
	for {
		rand := rand.Intn(100)
		fmt.Println(rand)
		count, ok := <-hit

		if !ok {
			fmt.Printf("%s won\n", name)
			return
		}

		if rand != 0 && rand%13 == 0 {
			fmt.Printf("%s missed the ball\n", name)
			close(hit)
			return
		}

		count++
		fmt.Printf("%s hits the ball back %d\n", name, count)
		hit <- count
	}
}

func initializeStruct(t reflect.Type, v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		ft := t.Field(i)
		switch ft.Type.Kind() {
		case reflect.Map:
			f.Set(reflect.MakeMap(ft.Type))
		case reflect.Slice:
			f.Set(reflect.MakeSlice(ft.Type, 0, 0))
		case reflect.Chan:
			f.Set(reflect.MakeChan(ft.Type, 0))
		case reflect.Struct:
			initializeStruct(ft.Type, f)
		case reflect.Ptr:
			fv := reflect.New(ft.Type.Elem())
			initializeStruct(ft.Type.Elem(), fv.Elem())
			f.Set(fv)
		default:
		}
	}
}

// DeepAssignment dst and src should be the same type in different API version
// dst should be pointer type
func DeepAssignment(dst, src interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Fatal("Fail to covert object", r)
		}
	}()
	dstValue := reflect.ValueOf(dst)
	srcValue := reflect.ValueOf(src)
	if dstValue.Kind() != reflect.Ptr {
		log.Fatal("The dst must be Ptr")
	}
	dstValue = dstValue.Elem()
	deepAssignmentInternal(dstValue, srcValue, 0, "")
}

func deepAssignmentInternal(dstValue, srcValue reflect.Value, depth int, path string) {
	if dstValue.CanSet() {
		switch srcValue.Kind() {
		case reflect.Bool:
			dstValue.SetBool(srcValue.Bool())
		case reflect.String:
			dstValue.SetString(srcValue.String())
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
			dstValue.SetInt(srcValue.Int())
		case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
			dstValue.SetUint(srcValue.Uint())
		case reflect.Float64, reflect.Float32:
			dstValue.SetFloat(srcValue.Float())
		case reflect.Complex64, reflect.Complex128:
			dstValue.SetComplex(srcValue.Complex())
		case reflect.Ptr:
			if !srcValue.IsNil() {
				d := reflect.New(dstValue.Type().Elem())
				dstValue.Set(d)
				deepAssignmentInternal(dstValue.Elem(), srcValue.Elem(), depth+1, "")
			}
		case reflect.Slice:
			if !srcValue.IsNil() {
				d := reflect.MakeSlice(dstValue.Type(), srcValue.Len(), srcValue.Cap())
				dstValue.Set(d)
				for i := 0; i < srcValue.Len(); i++ {
					v := dstValue.Index(i)
					deepAssignmentInternal(v, srcValue.Index(i), depth+1, "")
					v.Set(v)
				}
			}
		case reflect.Array:
			d := reflect.New(dstValue.Type()).Elem()
			for i := 0; i < srcValue.Len(); i++ {
				v := reflect.New(srcValue.Index(i).Type()).Elem()
				deepAssignmentInternal(v, srcValue.Index(i), depth+1, "")
				d.Index(i).Set(v)
			}
			dstValue.Set(d)
		case reflect.Map:
			if !srcValue.IsNil() {
				d := reflect.MakeMap(dstValue.Type())
				for _, key := range srcValue.MapKeys() {
					v := reflect.New(srcValue.MapIndex(key).Type()).Elem()
					deepAssignmentInternal(v, srcValue.MapIndex(key), depth+1, "")
					d.SetMapIndex(key, v)
				}
				dstValue.Set(d)
			}
		case reflect.Struct:
			for i := 0; i < srcValue.NumField(); i++ {
				srcField := srcValue.Field(i)
				dstField := dstValue.FieldByName(srcValue.Type().Field(i).Name)
				if dstField.IsValid() && dstField.CanAddr() {
					deepAssignmentInternal(dstField, srcField, depth+1, "")
				}
			}
		default:
		}
	}
}

type myString string

func main() {

	/*
		var x float64 = 3.4
		ax := &x
		fmt.Println(ax)
		fmt.Println(&ax)

		xv := reflect.ValueOf(x)
		//fmt.Println(&xv)
		//xve := xv.Elem()
		fmt.Println(&xv)
		xv.SetFloat(5.0)
		fmt.Print(x)
		fmt.Println(ax)
		fmt.Println(*ax)

		var y float64 = 6.7
		xy := &y
		fmt.Println(y)
		fmt.Println(xy)
		fmt.Println(*xy)
		*xy = 5.9
		fmt.Println(y)
		fmt.Println(xy)
		fmt.Println(*xy)
	*/

	/*
		//src := []string{"abc"}
		dst := make([]string, 5, 5)
		//srcV := reflect.ValueOf(src)
		dstx := &dst
		dstxV := reflect.ValueOf(dstx)

		velm := dstxV.Elem()
		a := velm.Index(0)
		a.Set(reflect.ValueOf("abc"))
		//velm.Set(velm)

		fmt.Println(dst)
	*/
	vmsazs := []compute.VirtualMachine{}
	datvm, _ := ioutil.ReadFile("vm.json")
	if err := json.Unmarshal(datvm, &vmsazs); err != nil {
		panic(err)
	}
	vms := []azcompute.VirtualMachine{} //make([]azcompute.VirtualMachine, len(vmsazs), cap(vmsazs))
	DeepAssignment(&vms, vmsazs)
	fmt.Println(vms)

	vmsssazs := []compute.VirtualMachineScaleSet{}
	vmsss := []azcompute.VirtualMachineScaleSet{}
	datvmss, _ := ioutil.ReadFile("vmss.json")
	if err := json.Unmarshal(datvmss, &vmsssazs); err != nil {
		panic(err)
	}
	DeepAssignment(&vmsss, vmsssazs)
	fmt.Println(vmsss)
}
