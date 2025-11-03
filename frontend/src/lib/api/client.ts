const API_BASE = import.meta.env.DEV 
	? 'http://localhost:8080/api'
	: '/api';

export async function fetchActivities() {
	const res = await fetch(`${API_BASE}/activities`);
	if (!res.ok) throw new Error('Failed to fetch activities');
	return res.json();
}

export async function fetchTournaments() {
	const res = await fetch(`${API_BASE}/tournaments`);
	if (!res.ok) throw new Error('Failed to fetch tournaments');
	return res.json();
}

export async function fetchArtwork() {
	const res = await fetch(`${API_BASE}/art/random`);
	if (!res.ok) throw new Error('Failed to fetch artwork');
	return res.json();
}

export async function createActivity(activity: any) {
	const res = await fetch(`${API_BASE}/activities`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(activity)
	});
	if (!res.ok) throw new Error('Failed to create activity');
	return res.json();
}
