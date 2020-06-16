package master

import (
	"gotest.tools/assert/cmp"
	"k8s.io/apimachinery/pkg/types"
)

type Merger struct {
	oldMaster Master
	newMaster Master
	namespacedName types.NamespacedName
}

func (merger *Merger) MergeSettings()(master Master, mergedSettings, mergedScaleSettings bool){

	oldSettings := merger.oldMaster.settings
	newSettings := merger.newMaster.settings

	scaleEqual := cmp.DeepEqual(oldSettings.MasterScaleSettings, newSettings.MasterScaleSettings)()

	if scaleEqual.Success() {
		mergedScaleSettings = true
		oldSettings.MasterScaleSettings = newSettings.MasterScaleSettings
	}

	settingsEqual := cmp.DeepEqual(oldSettings.MasterClusterSettings, newSettings.MasterClusterSettings)()

	if settingsEqual.Success() {
		mergedSettings = true
		oldSettings.MasterClusterSettings = newSettings.MasterClusterSettings
	}

	master = NewMaster(merger.namespacedName, oldSettings)
	return
}