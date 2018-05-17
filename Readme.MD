# task-status-update


## Setting up to build go

See this https://steel-ventures.atlassian.net/wiki/spaces/VT/pages/8454186/Setup+development+environment

And look for the `GO specifics` section


In addition, make sure you have access to the repo https://github.com/veritone/go-veritone-api or practically ask Veritone DevOp to add you to the `InternalDevelopers` group



## Setting up locally

After setting up go,

```
go get github.com/quynhdang-vt/task-status-update
```

Assuming that you have access to a number of Veritone private repos, this command should clone this repo and dependencies to your local $GOPATH/src directory


Once that's cloned, to work with this tool, open a terminal window and ..

```
cd $GOPATH/src/github.com/quynhdang-vt/task-status-update
```


## To build:

```
make -f Makefile.inContainer dev
```

This will create the `task-status-update` executable

## Usage:

Make sure to prepare a `payload.json` file.  See [payload](#payload_json)

```
./task-status-update -payload payload.json
```


## payload.json<a name=payload_json></a>

| field name | Notes |
|------------|-------|
|recordings | JSON array of strings containing the recording ids to be modified }|
| engineId | engine Id or GUID executing the tasks that you want to modify status. |
| fromStatus | current status of the task, typically "pending" or "running" |
| toStatus | the new status to move the task to, typically "failed" or "complete" |
| token | API token -- must have enough priv, e.g. "recording:read, job:all" |
| veritoneApiBaseUrl | for Prod :  "https://api.veritone.com".  For uk-prod: "https://api.uk.veritone.com |




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


