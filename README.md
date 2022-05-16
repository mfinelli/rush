# rush

An ssh certificate client and server.

## server

The server component is responsible for generating and signing certificates
when they are requested either by a host or by a user.

**Note:** it is not a goal of rush to provide authentication or authorization,
it is better to rely on your API gateway or reverse proxy to provide access
control.
