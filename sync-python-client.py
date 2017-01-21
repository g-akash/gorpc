import requests
import json
import sys
import socket




class WordType:
	Noun = 1
	Verb = 2
	Adjective = 3

class Insert_args(object):
	def __init__(self,word,meaning,synonyms,word_type):
		self.Word = word
		self.Meaning = meaning
		self.Synonyms = synonyms
		self.Word_type = word_type

class Remove_args(object):
	def __init__(self,word):
		self.Word = word


class Lookup_args(object):
	def __init__(self,word):
		self.Word = word


class Dictionary:
	def __init__(self,url,headers):
		self.url = url
		self.headers = headers

	def InsertWord(self,word,meaning,synonyms,Type):
		insert_arg = Insert_args(word,meaning,synonyms,Type)
		reply = ""
		payload = {
			"method": "Dictionary.InsertWord",
			"params": [insert_arg.__dict__,reply],
			"jsonrpc": "2.0",
			"id": 0,
		}
		try:
			s=socket.socket(socket.AF_INET,socket.SOCK_STREAM)
			s.connect((self.url,5000))
			s.send(json.dumps(payload))
			
			err = s.recv(10000)
		except requests.exceptions.RequestException as e:
			print e
			return
		err = json.loads(err)
		reply = err['result']
		err = err['error']
		if err is not None:
			print err
		else:
			print reply

	def RemoveWord(self,word):
		remove_arg = Remove_args(word)
		reply = ""
		payload = {
			"method":"Dictionary.RemoveWord",
			"params":[remove_arg.__dict__,reply],
			"jsonrpc":"2.0",
			"id": 0,
		}
		try:
			s=socket.socket(socket.AF_INET,socket.SOCK_STREAM)
			s.connect((self.url,5000))
			s.send(json.dumps(payload))
			
			err = s.recv(10000)
		except requests.exceptions.RequestException as e:
			print e
			return
		err = json.loads(err)
		reply = err['result']
		err = err['error']
		if err is not None:
			print err
		else:
			print reply
	

	def LookupWord(self,word):
		lookup_arg = Lookup_args(word)
		reply = Insert_args("","",[],"")
		payload = {
			"method":"Dictionary.LookupWord",
			"params":[lookup_arg.__dict__,reply.__dict__],
			"jsonrpc":"2.0",
			"id": 0,
		}
		try:
			s=socket.socket(socket.AF_INET,socket.SOCK_STREAM)
			s.connect((self.url,5000))
			s.send(json.dumps(payload))
			
			err = s.recv(10000)
		except requests.exceptions.RequestException as e:
			print e
			return
		err = json.loads(err)
		reply = err['result']
		err = err['error']
		if err is not None:
			print err
		else:
			print reply['Word'],reply['Meaning'],json.dumps(reply['Synonyms']),
			if reply['Word_type']==1:
				print "Noun"
			elif reply['Word_type']==2:
				print "Verb"
			elif reply['Word_type']==3:
				print "Adjective"



def main():
	addr = sys.argv[1]
	dictionary = Dictionary(addr,{'User-Agent': 'Mozilla/5.0 (Windows NT 6.0; WOW64; rv:24.0) Gecko/20100101 Firefox/24.0','content-type': 'application/json'})
	while True:
		print "What is the operation you would like to perform? Choose one of these, add, delete, lookup"
		x=raw_input()
		if x=="add":
			word = raw_input()
			meaning = raw_input()
			synonyms = raw_input().split()
			string_type = raw_input()
			if string_type == "Noun":
				Type = 1
			elif string_type == "Verb":
				Type = 2
			elif string_type == "Adjective":
				Type = 3
			else:
				print "Wrong type given"
				continue
			dictionary.InsertWord(word,meaning,synonyms,Type)

		elif x=="delete":
			word = raw_input()
			dictionary.RemoveWord(word)
		elif x=="lookup":
			word = raw_input()
			dictionary.LookupWord(word)
		else:
			print "Could not recognize command"

if __name__ == "__main__":
	main()





