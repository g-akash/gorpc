/**
 *	author		:	akash garg
 *	roll no.		130070060
 *	
 *	synchronous client for rpc in go
 */


 package main


import (
 	"fmt"
 	"os"
 	"net"
 	"net/rpc"
 	//"errors"
 	"strings"
 	"bufio"
 ) 




type Insert_args struct{
	Word string
	Meaning string
	//string list and type to be defined here
	Synonyms []string
	Word_type string
}

type Remove_args struct{
	Word string
	//it could be entire word here
}

type Lookup_args struct{
	Word string
}
type Dictionary struct{
	client *rpc.Client
}

func (dictionary *Dictionary)InsertWord(word string, meaning string, synonyms []string, Type string){
	insert_arg := Insert_args{Word: word, Meaning: meaning, Synonyms: synonyms, Word_type: Type}
	reply := new(string)
	err := dictionary.client.Call("Dictionary.InsertWord",insert_arg,reply)
	if err!=nil{
		fmt.Println(err)
	} else {
		fmt.Println(*reply)
	}
	//do the reply work here

}

func (dictionary *Dictionary)RemoveWord(word string){
	remove_arg := Remove_args{Word:word}
	reply := new(string)
	err := dictionary.client.Call("Dictionary.RemoveWord",remove_arg,reply)
	if err != nil{
		fmt.Println(err)
	} else {
		fmt.Println(*reply)
	}
}

func (dictionary *Dictionary) LookupWord(word string){
	lookup_arg := Lookup_args{Word: word}
	reply := new(Insert_args)
	err := dictionary.client.Call("Dictionary.LookupWord",lookup_arg,reply)
	if err != nil{
		fmt.Println(err)
	} else {
		fmt.Println(*reply)
	}
}




func main(){
	addr := os.Args[1]+":5000"
	fmt.Println(addr)
	client ,err := net.Dial("tcp",addr)
	if err != nil{
		fmt.Println(err)
		return
	}
	var dictionary = Dictionary{client:rpc.NewClient(client)}
	
	scanner := bufio.NewScanner(os.Stdin)
	for true{
		fmt.Println("What is the operation you would like to perform? Choose one of these, add, delete, lookup")
		scanner.Scan()
		command := scanner.Text()
		if command == "add"{
			scanner.Scan()
			word := scanner.Text()
			scanner.Scan()
			meaning := scanner.Text()
			scanner.Scan()
			syns := scanner.Text()
			synonyms := make([]string,0)
			if syns !=""{
				synonyms = strings.Split(syns," ")
			} 
			
			scanner.Scan()
			Type := scanner.Text()
			fmt.Println(word,meaning,synonyms,Type)
			dictionary.InsertWord(word,meaning,synonyms,Type)


		}	else if command == "delete"{
			scanner.Scan()
			word := scanner.Text()
			dictionary.RemoveWord(word)
			
		}	else if command == "lookup"{
			scanner.Scan()
			word := scanner.Text()
			dictionary.LookupWord(word)
			
		}	else{
			fmt.Println("Could not recognize command.")
		}
	}
}