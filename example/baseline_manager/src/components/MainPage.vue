<template>
    <v-container class="fill-height">
        <v-responsive class="align-centerfill-height mx-auto">
            <v-card title="Option" class="px-5 ma-2">
                <template v-slot:append>
                    <v-btn density="compact" variant="plain" :icon="optionLocked ? 'mdi-lock' : 'mdi-lock-open-variant'"
                        @click="optionLocked = !optionLocked"></v-btn>
                </template>
                <v-select v-model="confStore.conf.option.output_format" label="Output format"
                    :items="output_format_items" :disabled="optionLocked"></v-select>
                <v-text-field v-model="confStore.conf.option.output_filename" label="Output filename"
                    :rules="[v => (!v.includes('/') && !v.includes('\\')) || 'Subdirectory not allowed']"
                    :disabled="optionLocked">
                </v-text-field>
                <v-combobox v-model="confStore.conf.option.output_metadata" label="Output metadata" chips multiple
                    :items="allMetadataKeys" :disabled="optionLocked"></v-combobox>
                <v-checkbox v-model="confStore.conf.option.output_risk_only" label="Output risk item only"
                    :disabled="optionLocked"></v-checkbox>
            </v-card>

            <v-card title="Profile" class="px-5 ma-2">
                <template v-slot:append>
                    <v-btn icon="mdi-plus" density="compact" variant="plain" :disabled="profileLocked || !canAddProfile"
                        @click="ShowAddProfile()"></v-btn>
                    <v-btn density="compact" variant="plain"
                        :icon="profileLocked ? 'mdi-lock' : 'mdi-lock-open-variant'"
                        @click="profileLocked = !profileLocked"></v-btn>
                </template>
                <v-empty-state v-if="Object.keys(confStore.conf.profile).length === 0"
                    text="No profile defined"></v-empty-state>
                <div v-for="_, key in confStore.conf.profile" class="ma-0 ga-5 d-flex align-center">
                    {{ key }}:
                    <v-text-field v-model="confStore.conf.profile[key]" density="compact"
                        :disabled="profileLocked"></v-text-field>
                    <v-btn icon="mdi-minus" density="compact" variant="plain" :disabled="profileLocked"
                        @click="ShowDeleteProfile(String(key))"></v-btn>
                </div>
            </v-card>

            <v-card title="Listor" class="px-5 ma-2">
                <template v-slot:append>
                    <v-text-field v-model="searchForListor" density="comfortable" placeholder="Search"
                        prepend-inner-icon="mdi-magnify" style="width: 300px;" clearable hide-details></v-text-field>
                    <v-menu>
                        <template v-slot:activator="{ props }">
                            <v-btn icon="mdi-dots-vertical" v-bind="props" density="comfortable"
                                variant="plain"></v-btn>
                        </template>
                        <v-list>
                            <v-list-item :disabled="selectedListor.length === 0"
                                @click="ShowDeleteListor(selectedListor)">
                                {{ selectedListor.length === 0 ? 'No listor selected' : 'Delete selected' }}
                            </v-list-item>
                        </v-list>
                    </v-menu>
                </template>
                <v-data-iterator v-model="selectedListor" :items="confStore.conf.listor" items-per-page="12"
                    :search="searchForListor">
                    <template v-slot:default="{ items, isSelected, toggleSelect, isExpanded, toggleExpand }">
                        <v-row>
                            <v-col v-for="(item, i) in items" :key="i" cols="12" sm="6" xl="4">
                                <v-card :title="item.raw.rs_type" class="px-5 pb-5 ma-2">
                                    <template v-slot:append>
                                        <v-checkbox :model-value="isSelected(item)" density="compact" hide-details
                                            @click="toggleSelect(item)"></v-checkbox>
                                        <v-btn :icon="isExpanded(item as any) ? 'mdi-code-braces' : 'mdi-view-list'"
                                            density="compact" variant="plain"
                                            @click="toggleExpand(item as any)"></v-btn>
                                        <v-menu>
                                            <template v-slot:activator="{ props }">
                                                <v-btn icon="mdi-dots-vertical" v-bind="props" density="compact"
                                                    variant="plain"></v-btn>
                                            </template>
                                            <v-list>
                                                <v-list-item @click="ShowDeleteListor([item.raw.id])">
                                                    Delete
                                                </v-list-item>
                                            </v-list>
                                        </v-menu>
                                    </template>
                                    <v-window :model-value="isExpanded(item as any) ? 1 : 0">
                                        <v-window-item>
                                            <v-table>
                                                <tr>
                                                    <td>Id:</td>
                                                    <td>
                                                        {{ item.raw.id }}
                                                        <v-tooltip v-if="IsListorDangling(item.raw)"
                                                            text="Invalid: Listor is dangling">
                                                            <template v-slot:activator="{ props }">
                                                                <v-icon v-bind="props" size="x-small"
                                                                    color="error">mdi-alert-circle</v-icon>
                                                            </template>
                                                        </v-tooltip>
                                                    </td>
                                                </tr>
                                                <tr style="height: 1px;"></tr>
                                                <tr>
                                                    <td>Cloud type:</td>
                                                    <td>{{ item.raw.cloud_type }}</td>
                                                </tr>
                                            </v-table>
                                        </v-window-item>
                                        <v-window-item>
                                            <highlightjs language="yaml" :code="item.raw.raw_conf"
                                                style="max-height: 10em; overflow-y: scroll;">
                                            </highlightjs>
                                        </v-window-item>
                                    </v-window>
                                </v-card>
                            </v-col>
                        </v-row>
                    </template>
                    <template v-slot:footer="{ page, pageCount, prevPage, nextPage }">
                        <div class="d-flex align-center justify-center px-4 pb-5">
                            <v-btn :disabled="page === 1" density="comfortable" icon="mdi-arrow-left" variant="tonal"
                                rounded @click="prevPage"></v-btn>
                            <div class="mx-2 text-caption">
                                Page {{ page }} of {{ pageCount }}
                            </div>
                            <v-btn :disabled="page >= pageCount" density="comfortable" icon="mdi-arrow-right"
                                variant="tonal" rounded @click="nextPage"></v-btn>
                        </div>
                    </template>
                </v-data-iterator>
            </v-card>

            <v-card title="Baseline" class="px-5 ma-2">
                <template v-slot:append>
                    <v-text-field v-model="searchForBaseline" density="comfortable" placeholder="Search"
                        prepend-inner-icon="mdi-magnify" style="width: 300px;" clearable hide-details></v-text-field>
                    <v-menu>
                        <template v-slot:activator="{ props }">
                            <v-btn icon="mdi-dots-vertical" v-bind="props" density="comfortable"
                                variant="plain"></v-btn>
                        </template>
                        <v-list>
                            <v-list-item :disabled="selectedBaseline.length === 0"
                                @click="ShowDeleteBaseline(selectedBaseline)">
                                {{ selectedBaseline.length === 0 ? 'No baseline selected' : 'Delete selected' }}
                            </v-list-item>
                        </v-list>
                    </v-menu>
                </template>
                <v-data-iterator v-model="selectedBaseline" :items="confStore.conf.baseline" items-per-page="12"
                    :search="searchForBaseline">
                    <template v-slot:default="{ items, isSelected, toggleSelect, isExpanded, toggleExpand }">
                        <v-row>
                            <v-col v-for="(item, i) in items" :key="i" cols="12" sm="6" xl="4">
                                <v-card class="px-5 pb-5 ma-2">
                                    <template v-slot:append>
                                        <v-checkbox :model-value="isSelected(item)" density="compact" hide-details
                                            @click="toggleSelect(item)"></v-checkbox>
                                        <v-btn :icon="isExpanded(item as any) ? 'mdi-code-braces' : 'mdi-view-list'"
                                            density="compact" variant="plain"
                                            @click="toggleExpand(item as any)"></v-btn>
                                        <v-menu>
                                            <template v-slot:activator="{ props }">
                                                <v-btn icon="mdi-dots-vertical" v-bind="props" density="compact"
                                                    variant="plain"></v-btn>
                                            </template>
                                            <v-list>
                                                <v-list-item @click="ShowDeleteBaseline([item.raw.id])">
                                                    Delete
                                                </v-list-item>
                                            </v-list>
                                        </v-menu>
                                    </template>
                                    <v-window :model-value="isExpanded(item as any) ? 1 : 0">
                                        <v-window-item>
                                            <v-table>
                                                <tr>
                                                    <td>Tag:</td>
                                                    <td>
                                                        <v-chip v-for="(tag, j) in item.raw.tag" :key="`tag_${i}_${j}`"
                                                            size="small" variant="outlined">
                                                            <v-tooltip :text="tag">
                                                                <template v-slot:activator="{ props }">
                                                                    <span v-bind="props"
                                                                        class="d-inline-block text-truncate"
                                                                        style="max-width: 5em;">{{ tag }}</span>
                                                                </template>
                                                            </v-tooltip>
                                                        </v-chip>
                                                    </td>
                                                </tr>
                                                <tr style="height: 1px;"></tr>
                                                <tr>
                                                    <td>Metadata:</td>
                                                    <td>
                                                        <v-chip v-for="(md_value, md_key, j) in item.raw.metadata"
                                                            :key="`md_${i}_${j}`" size="small" variant="outlined">
                                                            <v-tooltip :text="`${md_key}: ${md_value}`">
                                                                <template v-slot:activator="{ props }">
                                                                    <span v-bind="props"
                                                                        class="d-inline-block text-truncate"
                                                                        style="max-width: 5em;">
                                                                        {{ md_key }}: {{ md_value }}
                                                                    </span>
                                                                </template>
                                                            </v-tooltip>
                                                        </v-chip>
                                                    </td>
                                                </tr>
                                                <tr style="height: 1px;"></tr>
                                                <tr>
                                                    <td>Cloud type:</td>
                                                    <td>
                                                        {{ GetBaselineCheckerCloudType(item.raw).join(', ') }}
                                                    </td>
                                                </tr>
                                                <tr style="height: 1px;"></tr>
                                                <tr>
                                                    <td>Listor ids</td>
                                                    <td>
                                                        {{ GetBaselineCheckerListorIdStr(item.raw).join(', ') }}
                                                        <v-tooltip v-if="IsInvalidListorRefInBaseline(item.raw)"
                                                            :text="`Invalid: The following id not exist: ${GetInvalidListorRefInBaseline(item.raw).join(', ')}`">
                                                            <template v-slot:activator="{ props }">
                                                                <v-icon v-bind="props" size="x-small"
                                                                    color="error">mdi-alert-circle</v-icon>
                                                            </template>
                                                        </v-tooltip>
                                                    </td>
                                                </tr>
                                            </v-table>
                                        </v-window-item>
                                        <v-window-item>
                                            <highlightjs language="yaml" :code="item.raw.raw_conf"
                                                style="max-height: 10em; overflow-y: scroll;">
                                            </highlightjs>
                                        </v-window-item>
                                    </v-window>
                                </v-card>
                            </v-col>
                        </v-row>
                    </template>
                    <template v-slot:footer="{ page, pageCount, prevPage, nextPage }">
                        <div class="d-flex align-center justify-center px-4 pb-5">
                            <v-btn :disabled="page === 1" density="comfortable" icon="mdi-arrow-left" variant="tonal"
                                rounded @click="prevPage"></v-btn>
                            <div class="mx-2 text-caption">
                                Page {{ page }} of {{ pageCount }}
                            </div>
                            <v-btn :disabled="page >= pageCount" density="comfortable" icon="mdi-arrow-right"
                                variant="tonal" rounded @click="nextPage"></v-btn>
                        </div>
                    </template>
                </v-data-iterator>
                <div class="d-flex align-center">
                    Export selected baseline only
                    <v-switch v-model="exportAll" color="primary" class="mx-2" hide-details></v-switch>
                    Export all baseline and listor
                </div>
            </v-card>
        </v-responsive>
    </v-container>

    <!-- FAB -->
    <v-menu location="top">
        <template v-slot:activator="{ props }">
            <v-fab icon="mdi-menu" class="position-fixed right-0 bottom-0 ma-15" v-bind="props"
                :loading="exporting"></v-fab>
        </template>
        <v-list density="compact">
            <v-list-item title="Toggle theme" prepend-icon="mdi-theme-light-dark" @click="ToggleTheme()"></v-list-item>
            <v-divider></v-divider>
            <v-list-item title="Import file" prepend-icon="mdi-import" @click="ShowImport()"></v-list-item>
            <v-list-item title="Export file" prepend-icon="mdi-export" @click="TryExport()"></v-list-item>
            <v-divider></v-divider>
            <v-list-item title="Reset localStorage" prepend-icon="mdi-firework"
                @click="showDlgResetLocalStorage = true"></v-list-item>
        </v-list>
    </v-menu>

    <!-- Dialog - Add profile -->
    <v-dialog v-model="showDlgAddProfile" width="auto">
        <v-card width="400" title="Add profile">
            <v-card-text>
                <v-select v-model="dlgAddProfileKey" label="Cloud type" :items="profileKeyToAdd"></v-select>
            </v-card-text>
            <template v-slot:actions>
                <v-spacer></v-spacer>
                <v-btn class="ms-auto" text="Ok" :disabled="dlgAddProfileKey === ''"
                    @click="AddProfile(dlgAddProfileKey)"></v-btn>
                <v-btn class="ms-auto" text="Cancel" @click="showDlgAddProfile = false"></v-btn>
            </template>
        </v-card>
    </v-dialog>

    <!-- Dialog - Delete profile -->
    <v-dialog v-model="showDlgDeleteProfile" width="auto">
        <v-card width="400" title="Delete profile">
            <v-card-text>
                <span>Are you sure to delete profile </span>
                <span class="text-warning">{{ dlgDeleteProfileKey }}</span>
                <span>?</span>
            </v-card-text>
            <template v-slot:actions>
                <v-spacer></v-spacer>
                <v-btn class="ms-auto" text="DELETE" color="warning"
                    @click="DeleteProfile(dlgDeleteProfileKey)"></v-btn>
                <v-btn class="ms-auto" text="Cancel" @click="showDlgDeleteProfile = false"></v-btn>
            </template>
        </v-card>
    </v-dialog>

    <!-- Dialog - Delete listor -->
    <v-dialog v-model="showDlgDeleteListor" width="auto">
        <v-card width="400" title="Delete listor">
            <v-card-text>
                <div>Are you sure to delete following listor?</div>
                <div class="text-warning">{{ listorIdToDelete.filter(i => String(i)).join(', ') }}</div>
            </v-card-text>
            <template v-slot:actions>
                <v-spacer></v-spacer>
                <v-btn class="ms-auto" text="DELETE" color="warning" @click="DeleteListor(listorIdToDelete)"></v-btn>
                <v-btn class="ms-auto" text="Cancel" @click="showDlgDeleteListor = false"></v-btn>
            </template>
        </v-card>
    </v-dialog>

    <!-- Dialog - Delete baseline -->
    <v-dialog v-model="showDlgDeleteBaseline" width="auto">
        <v-card width="400" title="Delete baseline">
            <v-card-text>
                <div>
                    <span>Are you sure to delete </span>
                    <span class="text-warning">{{ baselineIdToDelete.length }}</span>
                    <span> baseline?</span>
                </div>
                <v-checkbox v-model="deleteBaselineWithListor" label="Also delete dangling listor"></v-checkbox>
            </v-card-text>
            <template v-slot:actions>
                <v-spacer></v-spacer>
                <v-btn class="ms-auto" text="DELETE" color="warning"
                    @click="DeleteBaseline(baselineIdToDelete)"></v-btn>
                <v-btn class="ms-auto" text="Cancel" @click="showDlgDeleteBaseline = false"></v-btn>
            </template>
        </v-card>
    </v-dialog>

    <!-- Dialog - Reset localStorage -->
    <v-dialog v-model="showDlgResetLocalStorage" width="auto">
        <v-card width="400" title="Reset localStorage">
            <v-card-text>
                <span class="text-warning">Are you sure to reset all data
                    stored in the localStorage? This command cannot be undone</span>
            </v-card-text>
            <template v-slot:actions>
                <v-spacer></v-spacer>
                <v-btn class="ms-auto" text="RESET" color="warning" @click="resetLocalStorage()"></v-btn>
                <v-btn class="ms-auto" text="Cancel" @click="showDlgResetLocalStorage = false"></v-btn>
            </template>
        </v-card>
    </v-dialog>

    <!-- Dialog - Import -->
    <v-dialog v-model="showDlgImport" width="auto" :persistent="importing">
        <v-card width="800" title="Import">
            <v-window v-model="stepImport">
                <v-window-item :value="1">
                    <v-card-text>
                        <v-file-input v-model="importFile" label="Select yaml configuration file to import" show-size
                            :accept="`${YAML_MIME},.conf,.yaml,.yml`" :disabled="importing"></v-file-input>
                    </v-card-text>
                </v-window-item>
                <v-window-item :value="2">
                    <v-card-text>
                        <v-alert v-if="textImportSuccess.length > 0" title="Success" type="success"></v-alert>
                        <v-alert v-if="textImportFail.length > 0" title="Fail" type="error"></v-alert>
                        <v-textarea :model-value="textImportSuccess.length > 0 ? textImportSuccess : textImportFail"
                            rows="5" no-resize readonly></v-textarea>
                    </v-card-text>
                </v-window-item>
            </v-window>
            <v-divider></v-divider>
            <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn v-if="stepImport == 2" text="Back" variant="text" @click="stepImport--"></v-btn>
                <v-btn v-if="stepImport == 1" text="Import" color="primary" variant="flat"
                    :disabled="importFile === null" :loading="importing" @click="Import()"></v-btn>
                <v-btn v-if="stepImport == 2" text="Done" variant="text" @click="showDlgImport = false"></v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>

    <!-- Dialog - Export list empty -->
    <v-dialog v-model="showDlgExportEmpty" width="auto">
        <v-card width="400" title="Export">
            <v-card-text>
                <span v-if="exportAll">No baseline and listor to be exported</span>
                <span v-else>No baseline selected</span>
            </v-card-text>
        </v-card>
    </v-dialog>

    <!-- Dialog - Export confirm -->
    <v-dialog v-model="showDlgExportConfirm" width="auto">
        <v-card width="400" title="Export">
            <v-card-text>
                <div>There is baseline with non-existent listor id.</div>
                <div>Export anyway?</div>
            </v-card-text>
            <template v-slot:actions>
                <v-spacer></v-spacer>
                <v-btn class="ms-auto" text="EXPORT" color="warning" @click="Export()"></v-btn>
                <v-btn class="ms-auto" text="Cancel" @click="showDlgExportConfirm = false"></v-btn>
            </template>
        </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useTheme } from 'vuetify'
