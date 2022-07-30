# Blocky Logs

Blocky Logs is a utility to show logs from [blocky DNS](https://0xERR0R.github.io/blocky/)

At this moment it's just working if you logging to MySQL database.

## Features

- **Show by event type**
  - By default its show events of type "BLOCKED (porn)"
  - To show all events user --event "all"
  
- **Show events by IP**
  - Show all logs from last 24 hours from specified IP addresses
  
## Quick start

If you run blocky-log it's will generate a sample configuration file.

The configuration file accepts multiple hosts.

Configuration file can be placed in `/etc` directory, current directory or in `~/etc` directory.

```YAML
hours: 24
blockys:
  - database:
      database: blocky
      host: 192.168.1.1
      password: yourpassword
      port: 3306
      type: mysql
      username: alexandre
    name: Home
```

## Output example

<p align="center">
  <img height="600" src="https://github.com/alexandrepossebom/blocky-log/raw/main/docs/sample.jpg">
</p>

## Contribution

Issues, feature suggestions and pull requests are welcome!

## Disclaimer

Is better badly done than not done.
It's not perfect but it's working.
