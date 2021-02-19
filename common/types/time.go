//@Title time.go
//@Description 
//@Author bertang
//@Created bertang 2021/2/17 6:10 下午
package types

import (
    "bytes"
    "database/sql/driver"
    "github.com/bertang/bert/common/config/application"
    jsoniter "github.com/json-iterator/go"
    "time"
)

type BertTime struct {
    time.Time
}

//MarshalJSON json
func (b BertTime) MarshalJSON() ([]byte, error) {
    bb := b.Time
    if bb.IsZero() {
        return []byte("null"), nil
    }

    return jsoniter.Marshal(bb.Format(application.GetAppConf().TimeFormat))
}

// Value insert timestamp into mysql need this function.
func (b BertTime) Value() (driver.Value, error) {
    var zeroTime time.Time
    t := b.Time
    if t.UnixNano() == zeroTime.UnixNano() {
        return nil, nil
    }
    return t, nil
}

func (b *BertTime) UnmarshalJSON(data []byte) error {
    data = bytes.ReplaceAll(data, []byte("\""), []byte(""))
    if string(data) == "" || len(data) == 0 {
        return nil
    }
    temp, err := time.ParseInLocation(application.GetAppConf().TimeFormat, string(data), time.Local)
    if err != nil {
        return err
    }
    b.Time = temp
    return nil
}
