package main

import (
	"context"
	"fmt"
	"log"

	//"github.com/honcao/golang/pkg/apiv2"
	"math/rand"
	"reflect"
	"sync"
	"time"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/honcao/golang/pkg/apiv1"
	"github.com/honcao/golang/pkg/armhelper"
	azarmhelper "github.com/honcao/golang/pkg/armhelper/azure"
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

func DeepAssignmentWrapper(dst, src interface{}) {

	dstValue := reflect.ValueOf(dst)
	srcValue := reflect.ValueOf(src)

	if dstValue.Kind() != reflect.Ptr {
		log.Fatal("dst is not pointer to stuct")
	}
	dstValue = dstValue.Elem()
	if dstValue.Kind() != reflect.Struct {
		fmt.Println(dstValue.Kind())
		log.Fatal("dst is not pointer to stuct")
	}
	if srcValue.Kind() != reflect.Struct {
		log.Fatal("src is not stuct")
	}
	//initializeStruct(dstValue.Type(), dstValue)
	DeepAssignment(dstValue, srcValue, 0, "")
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

func DeepAssignment(dstValue, srcValue reflect.Value, depth int, path string) {
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
			d := reflect.New(dstValue.Type().Elem())
			dstValue.Set(d)
			DeepAssignment(dstValue.Elem(), srcValue.Elem(), depth+1, "")
		case reflect.Slice:
			d := reflect.New(dstValue.Type()).Elem()
			for i := 0; i < srcValue.Len(); i++ {
				v := reflect.New(srcValue.Index(i).Type()).Elem()
				DeepAssignment(v, srcValue.Index(i), depth+1, "")
				d = reflect.Append(d, v)
			}
			dstValue.Set(d)
		case reflect.Array:
			d := reflect.New(dstValue.Type()).Elem()
			for i := 0; i < srcValue.Len(); i++ {
				v := reflect.New(srcValue.Index(i).Type()).Elem()
				DeepAssignment(v, srcValue.Index(i), depth+1, "")
				d.Index(i).Set(v)
			}
			dstValue.Set(d)
		case reflect.Map:
			d := reflect.MakeMap(dstValue.Type())
			for _, key := range srcValue.MapKeys() {
				v := reflect.New(srcValue.MapIndex(key).Type()).Elem()
				fmt.Println(srcValue.MapIndex(key))
				DeepAssignment(v, srcValue.MapIndex(key), depth+1, "")
				d.SetMapIndex(key, v)
			}
			dstValue.Set(d)
		case reflect.Struct:
			for i := 0; i < srcValue.NumField(); i++ {
				srcField := srcValue.Field(i)
				fmt.Println(srcValue.Type().Field(i).Name)
				dstField := dstValue.FieldByName(srcValue.Type().Field(i).Name)
				if dstField.IsValid() {
					DeepAssignment(dstField, srcField, depth+1, "")
				}
			}
		default:
		}
	}
}

type myString string

