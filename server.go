/**
*	author		:	akash garg
*	roll no.	:	130070060
*	
*	server side code for rpc
*/	

package main

import(
	"errors"
	"log"
	"net"
	//"os"
	"fmt"
	"strings"
	"sync"
	"net/rpc"
	"net/http"
	"time"
	"net/rpc/jsonrpc"
)

type WordType uint8

const (
	Noun WordType = iota + 1
	Verb
	Adjective
)


type Dict_Word struct{
	Word string
	Meaning string
	//list of words
	Synonyms []*Dict_Word
	//enum is to be used here for type of the word
	Word_type WordType
}

type Insert_args struct{
	Word string
	Meaning string
	//string list and type to be defined here
	Synonyms []string
	Word_type WordType
}

type Remove_args struct{
	Word string
	//it could be entire word here
}

type Lookup_args struct{
	Word string
}

type Dictionary struct{
	dictionary map[string]*Dict_Word
	lock sync.Mutex
	}


var dictionary Dictionary

func otherErrorFound(args *Insert_args)bool{
	if args.Meaning=="" || args.Word == ""{
		return true
	}
	if args.Word_type<=0 || args.Word_type>3{
		return true
	}
	
	return false
}

func remove_synonyms(dict *Dictionary,word string, syn string){
	var new_synonyms []*Dict_Word
	for i:=0; i<len(dict.dictionary[word].Synonyms);i++{
		if dict.dictionary[word].Synonyms[i].Word!=syn{
			new_synonyms=append(new_synonyms,dict.dictionary[word].Synonyms[i])
		}
	}
	dict.dictionary[word].Synonyms=new_synonyms
	return
}


func (dictionary *Dictionary)InsertWord(args *Insert_args, reply *string) error {
	//store the arguments into appropriate variables
	fmt.Println("Request arrived for inserting new word ",*args)
	word := args.Word
	meaning := args.Meaning
	synonyms := args.Synonyms
	word_type:= args.Word_type
	var err error
	err = nil
	//at this point we have stored all the args into the appropriate variables.
	//get a lock
	dictionary.lock.Lock()
	_,ok := dictionary.dictionary[word]
	if ok{
		err = errors.New("AlreadyExists")
	} else if otherErrorFound(args) {
		err = errors.New("OtherServerSideError")
	} else {
		
		err = nil
		for i:=0;i<len(synonyms);i++{
			_,ok := dictionary.dictionary[synonyms[i]]
			if !ok{
				err = errors.New("OtherServerSideError")
				
				break;
			}
		}
		
		if err == nil{
			dictionary.dictionary[word] = &Dict_Word{Word:word,Meaning:meaning,Synonyms:nil,Word_type:word_type}
			var syn_list []*Dict_Word
			//found,ok = dictionary.dictionary[word]
			for i:=0;i<len(synonyms);i++{
				syns:=dictionary.dictionary[synonyms[i]]
				syns.Synonyms=append(dictionary.dictionary[synonyms[i]].Synonyms,dictionary.dictionary[word])
				dictionary.dictionary[synonyms[i]]=syns
				//syn_word,_ := dictionary.dictionary[synonyms[i]]
				syn_list=append(syn_list,dictionary.dictionary[synonyms[i]])
				}
			dictionary.dictionary[word] = &Dict_Word{Word:word,Meaning:meaning,Synonyms:syn_list,Word_type:word_type}

			*reply = "Success Inserting "+word
			}
		}

		
	
	//release the lock
	dictionary.lock.Unlock()
	if err!=nil{
		fmt.Println("Final result of the process ",err)
	} else {
		fmt.Println("Final result of the process ", *reply)
	}
	return err
}


func (dictionary *Dictionary)RemoveWord(args *Remove_args, reply *string) error {
	//store the args in appropriate variables
	fmt.Println("Request arrived for removing ",*args)
	word := args.Word
	var err error
	//get a lock
	dictionary.lock.Lock()

	//do the removal
	found,ok := dictionary.dictionary[word]
	if !ok{
		err = errors.New("UnknownWord")

	} else {
		synonyms := found.Synonyms
		synonyms_exist := true
		for i:=0; i<len(found.Synonyms);i++{
			_,ok:=dictionary.dictionary[found.Synonyms[i].Word]
			synonyms_exist = synonyms_exist && ok
		}
		if !synonyms_exist{
			err = errors.New("OtherServerSideError")
		} else {
			for i:=0; i<len(found.Synonyms);i++{
				remove_synonyms(dictionary,synonyms[i].Word,word)
			}
			delete(dictionary.dictionary,word)

		}
	}
	//sleep for sometime

	time.Sleep(time.Second*5)
	//release the lock
	//word = "Success Removing "+word
	*reply = "Success Removing "+word
	dictionary.lock.Unlock()
	if err!=nil{
		fmt.Println("Final result of the process ",err)
	} else {
		fmt.Println("Final result of the process ", *reply)
	}

	return err

}

func (dictionary * Dictionary)LookupWord(args *Lookup_args, reply *Insert_args) error {
	//store the args in appropriate variables
	fmt.Println("Request arrived for finding ",*args)
	word := args.Word
	var err error
	//get a lock
	dictionary.lock.Lock()

	dict_word,ok := dictionary.dictionary[word]
	if ok{
		reply.Word = dict_word.Word
		reply.Meaning = dict_word.Meaning
		reply.Word_type = dict_word.Word_type
		for i:=0;i<len(dict_word.Synonyms);i++{
			reply.Synonyms=append(reply.Synonyms,dict_word.Synonyms[i].Word)
		}
		err = nil
	} else {
		err = errors.New("UnknownWord")
		reply = nil
	}
	//do the work


	//release the lock
	dictionary.lock.Unlock()
	if err!=nil{
		fmt.Println("Final result of the process ",err)
	} else {
		fmt.Println("Final result of the process ", *reply)
	}

	return err

}


func http_listener(){
	//fmt.Println(rpc.DefaultRPCPath,"  ",rpc.DefaultDebugPath)
	server := rpc.NewServer()
	server.RegisterName("Dictionary",&dictionary)
	server.HandleHTTP(rpc.DefaultRPCPath,rpc.DefaultDebugPath)

	for true{
		l, e := net.Listen("tcp", ":6000")
		//fmt.Println(l,e)
		if e != nil {
			log.Fatal("there was an error in listening for http connection on port 6000:", e)
			return
		} else {
		fmt.Println("Started listening for new http connections.")
		}
			
		err := http.Serve(l, nil)
		if err!=nil{
			fmt.Println("Error serving connection.")
			continue
		}

		fmt.Println("Serving new connection.")
	}

}

func tcp_listener(){
	server := rpc.NewServer()
	server.RegisterName("Dictionary",&dictionary)
	//do something here so that we can accept both http and tcp connections

	
	l,e := net.Listen("tcp",":5000")
	//fmt.Println(l,e)
	if e!=nil{
		log.Fatal("there was an error in listening for tcp connection on port 5000:", e)
		return
	} else {
		fmt.Println("Started listening for new tcp connections.")
	}
	for true{

		conn, err := l.Accept()
		fmt.Println("Accepted a new tcp connection")
		if err != nil{
			fmt.Println("Cannot serve this connection. Will wait for new connection.")
			continue
		}
		fmt.Println("Serving new connection.")
		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
	

}


func main() {
	dictionary = *new(Dictionary)
	dictionary.dictionary = make(map[string]*Dict_Word)
	fmt.Println(strings.Split("Starting the server"," "))
	//start the two threads here
	go http_listener()
	go tcp_listener()
	for true{
		time.Sleep(time.Second*5)
	}
	
}



