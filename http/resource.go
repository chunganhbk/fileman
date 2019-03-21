package http

import (
	"fmt"
	"github.com/raedahgroup/fileman/errors"
	"github.com/raedahgroup/fileman/files"
	"github.com/raedahgroup/fileman/fileutils"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var resourceGetHandler = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	file, err := files.NewFileInfo(files.FileOptions{
		Path:    r.URL.Path,
		Modify:  d.user.Perm.Modify,
		Expand:  true,
		Fs:      d.user.Fs,
		Checker: d,
	})
	if err != nil {
		return errToStatus(err), err
	}

	if file.IsDir {
		file.Listing.Sorting = d.user.Sorting
		file.Listing.ApplySort()
		return renderJSON(w, r, file)
	}

	if checksum := r.URL.Query().Get("checksum"); checksum != "" {
		err := file.Checksum(checksum)
		if err == errors.ErrInvalidOption {
			return http.StatusBadRequest, nil
		} else if err != nil {
			return http.StatusInternalServerError, err
		}

		// do not waste bandwidth if we just want the checksum
		file.Content = ""
	}

	return renderJSON(w, r, file)
})

var resourceDeleteHandler = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	if r.URL.Path == "/" || !d.user.Perm.Delete {
		return http.StatusForbidden, nil
	}
	err := d.user.Fs.RemoveAll(r.URL.Path)
	if err != nil {
		return errToStatus(err), err
	}

	return http.StatusOK, nil
})

var resourcePostPutHandler = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	if !d.user.Perm.Create && r.Method == http.MethodPost {
		return http.StatusForbidden, nil
	}

	if !d.user.Perm.Modify && r.Method == http.MethodPut {
		return http.StatusForbidden, nil
	}

	defer func() {
		io.Copy(ioutil.Discard, r.Body)
	}()

	// For directories, only allow POST for creation.
	if strings.HasSuffix(r.URL.Path, "/") {
		if r.Method == http.MethodPut {
			return http.StatusMethodNotAllowed, nil
		}

		err := d.user.Fs.MkdirAll(r.URL.Path, 775);

		return errToStatus(err), err
	}

	if r.Method == http.MethodPost && r.URL.Query().Get("override") != "true" {
		if _, err := d.user.Fs.Stat(r.URL.Path); err == nil {
			return http.StatusConflict, nil
		}
	}

	file, err := d.user.Fs.OpenFile(r.URL.Path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0775)
	if err != nil {
		return errToStatus(err), err
	}
	defer file.Close()

	_, err = io.Copy(file, r.Body)
	if err != nil {
		return errToStatus(err), err
	}
	// Gets the info about the file.
	info, err := file.Stat()
	if err != nil {
		return errToStatus(err), err
	}

	etag := fmt.Sprintf(`"%x%x"`, info.ModTime().UnixNano(), info.Size())
	w.Header().Set("ETag", etag)
	return errToStatus(err), err
})

var resourcePatchHandler = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	src := r.URL.Path
	dst := r.URL.Query().Get("destination")
	action := r.URL.Query().Get("action")
	dst, err := url.QueryUnescape(dst)

	if err != nil {
		return errToStatus(err), err
	}

	if dst == "/" || src == "/" {
		return http.StatusForbidden, nil
	}

	switch action {
	case "copy":
		if !d.user.Perm.Create {
			return http.StatusForbidden, nil
		}
	case "rename":
	default:
		action = "rename"
		if !d.user.Perm.Rename {
			return http.StatusForbidden, nil
		}
	}
	if action == "copy" {
		err = fileutils.Copy( src, dst)
	}else {
		err = os.Rename(src, dst)
	}
	return errToStatus(err), err
})
