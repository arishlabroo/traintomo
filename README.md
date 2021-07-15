# traintomo

## endpoints
### To Post Schedule

`/postschedule`

###### request body
```json 
{
    "Name":"ABCD",
    "Schedule":["01:20 AM", "04:15PM", "11:59 PM"]
}
```

### To Get Next Conflict
`/getnextconflict?aftertime={time}`

###### response body
```json 
{
    "Conflict":"04:15 PM"
}
```
 