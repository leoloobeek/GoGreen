# GoGreen

## Introduction
This project was created to bring environmental (and HTTP) keying to scripting languages. As its common place to use PowerShell/JScript/VBScript as an initial vector of code execution, as a result of phishing or lateral movement, I see value of the techniques for these languages.

If you haven't heard of or used [Ebowla](https://github.com/Genetic-Malware/Ebowla) before, stop what your doing and follow that link. [Josh Pitts](https://twitter.com/midnite_runr) and [Travis Morrow](https://twitter.com/wired33) put together an awesome project that provides several environmental keying options. I'd recommend reading through the slides or watching the presentations on the Github page for Ebowla before moving on.

## GoGreen's Features

GoGreen offers environmental keying and HTTP keying to protect your payloads. 

- Environmental Keying
  - Encrypt payloads based on environment variables (username, dnshostname, architecture, etc.)
  - Encrypt payloads based on a certain file path which exists on the target (C:\users\username\Desktop\Some Link.lnk)
- HTTP Keying
  - Visits a web page, hashes the HTML response, uses that as the key
  - Taken from: https://cybersyndicates.com/2015/06/veil-evasion-aes-encrypted-httpkey-request-module/
  
At this time, Environmental Keying options are required, it is not possible to have HTTP keying without Environmental Keying options selected.

