---
sidebar_position: 1
---

# File storages

By default we support only file storage. Later on we will provide storages to S3, etc.  
It could be used like here:
```go
if storage == nil {
  storage = NewFsStorage()
}
ret := make([]string, 0)
var filename string
// files from POST request
for _, file := range files {
  f, err := file.Open()
  if err != nil {
    return err
  }
  bytecontent := make([]byte, file.Size)
  _, err = f.Read(bytecontent)
  if err != nil {
    return err
  }
  filename, err = storage.Save(&FileForStorage{
    Content:           bytecontent,
    PatternForTheFile: "*." + strings.Split(file.Filename, ".")[1],
    Filename:          file.Filename,
  })
  if err != nil {
    return err
  }
  err = f.Close()
  if err != nil {
    return err
  }
  ret = append(ret, filename)
}
```
