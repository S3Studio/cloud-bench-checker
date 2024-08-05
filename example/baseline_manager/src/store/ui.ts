import { defineStore } from 'pinia'

const useUIStore = defineStore('ui', {
    state: () => {
        return {
            themeName: 'dark',
            exportAll: false,
            optionLocked: false,
            profileLocked: false,
            deleteBaselineWithListor: false,
        }
    }
})

const prepareUIStore = () => {
    const key = 'ui-store'
    const instance = useUIStore()
    instance.$subscribe((_, state) => {
        localStorage.setItem(key, JSON.stringify({ ...state }))
    })
    const old = localStorage.getItem(key)
    if (old) {
        instance.$state = JSON.parse(old)
    }
}

export {useUIStore, prepareUIStore}
