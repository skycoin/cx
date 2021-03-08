this is first draft wrapper libary for cx in go Programmning Language.

Steps for writting wrapper libary for cx in go Programmning Language.


use top to down approach for writting package i.e start with Opcode.

1. The first step is add package entry into package varialble CorePackages so cx can import.

2. define Opcode for all the import variable in init fucntion.

for example : 	Op(OP_HTTP_SERVE, "http.Serve", opHTTPServe, In(ASTR), Out(ASTR))


opcode used to define statement of experssion.

this will define the opcode which we can use in cx program 
it also define the imput and output parameter.

3. create new file with op_packagename.go

a. In this file in init function define all the init function which will cerate package with MakePackage.

b.define all the vaiable which you want to use as struct also it substruct.


4. Define all the method which you want to use in program.


note always use base type when you are calling go Programmning Language method instead of nterface.

example : https://github.com/skycoin/cx/blob/develop/cx/op_http.go
