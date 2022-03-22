# Pharmacy automatization

This project is made for university work. This microservice for working with pharmacies.

## Layout of project:
- cmd - Small main function(s)
- internal - All business logic 
- pkg - Things than is too small to move on different repository, or just wrap over wrap :) May just copy to other projects
- configs - config for start up project
- .build - Docker
- api - any kind of API or description of project

Tasks:
- [x] Think about client error handling 
- [x] Move apperror & interceptors package to other repo
- [ ] Add Jaeger and OpenTracing

