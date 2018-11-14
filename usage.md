# Usage at this moment (5.11.18)

The webservice was introduced and has some basic functionality to use.
The current version is always automatically deployed at

    https://infogration.now.sh/

In order to get ALL the data which is saved you can use
    
    curl https://infogration.now.sh/content/all | python -m json.tool

where the last part is for pretty print and will return

    {
    "Status": "OK",
    "Data": [
        {
            "Content": {
                "Id": 1,
                "Question": "How is life these days?",
                "Answer": "So good"
            },
            "Language": {
                "Code": "en"
            },
            "CreatedAt": "1541455846"
        },
        {
            "Content": {
                "Id": 2,
                "Question": "Are 2 questions sufficient?",
                "Answer": "I do not think so!"
            },
            "Language": {
                "Code": "en"
            },
            "CreatedAt": "1541455846"
        },
        {
            "Content": {
                "Id": 3,
                "Question": "Are 3 questions sufficient?",
                "Answer": "I think so!"
            },
            "Language": {
                "Code": "en"
            },
            "CreatedAt": "1541455846"
        },
        {
            "Content": {
                "Id": 2,
                "Question": "2 preguntas son suficiente?",
                "Answer": "Creo que no!"
            },
            "Language": {
                "Code": "es"
            },
            "CreatedAt": "1541455846"
        }
    ]
    }


which includes some dummy data that I created there.

To get all instances of one language one can run

    https://infogration.now.sh/content/{lang}

For lang='es' for example we get

    {
    "status": "OK",
    "data": [
        {
            "content": {
                "Id": 2,
                "Question": "2 preguntas son suficiente?",
                "Answer": "Creo que no!"
            },
            "language": {
                "Code": "es"
            },
            "createdAt": "1541456019"
        }
    ]
    }

as an answer. If we enter something like

    https://infogration.now.sh/content/asdf

we get

    {"status":"Bad Request","data":[]}
