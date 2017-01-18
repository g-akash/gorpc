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
)


type Dict_Word struct{
	Word string
	Meaning string
	//list of words
	Synonyms []*Dict_Word
	//enum is to be used here for type of the word
	Word_type string
}

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
	dictionary map[string]*Dict_Word
	lock sync.Mutex
	}


var dictionary Dictionary

func otherErrorFound(args *Insert_args)bool{
	return false
}

func remove_synonyms(dict *Dictionary,syn string, word string){
	return
}


func (dictionary *Dictionary)InsertWord(args *Insert_args, reply *string) error {
	//store the arguments into appropriate variables
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
	return err
}


func (dictionary *Dictionary)RemoveWord(args *Remove_args, reply *string) error {
	//store the args in appropriate variables
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
	return err

}

func (dictionary * Dictionary)LookupWord(args *Lookup_args, reply *Insert_args) error {
	//store the args in appropriate variables
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
	return err

}


func http_listener(){
	server := rpc.NewServer()
	server.RegisterName("Dictionary",&dictionary)
	server.HandleHTTP("/","/debug")

	for true{
		l, e := net.Listen("tcp", ":6000")
		//fmt.Println(l,e)
		if e != nil {
			log.Fatal("listen error:", e)
		}

		http.Serve(l, nil)

	}

}

func tcp_listener(){
	server := rpc.NewServer()
	server.RegisterName("Dictionary",&dictionary)
	//do something here so that we can accept both http and tcp connections

	for true{
		l,e := net.Listen("tcp",":5000")
		//fmt.Println(l,e)
		if e!=nil{
			log.Fatal("there was an error in listening for connection", e)
		}
		server.Accept(l)
	}

}


func main() {
	dictionary = *new(Dictionary)
	dictionary.dictionary = make(map[string]*Dict_Word)
	fmt.Println(strings.Split("Server is up and running"," "))
	//start the two threads here
	go http_listener()
	go tcp_listener()
	for true{
		time.Sleep(time.Second*5)
	}
	

	// var reply string
	// var synonyms []string
	// var word Insert_args
	// word = *new(Insert_args)
	// var iarg Insert_args
	// iarg = Insert_args{word:"hello",meaning:"ohkay",synonyms:synonyms,word_type:"noun"}
	// var err error
	// dictionary.InsertWord(&iarg,&reply)
	// fmt.Println(word)
	// synonyms=append(synonyms,"hello")
	// iarg = Insert_args{word:"new_word",meaning:"ohkay",synonyms:synonyms,word_type:"noun"}
	// err=dictionary.InsertWord(&iarg,&reply)
	// fmt.Println(*dictionary.dictionary["hello"],*dictionary.dictionary["new_word"])
	// var rarg Remove_args
	// rarg = Remove_args{word:"hello"}
	// err=dictionary.RemoveWord(&rarg,&reply)
	// if err!=nil{
	// 	fmt.Println(err)
	// }
	// fmt.Println(reply+"jjjjjjj")
	// var larg Lookup_args
	// larg = Lookup_args{word:"new_word"}
	// err=dictionary.LookupWord(&larg,&word)
	// if err!=nil{
	// 	fmt.Println(err)
	// }
	// fmt.Println(word)
	// fmt.Println(dictionary)
}