import { storeToRefs } from 'pinia'
import yaml from 'js-yaml'

import { useUIStore, prepareUIStore } from '@/store/ui'
import { useConfStore, prepareConfStore, output_format_items, profile_key_items, Conf, emptyConf } from '@/store/conf';

const { themeName, exportAll, optionLocked, profileLocked, deleteBaselineWithListor } = storeToRefs(useUIStore())
const confStore = useConfStore()

// defined in https://www.iana.org/assignments/media-types/application/yaml
const YAML_MIME = 'application/yaml'
const theme = useTheme()

onMounted(() => {
    prepareUIStore()
    prepareConfStore()
    theme.global.name.value = themeName.value
})

const ToggleTheme = () => {
    themeName.value = theme.global.current.value.dark ? 'light' : 'dark'
    theme.global.name.value = themeName.value
}

const allMetadataKeys = computed(() => {
    const allMetadataKeys = confStore.conf.baseline.map(
        b => Object.keys(b.metadata))
        .reduce((prev, cur) => [...prev, ...cur], [])

    return [...new Set(allMetadataKeys)]
})

const searchForListor = ref("")
const searchForBaseline = ref("")
const selectedListor = ref([])
const selectedBaseline = ref([])

const allRefListorId = computed(() => {
    const allRefId = confStore.conf.baseline.map(b => b.checker)
        .reduce((prev, cur) => [...prev, ...cur], [])
        .map(c => c.listor)
        .reduce((prev, cur) => [...prev, ...cur], [])

    return [...new Set(allRefId)]
})

