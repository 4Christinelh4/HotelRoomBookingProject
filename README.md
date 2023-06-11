### Session management
Session is a logical connection between a server and a client. 
It persists over a period of time. 
When a user logs in the server, the server assigns a unique 
session identifier to the user. The user transfers the key when 
it makes requests to the server subsequently, before the 
session expires. 

In this project, `github.com/alexedwards/scs` is employed 
for session management. 

### Middleware
#### CSRF protection
When someone is clicking an img or something on web1, the attacker "injects" 
some API of the target website in the img or link. Then, when the img or link is clicked, 
requests are sent to the target website.

`CSRF token` can be used to prevent this.
1. The web application generates a unique CSRF token and associates it with 
the user's session.
2. The CSRF token is included in the HTML form or added to 
JavaScript-generated requests.
3. When the user submits a form or triggers an action, the CSRF token is sent along with the request.
4. The web application checks if the CSRF token received matches the one associated 
with the user's session. If they don't match or the token is missing, the request is considered invalid, and the action is not performed.


