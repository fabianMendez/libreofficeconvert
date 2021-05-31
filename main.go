package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func convert(src io.Reader, dst io.Writer, extension string) error {
	tmpDir, err := os.MkdirTemp(os.TempDir(), "*")
	if err != nil {
		return fmt.Errorf("creating temp dir: %w", err)
	}
	log.Println("Created temp directory:", tmpDir)

	defer func() {
		err = os.RemoveAll(tmpDir)
		if err != nil {
			log.Printf("deleting temp directory: %v\n", err)
		}
		log.Println("Deleted temp directory:", tmpDir)
	}()

	srcFilename := filepath.Join(tmpDir, "src")
	srcFile, err := os.OpenFile(srcFilename, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not create temp file: %w", err)
	}
	defer srcFile.Close()
	log.Println("Source filename:", srcFilename)

	_, err = io.Copy(srcFile, src)
	if err != nil {
		return fmt.Errorf("could not copy source file: %w", err)
	}

	libreofficePath := os.Getenv("LIBREOFFICE_PATH")
	if libreofficePath == "" {
		libreofficePath = "libreoffice"
	}

	cmd := exec.Command(libreofficePath, "--headless", "--convert-to", extension, "--outdir", tmpDir, srcFilename)
	log.Println(cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("executing command: %s - %w", cmd.String(), err)
	}

	srcCleanname := strings.TrimSuffix(srcFilename, filepath.Ext(srcFilename))
	dstFilename := srcCleanname + "." + extension
	log.Println("Destination filename:", dstFilename)

	dstFile, err := os.Open(dstFilename)
	if err != nil {
		return fmt.Errorf("could not open destination file: %w", err)
	}

	_, err = io.Copy(dst, dstFile)
	if err != nil {
		return fmt.Errorf("could not copy destination file: %w", err)
	}

	return nil
}

func respondWith(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(`{"error" "` + message + `"}`))
}

func convertHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL)

	extension := path.Base(r.URL.Path)

	if extension == "" {
		respondWith(w, http.StatusBadRequest, "destination extension is required")
		return
	}

	tmpfile, err := os.CreateTemp("", "*")
	if err != nil {
		log.Println(err)
		respondWith(w, http.StatusInternalServerError, "could not create temp file")
		return
	}
	defer func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}()

	err = convert(r.Body, tmpfile, extension)
	if err != nil {
		log.Println(err)
		respondWith(w, http.StatusInternalServerError, "could not generate library")
		return
	}

	tmpfileread, err := os.Open(tmpfile.Name())
	if err != nil {
		log.Println(err)
	}
	_, err = io.Copy(w, tmpfileread)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	http.HandleFunc("/", convertHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on port %s\n", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