const IsListorDangling = (item: { id: number }) => !allRefListorId.value.includes(item.id)

const GetBaselineCheckerCloudType = (item: {
    checker: { cloud_type: string }[]
}) => {
    return [...new Set(item.checker.map(
        (val: { cloud_type: string; }) => val.cloud_type
    ))]
}

const GetBaselineCheckerListorId = (item: {
    checker: { listor: number[] }[]
}) => {
    const allId = item.checker.map((val: { listor: number[] }) => val.listor)
        .reduce((prev, cur) => [...prev, ...cur], [])

    return [...new Set(allId)]
}

const GetBaselineCheckerListorIdStr = (item: {
    checker: { listor: number[] }[]
}) => GetBaselineCheckerListorId(item).map(id => String(id))

const allListorId = computed(() => {
    const allId = confStore.conf.listor.map(l => l.id)
    return [...new Set(allId)]
})

const GetInvalidListorRefInBaseline = (item: {
    checker: { listor: number[] }[]
}) => {
    const refListId = GetBaselineCheckerListorId(item)
    return refListId.filter(id => !allListorId.value.includes(id))
}

const IsInvalidListorRefInBaseline = (item: {
    checker: { listor: number[] }[]
}) => GetInvalidListorRefInBaseline(item).length !== 0

// Dialog - Add profile
const showDlgAddProfile = ref(false)
const dlgAddProfileKey = ref('')
const profileKeyToAdd = computed(() => {
    const curKeys = Object.keys(confStore.conf.profile)
    return profile_key_items.filter(item => !(curKeys.includes(item)))
})
const canAddProfile = computed(() => profileKeyToAdd.value.length > 0)

