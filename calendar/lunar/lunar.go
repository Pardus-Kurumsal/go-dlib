/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package lunar

import (
	"pkg.deepin.io/lib/calendar/util"
	"time"
)

// Calendar 保存公历年内计算农历所需的信息
type Calendar struct {
	Year              int            // 公历年份
	SolarTermJDs      *[25]float64   // 相关的 25 节气 北京时间 儒略日
	SolarTermTimes    *[25]time.Time // 对应 SolarTermJDs 转换为 time.Time 的时间
	NewMoonJDs        *[15]float64   // 相关的 15 个朔日 北京时间 儒略日
	Months            []*Month       // 月
	solarTermYearDays []int          // 十二节的 yearDay 列表
}

// Month 保存农历月信息
type Month struct {
	Name     int       // 农历月名
	Days     int       // 本月天数
	ShuoJD   float64   // 本月朔日时间 北京时间 儒略日
	ShuoTime time.Time // 本月朔日时间 北京时间
	IsLeap   bool      // 是否为闰月
}

func newCalendar(year int) *Calendar {
	cc := &Calendar{Year: year}
	cc.Months = make([]*Month, 0, 13)
	cc.calcProcData()
	cc.fillMonths()
	cc.calcLeapMonth()
	return cc
}

var ccCache = make(map[int]*Calendar)

// New 从缓存获取 Calendar 对象，没有则先创建
func New(year int) *Calendar {
	if cc, ok := ccCache[year]; ok {
		return cc
	} else {
		cc := newCalendar(year)
		ccCache[year] = cc
		return cc
	}
}

func (cc *Calendar) calcProcData() {
	//计算从上一年冬至开始到当年冬至全部25个节气
	cc.SolarTermJDs = get25SolarTermJDs(cc.Year-1, DongZhi)
	var solarTermTimes [25]time.Time
	for i := 0; i < 25; i++ {
		solarTermTimes[i] = util.GetDateTimeFromJulianDay(cc.SolarTermJDs[i])
	}
	cc.SolarTermTimes = &solarTermTimes

	var solarTermYearDays []int
	for i := 1; i < 25; i += 2 {
		yd := cc.SolarTermTimes[i].YearDay()
		solarTermYearDays = append(solarTermYearDays, yd)
	}
	cc.solarTermYearDays = solarTermYearDays

	tmpNewMoonJD := getNewMoonJD(util.JDBeijingTime2UTC(cc.SolarTermJDs[0]))
	if tmpNewMoonJD > cc.SolarTermJDs[0] {
		tmpNewMoonJD -= 29.53
	}
	cc.NewMoonJDs = get15NewMoonJDs(tmpNewMoonJD)

}

func get25SolarTermJDs(year, start int) *[25]float64 {
	// 从某一年的某个节气开始，连续计算25个节气，返回各节气的儒略日
	// year 年份
	// start 起始的节气
	// 返回 25 个节气的 儒略日北京时间
	stOrder := start
	var list [25]float64
	for i := 0; i < 25; i++ {
		jd := GetSolarTermJD(year, stOrder)
		list[i] = util.JDUTC2BeijingTime(jd)
		if stOrder == DongZhi {
			year++
		}
		stOrder = (stOrder + 1) % 24
	}
	return &list
}

func get15NewMoonJDs(jd float64) *[15]float64 {
	// 计算从某个时间之后的连续15个朔日
	// 参数: jd 开始时间的 儒略日
	// 返回 15个朔日时间 数组指针 儒略日北京时间
	var list [15]float64
	for i := 0; i < 15; i++ {
		newMoonJD := getNewMoonJD(jd)
		list[i] = util.JDUTC2BeijingTime(newMoonJD)
		// 转到下一个最接近朔日的时间
		jd = newMoonJD + 29.53
	}
	return &list
}

func deltaDays(t1, t2 time.Time) int {
	// 计算两个时间相差的天数
	// t2 > t1 结果为正数
	date1 := time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.UTC)
	date2 := time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, time.UTC)
	dd := int((date2.Unix() - date1.Unix()) / 86400)
	// fmt.Printf("%v <=> %v dd %v\n", t2, t1, dd)
	return dd
}

