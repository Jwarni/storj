// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

<template>
	<div id="app">
		<router-view/>
		<!-- Area for displaying notification -->
		<NotificationArea />
	</div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import NotificationArea from '@/components/notifications/NotificationArea.vue';

@Component({
	components: {
		NotificationArea
	}
})

export default class App extends Vue {}
</script>

<style lang="scss">

@font-face {
	font-family: "montserrat_regular";
	src: url("../static/fonts/montserrat_regular.ttf");
}
@font-face {
	font-family: "montserrat_medium";
	src: url("../static/fonts/montserrat_medium.ttf");
}
@font-face {
	font-family: "montserrat_bold";
	src: url("../static/fonts/montserrat_bold.ttf");
}

a{
  	cursor: pointer;
}

input,
textarea {
	font-family: inherit;
	font-weight: 600;
	border: 1px solid rgba(56, 75, 101, 0.4);
}

input:hover,
textarea:hover {
  	border-color: #737791 !important;
}
</style>
