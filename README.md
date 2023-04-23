# Is-connected

Checks internet status. If connection is lost once it's back send's you a notification via discord if that time was greater than 1 minute.

## Deamon

```=bash
sudo cp is-connect.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl start is-connect
```