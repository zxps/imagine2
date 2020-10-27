# Sindel Imagine2 microservice for files and images uploading and processing

- http://localhost:8385/save         - POST Upload file to storage (JSON response)

- http://localhost:8385/save_base64  - POST Upload file from base64 data. Update existing or create new file if not exixsts (JSON response)

- http://localhost:8385/show         - GET To get file by id or path ( File response  )

- http://localhost:8385/file         - GET file information by id or path ( JSON response )

- http://localhost:8385/image_resize - GET Get resized image file (id,width,height) ( File response )

- http://localhost:8385/image_resize_save - POST Resize and save resized image (id,width,height) ( JSON response )

- http://localhost:8385/image_crop   - GET Get cropped image file (id,x,y,width,height) ( File response )

- http://localhost:8385/image_crop_save - POST Crop and save cropped image (id,x,y,width,height) ( JSON response )

- http://localhost:8385/invalidate_cache   - POST Invalidate all image files transformations ( JSON response )