const ShowAddProfile = () => {
    dlgAddProfileKey.value = ''
    showDlgAddProfile.value = true
}

const AddProfile = (key: keyof Conf['profile']) => {
    confStore.conf.profile[key] = ""
    showDlgAddProfile.value = false
}

// Dialog - Delete profile
const showDlgDeleteProfile = ref(false)
const dlgDeleteProfileKey = ref('')

const ShowDeleteProfile = (key: string) => {
    dlgDeleteProfileKey.value = key
    showDlgDeleteProfile.value = true
}

const DeleteProfile = (key: keyof Conf['profile']) => {
    delete confStore.conf.profile[key]
    showDlgDeleteProfile.value = false
}

// Dialog - Delete listor
const showDlgDeleteListor = ref(false)
const listorIdToDelete = ref([] as number[])

const ShowDeleteListor = (idToDelete: number[]) => {
    listorIdToDelete.value.splice(0, listorIdToDelete.value.length, ...idToDelete)
    showDlgDeleteListor.value = true
}

const DeleteListor = (idList: number[]) => {
    confStore.conf.listor = confStore.conf.listor.filter(
        l => !idList.includes(l.id)
    )
    selectedListor.value = selectedListor.value.filter(
        l => !idList.includes(l)
    )
    showDlgDeleteListor.value = false
}

