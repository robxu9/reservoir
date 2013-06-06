package main

type Metadatable interface {
	GetMetadata() map[string]string
}