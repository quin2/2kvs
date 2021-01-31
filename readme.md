#calmDB

a simple key-key-value database written in Go. Designed to be concurrent and fast, with a web-first HTTP frontend. The (very rough) paper outlining my motivations and work are included in this repo.

## examples 

```
INSERT quinnvinlove name quinn

DELETE quinnvinlove name

SELECT quinnvinlove
SELECT quinnvinlove name
```

## issues
* needs better API docs
* only supports string for now for value. allow byte object?
* in-memory only at the moment
* atomic delete no longer supported. should it be supported?
* need way to clean out tombstone cache every so often, try gocron or similar?
* make tombstoning faster
* add mutex lock to map altering to make concurrent
* (any) operation not supported for k2, but maybe I could make our hack more elegant
* remove support for edit!
* fix issue where deleted keys can't be added again