func (cc *Calendar) fillMonths() {
	//采用夏历建寅，冬至所在月份为农历11月
	yuejian := 11
	for i := 0; i < 14; i++ {
		info := new(Month)
		if yuejian <= 12 {
			info.Name = yuejian
		} else {
			info.Name = yuejian - 12
		}
		info.ShuoJD = cc.NewMoonJDs[i]
		info.ShuoTime = util.GetDateTimeFromJulianDay(info.ShuoJD)
		nextShuoJD := cc.NewMoonJDs[i+1]
		info.Days = int(nextShuoJD+0.5) - int(info.ShuoJD+0.5)
		cc.Months = append(cc.Months, info)
		yuejian++
	}
}

func (cc *Calendar) calcLeapMonth() {
	// 根据节气计算是否有闰月，如果有闰月，根据农历月命名规则，调整月名称
	if int(cc.NewMoonJDs[13]+0.5) <= int(cc.SolarTermJDs[24]+0.5) {
		// 第13月的月末没有超过冬至，说明今年需要闰一个月
		i := 1
		for i < 14 {
			if int(cc.NewMoonJDs[i+1]+0.5) <= int(cc.SolarTermJDs[2*i]+0.5) {
				/* cc.NewMoonJDs[i + 1] 是第i个农历月的下一个月的月首
				   本该属于第i个月的中气如果比下一个月的月首还晚，或者与下个月的月首是同一天（民间历法），则说明第 i 个月没有中气, 是闰月 */
				break
			}
			i++
		}
		if i < 14 {
			// 找到闰月
			// fmt.Println("找到闰月 ", i)
			cc.Months[i].IsLeap = true
			// 对后面的农历月调整月名
			for i < 14 {
				cc.Months[i].Name--
				i++
			}
		}
	}
}

// SolarDayToLunarDay 指定年份内公历日期转换为农历日
func (cc *Calendar) SolarDayToLunarDay(month, day int) *Day {
	dt := time.Date(cc.Year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	yd := dt.YearDay()

	// 求月地支
	monthZhi := 0
	for monthZhi < len(cc.solarTermYearDays) {
		if yd >= cc.solarTermYearDays[monthZhi] {
			monthZhi++
		} else {
			break
		}
	}

	// 求农历月份和日
	var lunarMonth *Month
	var lunarDay int
	for _, m := range cc.Months {
		dd := deltaDays(m.ShuoTime, dt) + 1
		if 1 <= dd && dd <= m.Days {
			lunarMonth = m
			lunarDay = dd
			break
		}
	}
	// 求二十四节气
	solarTerm := -1
	solarTermInfos := cc.getMonthSolarTerms(month)
	if day == solarTermInfos[0].Day {
		solarTerm = solarTermInfos[0].SolarTerm
	} else if day == solarTermInfos[1].Day {
		solarTerm = solarTermInfos[1].SolarTerm
	}

	return &Day{
		Year:      cc.Year,
		Day:       lunarDay,
		Month:     lunarMonth,
		MonthZhi:  monthZhi,
		SolarTerm: solarTerm,
	}
}

type solarTermInfo struct {
	Day       int
	SolarTerm int
}

func (cc *Calendar) getSolarTermInfo(index int) *solarTermInfo {
	dt := cc.SolarTermTimes[index]
	day := dt.Day()
	stIndex := (index + DongZhi) % 24
	return &solarTermInfo{day, stIndex}
}

func (cc *Calendar) getMonthSolarTerms(month int) [2]*solarTermInfo {
	var list [2]*solarTermInfo
	index := 2*month - 1
	list[0] = cc.getSolarTermInfo(index)
	list[1] = cc.getSolarTermInfo(index + 1)
	return list
}

var (
	// 十二生肖
	Animals = []string{"鼠", "牛", "虎", "兔", "龙", "蛇", "马", "羊", "猴", "鸡", "狗", "猪"}
	// 天干
	Gan = []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
	// 地支
	Zhi = []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}
)

// GetYearZodiac 计算年份的生肖
func GetYearZodiac(year int) string {
	return Animals[(year-4)%12]
}

func cyclical(num int) string {
	// 将数字转换为天干地支表示
	return Gan[num%10] + Zhi[num%12]
}

// GetYearGanZhi 计算年份的干支
func GetYearGanZhi(year int) string {
	return cyclical(year - 1864)
}

// GetDayGanZhi 计算日干支
func GetDayGanZhi(year, month, day int) string {
	unixTime := time.Date(int(year), time.Month(month), int(day),
		0, 0, 0, 0, time.UTC).Unix()
	dayCyclical := int(unixTime/86400) + 29219 + 18
	return cyclical(dayCyclical)
}
