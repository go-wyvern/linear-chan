## linear-chan
    
## linear-chan is a tool for managing company Automated Deployment..

#### Author 
    
    name  :  wyvern.wu.
    email :  wyvern.wu@aliyun.com

#### Installation

    go get github.com/go-wyvern/linear-chan
    go install github.com/go-wyvern/linear-chan
    
#### Usage:

    linear-chan command [arguments]

##### The commands are:

    version                        show the linear-chan version
    init                           add config to /etc/linear-chan.conf
    up [projectname]               Automated Deployment server
    auth [servername]              Send public key to server
    ssh [servername]               SSH to server
    send [filenames] [servername]  send files to server
    create [projectname]           read /etc/linear-chan.conf and add tag(project) to mysql
    delete [projectname]           delete tag(project) from mysql
    update [projectname]           read /etc/linear-chan.conf and update tag(project) to mysql
    ls [servername]                list server message

#### If need config file pls use -f config file path in the End

#### Use "linear-chan help [command]" for more information about a command.

#### Additional help topics:

    version                        show the linear-chan version
    init                           add config to /etc/linear-chan.conf
    up [projectname]               Automated Deployment server
    auth [servername]              Send public key to server
    ssh [servername]               SSH to server
    send [filenames] [servername]  send files to server
    create [projectname]           read /etc/linear-chan.conf and add tag(project) to mysql
    delete [projectname]           delete tag(project) from mysql
    update [projectname]           read /etc/linear-chan.conf and update tag(project) to mysql
    ls [servername]                list server message

#### Use "linear-chan help [topic]" for more information about that topic.

#### Notice

    first use please [linear-chan init] first, it will create config files in /etc/linear-chan.ini and /etc/linear-chan.d/projects.ini
