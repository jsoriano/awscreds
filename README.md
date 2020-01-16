Simple helper script to convert output of `aws sts` to usable credentials
sections into `~/.aws/credentials`.

Installation:
```
go get -u github.com/jsoriano/awscreds
```

To generate a new MFA token and add it to a section called `newmfa`, run:
```
aws sts get-session-token --serial-number <mfa-device-serial-number> --token-code <token-code> | awscreds -section newmfa -u
```
