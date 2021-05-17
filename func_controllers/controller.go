package controller

import (
	"fmt"
	"github.com/openfaas/faas-cli/proxy"
	"io/ioutil"
	"github.com/Jeffail/gabs"
	"os"
	"strings"

	//"net/http"

)

func InitState(nameJson string, gateway string, name string, bytesIn *[]byte, contentType string, query []string, headers []string, async bool, httpMethod string, tlsInsecure bool, namespace string) error{
	js, _ := ioutil.ReadFile(nameJson)
	jsonParsed, err := gabs.ParseJSON(js)
	fmt.Println(nameJson)
	fmt.Println(string(*bytesIn))
	if err != nil {
		return err
	}
	input := string(*bytesIn)
	input = strings.Trim(input, "\n")
	startFunc := jsonParsed.Path("StartFunction").Data().(string)
	fmt.Println(startFunc)
	fmt.Printf("input : %s\n",input)
	states, err := jsonParsed.S("States").ChildrenMap()
	if err != nil {
		return err
	}
	nextFunc := startFunc
	count:= 0
	response := bytesIn
	for {
		fmt.Printf("/////// %v //////////\n", count)
		funcJson, err := gabs.ParseJSON(states[nextFunc].Bytes())
		if err != nil {
			return err
		}
		currFunc := nextFunc
		fmt.Printf("/////// %s ///////////\n", currFunc)
		if funcJson.Path("Type").Data().(string) == "Fail"{
			return nil
		}
		if funcJson.Exists("ResultPath") {

			response, err = proxy.InvokeFunction(gateway, currFunc, response, contentType, query, headers, async, httpMethod, tlsInsecure, namespace)
			//*response = []byte(strings.Trim(string(*response), "\\n"))
		} else{
			response, err = proxy.InvokeFunction(gateway, currFunc, bytesIn, contentType, query, headers, async, httpMethod, tlsInsecure, namespace)
		}
		if err != nil {
			os.Stdout.Write([]byte(err.Error()))
			if funcJson.Exists("Catch") {
				nextFunc, err = Catch(funcJson)
				if err != nil{
					return fmt.Errorf("Wrong format json file: catch\n")
				}
			}else{
				return err
			}
		} else {
			fmt.Printf("///////else %s //////////\n", currFunc)
			if response != nil {
				os.Stdout.Write(*response)
			}
			//fmt.Printf("current func: %v, ", currFunc)
			funcType := funcJson.Path("Type").Data().(string)
			if funcType == "Choice" {
				nextFunc, err = Choice(funcJson, input)
				if err != nil {
					panic(err)
				}
			} else if funcType == "Task" {
				nextFunc, err = Task(funcJson)
				if err != nil{
					return err
				}
				if nextFunc == "" {
					break
				}
			}
		}
		count++
	}
	return nil
}