func main() {

	env, err := azure.EnvironmentFromName("AzurePublicCloud")
	if err != nil {
		panic(err)
	}

	armclient, _ := azarmhelper.NewAzureClientWithClientSecret(env, "9ee2ec52-83c0-405e-a009-6636ead37acd", "72f988bf-86f1-41af-91ab-2d7cd011db47", "85115f84-ef7b-4ddb-b44d-b3a9d3b1990d", "Jw5Sx64ANZoU0bFbxmzoLZxC3U01Kn3PiKAQB3Aa+ZU=")

	rg, err := armclient.EnsureResourceGroup(context.Background(), "honcaorg1", "eastus", nil)
	g := armhelper.Group{}
	DeepAssignmentWrapper(&g, *rg)
	fmt.Println(err)
	fmt.Printf("%T\n", *rg.Properties)
	fmt.Println(*rg.Name)

	s := "abc"
	i := -100
	src := apiv1.TypeB{
		Name:        "abc",
		Age:         -15,
		Uint:        100,
		Float32:     1.9005,
		StringPtr:   &s,
		IntPtr:      &i,
		IntSlice:    []int{1, 2, 3, 4},
		StringSlice: []string{"a1", "b", "c", "d"},
		StructPtr:   &apiv1.TypeD{DName: "dname"},
		TypeCSlice: []apiv1.TypeC{
			apiv1.TypeC{
				CName: "c1",
			},
			apiv1.TypeC{
				CName: "c2",
			},
		},
		StringArray:  [5]string{"aa", "bb", "cc", "dd", "ee"},
		StringMap:    map[string]string{"abc": "bcd", "def": "fgh"},
		StringMapMap: map[string]map[string]string{"abc": {"abci": "bcd"}, "def": {"defi": "fgh"}},
	}
	dst := apiv1.TypeB{}
	DeepAssignmentWrapper(&dst, src)
	fmt.Println(src)
	fmt.Println(dst)
	fmt.Println(*dst.StructPtr)

	/*
		var src string
		var dst myString
		src = "test"
		DeepAssignmentWrapper(&dst, src)
		fmt.Println(dst)

		uv1 := apiv1.User{
			Name: "Hongbin",
			Age:  38,
			Test: &src,
		}
		t := reflect.TypeOf(uv1)
		v := reflect.ValueOf(&uv1)
		fmt.Println(t.String())
		fmt.Println(t.Field(1).Name)
		fmt.Println(v.String())
		fmt.Println(v.Type() == t)
		fmt.Println(v.Elem().Field(1))
		v.Elem().Field(1).SetInt(100)
		fmt.Println(v.Elem().Field(1))
		fmt.Println(v.Elem().Type().Field(1).Name)

		e := reflect.ValueOf(&uv1).Elem()
		fmt.Println(e.Type())
		fmt.Println(e.Kind())
		for i := 0; i < e.NumField(); i++ {
			varName := e.Type().Field(i).Name
			varType := e.Type().Field(i).Type
			varValue := e.Field(i).Interface()
			fmt.Printf("%v %v %v\n", varName, varType, varValue)
		}
	*/
	/*
		uv2 := apiv2.User{
			Name: "Hongbin",
			Age:  "38",
		}
		e2 := reflect.ValueOf(&uv2).Elem()
		for i := 0; i < e2.NumField(); i++ {
			varName := e2.Type().Field(i).Name
			varType := e2.Type().Field(i).Type
			e2.Field(i)
			varValue := e2.Field(i)
			fmt.Printf("%v %v %v\n", varName, varType, varValue)
		}
		fmt.Println(e2)
	*/
	/*
			w := &sync.WaitGroup{}
			ch := make(chan int)

			w.Add(2)
			go func() {
				Play("Hongbin", ch)
				w.Done()
			}()
			go func() {
				Play("Hongliang", ch)
				w.Done()
			}()
			ch <- 0

		w.Wait()
	*/
	/*
		fn := decorator.PreparerFunc(func(r *http.Request) (*http.Request, error) {
			r.Method = "PUT"
			return r, nil
		})

		rep, _ := fn.Prepare(&http.Request{})
		fmt.Println(rep.Method)

		p := decorator.CreatePreparer(decorator.SetMethod("PUT"))
		r, _ := p.Prepare(&http.Request{})
		fmt.Println(r.Method)
	*/
	/*
		s3, _ := exam.New("Hongbin", 100)
		r, _ := s3.Test(456)
		a := exam.Tester(s3)
		fmt.Println(a.Test(1234))
		fmt.Println(r)
		env, _ := azure.EnvironmentFromName("AZUREPUBLICCLOUD")
		bytes, _ := json.Marshal(env)
		s1 := string(bytes)

		fmt.Println(strings.Replace(s1, "\"", "\\\"", -1))

		s := make(chan string, 1)

		go func() {
			time.Sleep(200 * time.Millisecond)
			s <- "abc"
		}()

		select {
		case v := <-s:
			fmt.Println(v)
		case t := <-time.After(100 * time.Millisecond):
			fmt.Println(t)
			fmt.Println("timeout")
		}
	*/
	/*
		i := make([]int, 10)
		for _, value := range i {
			fmt.Println(value)
		}
		fmt.Printf("%T\n", i)
		env, err := azure.EnvironmentFromName("AzurePublicCloud")
		if err != nil {
			panic(err)
		}

		armclient, err := azarmhelper.NewAzureClientWithClientSecret(env, "9ee2ec52-83c0-405e-a009-6636ead37acd", "72f988bf-86f1-41af-91ab-2d7cd011db47", "85115f84-ef7b-4ddb-b44d-b3a9d3b1990d", "y9liZF65vOyPgpjqJLUnOnjRRH7i4rCA+EMhPAM4dac=")

		rg, _ := armclient.EnsureResourceGroup(context.Background(), "honcaorg1", "eastus", nil)
		fmt.Printf("%T\n", *rg.Properties)
		fmt.Println(*rg.Name)

		azsarmclient, err := azsarmhelper.NewAzureClientWithClientSecret(env, "9ee2ec52-83c0-405e-a009-6636ead37acd", "72f988bf-86f1-41af-91ab-2d7cd011db47", "85115f84-ef7b-4ddb-b44d-b3a9d3b1990d", "y9liZF65vOyPgpjqJLUnOnjRRH7i4rCA+EMhPAM4dac=")

		rg2, err := azsarmclient.EnsureResourceGroup(context.Background(), "honcaorg2", "westus", nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(*rg2.Name)
	*/
	/*
		os.Setenv("AZURE_TENANT_ID", "72f988bf-86f1-41af-91ab-2d7cd011db47")
		os.Setenv("AZURE_CLIENT_ID", "85115f84-ef7b-4ddb-b44d-b3a9d3b1990d")
		os.Setenv("AZURE_CLIENT_SECRET", "y9liZF65vOyPgpjqJLUnOnjRRH7i4rCA+EMhPAM4dac=")
		subID := "9ee2ec52-83c0-405e-a009-6636ead37acd"

		groupClient := resources.NewGroupsClient(subID)

		authorizer, err := auth.NewAuthorizerFromEnvironment()
		if err == nil {
			groupClient.Authorizer = authorizer
		}

		groupClient.CreateOrUpdate(
			context.Background(),
			"honcaogoclientrg1",
			resources.Group{
				Location: to.StringPtr("eastus"),
			},
		)
	*/

	/*
		// create a VirtualNetworks client
		vnetClient := network.NewVirtualNetworksClient(subID)

		// create an authorizer from env vars or Azure Managed Service Idenity
		authorizer, err := auth.NewAuthorizerFromEnvironment()
		if err == nil {
			vnetClient.Authorizer = authorizer
		}

		// call the VirtualNetworks CreateOrUpdate API
		vnetClient.CreateOrUpdate(context.Background(),
			"<resourceGroupName>",
			"<vnetName>",
			network.VirtualNetwork{
				Location: to.StringPtr("<azureRegion>"),
				VirtualNetworkPropertiesFormat: &network.VirtualNetworkPropertiesFormat{
					AddressSpace: &network.AddressSpace{
						AddressPrefixes: &[]string{"10.0.0.0/8"},
					},
					Subnets: &[]network.Subnet{
						{
							Name: to.StringPtr("<subnet1Name>"),
							SubnetPropertiesFormat: &network.SubnetPropertiesFormat{
								AddressPrefix: to.StringPtr("10.0.0.0/16"),
							},
						},
						{
							Name: to.StringPtr("<subnet2Name>"),
							SubnetPropertiesFormat: &network.SubnetPropertiesFormat{
								AddressPrefix: to.StringPtr("10.1.0.0/16"),
							},
						},
					},
				},
			})
	*/
}
