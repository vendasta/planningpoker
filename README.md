# planningpoker
Planning Poker

# Create Session
Create Session <url>
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

# Join Session
Join Session <url>
```json
{
    'participant_id': 'Jesse Redl',
    'session_id': '1234'
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

# New Prompt
New Prompt <url>
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

# Vote
New Vote <url>
```json
{
  'prompt_id': 'abc',
  'vote': '1',
}
```
Headers
- content-type = 'application/json'
- accept = 'application/json'
- authorization = 'bearer <token>'

# Close Session
Close Session <url>
```json
{
  'session_id': '1234',
}
```
Headers
- content-type = 'application/json'
- accept = 'application/json'
- authorization = 'bearer <token>'

# 
