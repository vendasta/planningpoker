# planningpoker
Planning Poker

Available Votes
`0`, `1`, `2`, `3`, `5`, `8`, `13`, `20`, `40`, `100`, `?`, `☕️`

# Create Session
Create Session URL `/session/create`
```json
{ 'participant_id': 'Dale Hopkins'}
```
Headers
 - content-type = 'application/json'
 - accept = 'application/json'

Return
HTTP 201 Created
```json
{
  'session_id': '1234',
  'token': 'opaque access token for Dale'
}
```

Sample CURL
```bash
curl 10.50.11.231:9000/session/create \
  -H "content-type: application/json" \
  -d '{"participant_id":"Dale Hopkins"}' \
  -X POST
```

# Join Session
Join Session URL: `/session/<session_id>/join`
```json
{
    'participant_id': 'Jesse Redl',
}
```
Headers
- content-type = 'application/json'
- accept = 'application/json'

Return
HTTP 200 OK
```json
{
  'token': 'opaque access token for Jesse'
}
```

```
HTTP 404 Not Found
(No open session with that ID was found)

HTTP 409 Conflict
(Another participant with that ID is already in the session)
```

Sample CURL
```bash
curl 10.50.11.231:9000/session/spicy%20bat/join \
  -H "content-type: application/json" \
  -d '{"participant_id":"Jesse Redl"}' \
  -X POST 
```

# Wait for Prompt
Wait for Prompt URL `/session/<session_id>/prompt/wait`
```json
{
  'last_prompt_id': 'abc'
}
```
Return
HTTP 200 OK
```json
{
  'prompt_id': 'def',
  'prompt': 'Add support for OpenID Connect to the Application's Authorization flow',
}
```
HTTP 204 No Content
(No new prompt) received before timeout

HTTP 404 Not Found
(Session was not found or was closed)

# Create Prompt
Create Prompt /session/<session_id>/prompt/create
```json
{
  'prompt': 'Add support for OpenID Connect to the Application\'s Authorization flow',
}
```
Headers
- content-type = 'application/json'
- accept = 'application/json'
- authorization = 'bearer <token>'


Return
HTTP 201 Created
{
    'prompt_id': 'abc'
}

Sample CURL
```bash
curl 10.50.11.231:9000/session/lazy%20rat/prompt/create \
  -H "content-type: application/json" \
  -H "authorization: Bearer 99222f4d-37b9-4171-89d5-d4eb82346661" \
  -d '{"prompt":"What is your favourite card?"}' \
  -X POST
```

# Vote
New Vote /session/<session_id>/prompt/<prompt_id>/vote/submit
```json
{
  'vote': '1',
}
```
Headers
- content-type = 'application/json'
- accept = 'application/json'
- authorization = 'bearer <token>'

Sample CURL
```bash
curl 10.50.11.231:9000/session/lazy%20rat/prompt/03727afc-6c28-470a-8ed9-2eb78783397c/vote/submit \
  -H "content-type: application/json" \
  -H "authorization: Bearer 99222f4d-37b9-4171-89d5-d4eb82346661" \
  -d '{"vote":"1"}' \
  -X POST
```

# Watch Vote 
URL `/session/<session_id>/prompt/<prompt_id>/watch`
```json
{
}
```
```json
{

}
```
Headers
- content-type = 'application/json'
- accept = 'application/json'
- authorization = 'bearer <token>'

# Close Session
Close Session /session/<session_id>/close
```json
{
}
```
Headers
- content-type = 'application/json'
- accept = 'application/json'
- authorization = 'bearer <token>'

Return
HTTP 204 No Content
(the session was closed)
HTTP 404 Not Found
(the session is not currently open)
HTTP 401 Unauthorized
(you are not the owner of the session)