// Dialog - Delete baseline
const showDlgDeleteBaseline = ref(false)
const baselineIdToDelete = ref([] as number[])

const ShowDeleteBaseline = (idToDelete: number[]) => {
    baselineIdToDelete.value.splice(0, baselineIdToDelete.value.length, ...idToDelete)
    showDlgDeleteBaseline.value = true
}

const DeleteBaseline = (idList: number[]) => {
    const idListorOfDeleteBaseline = [...new Set(
        confStore.conf.baseline.filter(l => idList.includes(l.id))
            .map(b => b.checker)
            .reduce((prev, cur) => [...prev, ...cur], [])
            .map(c => c.listor)
            .reduce((prev, cur) => [...prev, ...cur], [])
    )]
    confStore.conf.baseline = confStore.conf.baseline.filter(
        l => !idList.includes(l.id)
    )
    selectedBaseline.value = selectedBaseline.value.filter(
        l => !idList.includes(l)
    )

    if (deleteBaselineWithListor.value) {
        const idListorToDelete = idListorOfDeleteBaseline.filter(
            i => IsListorDangling({ id: i })
        )
        DeleteListor(idListorToDelete)
    }
    showDlgDeleteBaseline.value = false
}

// Dialog - Reset LocalStorage
const showDlgResetLocalStorage = ref(false)

