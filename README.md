### Interactio Test-Task

##### How to use this API:

---

| Endpoint    | Method | Notes                                  |
|-------------|--------|----------------------------------------|
| /           | GET    | A simple greeting                      |
| /events     | GET    | Returns a list of all available events |
| /events     | POST   | Creates a new event                    |
| /events/:id | GET    | Returns the event with the given id.   |

_For the events/:id table, you may specify the desired audio and video quality.
If they are not specified, it will return the quality specified by the server defaults._

##### How to start the server: 

---

The server should be started with ```go build```. When running ```go run test-task```, you may
specify the following flags:
-  ```-p```: allows the user to specify the port the server listens to
-  ```-m```: allows the user to specify the maximum number of invitees an event can have
-  ```-v```: allows the user to specify the default video quality served to clients
-  ```-a```: allows the user to specify the default audio quality served to clients

Created by Alex Kalpakidis