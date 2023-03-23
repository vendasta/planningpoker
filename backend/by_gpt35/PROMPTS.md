# Coding with GPT (OpenAI)

If you haven't already registered a free account with OpenAI, please do so here: https://chat.openai.com/

# Prompts (use a language you are comfortable with)

# Create Session

## Prompt 1
We are building a planning poker game. The create session api accepts the following json:

```
{
'participant_id': '1234'
}
```

The api should return a session_id and an access token that is a random alphanumeric string of 6 characters.

```
{
'session_id': '1234',
'token': 'QN7834'
}
```

Build a REST api that implements this in golang and the gorilla mux router.


## Prompt 2
Persist the session id into memory via cookies so that multiple participants can join the session. Cookies should be handled via the gorilla securecookie package.

# Join Session

## Prompt 3
The join session handler should accept the session id as part of its request path: /session/<session_id>/join

## Prompt 4
The join session handler should accept the following json:

```
{ 'participant_id': 'Dale Hopkins'}
```

The api should return a new session id for the joining participant.

```
{
'session_id': '1234',
'token': 'QN7834'
}
```

### Tweaks
Stored the joining participant's session id in a cookie.
Ensured join session had its own request and response objects. 

# Wait For Prompt

## Prompt 5
After a user has joined a session, they will be redirected in the frontend where they will wait for a prompt.

The handler url should be /session/<session_id>/prompt/wait

The api would accept the following json:
```
{
  'last_prompt_id': 'abc'
}

And it would return the following response. 
```
{
'prompt_id': 'def',
'prompt': 'Add support for OpenID Connect to the Application's Authorization flow',
}

A 404 should be returned if the session is not found.

## Prompt 6
The prompt will have be written to a new line delimited json filed that stores the prompt in the following format. Update the promptWait handler to read from prompts.json

```
{
        SessionID: "",
		PromptID: "def",
		Prompt:   "Add support for OpenID Connect to the Application's Authorization flow"
}
```

## Create Prompt

## Prompt 7

Generate a New Prompt handler for the following url /session/<session_id>/prompt/create
```
{
  'prompt': 'Add support for OpenID Connect to the Application\'s Authorization flow',
}
```

The session_id should exist in the cookie store. The Prompt would be written to a prompts.json file that is new line delimited json that uses the Prompt struct defined previously. 

### Tweaks
Created PromptRequest object

## New Vote

## Prompt 8

Build a vote handler `/session/<session_id>/prompt/<prompt_id>/vote` that validates the session and prompt id and records the results to votes.json

The api payload would would accept one of the following votes: `0`, `1`, `2`, `3`, `5`, `8`, `13`, `20`, `40`, `100`, `?`, `☕️`

The payload would be:

```json
{
  'vote': '1',
}
```

## Watch for Votes

## Prompt 9

After a vote is place the frontend will redirect the user to the watch for votes endpoint. This endpoint would return a response of the votes that have been made: URL `/session/<session_id>/prompt/<prompt_id>/watch`

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

## Close Session

## Prompt 10

Generate the close session handler. This should delete the session no longer allows votes to be made or viewed.

The route should be `/session/<session_id>/close`


## Final Tweaks

- Moved the handlers into their own files.
- The longer the prompt got the more context was lost. Needed to instruct the session handling
- The session handling and auth got a bit messy as the generation got longer. Fixed by hand
