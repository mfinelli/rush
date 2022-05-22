# rush

An ssh certificate client and server.

## server

The server component is responsible for generating and signing certificates
when they are requested either by a host or by a user.

### authentication

Rush currently supports basic authentication via an `htpasswd` file. You can
generate one and add new users like so (you'll be prompted for the password):

```shell
htpasswd -c /etc/rush/htpasswd username
```
