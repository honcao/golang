package main

import (
	"errors"
	"fmt"
	//"github.com/honcao/golang/pkg/apiv2"
	"math/rand"
	"reflect"
	"sync"
	"time"

	"github.com/honcao/golang/pkg/apiv1"
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

func DeepAssignmentWrapper(dst, src interface{}) (err error) {

	dstValue := reflect.ValueOf(dst)
	srcValue := reflect.ValueOf(src)
	if dstValue.Kind() != reflect.Ptr {
		return errors.New("dst is not pointer to stuct")
	}
	dstValue = dstValue.Elem()
	if dstValue.Kind() != reflect.Struct {
		return errors.New("dst is not pointer to stuct")
	}
	if srcValue.Kind() != reflect.Struct {
		return errors.New("src is not pointer to stuct")
	}
	DeepAssignment(dstValue, srcValue, 0, "")
	return nil
}

func DeepAssignment(dstValue, srcValue reflect.Value, depth int, path string) (err error) {

	switch srcValue.Kind() {
	case reflect.String:
		dstValue.SetString(srcValue.String())
	case reflect.Int:
		dstValue.SetInt(srcValue.Int())
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
	return nil
}

type myString string

func main() {

	src := apiv1.TypeB{
		Name: "abc",
		Age:  15,
	}
	dst := apiv1.TypeB{}
	err := DeepAssignmentWrapper(&dst, src)
	fmt.Println(err)
	fmt.Println(dst)

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
