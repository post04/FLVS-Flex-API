# FLVS Flex API
 An API for FLVS Flex because FLVS has no public API.

# Why
FLVS has no public API that students can use, when I was a student at FLVS I made this so that I didn't need to take the time to login every time I wanted to check my grades or when the last time I submitted an assignment was.\
This can also be useful for when you aren't home but would like to check your grades, you can use this API for something like discord to automatically update your grades in a message (like I did).

# Information
APIkeys need to be regenerated every ~1 hour or if you login on your client, for this reason I would just do something like this (what I did)
```go
infoUnFormatted, err := schoolapi.GetClasses(userid)
if err != nil {
	apiKey = schoolapi.GetAPIKey(username, password, true)
	infoUnFormatted, err = schoolapi.GetClasses(userid)
	if err != nil {
		os.Exit(2)
		return
	}
}
```

# Disclaimer
This is **NOT** an official or supported FLVS project and any action that FLVS may deem necessary to take on your account is not my fault. I didn't see anything in the FLVS TOS about this but **if** something does happen it is not my fault.

# Note
I no longer attend FLVS and cannot actually test this project anymore, if it stops working then you can fix it yourself and submit a PR.
