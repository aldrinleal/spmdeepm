# spm de epm

A Monitor I use for my lovely telco in colombia. Basically, it just issues a ping (to, say, your telco gateway) and reports it to CloudWatch. Thus, you're able to plot metrics and issue alarms. Nice, isn't it?

Even nicer: It works from a Raspberry Pi at home.

## Installation:

```shell
go get github.com/aldrinleal/spmdeepm
install -o root -m 6755 `which spmdeepm` /usr/s
```

Then find your gateway (hint: traceroute -n 8.8.8.8). And write a shell script to launch it via crontab e.g.:

```shell
#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export AWS_ACCESS_KEY_ID=<yourawskey>
export AWS_SECRET_ACCESS_KEY=<yourawssecret>

exec /usr/bin/spmdeepm $*
```

then on crontab:

```crontab
*/1 * * * * /home/pi/bin/spmdeepm.sh 200.24.33.103 2>&1 | xargs -i_ logger _
*/1 * * * * /home/pi/bin/spmdeepm.sh 8.8.8.8 2>&1 | xargs -i_ logger _
```

Works best if you have an RTC clock on your Pi. Or on your network. :)

