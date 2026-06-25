<script setup lang="ts">
    import { ref, onMounted } from 'vue';
    import { type GameAsset, parseGameAsset } from '@/types';

    // need run id and player id to get data from server
    // could be better to batch fetch but want to use many places
    const props = defineProps<{
        RunId: string,
        PlayerId: string,
    }>();

    let weaponData = ref<GameAsset[]>([]);

    const fetchWeaponData = () => {
        fetch(`/api/getItemOverview?PlayerId=${props.PlayerId}&RunId=${props.RunId}`).then(resp => {
            if (resp.ok) {
                return resp.json()
            }
        }).then(json => {
            if (Array.isArray(json)) {
                    weaponData.value = json.map(item => parseGameAsset(item));
                } else {
                    weaponData.value = [parseGameAsset(json)];
                }
        })
    }

    onMounted(() => {
        fetchWeaponData();
    });

    const assetPath = "/assets/"

</script>

<template>
    <div class="woc">
        <div class="equipped" v-if="weaponData.length != 0">
            <div v-for="wd in weaponData" class="asset">
                <img :src="assetPath + wd.Asset" width="100%" height="90%">
                <p>{{ wd.Name }}</p>
            </div>
        </div>
        <div v-else>
            <p>Loading</p>
        </div>
    </div>

</template>

<style>
    .equipped {
        display: flex;
        flex-direction: row;
        height: 100%;
        width: 30%;
        gap: 3%;
    }

    .asset {
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        border: solid thin grey; 
    }
    .woc {
        width: 65%;
    }
</style>