# go-echo
Help to identify HTTP routing issues in Cloud Foundry.

## How does it work?

Sends a HTTP POST request from ```client``` through to the ```backend app``` and expects the same result to be returned:

```client``` <-> ```frontend app``` <-> ```backend app```

If the `client` does not receive the same data it POSTed then there is a potential issue.

## Setup
- Push the backend app:
```
cf push echo-backend -i 10 -m 64M -b go_buildpack
```

- Push the frontend app with no-start:
```
cf push echo-frontend -i 10 -m 64M -b go_buildpack --no-start
```

- Set the env variable for the frontend app that points to URL of backend app:
```
cf set-env echo-frontend BACKEND_URL http://echo-backend.example.com
```

- Start the frontend app:
```
cf start echo-frontend
```

- Verify everything is working:
```
curl -s echo-frontend.example.com/frontend -d "ECHO"

// You can replace "ECHO" with whatever string you want.
// The app will return whatever is sent to it.
```

- Response should look like:
```
2017-08-10T09:19:30 "echo-frontend.example.com/frontend" | 200 | 2017-08-10T09:19:30 "echo-backend.example.com/backend" ECHO
```

- Set the frontend URL (needed for ```client.sh```):
```
export FRONTEND_URL=http://echo-frontend.example.com
```

- Kick off the script:
```
nohup ./client.sh &
```

- Monitor the log file:
```
tail -f client.log
```
 
## Example of potential routing issue:
```
2017-08-10T09:19:30 "echo-frontend.example.com/frontend" | 200 | 2017-08-10T09:19:30 "echo-backend.example.com/backend" PIVOTAL
2017-08-10T09:19:32 "echo-frontend.example.com/frontend" | 200 | 2017-08-10T09:19:32 "echo-backend.example.com/backend" PIVOTAL
2017-08-10T09:19:33 "echo-frontend.example.com/frontend" | 200 | 2017-08-10T09:19:33 "echo-backend.example.com/backend" PIVOTAL
2017-08-10T09:19:35 "echo-frontend.example.com/frontend" | 404 | <!DOCTYPE HTML PUBLIC "-//IETF//DTD HTML 2.0//EN">
<html><head>
<title>404 Not Found</title>
</head><body>
<h1>Not Found</h1>
<p>The requested URL /backend was not found on this server.</p>
</body></html>
```
