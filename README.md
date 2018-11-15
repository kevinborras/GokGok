# GokGok
A tool to check domains

## Requirements

+ Nmap from [Nmap Official Webpage](https://nmap.org/)

***

## Dependencies
```bash
github.com/fatih/color
github.com/integrii/flaggy
github.com/tomsteele/go-nmap
```
## Usage

```bash
$ go run gokgok.go -h
Gok-Gok - Is something there?

  Flags:
    -h --help  Displays help with available flag, subcommand, and positional value parameters.
    -v --version  Print version
    -t --targetList  File with targets to be checked
    -s --scanThem  Scan the the targets with Nmap
```
+ Check if hosts are up or down

```bash
$ go run gokgok.go -t domainTest.txt
 [+] SUCCESS:  Host scanme.nmap.org is alive with 200 OK
 [+] SUCCESS:  Host google.com is alive with 200 OK
 [+] SUCCESS:  Host wikipedia.org is alive with 200 OK
 [-] ERROR:  No such host or it's down: icsdfsdfsd.es
 [-] ERROR:  No such host or it's down: werwerwer.com
```

+ Launch nmap over alive hosts (*Scans will be saved in the nmapResults directory *)

```bash
$ go run gokgok.go -t domainTest.txt -s
 [+] SUCCESS:  Host scanme.nmap.org is alive with 200 OK
 [-] ERROR:  No such host or it's down: icsdfsdfsd.es
 [-] ERROR:  No such host or it's down: werwerwer.com
 [i] INFO:  Executing command nmap -sV -P0 45.33.32.156 -oX result_45.33.32.156.xml
 [+] SUCCESS:  Scan done
```

## TODO

+ Parse nmap files
+ Display the results
