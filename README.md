# Go-Cowpoke

 __NOTE__: This is not production ready!

A re-implementation of [cowpoke](https://github.com/LeanKit-Labs/cowpoke) in Go

## Requirements
  * Go 1.8 (for properly encoding URI parts)
  * [GB](https://getgb.io/)

## Building and Running

  * ```gb build```
  * ```bin/cowpoke```

## Dependencies

Dependencies should be added to __/vendor__ see [vendoring in Go](https://blog.gopheracademy.com/advent-2015/vendor-folder/)
  * ```gb vendor fetch <dep>```
  * copied manually

## Configuration

The service relies solely on environment variables. If the service is running in anything other than
release mode (__GIN_MODE=release__) the service will panic if a __.env__ file is not found in the project route.
This is to prevent the service from using configuration outside of env vars in say a Docker container.

The following table lists environment variables for the service

| Name  | Required?  | Notes
|---|---|---|---|---|
| HOST_PORT | YES | HTTP port the server listens on |
| RANCHER_USER_KEY  | YES | The account api key the service uses to preform upgrades |
| RANCHER_USER_SECRET | YES | The account api secret the service uses to preform upgrades |
| RANCHER_URL | YES | The URL to the Rancher server (http://someUrl:8080) |
| API_KEY| NO | Used to authenticate calls to the service. If this value is set, callers should supply the token via an HTTP Bearer header |
| SLACK_TOKEN | NO | Authentication token for optional slack notifications |
| SLACK_CHANNELS | NO | Comma delimited list of slack channels |

Example .env file

```
HOST_PORT=8080
RANCHER_USER_KEY="some key"
RANCHER_USER_SECRET="some key"
...
```

## Differences with the original implementation

There are a few key differences in this implementation, at a high level:

* configuration
  * handled strictly via environment variables
  * local configuration is handled via a .env file
* upgrading a stack
  * the request parameters required have been simplified
  * the response body is more consistent, and let's callers know the status of every stack upgrade attempt
  * the service gets template version files directly from Rancher instead of Github
  * the request is blocking. this was done to make the reply data more consistent with what actually happens in Rancher.
    The intent is to replace this endpoint entirely and redesign it as an async job queue.

## API

### /api/_status -> Health Check

__GET__

Response
```
"service up"
```

### /api/stack ->  Upgrade stack(s) built from the specified template to the provided template version

__PATCH__

Request (all fields required)
```
{
  "catalog": "some catalog",
  "template": "my template",
  "templateVersion": "2.0.0" //from the 'version' attribute in the template's rancher-compose.yml
}
```

Response
```
{
  "msg": "results from upgrading stack(s)",
  "results": [
    {
      "name": "cowpoke",
      "environment": "1a5",
      "upgradedTo": "cowpoke",
      "error": ""
    },
    {
       "name": "cowpoke",
       "environment": "1a6",
       "upgradedTo": "cowpoke",
       "error": ""
    },
    {
       "name": "cowpoke",
       "environment": "1a7",
       "upgradedTo": "cowpoke",
       "error": "something bad happened!"
    }
  ]
}
```

__Notes on Upgrading Stacks__

This endpoint will upgrade a stack if the following conditions are met

  * the service has access to an environment (i.e. account access key env vars)
  * a stack was created from a template (i.e. an externalId of catalog://catalog:template:version)
  * the stack was created from the template provided in the request
  * the version portion of the stack's external id is "less than" the version of the template in the request

__Assumptions and Opinions on Catalog Structure__

For the service to be able to upgrade it needs to differentiate between version. It currently determines this
by assuming that each version of a template is put in a directory structure where each directory is an auto incrementing int

Example Catalog:

```
  0/
    rancher-compose.yml
    docker-compose.yml
  1/
    rancher-compose.yml
    docker-compose.yml
  ...
```

This structure will cause each version to have an id of

* someCatalog:someTemplate:0
* someCatalog:someTemplate:1







