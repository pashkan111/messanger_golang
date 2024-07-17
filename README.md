Хочу сделать просто мессенджер, чтобы попробовать конкурентность в ГО


Stack: Go, Postgres, Redis

Postgres schema:

User:
    user_id: int
    username: str unique
    chats: array[int]

Chat:
    chat_id: int
    name: str
    participants: array[int]

Message:
    message_id: int
    text: str
    created_at: datetime
    chat_id int FK Chat
    author_id int FK User


User login or register in a system. He is sent a token

AUTH:
    /register
    /login

MESSANGER:
    /chats - returns all chats
    WS /chat/{id} - return last 20 messages. Make a pagination

Создание чата:

1. Создается запись в табл Chat и добавляется поле юзерс
2. User.chats добавляется айди чата
3. В редисе создается сет со всеми участниками чата и при отправке сообщения в чат, всем участникам группы уходит это сообщение


При переходе по ВС роуту /chat/{id}:
1. устанавливается соединение с сервером. Загружаются сообщения. Использовать редис стримс. 
2. В канал онлайн в редисе добавляется айди пользователя
3. Запускается горутина, которая прослушивает канал редис стримс и
    обрабатывает сообщения
4. В редис стримс отправлять событие типа         
        Message:
            message_id
            text
            created_at
            chat_id
            author_id


5. Клиентские События: 
    MessageCreated
    MessageUpdated
    MessageDeleted

    Будет ВС хэндлер, который получает сообщение от клиента, валидирует его
    и отправляет в очередь

    В то же время, запущенный консумер вычитывает сообщения из очереди
    и обрабатывает их и если это сообщение от текущего юзера, то отправляет
    на фронт подтверждение о сохранении. Если не от текущего, то отправляем на фронт это сообщение. 

    В случае если сообщение от текущего юзера, параллельно запускать горутину, которая добавит сообщение в БД

    Структура: 
        Message:
            message_id
            text
            created_at
            chat_id
            author_id

Реализация

1. Сделать логику создания пользователя. 
    /register/
    /login/

Для регистрации и логина использовать номер тлф и пароль
Использовать токен

2. Написать круд для чатов
3. Написать круд для сообщений


Есть страница с поиском участников
Там же можно создать группу (при создании группы, не нужна проверка на наличие такой же группы)
And it is allowed to write a message any person.
When message is sent chat creation happens under the hood. If such chat exists
all messages will be attached to it, else chat will be created
