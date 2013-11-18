=======
Espresso Webserver
========

Espresso is a small, lightweight webserver that parses pages written in the Lua Scripting Language. It's functionality is similiar to PHP. It is built by Eric C. and River B. We don't really maintain this anymore, but the source is availiable for
anyone who wishes to make use of it.

## Building

### Requirements

Requirements:
	Lua 5.2 Library

### Win32

Don't even bother. GoLua does not currently support Windows Targets.

### UNIX (including Linux and OSX)

```
cd Espresso
export GOPATH=$PWD
go build
./Espresso
```

### CUDA

Why the hell are you building Espresso on CUDA? 

##Running Espresso

To run Espresso, type the following command
./Espresso

The server will bind to port 4000 by default.
To change this, edit config.json. Note that you
cannot bind to port 80 without running it at root.
That's not a very good idea by the way.

##Using Espresso

You can browse the api.go file or the webroot folder to get an idea of how the functions work.

##How to prepare an egg

First, open the egg containment unit and withdraw a single egg. Be careful not to withdraw more then
one egg as they could cause confusion. Inspect the egg for any signs of tampering. It is important to
note that everyone is out to get you. Do not trust anyone, not even the egg.

Turn on your stove. Check if your stove is on by placing your hand above the stove. Do not directly touch the stove. When you sense the temperature is exactly 400.7K, place a pan on the stove, making sure the open end is facing upwards. Wait for 30.4 seconds. Lifting the egg directly above the panside, use two newtons of force to breach it's protective shielding. If
two newtons fail to compromise the integrity of the shell, a deadblow hammer will suffice.

Using a spoon, scoop up the contents of the egg and ingest. Throw the shells in the pan and let simmer for five minutes.
Do not eat the shells. Dispose of the stove and clean up the pan, wondering how bored a college student had to be to
write this all up in the middle of class.

