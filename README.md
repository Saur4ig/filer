# Sender and uploader large files (could be to 100gb)
## base description
    listner - webserver with 2 endpoints:
       /save - saves a file(s), chunk comes -> waits till all parts will be downloaded,
       then creates a file with a unique values and combines all part of files into one, then removes all rem parts and folders
       /get/{fileName.txt} - returns file(if exists) from files/unique
    
    removings duplicates -> ideally paste all string from file(s) to some db, for example, clichouse, then download all unique strings to file,
    coz it's too long to create solution like this - I`ve created an interface Database and used just map for storing all data
## usage
### listener
    make run_listener
### sender
    make run_sender

## api
* **URL**

  `/save`

* **Method:**

   `PUT`

* **Description:**

    if file size more than 5mb -> file should be chunked to several parts!
    _Sender_ was created for fast file sending! `Compile sender binary and paste file name using -flag=<fileName>,
    or use` `make run_sender` `to send with default host, port, endpoint, if u want to change - paste them using flags` 

---
* **URL**

  `/get/{fileName.txt}`

* **Method:**

   `GET`

* **Description:**

    Just a mux static file route.
---

The best I could come up with - to chunk file to a parts `5-50-500mb` or less and receive them in several threads 