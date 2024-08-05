import { defineStore } from 'pinia'

export const output_format_items = ['csv', 'json']
export const profile_key_items = ['tencent', 'aliyun', 'k8s', 'azure'].sort()
export const cloud_type_items = ['tencent_cloud', 'tencent_cos', 'aliyun', 'aliyun_oss', 'k8s', 'azure'].sort()

export type Conf = {
    option: {
        output_format: typeof output_format_items[number]
        output_filename: string,
        output_metadata: string[],
        output_risk_only: boolean,
    },
    profile: {
        [key: typeof profile_key_items[number]]: string,
    },
    listor: {
        id: number,
        cloud_type: typeof cloud_type_items[number],
        rs_type: string,
        raw_conf: string,
    }[],
    baseline: {
        id: number,
        tag: string[],
        metadata: { [key: string]: string },
        checker: {
            cloud_type: typeof cloud_type_items[number],
            listor: number[],
        }[],
        raw_conf: string,
    }[],
}

export const emptyConf: Conf = {
    option: {
        output_format: 'csv',
        output_filename: 'test',
        output_metadata: [],
        output_risk_only: true,
    },
    profile: {},
    listor: [],
    baseline: [],
}

const useConfStore = defineStore('conf', {
    state: () => {
        return { conf: JSON.parse(JSON.stringify(emptyConf)) as Conf } // make a deep clone
    }
})

const prepareConfStore = () => {
    const key = 'conf-store'
    const instance = useConfStore()
    instance.$subscribe((_, state) => {
        localStorage.setItem(key, JSON.stringify({ ...state }))
    })
    const old = localStorage.getItem(key)
    if (old) {
        instance.$state = JSON.parse(old)
    }
}

export { useConfStore, prepareConfStore }
