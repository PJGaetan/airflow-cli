# airflow-cli

Cli to interact with you airflow instance using the REST API when you can't ssh into the airflow server.


The code as been greatly inspired by [glab](https://gitlab.com/gitlab-org/cli/-/tree/main) and [jira-cli](https://github.com/ankitpokhrel/jira-cli).

## Todo 
- Choose a name
- Reoganise required module
- Wrap error and sys exit [x]
- Correct Select descriptions servey [x]
- Correctly escape with crtl+d [x]
- Write quickstart README [x]
- Output DAG with [diagon](https://github.com/ArthurSonzogni/Diagon) and [natural C biding](https://pkg.go.dev/cmd/cgo)

## Quickstart

### Create a new profile 

Create a profile using :

```sh
airflow-cli profile create
```

you can choose between two kind of authentication :
- user/password
- jwt auth

In case of a jwt auth you can choose to use a shell command to retrieve the jwt token instead of directly providing it.

### Dag status

Get a brief overview of your dags status with :

```sh
airflow-cli dag-run status
```

## Contribute

### Test

Use test makefile to spin up any airflow version.

```sh
 make up VERSION=2.5.3
 ```

