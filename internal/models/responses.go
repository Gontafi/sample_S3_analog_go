package models

import "encoding/xml"

type ListAllMyBucketsResult struct {
	XMLName xml.Name `xml:"ListAllMyBucketsResult"`
	Buckets Buckets  `xml:"Buckets"`
}

type Buckets struct {
	Bucket []Bucket `xml:"Bucket"`
}

type Bucket struct {
	CreationDate string `xml:"CreationDate"`
	Name         string `xml:"Name"`
}

type ErrResponse struct {
	Error struct {
		Code    int    `xml:"Code"`
		Message string `xml:"Message"`
	} `xml:"Error"`
}
