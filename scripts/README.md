# dp-identity-api scripts

This directory contains several scripts, many of which are used for:

- exporting users (and groups) from zebedee and/or Cognito
- importing users into Cognito
- [miscellaneous](#miscellaneous-scripts)

## Miscellaneous scripts

### collection_users

The `collection_users` script converts users' details from collections
(e.g. creator IDs in the collection JSON)
into more human-readable names and emails (from Cognito).

For more detail, run `./collection_users --help` from this directory.

Typical usage might be to summarise the users of the latest pre-publish collections;
this can be done by running:

```shell
$ ./collection_users prod --summary -NT
# collection_summary line for a given collection
#          cognito details (JSON) for the above collection
...
```
