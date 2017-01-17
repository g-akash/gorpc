/**
 *	author		:	akash garg
 *	roll no.		130070060
 *	
 *	synchronous client for roc in python
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


type Word struct{
	word string
	meaning string
	//list of words
	//enum is to be used here for type of the word
	word_type string
}

type Insert_args struct{
	word string
	meaning string
	//string list and type to be defined here
	synonyms []string
	Type string
}

type Remove_args struct{
	word string
	//it could be entire word here
}

type Lookup_args struct{
	word string
}


type Dictionary struct{
	client *rpc.Client
}

func (dictionary *Dictionary)InsertWord(word string, meaning string, synonyms []string, Type string){
	insert_arg := Insert_args{word: word, meaning: meaning, synonyms: synonyms, Type: Type}
	var reply string
	err := dictionary.client.Call("Dictionary.InsertWord",insert_arg,&reply)
	if err!=nil{
		log.Fatal("error in contacting server",err)
	}
	//do the reply work here

}




func main(){
	var dictionary = New(Dictionary)
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
			synonyms := strings.Split(scanner.Text()," ")
			scanner.Scan()
			Type := scanner.Text()
			
			dictionary.InsertWord(word,meaning,synonyms,Type)


		}	else if command == "delete"{
			scanner.Scan()
			word := scanner.Text()
			remove_arg := Remove_args{word:word}
			fmt.Println(remove_arg)
		}	else if command == "lookup"{
			scanner.Scan()
			word := scanner.Text()
			lookup_arg := Lookup_args{word:word}
			fmt.Println(lookup_arg)
		}	else{
			fmt.Println("Could not recognize command.")
		}
	}
}