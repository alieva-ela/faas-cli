package controller

import (
	"fmt"
	"github.com/Jeffail/gabs"
)

func Catch (funcJson *gabs.Container) (string, error) {
	children, err := funcJson.S("Catch").Children()
	if err != nil {
		return "", err
	}
	var nextFunc = children[0].Path("Next").Data().(string)

	return nextFunc, nil

}

func Choice(funcJson *gabs.Container, input string) (string, error){
	children, err := funcJson.S("Choices").Children()
	if err != nil {
		return "", err
	}
	wasFound := false
	var nextFunc = ""
	// в зависимости от входных данных (если успешно завершилась) смотрим next
	for _, child := range children {
		if input == child.Path("StringEquals").Data().(string) {
			nextFunc, wasFound = child.Path("Next").Data().(string), true
			fmt.Printf("next func: %v\n", nextFunc)
			break
		}
	}
	if wasFound{
		return nextFunc, nil
	}else{
		return "", err
	}
}

func Task(funcJson *gabs.Container) (string, error){
		// выполняем функцию currFunc и в зависимости от завершения смотрим catch

		// если успешно завершилась смотрим next
		if funcJson.Exists("End") {
			//fmt.Println("Exit")
			return "", nil
		} else if  funcJson.Exists("Next") {
			nextFunc := funcJson.Path("Next").Data().(string)
			//fmt.Printf("next func: %v\n", nextFunc)
			return nextFunc ,nil
		}

		return "", fmt.Errorf("Wrong format json file")
}