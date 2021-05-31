# Libreoffice Convert

Simple http server to convert files using libreoffice

## Install

```
go install github.com/fabianMendez/libreofficeconvert@latest
```

## Usage


### Run server

Run a local http server on port `1234` (optional) and specify the path to the libreoffice binary (optional).

```
PORT=1234 LIBREOFFICE_PATH=libreoffice7.0 libreofficeconvert
```

### Generate using CURL

You can use curl or any other http client to make a `POST` request with the document to convert in the body
of the request and the destination extension in the path of the URL (`pdf` in this case):

```
curl  --data-binary @/source/document.docx http:/localhost:1234/pdf --output /destination/document.pdf
```
