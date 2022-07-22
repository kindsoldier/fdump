/*
 * Copyright 2020 Oleg Borodin  <borodin@unix7.org>
 */

package dscron

import (
    "regexp"
    "strings"
    "strconv"
    "time"
)


func Match(mdays, wdays, hours, mins string, ts time.Time) bool {
    hour    := ts.Hour()
    min     := ts.Minute()
    mday    := ts.Day()
    wday    := int(ts.Weekday())

    mapMdays := Expander(mdays, 1, 31) // 1..31
    if mapMdays[mday] != true {
        return false
    }

    mapWdays := Expander(wdays, 1, 7)  // 1..7
    if mapWdays[wday] != true {
        return false
    }

    mapHours := Expander(hours, 0, 23) // 0..23
    if mapHours[hour] != true {
        return false
    }

    mapMins := Expander(mins, 0, 59)  // 0..59
    if mapMins[min] != true {
        return false
    }
    return true
}


func Expander(items string, min, max int) map[int]bool  {

    items = strings.ReplaceAll(items, "--", "-")
    items = strings.ReplaceAll(items, ",,", ",")
    items = strings.ReplaceAll(items, "//", "/")


    field := make(map[int]bool)

    for i := min; i <= max; i++ {
        field[i] = false
    }

    for _, item := range strings.Split(items, ",") {
        item := strings.Trim(item, " ")

        if len(item) == 0 {
            continue
        }

        var match bool

        match, _ = regexp.MatchString(`^\*$`, item)
        if match {
            for i := min; i <= max; i++ {
                field[i] = true
            }
        }

        match, _ = regexp.MatchString(`^[0-9]+$`, item)
        if match {
            num, _ := strconv.Atoi(item)
            if num >= min && num <= max {
                field[num] = true
            }
        }

        match, _ = regexp.MatchString(`^\*/[0-9]+$`, item)
        if match {
            arr := strings.Split(item, "/")
            period, _ := strconv.Atoi(arr[1])

            for i := min; i <= max; i++ {
                rep := i % period
                if rep == 0 {
                    field[i] = true
                }
            }
        }

        match, _ = regexp.MatchString(`^[0-9]+-[0-9]+$`, item)
        if match {
            arr := strings.Split(item, "-")
            begin, _ := strconv.Atoi(arr[0])
            end, _ := strconv.Atoi(arr[1])

            for i := begin; i < end + 1 && i <= max; i++ {
                field[i] = true
            }
        }
    }

    return field
}