## Using GoGreen
GoGreen is a Golang project (sorry). The Release section has executables for Windows, Mac and Linux. To use the source go [here](https://golang.org/doc/install) and get setup with Go. Then you just need the following to download and use the project:
```
go get github.com/leoloobeek/GoGreen
cd $GOPATH/src/github.com/leoloobeek/GoGreen
```

GoGreen generates the final payload code based on templates, if you want to edit them (BYOO - Bring Your Own Obfuscation) the files are within the `data` folder. The default configuration file is within the root directory and examples of other configuration files are within the `examples` directory. To run you have two options:

Run with default config (in root directory): 
```
go run gogreen.go
```
Run with specific config:
```
go run gogreen.go -config path/to/config.json
```

## How GoGreen's Env Keying Works
As mentioned, GoGreen uses file path and/or environment variables as inputs for decrypting the payload. This is done by the payload by spidering/walking the directory tree from a start directory (such as C:\Program Files) and harvesting all file paths found. Next, it reads in all environment variables specified in the code (which is chosen by you in the config file), and creates all possible [combinations](https://en.wikipedia.org/wiki/Combination). Then goes the following steps:

1. Loop through all file paths
2. For each file path, loop through each combination of env vars
3. Build a key with: sha512(\[file path\]\[env var combo\])
4. Try first 32 chars of sha512 hash above as key
5. If error, or sha512(decrypted) doesn't match payload hash, move onto next key

## How GoGreen's HTTP Keying Works
If you plan to use HTTP Keying, set HttpKeyUrl to the URL of the site you control. The page you will use for the HTTP Keying must be up and running, and accessible from the machine you're running GoGreen on. 

It's important to note, that if HTTP Keying is used the following execution flow of the payload is as follows:
1. GoGreen uses environmental keying (directory/env vars) to decrypt and execute code
2. The code decrypted will handle the HTTP keying, which is executed
3. The new executed code will go out to HttpKeyUrl and hash the page
4. The new executed code will use the hash and decrypt and execute your payload

This flow allows us to hide the fact we will be making HTTP connections (and hide New-Object Net.Webclient, etc.).

#### A note about the HttpKey
Since there's a lot of factors here (different langs, charsets, etc.), GoGreen will write out a "tester file" when using HTTP Keying. Run this file on one of your systems to ensure JScript/VBScript/PowerShell will obtain the same key that Golang/GoGreen obtained. If there's special characters in the source (such as copyright symbol) there's a chance it won't match. If you are dead set on using that site, you can specify a custom HttpKey with `go run gogreen.go -httpkey <httpkey>`. 

To avoid headaches here, its recommended to use a simple web page. Even the standard Microsoft 403 page would be benign and something that would be easy to use.

## Deployment Considerations
If you have a StartDir of C:\, the payload will take a very long time to walk the directory tree and harvest all the paths. It will then take a very long time to try each and every path to decrypt. It will also make it harder for a defender to determine which file/folder path is used as the key. The same goes for using many env vars (e.g. 14 env vars = 16383 combos). So there are trade offs.

Here's a super basic example of a config on a Windows 10, Intel i5, 8GB RAM:
- C:\Program Files (x86) with depth of 2: 2687 file/folder paths
- All combos of 5 env vars: 31
- Total number of keys that will try and decrypt: 83,297
- Time: 1 minute and 3 seconds to run through all

## The Config
To keep things easy without requiring outside dependencies, the config files are JSON. Here are what the options are, what they mean, and more.

#### Language
This is the language of the payload. The options are JScript, VBScript, or PowerShell.

#### WSHAutoVersion
Either 'Yes' or 'No'. Only relevant for JScript and VBScript. Some COM objects used in the JScript/VBScript payloads won't work on Windows 10 (without .NETv2) without setting this to 'Yes'. Long story -- https://twitter.com/tiraniddo/status/864473926218522625

#### StartDir
This needs to be set to something to use a file path as part of your encryption key. If its blank, the final payload will remove all code pertaining to using file paths as the key. 

The StartDir will be where the payload will start spidering/walking the directory structure (recursively). Every file and folder path within the directory tree from StartDir will be harvested and tried as the decryption key. Obviously the further away the StartDir and file path is, the longer the script will take to execute. It will also be harder for the DFIR to find which file/folder path is actually used as the key.

#### Depth
If StartDir is specified, depth is how far down the tree it will go. If you leave it empty, it will assign a large depth of 100000. To give you an idea of how many file paths you can end up with, below are the numbers based on my personal Win10 system. This would be the number of decryption keys the payload will try, 124740 will take a while. (Definitely if coupled with env vars).
- 124740 (depth=0)
- 718 (depth=1)
- 2687 (depth=2)

#### PathKey
This is what will actually be used as part of the encryption key. It must fall under the StartDir and within the depth specified. If you leave it empty, and StartDir is still specified, it won't be used as part of the key and the payload will still walk the StartDir. This edge case can be used to attempt to throw off defenders.

#### EnvVars
A dictionary of key:value pairs of various environmental variables to include as keys. You can have as many or zero of these. To not use any env vars within the code, do: `"EnvVars": {}`. 

The payload will then pull all environment variables on the target when executed and create [combinations](https://en.wikipedia.org/wiki/Combination) of the variables (not permutations) and will be used as part of the key. If specifying an env variable but setting the value as blank, it will be included in all combinations but not actually used as part of the key. Again, this option is there to try to throw off defenders.

See examples/envonly-config.json to see what a Env Var only payload config would look like.

#### Payload and PayloadPath
If your payload is short, throw the code right into Payload. Otherwise specify the file containing your payload in PayloadPath.

#### MinusBytes
This is used in [Ebowla](https://github.com/Genetic-Malware/Ebowla) and a great idea. To ensure valid data is decrypted we have to hash the payload and using MinusBytes essentially hashes (payload - MinusBytes). If you're using the same exact payload but don't want the same hash string in the final source code, switch this up a bit.

#### HttpKeyUrl
The URL of the site you control to hash and use as the encryption key.

#### HttpKeyUA
The User Agent you want to use for the HTTP request from the target. Helpful if you are using [Apache ModRewrite rules](https://bluescreenofjeff.com/2016-04-12-combatting-incident-responders-with-apache-mod_rewrite/). 

#### HttpKeyRetry
The HTTP keying will sleep 30 seconds, and retry <HttpKeyRetry> times before exiting. This is helpful if you want to wait to deploy the page for your HTTP key.

# Thanks!
I just started this tool for my own use and to better understand environmental keying. Huge thanks to [Josh Pitts](https://twitter.com/midnite_runr) and [Travis Morrow](https://twitter.com/wired33) for the technique and there's a lot more options/capabilities with [Ebowla](https://github.com/Genetic-Malware/Ebowla). I took tons of the ideas and code from that project for this one. Check it out!

Also thanks to the following:
- James Forshaw [@tiraniddo](https://twitter.com/tiraniddo) as I took some code from https://github.com/tyranid/DotNetToJScript/
- Will Schroeder [@harmj0y](https://twitter.com/harmj0y) and whoever else wrote the Empire PowerShell agent code
- Alex Rymdeko-harvey [@Killswitch_GUI](https://twitter.com/Killswitch_GUI) and Chris Truncer [@christruncer](https://twitter.com/christruncer) for [HttpKey](https://cybersyndicates.com/2015/06/veil-evasion-aes-encrypted-httpkey-request-module/) idea
