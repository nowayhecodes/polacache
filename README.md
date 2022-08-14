<div align="center">
    <img src="./asset/polacache.png"  width="300" alt="iters" />
</div>

<br />
<div align="center">
    <h1>Polacache</h1>
</div>

<div align="center" style="margin-top: -4rem;">
    <h3>A deadly simple and thread-safe map cache.</h3>
</div>
<br />

#### What?

Polacache is a deadly simple and thread-safe map cache. 
In it's constructor, you set a cleanupInterval, which launchs a goroutine to perform the cleanup loop.

#### Why?

For the fun

#### How? 

```go

package main

import (
    "time"
    
    pc "github.com/nowayhecodes/polacache"
)

func main() {
    cache := pc.New(1 * time.Minute)

    exampleItem := pc.Item{
        Key:   "example",
        Value: 42,
    }

    cache.Set(exampleItem, time.Now().Add(1*time.Hour).Unix())
    cache.Get(exampleItem.Key)
    cache.Delete(exampleItem.Key)

}

```
