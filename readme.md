#2kvs

a simple flat-file database written in Go 

## examples 

```
INSERT quinnvinlove name quinn

DELETE quinnvinlove
DELETE quinnvinlove name

SELECT quinnvinlove
SELECT quinnvinlove name
```

## issues
* multiple values with the same key combo are just added, rather than swapped in. use a hash table?
* only supports string for now 
* in-memory only at the moment. 