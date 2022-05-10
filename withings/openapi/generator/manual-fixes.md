## Fix swagger file
- Add server to oauth2 authorize endpoint.
- Add ?action=operationName to each path key, even though this is disallowed by OAS3.
- Remove format fields
- Change timezone, macaddress and various to string/integer
- Add items.type=string and "format": "json" to night_events
- Replace all `format": "timestamp",` with `format": "int64",`

## Fix build output
- Change `v.value = nil` to `v.value = SleepSummaryObject{}`
- Delete go.mod