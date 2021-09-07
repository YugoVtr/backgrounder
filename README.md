## requirements
- go 1.17
- gnome 40.4.0
- cron jobs 

## build
```bash
make build
```

## crontab on arch

```bash 
crontab -e 
```
```
*/1 * * * * $BACKGROUNDER/bin/backgrounder -config $BACKGROUNDER/settings.json
```

`$BACKGROUNDER` is path to project
