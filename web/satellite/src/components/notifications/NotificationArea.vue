// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

<template>
    <div id="notificationArea" class="notification-container" v-if="currentNotification" >
        <Notification :type="currentNotification.type" :message="currentNotification.message" />
    </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import Notification from '@/components/notifications/Notification.vue';
import { NOTIFICATION_TYPES } from '@/utils/constants/notification';

@Component({
    computed: {
        // Computed value for current notification depends on store.state value
        currentNotification: function(): Notification {
            return this.$store.getters.currentNotification;
        }
    },
    components: {
        Notification,
    }
})

export default class NotificationArea extends Vue {}
</script>

<style scoped lang="scss">
    .notification-container {
        height: 98px;
        max-width: 80%;
        width: 100%;
        background-color: #fff;
        display: flex;
        position: fixed;
        bottom: 50px;
        right: 30px;
        align-items: center;
        justify-content: space-between;
        box-shadow: 0px 12px 24px rgba(175, 183, 193, 0.4);
        border-radius: 6px;
        z-index: 999;
    }
</style>
