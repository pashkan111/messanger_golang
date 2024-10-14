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


TODO:
1. Write an interface of PubSub
2. Write redis impl
3. Write consumer
4. Write func that would check 
