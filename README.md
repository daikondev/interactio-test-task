### Interactio Test-Task

##### How to use this API:

---

| Endpoint    | Method | Notes                                  |
|-------------|--------|----------------------------------------|
| /           | GET    | A simple greeting                      |
| /events     | GET    | Returns a list of all available events |
| /events     | POST   | Creates a new event                    |
| /events/:id | GET    | Returns the event with the given id.   |

_For the events/:id table, you may specify the desired audio and video quality
using the videoQuality and audioQuality query strings.
If they are not specified, it will return the quality specified by the server defaults._

##### How to start the server: 

---

1. Run the ```Go build``` command.
2. Start the server with ```Go run test-task```.

##### Server defaults: 

- Max Invitees: the maximum number of invitees an event can have. Default value is 100 invitees. To change this setting, use the ```-m``` flag with the number of invitees you wish to specify.
- Default Video Quality: the default video quality served to clients. Default value is 720p. To change this setting, use the ```-v``` flag with the new default quality.
- Default Audio Quality: the default audio quality served to clients. Default value is mid. To change this setting, use the ```-a``` flag with the new default quality.

Created by Alex Kalpakidis