# ChitChat
Chitchat mandatory 3

How to run the program

Step 1:
CD in to the Server folder from your CLI and type $ go run server.go

Step 2: 
Open a second terminal where you CD into the client folder. From there you type $ go run client.go

Step 3:
You can now either start up multiple terminals to emulate a chat server.

Step 4: To exit the connection from client-side, type '.exit' in your terminal. To shutdown the server, type '.shutdown' in your terminal.


Tips & Tricks:
If you want to join the server from multiple computers on the same network, you have to change your target ip to  <Wifi/Internet ip-address>:5050
and also change the server ip listener from "localhost:5050" -> 0.0.0.0:5050

Have fun chatting!