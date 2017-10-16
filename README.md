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
GoGreen is a Golang project (sorry). Go [here](https://golang.org/doc/install) and get setup with Go. Then you just need the following to download and use the project:
```
go get github.com/leoloobeek/GoGreen
cd $GOPATH/src/github.com/leoloobeek/GoGreen
```

GoGreen generates the final payload code based on templates, if you want to edit them (BYOO - Bring Your Own Obfuscation :D) the files are within the `data` folder. The default configuration file is within the root directory and examples of other configuration files are within the `examples` directory. To run you have two options:

Run with default config (in root directory): 
```
go run gogreen.go
```
Run with specific config:
```
go run gogreen.go -config path/to/config.json
```
