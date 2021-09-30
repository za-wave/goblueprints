package main

import (
	"errors"
	"io/ioutil"
	"path/filepath"
)

// ErrNoAvatarURL is the error that is returned when the
// Avatar instance is unable to provide an avatar URL.
var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL")

// Avatar represents types capable of representing
// user profile pictures.
type Avatar interface {
	// GetAvatarURL gets the avatar URL for the specified client,
	// or returns an error if something goes wrong.
	// ErrNoAvatarURL is returned if the object is unable to get
	// a URL for the specified client.
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (_ AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

type GravatarAvatar struct{}

var UseGravatar GravatarAvatar

func (_ GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			return "//www.gravatar.com/avatar/" + useridStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

func (_ FileSystemAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			files, err := ioutil.ReadDir("avatars")
			if err == nil {
				for _, file := range files {
					if file.IsDir() {
						continue
					}
					match, _ := filepath.Match(useridStr+"*", file.Name())
					if match {
						return "/avatars/" + file.Name(), nil
					}
				}
			}
		}
	}
	return "", ErrNoAvatarURL
}

// func (FileSystemAvatar) GetAvatarURL(c *client) (string, error) {
// 	files, err := ioutil.ReadDir("avatars")
// 	if err != nil {
// 		return "", ErrNoAvatarURL
// 	}
// 	for _, file := range files {
// 		if file.IsDir() {
// 			continue
// 		}
// 		fname := file.Name()
// 		if c.userData["userid"] == strings.TrimSuffix(fname, filepath.Ext(fname)) {
// 			return "/avatars/" + fname, nil
// 		}
// 	}
// 	return "", ErrNoAvatarURL
// }
