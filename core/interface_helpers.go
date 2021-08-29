package core

import (
	"net/url"
	"reflect"
)

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func GenerateNumberSequence(start int, stop int) <- chan int {
	chnl := make(chan int)
	go func() {
		i := start
		for true {
			chnl <- i
			if i == stop {
				close(chnl)
				break
			}
			i += 1
		}
	}()
	return chnl
}

func CloneNetUrl(url1 *url.URL) *url.URL {
	clonedUrl := &url.URL{
		Scheme: url1.Scheme,
		Opaque: url1.Opaque,
		Host: url1.Host,
		Path: url1.Path,
		RawPath: url1.RawPath,
		ForceQuery: url1.ForceQuery,
		RawQuery: url1.RawQuery,
		Fragment: url1.Fragment,
		RawFragment: url1.RawFragment,
	}
	return clonedUrl
}

func GetID(m reflect.Value) uint {
	if m.Kind() == reflect.Ptr {
		return uint(m.Elem().FieldByName("ID").Uint())
	}
	return uint(m.FieldByName("ID").Uint())
}

func Remove(items []string, item string) []string {
	newitems := []string{}

	for _, i := range items {
		if i != item {
			newitems = append(newitems, i)
		}
	}

	return newitems
}
