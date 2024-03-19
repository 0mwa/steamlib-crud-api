# POST
## Publisher/Developer 
IN
```json 
{
  "name": "Valve",
  "country": "USA"
}
```
OUT
```json
{
  "error": "err"
}
```

# Games
IN
```json 
{
  "name": "TF2",
  "img": "img.png",
  "description": "cool",
  "rating": "10",
  "developer_id": "1",
  "publisher_id": "1",
  "steam_id": "440"        
  
}
```
OUT
```json
{
  "error": "err"
}
```

# GET
## Publisher/Developer 
IN
```json 
{
}
```
OUT
```json
{
  "name": "Valve",
  "country": "USA",
  "error": "err"
}
```
## Games
IN
```json 
{
}
```
OUT
```json
{
  "name": "TF2",
  "img": "img.png",
  "description": "cool",
  "rating": "10",
  "developer_id": "1",
  "publisher_id": "1",
  "steam_id": "440",
  "error": "err"
}
```

# PUT
## Publisher/Developer
IN
```json 
{
  "name": null,
  "country": "USA"
}
```
OUT
```json
{
  "error": "err"
}
```
## Games
IN
```json 
{
  "name": null,
  "img": "img.png",
  "description": null,
  "rating": "10",
  "developer_id": null,
  "publisher_id": null,
  "steam_id": "440"        
  
}
```
OUT
```json
{
  "error": "err"
}
```
# DELETE
## Publisher/Developer
IN
```json 
{
  "developer_id": "1"
}
```
OUT
```json
{
  "error": "err"
}
```
## Games
IN
```json 
{
  "steam_id": "440"
}
```
OUT
```json
{
  "error": "err"
}
```