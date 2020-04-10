# go-crond Openshift Ready!

A cron daemon written in golang



## Features

- system crontab (with username inside)
- user crontabs (without username inside)
- run-parts support
- Logging to STDOUT and STDERR
- Keep current environment (eg. for usage in Docker containers)
- Supports Linux, MacOS, ARM/ARM64 (Rasbperry Pi and others)
- Email notification

## Usage

```
Usage:
  go-crond

Application Options:
      --threads=            Number of parallel executions (default: 20)
      --configfile=         Include configfile email notifications
      --default-user=       Default user (default: root)
      --include=            Include files in directory as system crontabs (with user)
      --allow-unprivileged  Allow daemon to run as non root (unprivileged) user
  -v, --verbose             verbose mode
  -h, --help                show this help message
```

### Examples

Run crond with a system crontab:

    go-crond examples/crontab


Run crond with user crontabs (without user in it) under specific users:

    go-crond \
        root:examples/crontab-root \ 
        guest:examples/crontab-guest


Run crond with auto include of /etc/cron.d and script execution of hourly, weekly, daily and monthly:

    go-crond \
        --include=/etc/cron.d \
        --run-parts-hourly=/etc/cron.hourly \
        --run-parts-weekly=/etc/cron.weekly \
        --run-parts-daily=/etc/cron.daily \
        --run-parts-monthly=/etc/cron.monthly

Run crond with run-parts with custom time spec:

    go-crond \
        --run-parts=1m:/etc/cron.minute \
        --run-parts=15m:/etc/cron.15min

Run crond with run-parts with custom time spec and different user:

    go-crond \
        --run-parts=1m:application:/etc/cron.minute \
        --run-parts=15m:admin:/etc/cron.15min
        
        
## License
Open-sourced software licensed under the MIT license.
