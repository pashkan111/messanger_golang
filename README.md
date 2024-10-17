Actions when connecting to backend by ws:
    1. get token from auth header
    2. Check the token. 
    3. get user chats ids and make list of channels for consumer
    4. Run 2 goroutines.  First listens to messages from client
        and process them
        Second - runs a consumer that listens to events from other clients
    5. Run a consumer, pass ids of the chats and 
        the consumer starts listening channels of such chats.
        When he receives a message, he send it to frontend
    
WHEN consumer reads message it should only check user id of it
and if id does not match the current user id it should send this event to front

TODO 
1. create a function to handle messages from the websocket connection