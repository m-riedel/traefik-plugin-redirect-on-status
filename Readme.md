# Traefik Plugin: Redirect on Status

This Traefik Plugin provides a middleware to return a redirect to a specified domain when the next handler returns a specified error code. 
It behaves somewhat like the existing errors middleware.

## Configuration

| Key          | Allowed Values                                                                                                  | Required | Default               | Description                                                                                     |
|--------------|-----------------------------------------------------------------------------------------------------------------|----------|-----------------------|-------------------------------------------------------------------------------------------------|
| status       | The values are stored as List <br/> - Ranges: 200-299<br/> - enumeration: 200,201<br/>- single: 200             | yes      | -                     | Sets the repsonse stauts Codes of the down stream response, that the middleware does a redirect |
| redirectCode | 302, 303, 307 (All Temporary Redirect Status Codes)                                                             | no       | 307                   | Sets the status code of the response for the up stream                                          |
| method       | The values are stored as List <br/> - HTTP-Methods, like GET, POST, etc.                                        | no       | Acts for every Method | Sets that the middleware should only act for specific http methods                              |
| redirectUri  | A URI where to redirect to. May either be a FQDN or a relative Path. (Must align with the Location-Header spec) | yes      | -                     | Sets the value of the Location-Header. Therefore where the caller will be redirected to.        |