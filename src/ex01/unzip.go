package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func compareCredentialFields(args []string) {
	if len(args) < 2 {
		return
	}
	if args[0] == "adminLogin:" {
		adminLogin = args[1]
	} else if args[0] == "adminPassword:" {
		adminPassword = args[1]
	} else if args[0] == "SQLcommand:" {
		SQLcommand = args[1]
	}
}

func readCredentials() error {
	fileBytes, err := os.ReadFile("admin_credentials.txt")
	if err != nil {
		return err
	}
	fileSplitted := strings.Split(string(fileBytes), "\n")
	for _, value := range fileSplitted {
		args := strings.Fields(value)
		compareCredentialFields(args)
	}
	if len(adminLogin) == 0 {
		adminLogin = "bulat"
	}
	if len(adminPassword) == 0 {
		adminPassword = ""
	}
	if len(SQLcommand) == 0 {
		return errors.New("missing sql command to create table")
	}
	return nil
}

func unzipStaticFiles() error {
	dst, err := filepath.Abs("./")
	if err != nil {
		return err
	}
	archive, err := zip.OpenReader("static_files.zip")
	if err != nil {
		return err
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(dst, f.Name)
		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			return errors.New("wrong filepath")
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}
		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		fileInArchive, err := f.Open()
		if err != nil {
			return err
		}
		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return err
		}
		defer dstFile.Close()
		defer fileInArchive.Close()
	}
	fmt.Println("finished unzip")
	return nil
}
