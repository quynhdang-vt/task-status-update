# task-status-update

## To build:

```
make -f Makefile.inContainer dev
```


## Usage:

```
./task-status-update -payload payload.json
```


## payload.json

Contains the ids of the recordings which contain tasks with engines stuck in a certain `status` and need to be moved to a new status
Token: need to be able to check jobs!!


Examples

```json
{
"recordings":["77951604","77951513","77860949"],
"engineId":"9084f686-cf02-41c5-adb5-d48fb6558d77",
"fromStatus":"pending",
"toStatus":"failed",
"token":"",
"veritoneApiBaseUrl": "https://api.veritone.com"
}
```


