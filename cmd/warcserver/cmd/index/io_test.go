package index

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/nlnwa/gowarcserver/pkg/index"
	log "github.com/sirupsen/logrus"
)

func TestParseFormat(t *testing.T) {
	tests := []struct {
		name       string
		format     string
		expected   reflect.Type
		errorState bool
	}{
		{
			"'cdx' results in CdxLegacy writer",
			"cdx",
			reflect.TypeOf((*index.CdxLegacy)(nil)),
			false,
		},
		{
			"'cdxj' results in CdxJ writer",
			"cdxj",
			reflect.TypeOf((*index.CdxJ)(nil)),
			false,
		},
		{
			"'db' results in CdxDb writer",
			"db",
			reflect.TypeOf((*index.CdxDb)(nil)),
			false,
		},
		{
			"'cdxpb' results in CdxPd writer",
			"cdxpb",
			reflect.TypeOf((*index.CdxPb)(nil)),
			false,
		},
		{
			"'cd' results in error",
			"cd",
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFormat(tt.format)
			if err != nil && !tt.errorState {
				t.Errorf("Unexpected failure: %v", err)
			} else if err == nil && tt.errorState {
				t.Errorf("Expected error parsing '%v', got type %T", tt.format, got)
			}

			if reflect.TypeOf(got) != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, got)
			}
		})
	}
}

// TODO: this was hard to write tests for and therefore ReadFile
//		 should probably be refactored
func TestReadFile(t *testing.T) {
	log.SetLevel(log.WarnLevel)
	// same as testdata/example.warc except removed gzip content because of illegal go str characters
	testFileContent := []byte(`WARC/1.0
WARC-Date: 2017-03-06T04:03:53Z
WARC-Record-ID: <urn:uuid:e9a0cecc-0221-11e7-adb1-0242ac120008>
WARC-Filename: temp-20170306040353.warc.gz
WARC-Type: warcinfo
Content-Type: application/warc-fields
Content-Length: 249

software: Webrecorder Platform v3.7
format: WARC File Format 1.0
creator: temp-MJFXHZ4S
isPartOf: Temporary%20Collection
json-metadata: {"title": "Temporary Collection", "size": 2865, "created_at": 1488772924, "type": "collection", "desc": ""}


WARC/1.0
WARC-Date: 2017-03-06T04:03:53Z
WARC-Record-ID: <urn:uuid:e9a0ee48-0221-11e7-adb1-0242ac120008>
WARC-Filename: temp-20170306040353.warc.gz
WARC-Type: warcinfo
Content-Type: application/warc-fields
Content-Length: 470

software: Webrecorder Platform v3.7
format: WARC File Format 1.0
creator: temp-MJFXHZ4S
isPartOf: Temporary%20Collection/Recording%20Session
json-metadata: {"created_at": 1488772924, "type": "recording", "updated_at": 1488773028, "title": "Recording Session", "size": 2865, "pages": [{"url": "http://example.com/", "title": "Example Domain", "timestamp": "20170306040348"}, {"url": "http://example.com/", "title": "Example Domain", "timestamp": "20170306040206"}]}


WARC/1.0
WARC-Target-URI: http://example.com/
WARC-Date: 2017-03-06T04:02:06Z
WARC-Type: response
WARC-Record-ID: <urn:uuid:a9c51e3e-0221-11e7-bf66-0242ac120005>
WARC-IP-Address: 93.184.216.34
WARC-Block-Digest: sha1:DR5MBP7OD3OPA7RFKWJUD4CTNUQUGFC5
WARC-Payload-Digest: sha1:G7HRM7BGOKSKMSXZAHMUQTTV53QOFSMK
Content-Type: application/http; msgtype=response
Content-Length: 975

HTTP/1.1 200 OK
Content-Encoding: gzip
Accept-Ranges: bytes
Cache-Control: max-age=604800
Content-Type: text/html
Date: Mon, 06 Mar 2017 04:02:06 GMT
Etag: "359670651+gzip"
Expires: Mon, 13 Mar 2017 04:02:06 GMT
Last-Modified: Fri, 09 Aug 2013 23:54:35 GMT
Server: ECS (iad/182A)
Vary: Accept-Encoding
X-Cache: HIT
Content-Length: 606
Connection: close



WARC/1.0
WARC-Type: request
WARC-Record-ID: <urn:uuid:a9c5c23a-0221-11e7-8fe3-0242ac120007>
WARC-Target-URI: http://example.com/
WARC-Date: 2017-03-06T04:02:06Z
WARC-Concurrent-To: <urn:uuid:a9c51e3e-0221-11e7-bf66-0242ac120005>
Content-Type: application/http; msgtype=request
Content-Length: 493

GET / HTTP/1.0
Host: example.com
Accept-Language: en-US,en;q=0.8,ru;q=0.6
Referer: https://webrecorder.io/temp-MJFXHZ4S/temp/recording-session/record/http://example.com/
Upgrade-Insecure-Requests: 1
Connection: close
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
Dnt: 1
Accept-Encoding: gzip, deflate, sdch, br



WARC/1.0
WARC-Type: request
WARC-Record-ID: <urn:uuid:e6e41fea-0221-11e7-8fe3-0242ac120007>
WARC-Target-URI: http://example.com/
WARC-Date: 2017-03-06T04:03:48Z
WARC-Concurrent-To: <urn:uuid:e6e395ca-0221-11e7-a18d-0242ac120005>
Content-Type: application/http; msgtype=request
Content-Length: 493

GET / HTTP/1.0
Host: example.com
Accept-Language: en-US,en;q=0.8,ru;q=0.6
Referer: https://webrecorder.io/temp-MJFXHZ4S/temp/recording-session/record/http://example.com/
Upgrade-Insecure-Requests: 1
Connection: close
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
Dnt: 1
Accept-Encoding: gzip, deflate, sdch, br

`)

	filepath := fmt.Sprintf("%s/test.warc", t.TempDir())
	file, err := os.Create(filepath)
	if err != nil {
		t.Fatalf("Failed to create testfile at '%s'", filepath)
	}
	// This is not strictly needed because of tmp, but to be platform agnostic it might be a good idea
	defer file.Close()

	_, err = file.Write(testFileContent)
	if err != nil {
		t.Fatalf("Failed to write to testfile at '%s'", filepath)
	}

	err = file.Sync()
	if err != nil {
		t.Fatalf("Failed to sync testfile at '%s'", filepath)
	}

	tests := []struct {
		writerFormat string
		writer       index.CdxWriter
	}{
		{
			"cdx",
			&index.CdxLegacy{},
		},
		{
			"cdxj",
			&index.CdxJ{},
		},
		{

			"cdxpd",
			&index.CdxPb{},
		},
		{
			"db",
			&index.CdxDb{},
		},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%T successfully indexes", tt.writer)
		t.Run(testName, func(t *testing.T) {
			c := &conf{
				filepath,
				tt.writerFormat,
				tt.writer,
			}
			c.writer.Init()
			defer c.writer.Close()

			err := ReadFile(c)
			if err != nil {
				t.Errorf("Unexpected failure: %v", err)
			}

		})
	}
}