const resetLocalStorage = () => {
    confStore.conf = JSON.parse(JSON.stringify(emptyConf)) // make a deep clone
    selectedListor.value.length = 0
    selectedBaseline.value.length = 0
    showDlgResetLocalStorage.value = false
}

// Dialog - Import
const showDlgImport = ref(false)
const stepImport = ref(1)
const importFile = ref(null as null | File)
const importing = ref(false)
const textImportSuccess = ref("")
const textImportFail = ref("")

const ShowImport = () => {
    stepImport.value = 1
    importing.value = false
    showDlgImport.value = true
}

const Import = () => {
    importing.value = true
    selectedListor.value.length = 0
    selectedBaseline.value.length = 0

    importFile.value?.text().then(s => {
        const doc = yaml.load(s)

        return new Promise<string[]>((resolve, reject) => {
            const importRes = ImportYamlDoc(doc)
            if (Array.isArray(importRes)) {
                resolve(importRes)
            } else {
                reject(importRes)
            }
        })
    }).then(r => {
        textImportSuccess.value = r.join('\n')
        textImportFail.value = ''
    }).catch(r => {
        textImportSuccess.value = ''
        textImportFail.value = String(r)
    }).finally(() => {
        importing.value = false
        stepImport.value = 2
    })
}

// return string[] as success; return string as fail
const ImportYamlDoc = (doc: any) => {
    if (typeof (doc) !== 'object') {
        return 'No data to import, please check the yaml configuration file'
    }

    const textSuccessList = [] as string[]

    // option
    const option = doc.option || undefined
    if (option === undefined) {
        textSuccessList.push('"option" not defined')
    } else if (typeof (option) !== 'object') {
        textSuccessList.push('Invalid "option" bypassed')
    } else if (optionLocked.value) {
        textSuccessList.push('"option" bypassed as it is locked')
    } else {
        for (const key in confStore.conf.option) {
            if ((option[key] || undefined) === undefined) {
                continue
            }
            if (typeof (Reflect.get(confStore.conf.option, key)) === typeof (option[key])) {
                Reflect.set(confStore.conf.option, key, option[key])
            } else {
                textSuccessList.push(`Invalid "option.${key}" bypassed`)
            }
        }
        textSuccessList.push('"option" imported')
    }

    // profile
    const profile = doc.profile || undefined
    if (profile === undefined) {
        textSuccessList.push('"profile" not defined')
    } else if (typeof (profile) !== 'object') {
        textSuccessList.push('Invalid "profile" bypassed')
    } else if (profileLocked.value) {
        textSuccessList.push('"profile" bypassed as it is locked')
    } else {
        for (const key of profile_key_items) {
            if ((profile[key] || undefined) === undefined) {
                continue
            }
            if (typeof (profile[key]) === "string") {
                confStore.conf.profile[key] = profile[key]
            } else {
                textSuccessList.push(`Invalid "profile.${key}" bypassed`)
            }
        }
        textSuccessList.push('"profile" imported')
    }

    // listor
    const idRemapped = {} as { [key: number]: number }
    const listor = doc.listor || undefined
    if (listor === undefined) {
        textSuccessList.push('"listor" not defined')
    } else if (!Array.isArray(listor)) {
        textSuccessList.push('Invalid "listor" bypassed')
    } else {
        const idAdded = [] as number[]
        const existListorId = confStore.conf.listor.map(l => l.id)
        const existListorRawConf = confStore.conf.listor.map(l => l.raw_conf)

        listor.forEach(v => {
            if (typeof (v) !== 'object') {
                return
            }

            let curId = undefined as undefined | number
            if (typeof (v.id) === 'number') {
                curId = v.id
            }
            delete v.id // remove id in raw conf
            const curRawConf = yaml.dump(v)

            let addId = curId
            if (curId === undefined || existListorId.includes(curId)) {
                // remap of id is needed
                let remapTarget = undefined as undefined | number

                // try to find the same listor which may have been imported and remapped the last time
                for (const confindex in existListorRawConf) {
                    if (existListorRawConf[confindex] === curRawConf) {
                        addId = undefined // donnot add to existing conf
                        remapTarget = confStore.conf.listor[confindex].id
                        break
                    }
                }

                if (remapTarget === undefined) {
                    // assign a new id
                    let newid = existListorId.length + listor.length + 1
                    while (existListorId.includes(newid)) {
                        newid++
                    }

                    remapTarget = newid
                    addId = remapTarget
                }

                if (curId !== undefined && curId !== remapTarget) {
                    idRemapped[curId] = remapTarget
                }
            }

            if (addId !== undefined) {
                // add to existing list
                const ct = v.cloud_type || ''
                const rt = v.rs_type || ''
                confStore.conf.listor.push({
                    id: addId,
                    cloud_type: typeof (ct) === 'string' ? ct : '',
                    rs_type: typeof (rt) === 'string' ? rt : '',
                    raw_conf: curRawConf,
                })
                existListorId.push(addId)
                existListorRawConf.push(curRawConf)
                idAdded.push(addId)
            }
        })

        textSuccessList.push(`${idAdded.length} listor added`)
        textSuccessList.push(`${Object.keys(idRemapped).length} listor remapped`)
    }

    // baseline
    const baseline = doc.baseline || undefined
    if (baseline === undefined) {
        textSuccessList.push('"baseline" not defined')
    } else if (!Array.isArray(baseline)) {
        textSuccessList.push('Invalid "baseline" bypassed')
    } else {
        const idAdded = [] as number[]
        let idDup = 0
        const existBaselineId = confStore.conf.baseline.map(b => b.id)
        const existBaselineRawConf = confStore.conf.baseline.map(b => b.raw_conf)

        baseline.forEach(v => {
            if (typeof (v) !== 'object') {
                return
            }

            const validChecker = [] as {
                cloud_type: string,
                listor: number[],
            }[]
            const checker = v.checker || undefined
            if (Array.isArray(checker)) {
                checker.forEach(c => {
                    const idListor = c.listor || undefined
                    const validListorId = [] as number[]
                    if (Array.isArray(idListor)) {
                        validListorId.push(...idListor.filter(id => typeof (id) === 'number'))
                        // remove id in raw conf
                        c.listor = idListor.filter(id => typeof (id) !== 'number')
                        if (c.listor.length === 0) {
                            delete c.listor
                        }
                    }

                    const cloud_type = c.cloud_type || undefined
                    if (cloud_type !== undefined || validListorId.length > 0) {
                        validChecker.push({
                            cloud_type: cloud_type || '',
                            listor: validListorId,
                        })
                    } else {
                        c.delete = true
                    }
                })
                // make sure the size of checker in conf is equal to the size of checker in yaml string of raw_conf
                v.checker = v.checker.filter((c: any) => c.delete !== true)
            } else {
                delete v.checker
            }

            // remap id of listor
            validChecker.forEach(c => {
                Object.entries(idRemapped).forEach(m => {
                    const from = Number(m[0])
                    const to = m[1]
                    for (const idindex in c.listor) {
                        if (c.listor[idindex] === from) {
                            c.listor[idindex] = to
                        }
                    }
                })
            })

            const curRawConf = yaml.dump(v)
            let addId = undefined as number | undefined

            for (const confindex in existBaselineRawConf) {
                if (existBaselineRawConf[confindex] === curRawConf) {
                    addId = confStore.conf.baseline[confindex].id

                    // merge listor id to existing baseline
                    const existingBaseline = confStore.conf.baseline[confindex]
                    if (validChecker.length != existingBaseline.checker.length) {
                        console.log('internal error, should not run to this line')
                        break
                    }

                    for (const checkerindex in existingBaseline.checker) {
                        const idListor = [...existingBaseline.checker[checkerindex].listor]
                        idListor.push(...validChecker[checkerindex].listor)
                        existingBaseline.checker[checkerindex].listor.length = 0
                        existingBaseline.checker[checkerindex].listor.push(...new Set(idListor))
                    }

                    break
                }
            }

            if (addId === undefined) {
                // // assign a new id and add to existing list
                let newid = existBaselineId.length + 1
                while (existBaselineId.includes(newid)) {
                    newid++
                }

                const tag = v.tag || ''
                const strTag = [] as string[]
                if (Array.isArray(tag)) {
                    strTag.push(...tag.filter(s => typeof (s) === 'string'))
                }
                const metadata = v.metadata || ''
                const objMetadata = {} as { [key: string]: string }
                if (typeof (metadata) === 'object') {
                    Object.entries(metadata).forEach(md => {
                        const mdkey = md[0]
                        const mdval = md[1]
                        if (typeof (mdval) === 'string') {
                            objMetadata[mdkey] = mdval
                        }
                    })
                }

                confStore.conf.baseline.push({
                    id: newid,
                    tag: strTag,
                    metadata: objMetadata,
                    checker: validChecker,
                    raw_conf: curRawConf,
                })
                existBaselineId.push(newid)
                existBaselineRawConf.push(curRawConf)
                idAdded.push(newid)
            } else {
                idDup++
            }
        })

        textSuccessList.push(`${idAdded.length} baseline added`)
        textSuccessList.push(`${idDup} duplicated baseline bypassed`)
    }

    return textSuccessList
}

