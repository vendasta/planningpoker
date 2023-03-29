# Planning Poker

Planning Poker is a technique used in agile software development to estimate the effort required to complete a task or a project. It involves a team of people who collectively estimate the amount of work required for a given task or project.

Here's how Planning Poker works:

- A facilitator presents the task or project to the team.
- Each team member is given a set of cards, usually numbered from 0 to 100 or higher.
- The team member who is most familiar with the task or project provides a brief description of what needs to be done.
- Each team member privately selects a card from their set of cards to represent the effort required to complete the task. They should not share their card with the rest of the team.
- When everyone has made their selection, all the cards are revealed simultaneously.
- The team members then discuss the reasons behind their estimates and work to come to a consensus estimate.
- If there is a large discrepancy in the estimates, the team discusses the reasons behind the discrepancy and may choose to re-estimate the task or project.
- The process is repeated until the team agrees on a final estimate.

Planning Poker is a valuable tool for agile teams because it encourages collaboration and discussion, and it helps to eliminate bias in the estimation process. It also ensures that all team members have a say in the estimation process, which leads to more accurate estimates.

# Your Task

Given the following API's your task is to build a frontend for the planning poker session. Feel free to work in teams and use any language of your choice! Note, that you should feel comfortable using tools like Chat GPT!


# Example Prompt

We are building a frontend for a planning poker game. Build an angular component that can create a new session and display the session id for others to join. 
The token returned will be used to authorize future requests for the participant. The server url to create the session is `/session/create`

The api expects the following json payload:
```json
{ 'participant_id': 'Dale Hopkins'}
```

The api response be a json object with the following structure:

```json
{
'session_id': '1234',
'token': 'opaque access token for Dale'
}
```

# Available Votes
The following votes will be used within the planning session game. 

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
 "votes": [
  {
   "participant_id": "Dale Hopkins",
   "vote": "1"
  },
  {
   "participant_id": "Jesse Redl",
   "vote": "3"
  }
 ]
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
