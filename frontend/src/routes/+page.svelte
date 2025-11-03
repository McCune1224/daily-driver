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

<main class="container">
	<header>
		<h1>Aperture Science Data Portal</h1>
		<p class="tagline">The Portal to Your Daily Progress</p>
	</header>

	{#if loading}
		<div class="loading">Loading portal systems...</div>
	{:else}
		<div class="grid">
			<!-- Activities Section -->
			<section class="card">
				<h2>Recent Activities</h2>
				{#if activities.length > 0}
					<ul class="activity-list">
						{#each activities as activity}
							<li>
								<span class="activity-type">{activity.activity_type}</span>
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
			<section class="card">
				<h2>Tournament History</h2>
				{#if tournaments.length > 0}
					<ul class="tournament-list">
						{#each tournaments as tournament}
							<li>
								<strong>{tournament.tournament_name}</strong>
								<span class="placement">#{tournament.placement}</span>
								<span class="game">{tournament.game}</span>
								<span class="date">{new Date(tournament.tournament_date).toLocaleDateString()}</span>
							</li>
						{/each}
					</ul>
				{:else}
					<p>No tournaments found</p>
				{/if}
			</section>

			<!-- Art Section -->
			<section class="card art-card">
				<h2>Featured Artwork</h2>
				{#if artwork}
					<div class="artwork">
						<h3>{artwork.title}</h3>
						<p class="artist">{artwork.artist}</p>
						<p class="date">{artwork.date_display}</p>
					</div>
				{:else}
					<p>No artwork available</p>
				{/if}
			</section>
		</div>
	{/if}
</main>

<style>
	:global(body) {
		font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		margin: 0;
		padding: 0;
		min-height: 100vh;
	}

	.container {
		max-width: 1200px;
		margin: 0 auto;
		padding: 2rem;
	}

	header {
		text-align: center;
		color: white;
		margin-bottom: 3rem;
	}

	h1 {
		font-size: 3rem;
		margin-bottom: 0.5rem;
		text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3);
	}

	.tagline {
		font-size: 1.2rem;
		opacity: 0.9;
	}

	.loading {
		text-align: center;
		color: white;
		font-size: 1.5rem;
		padding: 4rem;
	}

	.grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
		gap: 2rem;
	}

	.card {
		background: white;
		border-radius: 12px;
		padding: 2rem;
		box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
		transition: transform 0.3s ease;
	}

	.card:hover {
		transform: translateY(-5px);
	}

	h2 {
		color: #667eea;
		margin-top: 0;
		margin-bottom: 1.5rem;
		border-bottom: 2px solid #667eea;
		padding-bottom: 0.5rem;
	}

	.activity-list, .tournament-list {
		list-style: none;
		padding: 0;
		margin: 0;
	}

	.activity-list li, .tournament-list li {
		padding: 1rem;
		border-bottom: 1px solid #eee;
		display: flex;
		flex-wrap: wrap;
		gap: 1rem;
	}

	.activity-list li:last-child, .tournament-list li:last-child {
		border-bottom: none;
	}

	.activity-type, .game {
		background: #667eea;
		color: white;
		padding: 0.25rem 0.75rem;
		border-radius: 20px;
		font-size: 0.9rem;
	}

	.placement {
		font-weight: bold;
		color: #764ba2;
	}

	.artwork {
		text-align: center;
	}

	.artwork h3 {
		color: #333;
		margin-bottom: 0.5rem;
	}

	.artist {
		font-style: italic;
		color: #666;
	}

	.date {
		color: #999;
		font-size: 0.9rem;
	}
</style>
