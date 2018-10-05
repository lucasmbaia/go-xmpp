package utils

import (
  "encoding/xml"
  "regexp"
  "bytes"
)

const (
  WITHOUT_ENT_TAG = `<\/[^>]+>$`
  SELF_CLOSING    = `><\/[^>]+>$`
)

func MarshalSelfClosingTag(i interface{}) ([]byte, error) {
  var (
    body []byte
    err  error
    r    *regexp.Regexp
    dst  []byte
  )

  if body, err = xml.Marshal(i); err != nil {
    return body, err
  }

  r = regexp.MustCompile(SELF_CLOSING)

  dst = r.ReplaceAll(body, []byte(" />"))
  return dst, nil
}

func MarshalWithOutEndTag(i interface{}, replaceQuotes bool) ([]byte, error) {
  var (
    body []byte
    err  error
    r    *regexp.Regexp
    dst  []byte
  )

  if body, err = xml.Marshal(i); err != nil {
    return body, err
  }

  r = regexp.MustCompile(WITHOUT_ENT_TAG)

  dst = bytes.Replace(r.ReplaceAll(body, []byte("")), []byte("\""), []byte("'"), -1)
  return dst, nil
}
