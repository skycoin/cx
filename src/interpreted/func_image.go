package interpreted

import (
	"image"
)

var images map[string]image.Image = make(map[string]image.Image, 0)

func image_Decode(fileName *CXArgument) error {
	if err := checkType("image.Decode", "str", fileName); err == nil {
		name := string(*fileName.Value)

		// we'll ignore format for now
		if img, _, err := image.Decode(openFiles[name]); err == nil {
			images[name] = img
		} else {
			return err
		}

		return nil
	} else {
		return err
	}
}