// Dialog - Export
const showDlgExportEmpty = ref(false)
const showDlgExportConfirm = ref(false)
const exportBaseline = ref([] as number[])
const exportListor = ref([] as number[])
const exporting = ref(false)

const TryExport = () => {
    exportBaseline.value.length = 0
    exportListor.value.length = 0

    if (exportAll.value) {
        exportBaseline.value.push(...confStore.conf.baseline.map(b => b.id))
        exportListor.value.push(...confStore.conf.listor.map(l => l.id))
        if (exportBaseline.value.length + exportListor.value.length === 0) {
            showDlgExportEmpty.value = true
            return
        }
    } else {
        exportBaseline.value.push(...selectedBaseline.value)
        if (exportBaseline.value.length === 0) {
            showDlgExportEmpty.value = true
            return
        }

        const idListor = [] as number[]
        let needConfirm = false
        exportBaseline.value.forEach(idBaseline => {
            const baseline = confStore.conf.baseline.filter(b => b.id === idBaseline)[0]
            if (baseline === undefined) {
                console.log('internal error, should not run to this line')
                return
            }

            idListor.push(...GetBaselineCheckerListorId(baseline))

            if (IsInvalidListorRefInBaseline(baseline)) {
                needConfirm = true
            }
        })

        exportListor.value.push(...new Set(idListor))
        if (needConfirm) {
            showDlgExportConfirm.value = true
            return
        }
    }

    Export()
}

const Export = () => {
    showDlgExportConfirm.value = false
    exporting.value = true

    const listorData = confStore.conf.listor.filter(l => exportListor.value.includes(l.id))
        .map(l => {
            const data = yaml.load(l.raw_conf) as any
            data.id = l.id
            return data
        })
    const baselineData = confStore.conf.baseline.filter(b => exportBaseline.value.includes(b.id))
        .map(b => {
            const data = yaml.load(b.raw_conf) as any
            if (b.checker.length !== data.checker?.length) {
                console.log('internal error, should not run to this line')
                return data
            }

            for (const checkerindex in data.checker as any[]) {
                data.checker[checkerindex].listor = [...b.checker[checkerindex].listor]
            }

            return data
        })

    const exportData = {
        option: confStore.conf.option,
        profile: confStore.conf.profile,
        listor: listorData,
        baseline: baselineData,
    }

    const link = document.createElement('a')
    link.download = "baseline.conf"
    link.href = URL.createObjectURL(new Blob([yaml.dump(exportData)], { type: YAML_MIME }))
    link.click()
    URL.revokeObjectURL(link.href)

    exporting.value = false
}
</script>
