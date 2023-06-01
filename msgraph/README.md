# onedrive-API info
  - API explorer : https://developer.microsoft.com/en-us/graph/graph-explorer
  - Onedrive REST APIs: https://learn.microsoft.com/en-us/onedrive/developer/rest-api/?view=odsp-graph-online
  - Onedrive access  https://learn.microsoft.com/en-us/graph/api/resources/onedrive?view=graph-rest-1.0
  - sample upload code here for file > chunkSize using looping:
        https://github.com/microsoftgraph/msgraph-sdk-java-core/blob/dev/src/main/java/com/microsoft/graph/tasks/LargeFileUploadTask.java
    - the above code does not handle ranges properly though, it is changed in go code in this repo.

  - sample go SDK code here - https://learn.microsoft.com/en-us/graph/api/driveitem-list-children?view=graph-rest-1.0&tabs=go 

## NOTES
 - sdk-go support for upload session is not complete, it lack method to provide filename. this is separately added in this repo.

## TODOS
- remove non drive related code
- reformat code to have simpler APIs
- use other auth mechanism eg asking token from user or using client certificates etc.
- Candidate for pull request - filename support for go sdk upload session.