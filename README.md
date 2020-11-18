# Sindel Imagine2 microservice for files and images uploading and processing

- GET http://localhost:8385
  To get service stats.
  Returns JSON response.

- POST http://localhost:8385/upload [parameters: file={Filedata}]
  To upload file to storage.
  Returns JSON response.

- POST http://localhost:8385/save_base64 [parameters: data={Base64_Data}&file={filename}]
  To upload file from base64 data. Create or Update existing file
  Returns JSON response.

- GET http://localhost:8385/show [parameters: id={file_id}&transform={transform}]
  To Render file by id with transformations.
  Returns raw file response with proper content-type header.

- GET http://localhost:8385/file?id={file_id}
  To get file information by id.
  Returns JSON response.

- GET http://localhost:8385/files
  Get files.
  Returns JSON response.

- GET http://localhost:8385/delete?id={file_id}
  To delete file by id
  Returns JSON response.

- GET http://localhost:8385/render/{filepath:*}  
  GET Render file by path with or without transformations.
  Returns raw file response with proper content-type header.

- POST http://localhost:8385/invalidate_cache [parameters: id={file_id}]
  Invalidate all image files transformations. Removes cached files.
  Return JSON response.

