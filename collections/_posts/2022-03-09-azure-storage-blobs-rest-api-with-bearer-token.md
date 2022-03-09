---
layout: post
title:  "Access the Azure Storage REST API using an OAuth2 bearer token"
date:   2022-03-09 19:47:32 +0100
categories: azure
excerpt: I was getting 403 errors when using a valid bearer token for the REST API calls
#proccessors: pymd
---

I am accessing the Azure Storage REST API to list, read and modify blobs in a storage
container.

I was having problems accessing the endpoints using an OAuth2 token
generated using an Azure AD Service Principal's client ID and secret ID.

First get a fresh token:

```
$ export AZURE_TOKEN="$(curl -X POST \
    https://login.microsoftonline.com/${AZURE_TENANT_ID}/oauth2/token \
    -d client_id=${AZURE_CLIENT_ID} \
    -d client_secret=${AZURE_CLIENT_SECRET} \
    -d resource=https://storage.azure.com/ \
    -d grant_type=client_credentials \
    -o- \
    | jq -r .access_token)"
```

Then attempt to list blobs:

```
$ curl -i --oauth2-bearer "$AZURE_TOKEN" \
    https://mystorageaccount.blob.core.windows.net/?comp=list

HTTP/1.1 403 Server failed to authenticate the request. Make sure the value of Authorization header is formed correctly including the signature.
Content-Length: 438
Content-Type: application/xml
Server: Windows-Azure-Blob/1.0 Microsoft-HTTPAPI/2.0
x-ms-request-id: <redacted>
Date: Wed, 09 Mar 2022 18:54:23 GMT

<?xml version="1.0" encoding="utf-8"?>
<Error><Code>AuthenticationFailed</Code><Message>Server failed to authenticate the request. Make sure the value of Authorization header is formed correctly including the signature.
RequestId:<redacted>
Time:2022-03-09T18:54:24.3070933Z</Message><AuthenticationErrorDetail>Authentication scheme Bearer is not supported in this version.</AuthenticationErrorDetail></Error>
```

It appears that I was implicitly using an old version of the API by not sending the
a `x-ms-version` request header -- a version that did not support OAuth2 tokens.

I must have skimmed the official documentation too fast, because after I started writing
this post I found [a place documenting this][api-versioning]:

> Requests using an OAuth 2.0 token from Azure Active Directory (Azure AD). To authorize a request with Azure AD, pass the x-ms-version header on the request with a service version of 2017-11-09 or higher. For more information, see Call storage operations with OAuth tokens in Authorize with Azure Active Directory.

And [this whole page][api-versioning-2] also mentions it pretty explicitly:

> To call Blob, Queue and Table service operations using OAuth access tokens, pass the access token in the Authorization header using the Bearer scheme, and specify a service version of 2017-11-09 or higher [...]

Passing a header containing a newer version suddenly made all the difference.

```
$ curl --oauth2-bearer "$AZURE_TOKEN" -i \
    -H 'x-ms-version: 2017-11-09' \
    https://mystorageaccount.blob.core.windows.net/?comp=list

HTTP/1.1 200 OK
Transfer-Encoding: chunked
Content-Type: application/xml
Server: Windows-Azure-Blob/1.0 Microsoft-HTTPAPI/2.0
x-ms-request-id: <redacted>
x-ms-version: 2017-11-09
Date: Wed, 09 Mar 2022 19:20:32 GMT

<?xml version="1.0" encoding="utf-8"?><EnumerationResults ServiceEndpoint="https://mystorageaccount.blob.core.windows.net/"><Containers><Container><Name>someblobfolder</Name><Properties><Last-Modified>Mon, 21 Feb 2022 23:33:06 GMT</Last-Modified><Etag>"0x8D9F59282F28953"</Etag><LeaseStatus>unlocked</LeaseStatus><LeaseState>available</LeaseState><HasImmutabilityPolicy>false</HasImmutabilityPolicy><HasLegalHold>false</HasLegalHold></Properties></Container></Containers><NextMarker /></EnumerationResults>
```

I opened [a PR for an addition to the documentation][pr] to add a note even one more
place, but let's see if they deem it necessary. I'm not even sure myself anymore :)

Note to self: RTFM -- don't just skim the fine manual

## References
- <https://docs.microsoft.com/en-us/rest/api/storageservices/versioning-for-the-azure-storage-services#specifying-service-versions-in-requests>

[api-versioning]: https://docs.microsoft.com/en-us/rest/api/storageservices/versioning-for-the-azure-storage-services#specifying-service-versions-in-requests
[api-versioning-2]: https://docs.microsoft.com/en-us/rest/api/storageservices/authorize-with-azure-active-directory
[pr]: https://github.com/MicrosoftDocs/azure-docs/pull/89507
