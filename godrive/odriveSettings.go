package godrive

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type odriveSettings struct {
	MountPaths           []string
	Autounsyncthreshold  unsyncPeriod
	Placeholderthreshold placeholderSize
	TrashClean           trashCleanFrequency
	Xlthreshold          splitSize
}

//go:generate stringer -type=trashCleanFrequency -trimprefix=trashCleanFrequency
type trashCleanFrequency int

const (
	trashCleanFrequencyNever trashCleanFrequency = iota
	trashCleanFrequencyImmediately
	trashCleanFrequencyFifteen
	trashCleanFrequencyHour
	trashCleanFrequencyDay
)

//go:generate stringer -type=splitSize -trimprefix=splitSize
type splitSize int

const (
	splitSizeSmall splitSize = iota
	splitSizeMedium
	splitSizeLarge
	splitSizeXlarge
)

//go:generate stringer -type=unsyncPeriod -trimprefix=unsyncPeriod
type unsyncPeriod int

const (
	unsyncPeriodNever unsyncPeriod = iota
	unsyncPeriodDay
	unsyncPeriodWeek
	unsyncPeriodMonth
)

//go:generate stringer -type=placeholderSize -trimprefix=placeholderSize
type placeholderSize int

const (
	placeholderSizeNever placeholderSize = iota
	placeholderSizeSmall
	placeholderSizeMedium
	placeholderSizeLarge
	placeholderSizeAlways
)

func getTrashFrequencyIndexFromString(value string) trashCleanFrequency {
	for i := trashCleanFrequency(0); i < trashCleanFrequency(len(_trashCleanFrequency_index)-1); i++ {
		if strings.Contains(i.String(), value) {
			return i
		}
	}
	return -1
}

func trashCleanFrequencyElements() []string {
	elements := []string{}
	for i := trashCleanFrequency(0); i < trashCleanFrequency(len(_trashCleanFrequency_index)-1); i++ {
		elements = append(elements, i.String())
	}
	return elements
}

func getSplitSiteIndexFromString(value string) splitSize {
	for i := splitSize(0); i < splitSize(len(_splitSize_index)-1); i++ {
		if strings.Contains(i.String(), value) {
			return i
		}
	}
	return -1
}

func splitSizeElements() []string {
	elements := []string{}
	for i := splitSize(0); i < splitSize(len(_splitSize_index)-1); i++ {
		elements = append(elements, i.String())
	}
	return elements
}

func getPlaceholderSizeIndexFromString(value string) placeholderSize {
	for i := placeholderSize(0); i < placeholderSize(len(_placeholderSize_index)-1); i++ {
		if strings.Contains(i.String(), value) {
			return i
		}
	}
	return -1
}

func placeholderSizeElements() []string {
	elements := []string{}
	for i := placeholderSize(0); i < placeholderSize(len(_placeholderSize_index)-1); i++ {
		elements = append(elements, i.String())
	}
	return elements
}

func getUnsyncPeriodIndexFromString(value string) unsyncPeriod {
	for i := unsyncPeriod(0); i < unsyncPeriod(len(_unsyncPeriod_index)-1); i++ {
		if strings.Contains(i.String(), value) {
			return i
		}
	}
	return -1
}

func unsyncPeriodElements() []string {
	elements := []string{}
	for i := unsyncPeriod(0); i < unsyncPeriod(len(_unsyncPeriod_index)-1); i++ {
		elements = append(elements, i.String())
	}
	return elements
}

func LoadSettings() odriveSettings {
	viper.SetConfigName("godrive")
	viper.SetConfigType("toml")
	viper.AddConfigPath("/etc/godrive/")
	viper.AddConfigPath("$HOME/.appname")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	viper.SetDefault("mountPaths", []string{})
	viper.SetDefault("autounsyncthreshold", unsyncPeriodNever)
	viper.SetDefault("placeholderSize", placeholderSizeNever)
	viper.SetDefault("trashClean", trashCleanFrequencyNever)
	viper.SetDefault("splitSize", splitSizeXlarge)

	if err != nil { // Handle errors reading the config file
		fmt.Printf("Unable to read values from config. %s\n", err)
	}
	return parseFromViperToSettings()
}

func SaveSettings(settings odriveSettings) {
	parseSettingsToViper(settings)

	viper.SetConfigName("godrive")
	viper.SetConfigType("toml")
	viper.AddConfigPath("/etc/godrive/")
	viper.AddConfigPath("$HOME/.appname")
	viper.AddConfigPath(".")
	err := viper.WriteConfig()
	if err != nil {
		fmt.Printf("unable to save config file. %s\n", err)
	}
}

func parseFromViperToSettings() odriveSettings {
	return odriveSettings{
		MountPaths:           viper.GetStringSlice("mountPaths"),
		Autounsyncthreshold:  getUnsyncPeriodIndexFromString(viper.GetString("autounsyncthreshold")),
		Placeholderthreshold: getPlaceholderSizeIndexFromString(viper.GetString("placeholderthreshold")),
		TrashClean:           getTrashFrequencyIndexFromString(viper.GetString("trashClean")),
		Xlthreshold:          getSplitSiteIndexFromString(viper.GetString("splitSize")),
	}
}

func parseSettingsToViper(settings odriveSettings) {
	viper.Set("mountPaths", settings.MountPaths)
	viper.Set("autounsyncthreshold", settings.Autounsyncthreshold)
	viper.Set("placeholderSize", settings.Placeholderthreshold)
	viper.Set("trashClean", settings.TrashClean)
	viper.Set("splitSize", settings.Xlthreshold)
}
