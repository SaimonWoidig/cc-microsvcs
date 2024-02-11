# cc-microsvcs

An agent command-center for remote agent/server maintenance and control.  
Using the microservices architecture.

## Project structure

- `common` - code shared by all microservices, containing OTEL tracing, metrics and logs setup and other packages used by all microservices
- `service.auth` - authentication microservice

## Local development

Use a go.work file and include the required packages.
