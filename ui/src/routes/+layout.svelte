<script lang="ts">
	import { init, keycloak } from '$lib/keycloak';
	import { profile } from '$lib/stores';
	import { onMount } from 'svelte';

	onMount(async () => await init());
</script>

<nav>
	<a href="/">Home</a>
	<a href="/about">About</a>
	<a href="/graphql">GraphQL</a>
</nav>
{#if $profile}
	<button on:click={() => keycloak.logout()}>logout</button>
	<slot />
{:else}
	<div>not authorized</div>
{/if}
