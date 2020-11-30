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
* needs docs
* multiple values with the same key combo are just added, rather than swapped in. use a hash table?
* only supports string for now 
* in-memory only at the moment. 

switch to golang maps https://blog.golang.org/maps issue is we can't have two keys...
make one mega key? (then can't perform general lookups)

use key for k1, and that map 

or extend comparable to work with any data! type is key tuple


m = make(map[string]map[string]string)

issue with dual key: both values have to be eq, so we can't do something like k1 (any). no way to change this either in the spec 