<script lang="ts">
	import { onMount } from 'svelte';
	import { fetchActivities, fetchTournaments, fetchArtwork } from '$lib/api/client';

	let activities: any[] = [];
	let tournaments: any[] = [];
	let artwork: any = null;
	let loading = true;

	onMount(async () => {
		try {
			const [activitiesData, tournamentsData, artworkData] = await Promise.all([
				fetchActivities(),
				fetchTournaments(),
				fetchArtwork()
			]);

			activities = activitiesData;
			tournaments = tournamentsData;
			artwork = artworkData;
		} catch (error) {
			console.error('Failed to load data:', error);
		} finally {
			loading = false;
		}
	});
</script>

<main class="max-w-7xl mx-auto px-8 py-8">
	<header class="text-center mb-12">
		<h1 class="text-4xl font-bold mb-2">Aperture Science Data Portal</h1>
		<p class="text-lg opacity-75 tagline">The Portal to Your Daily Progress</p>
	</header>

	{#if loading}
		<div class="text-center text-xl py-16 loading">Loading portal systems...</div>
	{:else}
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
			<!-- Activities Section -->
			<section class="bg-white border border-black p-8 shadow-sm card">
				<h2 class="text-2xl font-bold mb-6 border-b border-black pb-2">Recent Activities</h2>
				{#if activities.length > 0}
					<ul class="list-none p-0 m-0 activity-list">
						{#each activities as activity}
							<li class="p-4 border-b border-gray-300 flex flex-wrap gap-4">
								<span class="bg-portal-blue text-white px-3 py-1 text-sm activity-type">{activity.activity_type}</span>
								<span class="activity-distance">{(activity.distance_meters / 1000).toFixed(2)} km</span>
								<span class="activity-date">{new Date(activity.activity_date).toLocaleDateString()}</span>
							</li>
						{/each}
					</ul>
				{:else}
					<p>No activities found</p>
				{/if}
			</section>

			<!-- Tournaments Section -->
			<section class="bg-white border border-black p-8 shadow-sm card">
				<h2 class="text-2xl font-bold mb-6 border-b border-black pb-2">Tournament History</h2>
				{#if tournaments.length > 0}
					<ul class="list-none p-0 m-0 tournament-list">
						{#each tournaments as tournament}
							<li class="p-4 border-b border-gray-300 flex flex-wrap gap-4">
								<strong>{tournament.tournament_name}</strong>
								<span class="font-bold text-portal-blue placement">#{tournament.placement}</span>
								<span class="bg-portal-blue text-white px-3 py-1 text-sm game">{tournament.game}</span>
								<span class="date">{new Date(tournament.tournament_date).toLocaleDateString()}</span>
							</li>
						{/each}
					</ul>
				{:else}
					<p>No tournaments found</p>
				{/if}
			</section>

			<!-- Art Section -->
			<section class="bg-white border border-black p-8 shadow-sm card art-card">
				<h2 class="text-2xl font-bold mb-6 border-b border-black pb-2">Featured Artwork</h2>
				{#if artwork}
					<div class="text-center artwork">
						<h3 class="text-xl mb-2">{artwork.title}</h3>
						<p class="italic text-gray-600 artist">{artwork.artist}</p>
						<p class="text-gray-500 text-sm date">{artwork.date_display}</p>
					</div>
				{:else}
					<p>No artwork available</p>
				{/if}
			</section>
		</div>
	{/if}
</main>


