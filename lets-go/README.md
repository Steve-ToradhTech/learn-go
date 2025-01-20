# Lets go personal docs

## Troubleshooting

### Mysql is throwing an error showing the ip address as the connecting host? 
The docker bridge network is passing an IP address instead of the hostname. Quickfix
is to just update mysql to create a generic web accounts with wildcard hostname access.
Could probably do some hostname resolution on the containers but I'm just lazy.
`CREATE USER 'web'@'%';`
`GRANT ALL PRIVILEGES ON *.* TO 'web'@'%' IDENTIFIED BY 'yourpassword';`
`ALTER USER 'web'@'%' IDENTIFIED BY 'pass';`