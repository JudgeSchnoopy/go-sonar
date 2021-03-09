# go-sonar
Detailed dependency management for tying microservices together

## Objective
Prevent breaking changes by running live integration tests from all servers that integrate with the changing code.

## Why?
I've deployed many microservices to each integrate with specific vendor APIs.  These were all appropriately disconnected and independent.

I began deploying a second layer of microservices for richer functionality.  They tied information together from different sources, made decisions, and pushed information to various services.  This layer was specifically designed to tie into multiple microservices.  But how can I trust my future self, or perhaps my coworkers, to not push a breaking change that could topple multiple downstream services?

This is the strictest method I could come up with.

## How?
Each microservice implements an endpoint, "/sonar".  This endpoint takes a POST that details the microservice name and new test deployment address.
The server checks a set of 'dependency checks', or integration tests categorized by dependency.  If any tests are defined for the microservice name, it runs those tests using the provided test deployment address.  If any tests fail, it responds with a detailed failure message on which endpoint is broken.

A central 'sonar' microservice will coordinate service registration and test request / result aggregation.

Any server implementing a /sonar endpoint can run a go-sonar client function.  This function reports in to the central Sonar server with the service's name and address.  That deployment will now be monitored by Sonar.

Any service, even those not monitored by Sonar, can send a request to Sonar with it's service name and address.  Sonar simply pings all registered services /sonar endpoint reporting that service name and address, and waits for the results.  Results will be marked as "Pass", "Fail", "None", or "Incomplete" (by default) for all registered services.  
Any service that does not have a dependency on the named service reports back "None".

After all servers respond (or a timeout is reached), Sonar responds to the requestor with a full dependency report : which services report dependencies on this service and whether their integration tests passed.
