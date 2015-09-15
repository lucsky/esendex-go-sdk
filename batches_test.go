package esendex

import (
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBatches(t *testing.T) {
	const (
		startIndex         = 0
		count              = 15
		totalCount         = 15
		id                 = "messagebatchid"
		uri                = "messagebatchuri"
		batchSize          = 1
		persistedBatchSize = 1
		accountReference   = "EXHEYEYE"
		createdBy          = "efiwewe@example.com"
		name               = "my cool batch"
	)

	var (
		status       = map[string]int{"submitted": 1}
		createdAt    = time.Date(2012, 1, 1, 12, 0, 0, 0, time.UTC)
		createdAtStr = "2012-01-01T12:00:00Z"
	)

	h := newRecordingHandler(`<?xml version="1.0" encoding="utf-8"?>
<messagebatches startindex="`+strconv.Itoa(startIndex)+`" count="`+strconv.Itoa(count)+`" totalcount="`+strconv.Itoa(totalCount)+`" xmlns="http://api.esendex.com/ns/">
 <messagebatch id="`+id+`" uri="`+uri+`">
  <createdat>`+createdAtStr+`</createdat>
  <batchsize>`+strconv.Itoa(batchSize)+`</batchsize>
  <persistedbatchsize>`+strconv.Itoa(persistedBatchSize)+`</persistedbatchsize>
  <status>
   <acknowledged>0</acknowledged>
   <authorisationfailed>0</authorisationfailed>
   <connecting>0</connecting>
   <delivered>0</delivered>
   <failed>0</failed>
   <partiallydelivered>0</partiallydelivered>
   <rejected>0</rejected>
   <scheduled>0</scheduled>
   <sent>0</sent>
   <submitted>1</submitted>
   <validityperiodexpired>0</validityperiodexpired>
   <cancelled>0</cancelled>
  </status>
  <accountreference>`+accountReference+`</accountreference>
  <createdby>`+createdBy+`</createdby>
  <name>`+name+`</name>
 </messagebatch>
</messagebatches>`, 200, map[string]string{})
	s := httptest.NewServer(h)
	defer s.Close()

	client := New("user", "pass")
	client.BaseURL, _ = url.Parse(s.URL)

	result, err := client.Batches()

	assert := assert.New(t)

	assert.Nil(err)

	assert.Equal("GET", h.Request.Method)
	assert.Equal("/v1.1/messagebatches", h.Request.URL.String())

	if user, pass, ok := h.Request.BasicAuth(); assert.True(ok) {
		assert.Equal("user", user)
		assert.Equal("pass", pass)
	}

	assert.Equal(startIndex, result.StartIndex)
	assert.Equal(count, result.Count)
	assert.Equal(totalCount, result.TotalCount)

	if assert.Equal(1, len(result.Batches)) {
		batch := result.Batches[0]

		assert.Equal(id, batch.ID)
		assert.Equal(uri, batch.URI)
		assert.Equal(createdAt, batch.CreatedAt)
		assert.Equal(batchSize, batch.BatchSize)
		assert.Equal(persistedBatchSize, batch.PersistedBatchSize)
		assert.Equal(status, batch.Status)
		assert.Equal(accountReference, batch.AccountReference)
		assert.Equal(createdBy, batch.CreatedBy)
		assert.Equal(name, batch.Name)
	}
}

func TestBatchesWithPaging(t *testing.T) {
	h := newRecordingHandler(`<?xml version="1.0" encoding="utf-8"?>
<messagebatches startindex="4" count="10" totalcount="200" xmlns="http://api.esendex.com/ns/">
</messagebatches>`, 200, map[string]string{})
	s := httptest.NewServer(h)
	defer s.Close()

	client := New("user", "pass")
	client.BaseURL, _ = url.Parse(s.URL)

	_, err := client.Batches(Page(5, 10))

	assert := assert.New(t)

	assert.Nil(err)

	assert.Equal("GET", h.Request.Method)
	assert.Equal("/v1.1/messagebatches", h.Request.URL.Path)

	if user, pass, ok := h.Request.BasicAuth(); assert.True(ok) {
		assert.Equal("user", user)
		assert.Equal("pass", pass)
	}

	query := h.Request.URL.Query()
	assert.Equal("5", query.Get("startindex"))
	assert.Equal("10", query.Get("count"))
}

func TestBatch(t *testing.T) {
	const (
		startIndex         = 0
		count              = 15
		totalCount         = 15
		id                 = "messagebatchid"
		uri                = "messagebatchuri"
		batchSize          = 1
		persistedBatchSize = 1
		accountReference   = "EXHEYEYE"
		createdBy          = "efiwewe@example.com"
		name               = "my cool batch"
	)

	var (
		status       = map[string]int{"submitted": 1}
		createdAt    = time.Date(2012, 1, 1, 12, 0, 0, 0, time.UTC)
		createdAtStr = "2012-01-01T12:00:00Z"
	)

	h := newRecordingHandler(`<?xml version="1.0" encoding="utf-8"?>
<messagebatch id="`+id+`" uri="`+uri+`" xmlns="http://api.esendex.com/ns/">
 <createdat>`+createdAtStr+`</createdat>
 <batchsize>`+strconv.Itoa(batchSize)+`</batchsize>
 <persistedbatchsize>`+strconv.Itoa(persistedBatchSize)+`</persistedbatchsize>
 <status>
  <acknowledged>0</acknowledged>
  <authorisationfailed>0</authorisationfailed>
  <connecting>0</connecting>
  <delivered>0</delivered>
  <failed>0</failed>
  <partiallydelivered>0</partiallydelivered>
  <rejected>0</rejected>
  <scheduled>0</scheduled>
  <sent>0</sent>
  <submitted>1</submitted>
  <validityperiodexpired>0</validityperiodexpired>
  <cancelled>0</cancelled>
 </status>
 <accountreference>`+accountReference+`</accountreference>
 <createdby>`+createdBy+`</createdby>
 <name>`+name+`</name>
</messagebatch>`, 200, map[string]string{})
	s := httptest.NewServer(h)
	defer s.Close()

	client := New("user", "pass")
	client.BaseURL, _ = url.Parse(s.URL)

	result, err := client.Batch(id)

	assert := assert.New(t)

	assert.Nil(err)

	assert.Equal("GET", h.Request.Method)
	assert.Equal("/v1.1/messagebatches/"+id, h.Request.URL.String())

	if user, pass, ok := h.Request.BasicAuth(); assert.True(ok) {
		assert.Equal("user", user)
		assert.Equal("pass", pass)
	}

	assert.Equal(id, result.ID)
	assert.Equal(uri, result.URI)
	assert.Equal(createdAt, result.CreatedAt)
	assert.Equal(batchSize, result.BatchSize)
	assert.Equal(persistedBatchSize, result.PersistedBatchSize)
	assert.Equal(status, result.Status)
	assert.Equal(accountReference, result.AccountReference)
	assert.Equal(createdBy, result.CreatedBy)
	assert.Equal(name, result.Name)
}

func TestCancelBatch(t *testing.T) {
	h := newRecordingHandler("", 204, map[string]string{})
	s := httptest.NewServer(h)
	defer s.Close()

	client := New("user", "pass")
	client.BaseURL, _ = url.Parse(s.URL)

	err := client.CancelBatch("batchid")

	assert := assert.New(t)

	assert.Nil(err)

	assert.Equal("DELETE", h.Request.Method)
	assert.Equal("/v1.1/messagebatches/batchid/schedule", h.Request.URL.String())

	if user, pass, ok := h.Request.BasicAuth(); assert.True(ok) {
		assert.Equal("user", user)
		assert.Equal("pass", pass)
	}
}
