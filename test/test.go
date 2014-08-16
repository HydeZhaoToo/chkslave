package test

import "log"

func check(err error){
	if err != nil{
		log.Fatalln(err)
	}
}
